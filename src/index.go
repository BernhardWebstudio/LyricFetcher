package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mikkyang/id3-go"
	"github.com/mikkyang/id3-go/v2"

	"net/http"
	"net/url"
)

type WikiLyric struct {
	Lyric string      `json:"lyric"`
	Err   interface{} `json:"err"`
}

func main() {
	args := os.Args[1:]

	// iterate given paths
	for _, reet := range args {
		// iterate files in path
		filepath.Walk(reet, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
				return err
			}
			// automatically skip dir etc. by not doing anything
			if filepath.Ext(path) == ".mp3" { // || filepath.Ext(path) == ".m4a"
				handleFile(path)
			} else {
				fmt.Println("Skipping " + path)
			}
			return nil
		})
	}
}

func handleFile(file string) {
	fmt.Println("Handling " + file)
	// Open file and parse tag in it.
	mp3File, err := id3.Open(file)
	if err != nil {
		fmt.Printf("Error while opening mp3 file: %v\n", err)
		return
	}
	defer mp3File.Close()

	// do not load lyrics if already present
	uslfs, fail := mp3File.Frame("USLT").(*v2.UnsynchTextFrame)
	if fail != false {
		fmt.Printf("Error while fetching USLT frame: %v\n", fail)
		return
	}

	// has lyrics
	if uslfs != nil && uslfs.String() != "" {
		return
	}

	// Read frames to fetch lyrics.
	artist := mp3File.Artist()
	title := mp3File.Title()
	if artist == "" {
		fmt.Println("Artist not set. Skipping.")
		return
	}
	// TODO: use filename as title
	if title == "" {
		fmt.Println("Title not set. Skipping.")
		return
	}

	lyric, err := loadLyric(artist, title)
	if err != nil {
		fmt.Printf("Error while loading lyric: %v\n", err)
		return
	}

	fmt.Println("setting lyrics on " + file)
	// Set lyrics frame.
	mp3File.AddFrames(lyric)
}

/**
* Get the lyrics object
 */
func loadLyric(artist, title string) (uslf *v2.UnsynchTextFrame, err error) {
	// TODO: handle failures
	// maybe by using different API: https://github.com/BharatKalluri/lyricfetcher
	lyric, err := loadLyricsFromWikia(artist, title)
	if err != nil {
		return uslf, err
	}

	desc := "Lyrics of " + title + " by " + artist

	if strings.TrimSpace(lyric) == "" {
		return uslf, errors.New("Empty " + desc + " found: " + lyric)
	}

	// TODO: detect language, e.g. with https://github.com/chrisport/go-lang-detector
	ft := v2.V23FrameTypeMap["USLT"]
	uslf = v2.NewUnsynchTextFrame(ft, desc, lyric)
	fmt.Println("Loaded " + desc)
	return uslf, err
}

func loadLyricsFromWikia(artist, title string) (lyrics string, err error) {
	url := "http://lyric-api.herokuapp.com/api/find/" + url.QueryEscape(artist) + "/" + url.QueryEscape(title)
	response, err := http.Get(url)
	if err != nil {
		return lyrics, err
	}
	defer response.Body.Close()
	// req.Response will contain a JavaScript Document element that can
	// for example be used with the js/dom package.
	var l WikiLyric
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return lyrics, err
	}
	err = json.Unmarshal(contents, &l)
	if err != nil {
		return lyrics, err
	}
	if l.Err == nil {
		fmt.Printf("%+v\n", l.Err)
		errorMessage, err := json.Marshal(l.Err)
		if nil != err {
			err = errors.New(string(errorMessage[:]))
		}
	}
	return l.Lyric, err
}