package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"net/http"

	"github.com/bogem/id3v2"
)

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
			if filepath.Ext(path) == ".mp3" {
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
		log.Fatal("Error while opening mp3 file: ", err)
	}
	defer tag.Close()

	// do not load lyrics if already present
	uslfs := tag.GetFrames(tag.CommonID("Unsynchronised lyrics/text transcription"))
	if uslfs != nil {
		return
	}

	// Read frames to fetch lyrics.
	lyric := loadLyric(tag.Artist(), tag.Title())

	// Set lyrics frame.
	tag.AddUnsynchronisedLyricsFrame(lyric)

	// Write it to file.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

/**
* Get the lyrics object
 */
func loadLyric(artist, title string) id3v2.UnsynchronisedLyricsFrame {
	lyric := loadLyricsFromWikia(artist, title)

	uslt := id3v2.UnsynchronisedLyricsFrame{
		Encoding:          id3v2.EncodingUTF8,
		Language:          "eng", // todo: detect language
		ContentDescriptor: "Lyrics of " + title,
		Lyrics:            lyric,
	}

	return uslt
}

func loadLyricsFromWikia(artist, title string) string {
	url := "http://lyric-api.herokuapp.com/api/find/" + artist + "/" + title
	req, err := http.Get(url)
	if err != nil {
		log.Fatal("Error while fetching lyrics: ", err)
		return ""
	}
	defer req.Body.Close()
	// req.Response will contain a JavaScript Document element that can
	// for example be used with the js/dom package.
	var l WikiLyric
	json.NewDecoder(req.Body).Decode(l)
	return l.lyric
}

type WikiLyric struct {
	lyric string
	err   string
}
