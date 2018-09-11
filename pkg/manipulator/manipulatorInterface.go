package manipulator

type manipulator interface {
	CanHandle(file string) (able bool)
	GetArtist(file string) (artist string, err error)
	GetSongname(file string) (songname string, err error)
	SetLyrics(file string, lyric string) (err error)
}
