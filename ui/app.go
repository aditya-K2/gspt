package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	ImgY int
	ImgW int
	ImgH int
	ImgX int
	Ui   *Application
)

type Application struct {
	App            *tview.Application
	Main           *interactiveView
	Navbar         *tview.Table
	SearchBar      *tview.Box
	ProgressBar    *tview.Box
	Root           *Root
	ImagePreviewer *tview.Box
}

func NewApplication() *Application {

	App := tview.NewApplication()
	Main := NewMain()
	pBar := tview.NewBox().SetBorder(true).SetTitle("PROGRESS").SetBackgroundColor(tcell.ColorDefault)
	searchbar := tview.NewBox().SetBorder(true).SetTitle("SEARCH").SetBackgroundColor(tcell.ColorDefault)
	SetCurrentView(PView)
	mains := NewInteractiveView()
	mains.Table.SetBorder(true)

	mains.SetContentFunc(GetCurrentView().Content)
	mains.SetContextKey(GetCurrentView().ContextKey())
	f := func() {
		GetCurrentView().ContextOpener(Main, mains.SelectionHandler)
	}
	mains.SetContextOpener(f)
	mains.SetContextHandler(GetCurrentView().ContextHandler)
	mains.SetExternalCapture(GetCurrentView().ExternalInputCapture)

	Navbar := tview.NewTable()
	imagePreviewer := tview.NewBox()

	imagePreviewer.SetBorder(true)

	Navbar.SetBackgroundColor(tcell.ColorDefault)
	imagePreviewer.SetBackgroundColor(tcell.ColorDefault)

	Navbar.SetBorder(true)
	Navbar.SetSelectable(true, false)

	done := func(s bool, err error) {
		if s {
			App.Draw()
		}
	}
	playlistNav, err := NewPlaylistNav(done)
	if err != nil {
		panic(err)
	}

	playlistNav.Table.SetBackgroundColor(tcell.ColorDefault)
	PlaylistActions = map[string]*Action{
		"playEntry": NewAction(playlistNav.PlaySelectEntry, nil),
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			Main.AfterContextClose(func() { App.SetFocus(mains.Table) })
			r, _ := playlistNav.Table.GetSelection()
			PView.SetPlaylist(&(*playlistNav.Playlists)[r])
			App.SetFocus(mains.Table)
			return nil
		}, nil),
	}
	playlistNav.MapActions(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(Navbar, 6, 3, false).
		AddItem(playlistNav.Table, 0, 6, false).
		AddItem(imagePreviewer, 9, 3, false)

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 17, 1, false).
		AddItem(mains.Table, 0, 4, false)

	searchBarFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchbar, 3, 1, false).
		AddItem(sNavExpViewFlex, 0, 1, false)

	MainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBarFlex, 0, 8, false).
		AddItem(pBar, 5, 1, false)

	Main.Primitive("Main", MainFlex)
	App.EnableMouse(true)
	App.SetRoot(Main.Root, true).SetFocus(playlistNav.Table)

	Ui = &Application{
		App:            App,
		Main:           mains,
		Navbar:         Navbar,
		SearchBar:      searchbar,
		ProgressBar:    pBar,
		Root:           Main,
		ImagePreviewer: imagePreviewer,
	}

	return Ui
}
