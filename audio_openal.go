// Copyright (c) 2019 Thomas MILLET. All rights reserved.
// Copyright (C) 1991 Free Software Foundation, Inc.

// +build !js

package audio

import (
	"encoding/binary"
	fmt "fmt"
	filepath "path/filepath"
	time "time"
	unsafe "unsafe"

	tge "github.com/thommil/tge"
	al "github.com/thommil/tge-mobile/exp/audio/al"

	vorbis "github.com/mccoyst/vorbis"
)

type plugin struct {
	isInit  bool
	runtime tge.Runtime
}

var nativeEndian binary.ByteOrder

func (p *plugin) Init(runtime tge.Runtime) error {
	if !p.isInit {
		p.runtime = runtime

		buf := [2]byte{}
		*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)
		switch buf {
		case [2]byte{0xAB, 0xCD}:
			nativeEndian = binary.BigEndian
		default:
			nativeEndian = binary.LittleEndian
		}

		err := al.OpenDevice()
		if err != nil {
			return err
		}
		al.SetListenerPosition(al.Vector{0, 0, 0})
		al.SetListenerOrientation(al.Orientation{Forward: al.Vector{0, 0, -1}, Up: al.Vector{0, 1, 0}})
		sources := al.GenSources(sourcePoolSize)
		for _, source := range sources {
			source.SetMaxGain(1.0)
			source.SetGain(0)
			b := bufferSourceNode{
				node: node{
					sources: []*sourceProxy{&sourceProxy{
						handle:    source,
						connected: false,
						gain:      1,
					}},
					to: make([]connectListener, 0, 1),
				},
			}
			sourcePool <- &b
		}
		return nil
	}
	return fmt.Errorf("Already initialized")
}

func (p *plugin) Dispose() {
	if p.isInit {
		al.CloseDevice()
		p.runtime = nil
	}
}

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

const sourcePoolSize = 100

var sourcePool = make(chan *bufferSourceNode, sourcePoolSize)

var destinationNodeSingleton = destinationNode{
	node: node{
		sources: make([]*sourceProxy, 0, sourcePoolSize),
		to:      nil,
	},
}

// GRAPH

type connectListener interface {
	onConnectStateChanged(connected bool, sources []*sourceProxy)
}

type node struct {
	sources []*sourceProxy
	to      []connectListener
}

func (n *node) Connect(to Node) Node {
	n.to = append(n.to, to.(connectListener))
	to.(connectListener).onConnectStateChanged(true, n.sources)
	return to
}

func (n *node) Disconnect(to Node) {
	for i, currentToNode := range n.to {
		if currentToNode == to.(connectListener) {
			n.to = append(n.to[:i], n.to[i+1:]...)
		}
	}
	to.(connectListener).onConnectStateChanged(false, n.sources)
}

func (n *node) onConnectStateChanged(connected bool, sources []*sourceProxy) {
	if connected {
		n.sources = append(n.sources, sources...)
	} else {
		for _, newSource := range sources {
			for i, currentSource := range n.sources {
				if newSource == currentSource {
					n.sources = append(n.sources[:i], n.sources[i+1:]...)
				}
			}
		}
	}
	for _, subNode := range n.to {
		subNode.onConnectStateChanged(connected, sources)
	}
}

// Buffer

type buffer struct {
	handle   al.Buffer
	duration time.Duration
}

func (b *buffer) Delete() {
	al.DeleteBuffers(b.handle)
}

// Source

type sourceProxy struct {
	handle    al.Source
	connected bool
	gain      float32
}

// BufferSourceNode

type bufferSourceNode struct {
	node
	buffer *buffer
}

