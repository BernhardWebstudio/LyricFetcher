package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"net/http"
	"net/url"

	"github.com/bogem/id3v2"
)

type WikiLyric struct {
	Lyric string `json:"lyric"`
	Err   string `json:"err"`
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
	tag, err := id3v2.Open(file, id3v2.Options{Parse: true})
	if err != nil {
		fmt.Printf("Error while opening mp3 file: %v", err)
		return
	}
	defer tag.Close()

	// do not load lyrics if already present
	hasLyrics := false
	uslfs := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
	for _, f := range uslfs {
		uslf, ok := f.(id3v2.UnsynchronisedLyricsFrame)
		if !ok {
			fmt.Println("Couldn't assert USLT frame. Skipping.")
			return
		}

		if uslf.Lyrics != "" {
			hasLyrics = true
		}
	}

	if hasLyrics {
		return
	}

	// Read frames to fetch lyrics.
	artist := tag.Artist()
	title := tag.Title()
	if artist == "" {
		fmt.Println("Artist not set. Skipping.")
		return
	}
	if title == "" {
		fmt.Println("Title not set. Skipping.")
		return
	}

	lyric, err := loadLyric(artist, title)
	if err != nil {
		fmt.Printf("Error while loading lyric: %v", err)
		return
	}

	// Set lyrics frame.
	tag.AddUnsynchronisedLyricsFrame(lyric)

	// Write it to file.
	err = tag.Save()
	if err != nil {
		fmt.Printf("Error while saving tag: %v", err)
		return
	}
}

/**
* Get the lyrics object
 */
func loadLyric(artist, title string) (lyricFrame id3v2.UnsynchronisedLyricsFrame, err error) {
	// TODO: handle failures
	// maybe by using different API: https://github.com/BharatKalluri/lyricfetcher
	lyric, err := loadLyricsFromWikia(artist, title)
	if err != nil || strings.TrimSpace(lyric) == "" {
		return lyricFrame, err
	}

	fmt.Println(lyric)

	// TODO: detect language, e.g. with https://github.com/chrisport/go-lang-detector
	uslt := id3v2.UnsynchronisedLyricsFrame{
		Encoding:          id3v2.EncodingUTF8,
		Language:          "eng",
		ContentDescriptor: "Lyrics of " + title,
		Lyrics:            lyric,
	}

	return uslt, nil
}

func loadLyricsFromWikia(artist, title string) (lyrics string, err error) {
	url := "http://lyric-api.herokuapp.com/api/find/" + url.QueryEscape(artist) + "/" + url.QueryEscape(title)
	response, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer response.Body.Close()
	// req.Response will contain a JavaScript Document element that can
	// for example be used with the js/dom package.
	var l WikiLyric
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(contents[:]))
	err = json.Unmarshal(contents, &l)
	if err != nil {
		return "", err
	}
	if l.Err == "" {
		fmt.Println(l.Err)
	}
	return l.Lyric, nil
}
