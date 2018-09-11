package manipulator

import (
	"fmt"
	"path/filepath"

	"github.com/mikkyang/id3-go"
	"github.com/mikkyang/id3-go/v2"
)

type Mp3Manipulator struct {
}

func (m Mp3Manipulator) HasLyrics(file string) (has bool, err error) {
	mp3File, err := m.openFile(file)
	if err != nil {
		return has, err
	}

	uslfs, fail := mp3File.Frame("USLT").(*v2.UnsynchTextFrame)
	if fail != false {
		fmt.Printf("Error while fetching USLT frame: %v\n", fail)
		return has, fmt.Errorf("Error while fetching USLT frame: %v", fail)
	}

	// has lyrics
	if uslfs != nil && uslfs.String() != "" {
		return true, err
	}
	// no lyrics, no err
	return false, err
}
func (m Mp3Manipulator) CanHandle(file string) (able bool, err error) {
	return filepath.Ext(file) == ".mp3", err
}
func (m Mp3Manipulator) GetArtist(file string) (artist string, err error) {
	mp3File, err := m.openFile(file)

	return mp3File.Artist(), err
}
func (m Mp3Manipulator) GetSongname(file string) (songname string, err error) {
	mp3File, err := m.openFile(file)

	return mp3File.Title(), err
}
func (m Mp3Manipulator) SetLyrics(file string, lyric string) (err error) {
	mp3File, err := m.openFile(file)

	ft := v2.V23FrameTypeMap["USLT"]
	desc := "Lyrics"
	uslf := v2.NewUnsynchTextFrame(ft, desc, lyric)

	mp3File.AddFrames(uslf)

	return err
}
func (m Mp3Manipulator) openFile(file string) (mp3File *id3.File, err error) {
	mp3File, err = id3.Open(file)
	defer mp3File.Close()
	return mp3File, err
}
