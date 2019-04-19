package main

import (
	"log"
	"sync"

	"github.com/strengine/Core/av"
	"github.com/strengine/Core/av/transcode"
	"github.com/strengine/Core/cgo/ffmpeg"

	"github.com/strengine/Core/av/avutil"

	"github.com/strengine/Core/av/pubsub"
	"github.com/strengine/Core/format"
	"github.com/strengine/Core/format/rtmp"
)

func init() {
	format.RegisterAll()
}

func main() {
	server := &rtmp.Server{}

	l := &sync.RWMutex{}
	// d := &sync.RWMutex{}
	type Channel struct {
		que *pubsub.Queue
	}
	channels := map[string]*Channel{}
	// DChannels := map[string]*Channel{}

	test1 := func(key string) {
		log.Println(key + " is started processing")
		l.Lock()
		ch := channels[key]
		if ch == nil {
			log.Println(key + " already delete processing close")
			return
		}
		l.Unlock()

		findcodec := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
			need = true
			dec, _ = ffmpeg.NewAudioDecoder(stream)
			enc, _ = ffmpeg.NewAudioEncoderByName("libfdk_aac")
			enc.SetSampleRate(stream.SampleRate())
			enc.SetChannelLayout(av.CH_STEREO)
			enc.SetBitrate(12000)
			enc.SetOption("profile", "HE-AACv2")
			return
		}

		trans := &transcode.Demuxer{
			Options: transcode.Options{
				FindAudioDecoderEncoder: findcodec,
			},
			Demuxer: ch.que,
		}
		outfile, _ := avutil.Create("out.ts")
		avutil.CopyFile(outfile, trans)

		outfile.Close()
		trans.Close()

		// avutil.
		// l.Lock()
		// d.Lock()

	}

	server.HandlePublish = func(conn *rtmp.Conn) {
		l.Lock()
		key := conn.URL.Query().Get("key")
		ch := channels[conn.URL.Path]
		if ch == nil {
			ch = &Channel{}
			ch.que = pubsub.NewQueue()
			if key == "" {
				log.Println("'key' is undefined in url path")
				l.Unlock()
				return
			}
			log.Println(key + " is connected")
			channels[key] = ch
		}
		l.Unlock()
		go test1(key)
		avutil.CopyFile(ch.que, conn)
		log.Println(key + " is disconnected")
		l.Lock()
		delete(channels, key)
		l.Unlock()
		ch.que.Close()
	}

	server.ListenAndServe()
}
