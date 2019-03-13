<h1 align="center">TGE-AUDIO - Audio plugin for TGE</h1>

 <p align="center">
    <a href="https://godoc.org/github.com/thommil/tge-audio"><img src="https://godoc.org/github.com/thommil/tge-audio?status.svg" alt="Godoc"></img></a>
    <a href="https://goreportcard.com/report/github.com/thommil/tge-audio"><img src="https://goreportcard.com/badge/github.com/thommil/tge-audio"  alt="Go Report Card"/></a>
</p>

Audio support for TGE runtime - [TGE](https://github.com/thommil/tge)

Based on :
 * [gomobile exp/audio/al](https://godoc.org/golang.org/x/mobile/exp/audio/al) - Copyright (c) 2009 The Go Authors. All rights reserved.
 * [OpenAL](https://github.com/golang/mobile) - Copyright (C) 1991 Free Software Foundation, Inc.

The Audio API is based on [WebAudio API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Audio_API) specifications ported to OpenAL for 
Desktop and Mobile targets.

## Targets
 * Desktop (Work In Progress)
 * Mobile  (Work In Progress)
 * Browser

## Dependencies
 * [TGE core](https://github.com/thommil/tge)

### Desktop & Mobile

In order to use this package on Linux desktop distros, you may need OpenAL library as an external dependency:

```
sudo apt-get install libopenal-dev
```

## Limitations
Only [StereoPannerNode](https://developer.mozilla.org/en-US/docs/Web/API/StereoPannerNode) and [GainNode](https://developer.mozilla.org/en-US/docs/Web/API/GainNode) are currently available, additional and custom [AudioNodes](https://developer.mozilla.org/en-US/docs/Web/API/AudioNode) and spacialization features are planned in future releases.

Music source is not optimized and should be improved in future releases.

Audio can't be started without user interaction on browsers, see details in [implementation](#implementation).

## Implementation
The [WebAudio API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Audio_API) offers a complete and comprehensive approach to audio implementation in 2D/3D scenes, tge-audio tries to benefit from it by exposing also OpenAL the same way.

<p align="center">
<img src="https://mdn.mozillademos.org/files/9713/WebAudioBasics.png"/>
</p>

See example at [AUDIO examples](https://github.com/Thommil/tge-examples/tree/master/plugins/tge-audio)

```golang
package main

import (
	"fmt"
	"time"

	tge "github.com/thommil/tge"
	audio "github.com/thommil/tge-audio"
)

type AudioApp struct {
	runtime         tge.Runtime
	audioInit       bool
	sampleData      []byte
	sampleBuffer    audio.Buffer
	stereoPanNode   audio.StereoPannerNode
	gainNode        audio.GainNode
	destinationNode audio.DestinationNode
	width           int32
	height          int32
}

func (app *AudioApp) OnCreate(settings *tge.Settings) error {
	fmt.Println("OnCreate()")
	return nil
}

func (app *AudioApp) OnStart(runtime tge.Runtime) error {
	fmt.Println("OnStart()")
	app.runtime = runtime

	runtime.Subscribe(tge.MouseEvent{}.Channel(), app.OnMouseEvent)
	runtime.Subscribe(tge.ResizeEvent{}.Channel(), app.OnResize)
	return nil
}

// InitAudio is start on mouse click only due to browsers restrictions
func (app *AudioApp) InitAudio() error {
	var err error

	//Buffer (should be done in a loading screen off course)
	fmt.Println("Loading bass.wav in buffer")
	if app.sampleBuffer, err = audio.CreateBuffer("bass.wav"); err != nil {
		return err
	}

	// Audio graph
	fmt.Println("Creating audio graph")
	if app.destinationNode, err = audio.CreateDestinationNode(); err != nil {
		return err
	}
	if app.gainNode, err = audio.CreateGainNode(); err != nil {
		return err
	}
	if app.stereoPanNode, err = audio.CreateStereoPannerNode(); err != nil {
		return err
	}
	// destination is connected at resume state only
	app.stereoPanNode.Connect(app.gainNode)
	app.audioInit = true

	app.gainNode.Connect(app.destinationNode)

	return nil
}

func (app *AudioApp) OnResume() {
	fmt.Println("OnResume()")
	// Open sound
	if app.audioInit {
		app.gainNode.Connect(app.destinationNode)
	}
}

func (app *AudioApp) OnResize(event tge.Event) bool {
	app.width = (event.(tge.ResizeEvent)).Width
	app.height = (event.(tge.ResizeEvent)).Height
	return false
}

func (app *AudioApp) OnMouseEvent(event tge.Event) bool {
	mouseEvent := event.(tge.MouseEvent)
	if mouseEvent.Type == tge.TypeDown && (mouseEvent.Button&tge.ButtonLeftOrTouchFirst != 0) {
		if !app.audioInit {
			if err := app.InitAudio(); err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
		}

		app.stereoPanNode.Pan(float32(mouseEvent.X)/float32(app.width)*2 - 1.0)
		app.gainNode.Gain(float32(app.height-mouseEvent.Y) / float32(app.height))

		fmt.Printf("Set pan to %v and gain to %v%%\n", (float32(mouseEvent.X)/float32(app.width)*2 - 1.0), float32(app.height-mouseEvent.Y)/float32(app.height)*100)

		if sourceNode, err := audio.CreateBufferSourceNode(app.sampleBuffer); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			sourceNode.Connect(app.stereoPanNode)
			sourceNode.Start(false, 0, 0)
		}
	}
	return false
}

func (app *AudioApp) OnRender(elaspedTime time.Duration, syncChan <-chan interface{}) {
	<-syncChan
}

func (app *AudioApp) OnTick(elaspedTime time.Duration, syncChan chan<- interface{}) {
	syncChan <- true
}

func (app *AudioApp) OnPause() {
	fmt.Println("OnPause()")
	// Close sound
	app.gainNode.Disconnect(app.destinationNode)
}

func (app *AudioApp) OnStop() {
	fmt.Println("OnStop()")
	app.runtime.Unsubscribe(tge.MouseEvent{}.Channel(), app.OnMouseEvent)
	app.runtime.Unsubscribe(tge.ResizeEvent{}.Channel(), app.OnResize)
}

func (app *AudioApp) OnDispose() {
	fmt.Println("OnDispose()")
}

func main() {
	tge.Run(&AudioApp{})
}
```

