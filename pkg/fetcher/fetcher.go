package fetcher

// Fetcher interface
type Fetcher interface {
	FetchLyrics(artist, title string) (lyrics string, err error)
}
