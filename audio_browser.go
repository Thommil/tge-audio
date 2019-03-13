// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build js

package audio

import (
	fmt "fmt"

	tge "github.com/thommil/tge"
)

type plugin struct {
	isInit bool
}

func (p *plugin) Init(runtime tge.Runtime) error {
	if !p.isInit {

		return nil
	}
	return fmt.Errorf("Already initialized")
}

func (p *plugin) Dispose() {
	p.isInit = false
}

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

func createBuffer(data []byte) (Buffer, error) {
	return nil, nil
}

func createBufferSourceNode(buffer Buffer) (BufferSourceNode, error) {
	return nil, nil
}

func createDestinationNode() (DestinationNode, error) {
	return nil, nil
}

func createStereoPannerNode() (StereoPannerNode, error) {
	return nil, nil
}

func createGainNode() (GainNode, error) {
	return nil, nil
}
