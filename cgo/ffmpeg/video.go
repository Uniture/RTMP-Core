package ffmpeg

/*
#include "ffmpeg.h"
int wrap_avcodec_decode_video2(AVCodecContext *ctx, AVFrame *frame, void *data, int size, int *got) {
	struct AVPacket pkt = {.data = data, .size = size};
	return avcodec_decode_video2(ctx, frame, got, &pkt);
}
*/
import "C"
import (
	"fmt"
	"image"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/nulla-go/core/av"
	"github.com/nulla-go/core/codec/h264parser"
)

type VideoDecoder struct {
	ff        *ffctx
	Extradata []byte
}

func (self *VideoDecoder) Setup() (err error) {
	ff := &self.ff.ff
	if len(self.Extradata) > 0 {
		ff.codecCtx.extradata = (*C.uint8_t)(unsafe.Pointer(&self.Extradata[0]))
		ff.codecCtx.extradata_size = C.int(len(self.Extradata))
	}
	if C.avcodec_open2(ff.codecCtx, ff.codec, nil) != 0 {
		err = fmt.Errorf("ffmpeg: decoder: avcodec_open2 failed")
		return
	}
	return
}

func fromCPtr(buf unsafe.Pointer, size int) (ret []uint8) {
	hdr := (*reflect.SliceHeader)((unsafe.Pointer(&ret)))
	hdr.Cap = size
	hdr.Len = size
	hdr.Data = uintptr(buf)
	return
}

type VideoFrame struct {
	Image image.YCbCr
	frame *C.AVFrame
}

func (self *VideoFrame) Free() {
	self.Image = image.YCbCr{}
	C.av_frame_free(&self.frame)
}

func freeVideoFrame(self *VideoFrame) {
	self.Free()
}

func (self *VideoDecoder) Decode(pkt []byte) (img *VideoFrame, err error) {
	ff := &self.ff.ff

	cgotimg := C.int(0)
	frame := C.av_frame_alloc()
	cerr := C.wrap_avcodec_decode_video2(ff.codecCtx, frame, unsafe.Pointer(&pkt[0]), C.int(len(pkt)), &cgotimg)
	if cerr < C.int(0) {
		err = fmt.Errorf("ffmpeg: avcodec_decode_video2 failed: %d", cerr)
		return
	}

	if cgotimg != C.int(0) {
		w := int(frame.width)
		h := int(frame.height)
		ys := int(frame.linesize[0])
		cs := int(frame.linesize[1])

		img = &VideoFrame{Image: image.YCbCr{
			Y:              fromCPtr(unsafe.Pointer(frame.data[0]), ys*h),
			Cb:             fromCPtr(unsafe.Pointer(frame.data[1]), cs*h/2),
			Cr:             fromCPtr(unsafe.Pointer(frame.data[2]), cs*h/2),
			YStride:        ys,
			CStride:        cs,
			SubsampleRatio: image.YCbCrSubsampleRatio420,
			Rect:           image.Rect(0, 0, w, h),
		}, frame: frame}
		runtime.SetFinalizer(img, freeVideoFrame)
	}

	return
}

func NewVideoDecoder(stream av.CodecData) (dec *VideoDecoder, err error) {
	_dec := &VideoDecoder{}
	var id uint32

	switch stream.Type() {
	case av.H264:
		h264 := stream.(h264parser.CodecData)
		_dec.Extradata = h264.AVCDecoderConfRecordBytes()
		id = C.AV_CODEC_ID_H264

	default:
		err = fmt.Errorf("ffmpeg: NewVideoDecoder codec=%v unsupported", stream.Type())
		return
	}

	c := C.avcodec_find_decoder(id)
	if c == nil || C.avcodec_get_type(id) != C.AVMEDIA_TYPE_VIDEO {
		err = fmt.Errorf("ffmpeg: cannot find video decoder codecId=%d", id)
		return
	}

	if _dec.ff, err = newFFCtxByCodec(c); err != nil {
		return
	}
	if err = _dec.Setup(); err != nil {
		return
	}

	dec = _dec
	return
}


type VideoEncoder struct{
	ff        *ffctx
}

func (self * VideoEncoder) Encode(img *VideoFrame)(pkt []byte, err error){
	//
}

func NewVideoEncoder(codecID av.CodecType) (enc * VideoEncoder, err error){
	
	// var CCodecID int
	// switch(codecID){
	// case av.H264:
	// 	CCodecID = C.AV_CODEC_ID_H264
	// 	break;
	// default:
	// 	err = fmt.Errorf("ffmpeg: cannot find video encoder codecID=%d", codecID)
	// 	return
	// }
	var id uint32

	switch codecID{
	case av.H264:
		id = C.AV_CODEC_ID_H264
		break
	default:
		err = fmt.Errorf("ffmpeg: wrong video coder=%d", codecID)
		return
	}

	enc = &VideoEncoder{}
	codec := C.avcodec_find_encoder(id)
	if codec == nil{
		err = fmt.Errorf("ffmpeg: cannot find video encoder codecID=%d", id)
		return
	}
	enc.ff, err = newFFCtxByCodec(codec)
	if err!=nil{
		return
	}
	enc.ff.ff 

	return
}
