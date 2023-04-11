package ui

import "github.com/gdamore/tcell/v2"

var (
	CurrentView  View
	playlistView = &PlaylistView{}
	albumView    = &AlbumView{}
	albumSView   = &AlbumsView{}
)

type View interface {
	Content() [][]Content
	ContextOpener(*Root, func(int))
	ContextKey() rune
	ContextHandler(start, end, sel int)
	ExternalInputCapture(e *tcell.EventKey) *tcell.EventKey
	Name() string
}

func SetCurrentView(v View) {
	CurrentView = v
}

func GetCurrentView() View {
	return CurrentView
}
