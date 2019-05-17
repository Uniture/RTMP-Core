package format

import (
	"github.com/nulla-go/core/av/avutil"
	"github.com/nulla-go/core/format/aac"
	"github.com/nulla-go/core/format/flv"
	"github.com/nulla-go/core/format/mp4"
	"github.com/nulla-go/core/format/rtmp"
	"github.com/nulla-go/core/format/rtsp"
	"github.com/nulla-go/core/format/ts"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(mp4.Handler)
	avutil.DefaultHandlers.Add(ts.Handler)
	avutil.DefaultHandlers.Add(rtmp.Handler)
	avutil.DefaultHandlers.Add(rtsp.Handler)
	avutil.DefaultHandlers.Add(flv.Handler)
	avutil.DefaultHandlers.Add(aac.Handler)
}
