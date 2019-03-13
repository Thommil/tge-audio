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
 * Desktop
 * Mobile
 * Browser

## Dependencies
 * [TGE core](https://github.com/thommil/tge)

### Desktop & Mobile

In order to use this package on Linux desktop distros, you may need OpenAL library as an external dependency:

```
sudo apt-get install libopenal-dev
```

## Limitations
Only [StereoPannerNode](https://developer.mozilla.org/en-US/docs/Web/API/StereoPannerNode) and [GainNode](https://developer.mozilla.org/en-US/docs/Web/API/GainNode) are currently available, custom [AudioNode](https://developer.mozilla.org/en-US/docs/Web/API/AudioNode) and spacialization features are planned in future releases.

Music source is not optimized and should be improved in future releases.

## Implementation
See example at [AUDIO examples](https://github.com/Thommil/tge-examples/tree/master/plugins/tge-audio)


```golang
TODO
```