func (n *bufferSourceNode) Start(delay, offset, duration float32, loop bool, loopStart, loopEnd float32) {
	if delay > 0 {
		delayDuration := time.Duration(delay * 1000000000)
		go func() {
			<-time.After(delayDuration)
			if n.sources[0].connected {
				n.sources[0].handle.SetGain(n.sources[0].gain)
			}
			al.PlaySources(n.sources[0].handle)
			n.sources[0].handle.Setf(0x1024, offset) // OFFSET
			b := *(n.buffer)
			durationDuration := b.duration
			if duration > 0 {
				durationDuration = time.Duration(duration * 1000000000)
			}
			<-time.After(durationDuration)
			if loop {
				if loopEnd == 0 {
					loopEnd = float32(b.duration.Nanoseconds()) / 1000000000
				}
				for true {
					al.PlaySources(n.sources[0].handle)
					n.sources[0].handle.Setf(0x1024, loopStart) // OFFSET
					<-time.After(time.Duration((loopEnd - loopStart) * 1000000000))
				}
			} else {
				n.Stop()
			}
		}()
	} else {
		if n.sources[0].connected {
			n.sources[0].handle.SetGain(n.sources[0].gain)
		}
		al.PlaySources(n.sources[0].handle)
		n.sources[0].handle.Setf(0x1024, offset) // OFFSET
		go func() {
			b := *(n.buffer)
			durationDuration := b.duration
			if duration > 0 {
				durationDuration = time.Duration(duration * 1000000000)
			}
			<-time.After(durationDuration)
			if loop {
				if loopEnd == 0 {
					loopEnd = float32(b.duration.Nanoseconds()) / 1000000000
				}
				for true {
					al.PlaySources(n.sources[0].handle)
					n.sources[0].handle.Setf(0x1024, loopStart) // OFFSET
					<-time.After(time.Duration((loopEnd - loopStart) * 1000000000))
				}
			} else {
				n.Stop()
			}
		}()
	}
}

func (n *bufferSourceNode) Stop() {
	al.StopSources(n.sources[0].handle)
	if n.buffer != nil {
		b := *(n.buffer)
		n.sources[0].gain = 1
		n.sources[0].handle.SetGain(0)
		n.sources[0].handle.UnqueueBuffers(b.handle)
		n.node.to = n.node.to[:0]
		n.buffer = nil
		sourcePool <- n
	}
}

func (n *bufferSourceNode) Play(loop bool) {
	if n.sources[0].connected {
		n.sources[0].handle.SetGain(n.sources[0].gain)
	}
	al.PlaySources(n.sources[0].handle)
	if loop {
		n.sources[0].handle.Seti(0x1007, 1) //LOOP
	}

}

func (n *bufferSourceNode) Pause() {
	al.PauseSources(n.sources[0].handle)
}

func (n *bufferSourceNode) Delete() {
	al.DeleteBuffers(n.buffer.handle)
}

// DestinationNode

type destinationNode struct {
	node
}

func (n *destinationNode) onConnectStateChanged(connected bool, sources []*sourceProxy) {
	for _, source := range sources {
		if connected {
			source.connected = true
			if source.handle.State() == al.Playing {
				source.handle.SetGain(source.gain)
			}
		} else {
			source.connected = false
			source.handle.SetGain(0)
		}
	}
}

// StereoPannerNode

type stereoPannerNode struct {
	node
	pan float32
}

func (n *stereoPannerNode) Pan(value float32) {
	n.pan = value
	for _, source := range n.sources {
		source.handle.SetPosition(al.Vector{n.pan, 0, 0})
	}
}

func (n *stereoPannerNode) onConnectStateChanged(connected bool, sources []*sourceProxy) {
	for _, source := range sources {
		source.handle.SetPosition(al.Vector{n.pan, 0, 0})
	}
	n.node.onConnectStateChanged(connected, sources)
}

// GainNode

type gainNode struct {
	node
	gain float32
}

func (n *gainNode) Gain(value float32) {
	n.gain = value
	for _, source := range n.sources {
		source.gain = n.gain
		if source.connected {
			source.handle.SetGain(source.gain)
		}
	}
}

func (n *gainNode) onConnectStateChanged(connected bool, sources []*sourceProxy) {
	for _, source := range sources {
		source.gain = n.gain
	}
	n.node.onConnectStateChanged(connected, sources)
}

// Factories

func createBuffer(path string) (Buffer, error) {
	alBuffer, duration, err := alBufferFromPath(path)
	if err != nil {
		return nil, err
	}

	outBuffer := buffer{
		handle:   alBuffer,
		duration: duration,
	}

	return &outBuffer, nil
}

