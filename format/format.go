package format

import (
	"github.com/nulla-go/Core/av/avutil"
	"github.com/nulla-go/Core/format/aac"
	"github.com/nulla-go/Core/format/flv"
	"github.com/nulla-go/Core/format/mp4"
	"github.com/nulla-go/Core/format/rtmp"
	"github.com/nulla-go/Core/format/rtsp"
	"github.com/nulla-go/Core/format/ts"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(mp4.Handler)
	avutil.DefaultHandlers.Add(ts.Handler)
	avutil.DefaultHandlers.Add(rtmp.Handler)
	avutil.DefaultHandlers.Add(rtsp.Handler)
	avutil.DefaultHandlers.Add(flv.Handler)
	avutil.DefaultHandlers.Add(aac.Handler)
}
