package manipulator

// Manipulator Interface
type Manipulator interface {
	CanHandle(file string) (able bool, err error)
	GetArtist(file string) (artist string, err error)
	GetSongname(file string) (songname string, err error)
	HasLyrics(file string) (has bool, err error)
	SetLyrics(file string, lyric string) (err error)
}
