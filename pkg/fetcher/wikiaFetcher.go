package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type wikiLyric struct {
	Lyric string      `json:"lyric"`
	Err   interface{} `json:"err"`
}

type wikiaFetcher struct {
}

func (w wikiaFetcher) FetchLyrics(artist, title string) (lyrics string, err error) {
	url := "http://lyric-api.herokuapp.com/api/find/" + url.QueryEscape(artist) + "/" + url.QueryEscape(title)
	response, err := http.Get(url)
	if err != nil {
		return lyrics, err
	}
	defer response.Body.Close()
	// req.Response will contain a JavaScript Document element that can
	// for example be used with the js/dom package.
	var l wikiLyric
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
