// Copyright (c) 2019 Thomas MILLET. All rights reserved.
// Copyright (C) 1991 Free Software Foundation, Inc.

// +build android ios

package audio

import (
	fmt "fmt"

	tge "github.com/thommil/tge"
	al "github.com/thommil/tge-mobile/exp/audio/al"
)

type plugin struct {
	isInit bool
}

func (p *plugin) Init(r tge.Runtime) error {
	if !p.isInit {
		return al.OpenDevice()
	}
	return fmt.Errorf("Already initialized")
}

func (p *plugin) Dispose() {
	if p.isInit {
		al.CloseDevice()
	}
}

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

func createBuffer(path string) (Buffer, error) {
	return nil, nil
}

func createBufferSourceNode(buffer Buffer) (BufferSourceNode, error) {
	return nil, nil
}

func createMediaElementSourceNode(path string) (MediaElementSourceNode, error) {
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
