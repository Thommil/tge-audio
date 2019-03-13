// Copyright (c) 2019 Thomas MILLET. All rights reserved.

package audio

import (
	fmt "fmt"

	tge "github.com/thommil/tge"
)

// -------------------------------------------------------------------- //
// Plugin definition
// -------------------------------------------------------------------- //

// Name name of the plugin
const Name = "audio"

var _pluginInstance = &plugin{}

func init() {
	tge.Register(_pluginInstance)
}

func (p *plugin) GetName() string {
	return Name
}

// -------------------------------------------------------------------- //
// API
// -------------------------------------------------------------------- //

// Buffer interface represents a short audio asset residing in memory, created from an audio file
// See https://developer.mozilla.org/en-US/docs/Web/API/AudioBuffer
type Buffer interface {
}

// Node interface is a generic interface for representing an audio processing module.
// See https://developer.mozilla.org/en-US/docs/Web/API/AudioNode
type Node interface {
	// Connect the node output to given node
	Connect(node Node) Node
	// Disconnect the node output from the given node
	Disconnect(node Node)
}

// BufferSourceNode interface represents an audio source consisting of in-memory audio data, stored in an AudioBuffer
// See https://developer.mozilla.org/en-US/docs/Web/API/AudioBufferSourceNode
type BufferSourceNode interface {
	Node
	// Start playing node, loopStart and loopEnd is given in seconds
	Start(loop bool, loopStart, loopEnd float32)
	// Stop playing node
	Stop()
}

// MediaElementSourceNode interface represents an external audio source for continuous play (music)
// See https://developer.mozilla.org/en-US/docs/Web/API/MediaElementAudioSourceNode
type MediaElementSourceNode interface {
	Node
}

// DestinationNode interface represents the end destination of an audio graph in a given context — usually the speakers of your device
// See https://developer.mozilla.org/en-US/docs/Web/API/AudioDestinationNode
type DestinationNode interface {
	Node
}

// StereoPannerNode interface represents a simple stereo panner node that can be used to pan an audio stream left or right
// See https://developer.mozilla.org/en-US/docs/Web/API/StereoPannerNode
type StereoPannerNode interface {
	Node
	// Pan the output from -1 left to 1 right
	Pan(value float32)
}

// GainNode interface represents a change in volume
// See https://developer.mozilla.org/en-US/docs/Web/API/GainNode
type GainNode interface {
	Node
	// Gain changes the output volume from 0 silence to 1 full
	Gain(value float32)
}

// CreateBuffer creates a Buffer from an assets path
func CreateBuffer(path string) (Buffer, error) {
	return createBuffer(path)
}

// CreateNode creates a new custom node, not implemented yet
func CreateNode() (Node, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// CreateBufferSourceNode creates a new BufferSourceNode, this method must be called each time you want to play a BufferSourceNode
func CreateBufferSourceNode(buffer Buffer) (BufferSourceNode, error) {
	return createBufferSourceNode(buffer)
}

// CreateMediaElementSourceNode creates a new MediaElementSourceNode from an assets path, not implemented yet
func CreateMediaElementSourceNode(path string) (MediaElementSourceNode, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// CreateDestinationNode creates a new DestinationNode to connect audio graph output
func CreateDestinationNode() (DestinationNode, error) {
	return createDestinationNode()
}

// CreateStereoPannerNode creates a new StereoPannerNode to connect in audio graph
func CreateStereoPannerNode() (StereoPannerNode, error) {
	return createStereoPannerNode()
}

// CreateGainNode creates a new GainNode to connect in audio graph
func CreateGainNode() (GainNode, error) {
	return createGainNode()
}
