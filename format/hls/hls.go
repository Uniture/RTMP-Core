package hls

import (
	"fmt"
	"os"

	"github.com/nulla-go/Core/av"
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
)

type HlsFile interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (n int, err error)
}

type hlsDestination interface {
	Mkdir(name string, perm os.FileMode) error
	IsExist(err error) bool
	Create(name string) (HlsFile, error)
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
	// Creating new stream's folder
	if err := self.filesys.Mkdir(name, os.ModeAppend); err != nil {
		return hlsFileSystemError{err, "error creating dir with name" + name}
	}
	// Creating master playlist
	if _, err := self.filesys.Create(name + masterPlaylistPath); err != nil {
		return hlsFileSystemError{err, "error creating file with path " + name + masterPlaylistPath}
	}
	fmt.Println("MasterPlaylist was created")
	return nil
}
