// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build android ios

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
