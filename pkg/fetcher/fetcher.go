package fetcher

type fetcher interface {
	LoadLyrics(artist, title string) (lyrics string, err error)
}
