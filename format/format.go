package format

import (
	"github.com/strengine/core/av/avutil"
	"github.com/strengine/core/format/aac"
	"github.com/strengine/core/format/flv"
	"github.com/strengine/core/format/mp4"
	"github.com/strengine/core/format/rtmp"
	"github.com/strengine/core/format/rtsp"
	"github.com/strengine/core/format/ts"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(mp4.Handler)
	avutil.DefaultHandlers.Add(ts.Handler)
	avutil.DefaultHandlers.Add(rtmp.Handler)
	avutil.DefaultHandlers.Add(rtsp.Handler)
	avutil.DefaultHandlers.Add(flv.Handler)
	avutil.DefaultHandlers.Add(aac.Handler)
}
