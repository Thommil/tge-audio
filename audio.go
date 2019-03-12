// Copyright (c) 2019 Thomas MILLET. All rights reserved.

package audio

import (
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
