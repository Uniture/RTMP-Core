package format

import (
	"github.com/strengine/Core/av/avutil"
	"github.com/strengine/Core/format/aac"
	"github.com/strengine/Core/format/flv"
	"github.com/strengine/Core/format/mp4"
	"github.com/strengine/Core/format/rtmp"
	"github.com/strengine/Core/format/rtsp"
	"github.com/strengine/Core/format/ts"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(mp4.Handler)
	avutil.DefaultHandlers.Add(ts.Handler)
	avutil.DefaultHandlers.Add(rtmp.Handler)
	avutil.DefaultHandlers.Add(rtsp.Handler)
	avutil.DefaultHandlers.Add(flv.Handler)
	avutil.DefaultHandlers.Add(aac.Handler)
}
