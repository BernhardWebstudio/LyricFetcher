package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BernhardWebstudio/LyricFetcher/pkg/fetcher"
	"github.com/BernhardWebstudio/LyricFetcher/pkg/manipulator"
)

// HandleFiles handles files
func HandleFiles(args []string) {
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
				HandleFile(path)
			} else {
				fmt.Println("Skipping " + path)
			}
			return nil
		})
	}
}

// HandleFile handles one specific given file
func HandleFile(file string) {
	fmt.Println("Handling " + file)
	// Find appropriate Manipulator
	success := false

	// TODO: add more manipulators here
	var manipulators = []manipulator.Manipulator{
		manipulator.Mp3Manipulator{},
	}

	for _, manipulator := range manipulators {
		success = tryHandleFile(file, manipulator)
		if success {
			break
		}
	}

	if success {
		fmt.Println("Handled " + file + " successfully")
	} else {
		fmt.Println("Problem occurred while handling  " + file)
	}

}

func tryHandleFile(file string, manipulator manipulator.Manipulator) (success bool) {
	can, err := manipulator.CanHandle(file)

	// only try if manipulator can handle this filetype
	if err != nil || can == false {
		return false
	}

	// do not load lyrics if already present
	has, err := manipulator.HasLyrics(file)
	if err != nil || has == false {
		fmt.Println(file + " already has Lyrics. Skipping.")
		return true
	}

	// Read frames necessary to fetch lyrics.
	artist, err := manipulator.GetArtist(file)
	if artist == "" || err != nil {
		return false
	}
	title, err := manipulator.GetSongname(file)
	// TODO: use filename as title
	if title == "" || err != nil {
		return false
	}

	lyric, success := LoadLyric(artist, title)
	if success != true {
		// fmt.Printf("Lyrics could not be loaded for title '%s' by '%s'\n", title, artist)
		return false
	}

	// Set lyrics frame.
	err = manipulator.SetLyrics(file, lyric)

	return err == nil
}

// LoadLyric traverses the loaders and loads until success or all through
func LoadLyric(artist, title string) (lyric string, success bool) {
	// Find appropriate fetcher
	success = false

	// TODO: add more fetcher here
	var fetchers = []fetcher.Fetcher{
		fetcher.WikiaFetcher{},
	}

	for _, fetcher := range fetchers {
		lyric, success = tryLoadingLyric(artist, title, fetcher)
		if success {
			break
		}
	}

	return lyric, success
}

func tryLoadingLyric(artist string, title string, fetcher fetcher.Fetcher) (lyric string, success bool) {
	lyric, err := fetcher.FetchLyrics(artist, title)

	if strings.TrimSpace(lyric) == "" {
		return "", false
	}

	return lyric, err == nil
}
