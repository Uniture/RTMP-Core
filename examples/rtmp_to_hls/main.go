package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/nulla-go/Core/av/avutil"

	"github.com/nulla-go/Core/av/pubsub"
	"github.com/nulla-go/Core/format"
	"github.com/nulla-go/Core/format/hls"
	"github.com/nulla-go/Core/format/rtmp"
)

type FileSystem struct {
}

//type File interface {
//	Write(b []byte) (n int, err error)
//	WriteString(s string) (n int, err error)
//}

func (fs FileSystem) Mkdir(name string, perm os.FileMode) (err error) {
	err = os.Mkdir(name, perm)
	return
}

func (fs FileSystem) IsExist(err error) (res bool) {
	res = os.IsExist(err)
	return
}

func (fs FileSystem) Create(name string) (io.Writer, error) {
	file, err := os.Create(name)
	return file, err
}

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

		hls := hls.NewHLSProcessing(&FileSystem{})
		err := hls.Pipe(key, ch.que.Latest())

		if err != nil {
			fmt.Println(err)
		}
		//findcodec := func(stream av.AudioCodecData, i int) (need bool, dec av.AudioDecoder, enc av.AudioEncoder, err error) {
		//need = true
		//dec, _ = ffmpeg.NewAudioDecoder(stream)
		//enc, err = ffmpeg.NewAudioEncoderByName("libfdk_aac")
		//if err != nil {
		//	log.Println("Encoder is undefined")
		//	return
		//}
		//enc.SetSampleRate(stream.SampleRate())
		//enc.SetBitrate(48000)
		//enc.SetChannelLayout(av.CH_STEREO)
		//enc.SetOption("profile", "HE-AACv2")
		//return
		//}
		//Options: transcode.Options{
		//		FindAudioDecoderEncoder: findcodec,
		//	},
		//	Demuxer: ch.que.Latest(),
		//}
		//outfile, _ := avutil.Create("out.ts")
		//avutil.CopyFile(outfile, ch.que.Latest())

		//outfile.Close()
		//		trans.Close()

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
