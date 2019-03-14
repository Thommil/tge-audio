// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build js

package audio

import (
	fmt "fmt"
	js "syscall/js"

	tge "github.com/thommil/tge"
)

type plugin struct {
	isInit   bool
	jsTge    *js.Value
	audioCtx *js.Value
}

func (p *plugin) Init(runtime tge.Runtime) error {
	switch host := runtime.GetHost().(type) {
	case *js.Value:
		jsTge := host.Get("tge")
		if jsTge == js.Undefined() || jsTge == js.Null() {
			return fmt.Errorf("javascript tge not found")
		}
		p.jsTge = &jsTge
	default:
		return fmt.Errorf("Runtime host must be a *syscall/js.Value")
	}
	return nil
}

func (p *plugin) Dispose() {
	p.isInit = false
	p.audioCtx = nil
}

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

type buffer struct {
	value *js.Value
}

type node struct {
	value *js.Value
}

func (n *node) Connect(to Node) Node {
	switch to.(type) {
	case *bufferSourceNode:
		n.value.Call("connect", *(to.(*bufferSourceNode).value))
	case *mediaElementSourceNode:
		n.value.Call("connect", *(to.(*mediaElementSourceNode).value))
	case *destinationNode:
		n.value.Call("connect", *(to.(*destinationNode).value))
	case *stereoPannerNode:
		n.value.Call("connect", *(to.(*stereoPannerNode).value))
	case *gainNode:
		n.value.Call("connect", *(to.(*gainNode).value))
	}
	return to
}

func (n *node) Disconnect(to Node) {
	switch to.(type) {
	case *bufferSourceNode:
		n.value.Call("disconnect", *(to.(*bufferSourceNode).value))
	case *mediaElementSourceNode:
		n.value.Call("disconnect", *(to.(*mediaElementSourceNode).value))
	case *destinationNode:
		n.value.Call("disconnect", *(to.(*destinationNode).value))
	case *stereoPannerNode:
		n.value.Call("disconnect", *(to.(*stereoPannerNode).value))
	case *gainNode:
		n.value.Call("disconnect", *(to.(*gainNode).value))
	}
}

type bufferSourceNode struct {
	node
}

func (n *bufferSourceNode) Start(offset, duration float32, loop bool) {
	n.value.Set("loop", loop)
	n.value.Call("start", 0, offset, duration)
}

func (n *bufferSourceNode) Stop() {
	n.value.Call("stop")
}

type mediaElementSourceNode struct {
	node
}

type destinationNode struct {
	node
}

type stereoPannerNode struct {
	node
}

func (n *stereoPannerNode) Pan(value float32) {
	n.value.Get("pan").Set("value", value)
}

type gainNode struct {
	node
}

func (n *gainNode) Gain(value float32) {
	n.value.Get("gain").Set("value", value)
}

func createContext() error {
	audioContextClass := js.Global().Get("AudioContext")

	if audioContextClass == js.Undefined() || audioContextClass == js.Null() {
		audioContextClass = js.Global().Get("webkitAudioContext")
	}
	if audioContextClass == js.Undefined() || audioContextClass == js.Null() {
		return fmt.Errorf("AudioContext not supported")
	}

	audioCtx := audioContextClass.New()

	if audioCtx == js.Undefined() || audioCtx == js.Null() {
		return fmt.Errorf("failed to instanciate AudioContext")
	}

	_pluginInstance.audioCtx = &audioCtx
	_pluginInstance.isInit = true

	return nil
}

func createBuffer(path string) (Buffer, error) {
	if !_pluginInstance.isInit {
		if err := createContext(); err != nil {
			return nil, err
		}
	}
	var err error
	doneState := make(chan bool)
	buffer := buffer{}

	onCreateAudioBufferCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[1] != js.Null() {
			err = fmt.Errorf(args[1].String())
			doneState <- false
		} else {
			buffer.value = &args[0]
			doneState <- true
		}
		return false
	})
	defer onCreateAudioBufferCallback.Release()

	_pluginInstance.jsTge.Call("createAudioBuffer", *(_pluginInstance.audioCtx), path, onCreateAudioBufferCallback)

	<-doneState

	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

func createBufferSourceNode(buf Buffer) (BufferSourceNode, error) {
	if !_pluginInstance.isInit {
		if err := createContext(); err != nil {
			return nil, err
		}
	}

	jsBufferSourceNode := _pluginInstance.audioCtx.Call("createBufferSource")

	if jsBufferSourceNode == js.Undefined() || jsBufferSourceNode == js.Null() {
		return nil, fmt.Errorf("failed to create JS AudioBufferSourceNode")
	}

	jsBufferSourceNode.Set("buffer", *(buf.(*buffer).value))

	node := &bufferSourceNode{}
	node.value = &jsBufferSourceNode

	return node, nil
}

func createDestinationNode() (DestinationNode, error) {
	if !_pluginInstance.isInit {
		if err := createContext(); err != nil {
			return nil, err
		}
	}

	jsDestinationNode := _pluginInstance.audioCtx.Get("destination")

	if jsDestinationNode == js.Undefined() || jsDestinationNode == js.Null() {
		return nil, fmt.Errorf("failed to create JS AudioDestinationNode")
	}

	node := &destinationNode{}
	node.value = &jsDestinationNode

	return node, nil
}

func createStereoPannerNode() (StereoPannerNode, error) {
	if !_pluginInstance.isInit {
		if err := createContext(); err != nil {
			return nil, err
		}
	}

	jsStereoPannerNode := _pluginInstance.audioCtx.Call("createStereoPanner")

	if jsStereoPannerNode == js.Undefined() || jsStereoPannerNode == js.Null() {
		return nil, fmt.Errorf("failed to create JS StereoPannerNode")
	}

	node := &stereoPannerNode{}
	node.value = &jsStereoPannerNode

	return node, nil
}

func createGainNode() (GainNode, error) {
	if !_pluginInstance.isInit {
		if err := createContext(); err != nil {
			return nil, err
		}
	}

	jsGainNode := _pluginInstance.audioCtx.Call("createGain")

	if jsGainNode == js.Undefined() || jsGainNode == js.Null() {
		return nil, fmt.Errorf("failed to create JS GainNode")
	}

	node := &gainNode{}
	node.value = &jsGainNode

	return node, nil
}
