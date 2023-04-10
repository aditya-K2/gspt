package ui

import "github.com/gdamore/tcell/v2"

var (
	CurrentView View
	PView       = &PlaylistView{}
)

type View interface {
	Content() [][]Content
	ContextOpener(*Main, func(int))
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
