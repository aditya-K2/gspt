package ui

import "github.com/gdamore/tcell/v2"

var (
	CurrentView        View
	playlistView       = &PlaylistView{}
	albumView          = &AlbumView{}
	albumsView         = &AlbumsView{}
	likedSongsView     = &LikedSongsView{}
	recentlyPlayedView = &RecentlyPlayedView{}
	topTracksView      = &TopTracksView{}
	artistView         = &ArtistView{}
	artistsView        = &ArtistsView{}
	searchView         = &SearchView{}
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