func createBufferSourceNode(b Buffer) (BufferSourceNode, error) {
	audioBuffer := b.(*buffer)
	bufferSource := <-sourcePool
	bufferSource.buffer = audioBuffer
	bufferSource.node.sources[0].handle.QueueBuffers(audioBuffer.handle)
	return bufferSource, nil
}

func createMediaElementSourceNodePath(path string) (MediaElementSourceNode, error) {
	alBuffer, duration, err := alBufferFromPath(path)
	if err != nil {
		return nil, err
	}

	audioBuffer := buffer{
		handle:   alBuffer,
		duration: duration,
	}

	return createMediaElementSourceNodeBuffer(&audioBuffer)
}

func createMediaElementSourceNodeBuffer(b Buffer) (MediaElementSourceNode, error) {
	audioBuffer := b.(*buffer)
	bufferSource := <-sourcePool
	bufferSource.buffer = audioBuffer
	bufferSource.node.sources[0].handle.QueueBuffers(audioBuffer.handle)
	return bufferSource, nil
}

func createMediaElementSourceNode(path string) (MediaElementSourceNode, error) {
	return createMediaElementSourceNodePath(path)
}

func createDestinationNode() (DestinationNode, error) {
	return &destinationNodeSingleton, nil
}

func createStereoPannerNode() (StereoPannerNode, error) {
	return &stereoPannerNode{
		pan: 0,
		node: node{
			sources: make([]*sourceProxy, 0, sourcePoolSize),
			to:      make([]connectListener, 0, 1),
		},
	}, nil
}

func createGainNode() (GainNode, error) {
	return &gainNode{
		gain: 1,
		node: node{
			sources: make([]*sourceProxy, 0, sourcePoolSize),
			to:      make([]connectListener, 0, 1),
		},
	}, nil
}

// -------------------------------------------------------------------- //
// Tooling
// -------------------------------------------------------------------- //

// Byte buffer array singleton, allocate 1mB at startup
var byteArrayBuffer = make([]byte, 0)
var byteArrayBufferExtendFactor = 1

func getByteArrayBuffer(size int) []byte {
	if size > len(byteArrayBuffer) {
		for (1024 * 1024 * byteArrayBufferExtendFactor) < size {
			byteArrayBufferExtendFactor++
		}
		byteArrayBuffer = make([]byte, (1024 * 1024 * byteArrayBufferExtendFactor))
	}
	return byteArrayBuffer[:size]
}

func int16ToBytes(values []int16) []byte {
	b := getByteArrayBuffer(2 * len(values))
	if nativeEndian == binary.LittleEndian {
		for i, v := range values {
			u := *(*int16)(unsafe.Pointer(&v))
			b[2*i+0] = byte(u)
			b[2*i+1] = byte(u >> 8)
		}
	} else {
		for i, v := range values {
			u := *(*int16)(unsafe.Pointer(&v))
			b[2*i+0] = byte(u >> 8)
			b[2*i+1] = byte(u)
		}
	}
	return b
}

func alBufferFromPath(path string) (al.Buffer, time.Duration, error) {
	fileExt := filepath.Ext(path)
	// Read, Decode and upload data
	switch fileExt {
	case ".ogg":
		// Read
		inData, err := _pluginInstance.runtime.GetAsset(path)
		if err != nil {
			return 0, 0, err
		}
		// Decode
		data, channels, sampleRate, err := vorbis.Decode(inData)
		if err != nil {
			return 0, 0, err
		}
		// Gen AL buffer
		alBuffer := al.GenBuffers(1)[0]
		format := al.FormatStereo16
		if channels == 1 {
			format = al.FormatMono16
		}
		// Upload
		alBuffer.BufferData(uint32(format), int16ToBytes(data), int32(sampleRate))

		return alBuffer, time.Duration((float32(len(data)) / float32(channels) / float32(sampleRate)) * 1000000000), nil
	default:
		return 0, 0, fmt.Errorf("audio extension not supported %s", fileExt)
	}
}
