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
	ContextOpener() func(*Root, func(int))
	ContextKey() rune
	ContextHandler() func(start, end, sel int)
	DisableVisualMode() bool
	ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey
	Name() string
}

func SetCurrentView(v View) {
	CurrentView = v
}

func GetCurrentView() View {
	return CurrentView
}
