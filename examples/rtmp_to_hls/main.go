package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"net/http"
	"github.com/nulla-go/core/av/avutil"

	"github.com/nulla-go/core/av/pubsub"
	"github.com/nulla-go/core/format"
	"github.com/nulla-go/core/format/hls"
	"github.com/nulla-go/core/format/rtmp"
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
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
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

		hls := hls.NewHLSProcessing(&FileSystem{}, "http://localhost:9000/")
		err := hls.Pipe(key, ch.que.Latest())

		if err != nil {
			fmt.Println(err)
		}
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

	http.Handle("/", http.FileServer(http.Dir(".")))

	go http.ListenAndServe(":9000", nil)

	server.ListenAndServe()
}
