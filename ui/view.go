package ui

import "github.com/gdamore/tcell/v2"

var (
	CurrentView        View
	playlistView       = NewPlaylistView()
	albumView          = NewAlbumView()
	albumsView         = NewAlbumsView()
	likedSongsView     = NewLikedSongsView()
	recentlyPlayedView = NewRecentlyPlayedView()
	topTracksView      = NewTopTracksView()
	artistView         = NewArtistView()
	artistsView        = NewArtistsView()
	searchView         = NewSearchView()
)

type View interface {
	Content() func() [][]Content
	DisableVisualMode() bool
	ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey
	VisualCapture() func(start, end int, e *tcell.EventKey) *tcell.EventKey
	Name() string
	OpenEntry()
}

func SetCurrentView(v View) {
	CurrentView = v
}
