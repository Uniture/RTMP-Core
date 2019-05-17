package hls

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nulla-go/core/av"
	"github.com/nulla-go/core/format/ts"
)

type hlsFileSystemError struct {
	Original error
	Comment  string
}

func (e hlsFileSystemError) Error() string {
	return fmt.Sprintf("%s: %s", e.Comment, e.Original)
}

const (
	masterPlaylistPath = "/master.m3u8"
	playlistPath       = "/stream0"
)

type HLSFile interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
}

type hlsDestination interface {
	Mkdir(name string, perm os.FileMode) error
	IsExist(err error) bool
	Create(name string) (io.Writer, error)
}

type hls struct {
	filesys hlsDestination
}

func NewHLSProcessing(filesys hlsDestination) (chls *hls) {
	chls = &hls{}
	chls.filesys = filesys
	return
}

func (self *hls) Pipe(name string, src av.Demuxer) error {
	var masterPL io.Writer
	var err error
	// Creating new stream's folder
	if err := self.filesys.Mkdir(name, os.ModePerm); err != nil {
		return hlsFileSystemError{err, "error creating dir with name" + name}
	}
	// Creating master playlist
	if masterPL, err = self.filesys.Create(name + masterPlaylistPath); err != nil {
		return hlsFileSystemError{err,
			"error creating file with path " + name + masterPlaylistPath}
	}

	fmt.Println("MasterPlaylist was created")

	//	if err := self.filesys.Mkdir(name+playlistPath, os.ModeAppend); err != nil {
	//		return hlsFileSystemError{err, "error creating dir " + name + playlistPath}
	//}

	self.processing(masterPL, name, src)
	return nil
}

func (self *hls) processing(masterPlaylist io.Writer, root string, src av.Demuxer) {
	fmt.Println("processing started")
	streamP := root + "/stream0"
	err := self.filesys.Mkdir(streamP, os.ModePerm)
	if err != nil {
		fmt.Println("Cannot create stream dir ", err.Error())
		return
	}
	var videoStreamNumber int8
	codecsData, err := src.Streams()
	if err != nil {
		fmt.Println("Failed get streams")
		return
	}
	for streamN, cData := range codecsData {
		fmt.Print(streamN)
		if cData.Type().IsAudio() {
			fmt.Print(" is a audio stream")
			fmt.Println()
		} else if cData.Type().IsVideo() {
			fmt.Print(" is a video stream")
			fmt.Println()
			videoStreamNumber = int8(streamN)
		} else {
			fmt.Print(" is undefined stream")
			fmt.Println()
		}
	}

	//	var limit time.Duration = 9 * time.Second
	counter := 9 * time.Second
	//	var created bool = false
	var fileCounter int64

	var currentBufer int8
	var bufCounter int64
	//	var buf1Counter int64 = 0
	buf0 := make([]av.Packet, 4096)
	buf1 := make([]av.Packet, 4096)
	for {
		var pkt av.Packet
		//var c chan av.Packet
		//c := make(chan av.Packet, 10)
		var err error
		if pkt, err = src.ReadPacket(); err != nil {
			if err == io.EOF {
				break
			}
			//			close(c)
			return
		}

		if currentBufer == 0 {
			buf0[bufCounter] = pkt
		} else {
			buf1[bufCounter] = pkt
		}
		bufCounter++

		//c <- pkt
		//packateCounter += 1
		//		fmt.Println("Packet sent to channel")
		if bufCounter%20 == 0 {
			fmt.Println(bufCounter)
		}
		//if !created {
		//	fmt.Println("Writing started")
		//	go self.writePktsToFile(streamP+"/stream0.ts", c, codecsData)
		//	created = true
		//}
		//fmt.Println(pkt.Idx)
		if pkt.Idx == videoStreamNumber {
			if pkt.Time >= counter {
				fmt.Println("Limit", counter.Seconds())
				currentStreamFileName := fmt.Sprintf(root+"/stream%v.ts", fileCounter)
				fileCounter++
				if currentBufer == 0 {
					go self.writePktsToFile(currentStreamFileName, codecsData, &buf0, bufCounter)
					currentBufer = 1
				} else {
					go self.writePktsToFile(currentStreamFileName, codecsData, &buf1, bufCounter)
					currentBufer = 0
				}
				bufCounter = 0

				counter += 9 * time.Second
			}
		}

		//fmt.Println(pkt.Time)
	}

}

func (self *hls) writePktsToFile(
	filePath string, streams []av.CodecData, pkts *[]av.Packet, length int64) {

	writer, err := self.filesys.Create(filePath)
	if err != nil {
		fmt.Println("Cannot creating file " + filePath)
		return
	}
	tsMuxer := ts.NewMuxer(writer)
	//for i, cData := range streams{
	//	err = tsMuxer.newStream(s.Data)
	//	if err != nil{
	//		fmt.Println("failed codec data ", i, err.Error())
	//	}
	//}

	err = tsMuxer.WriteHeader(streams)
	if err != nil {
		fmt.Println("Cannot write header ", err.Error())
		return
	}

	//var pktCpunter int64

	for i, pkt := range *pkts {
		if i > int(length) {
			break
		}
		err = tsMuxer.WritePacket(pkt)
		if err != nil {
			fmt.Println("Cannot write Packet")
		}
	}
	//	for _, pkt := range *pkts {

	//select {
	//case <-c:
	//	fmt.Println("Packate got")
	//}

	//	}

	//	for pkt := range c {
	//		}

	err = tsMuxer.WriteTrailer()
	if err != nil {
		fmt.Println("Cannot write trailer ", err.Error())
		return
	}

}
