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
	MainS          *tview.Pages
	Navbar         *tview.Table
	SearchBar      *tview.Box
	ProgressBar    *tview.Box
	Pages          *tview.Pages
	ImagePreviewer *tview.Box
}

func NewApplication() *Application {

	App := tview.NewApplication()
	Main := NewMain()
	pBar := tview.NewBox().SetBorder(true).SetTitle("PROGRESS").SetBackgroundColor(tcell.ColorDefault)
	mainS := tview.NewPages()
	searchbar := tview.NewBox().SetBorder(true).SetTitle("SEARCH").SetBackgroundColor(tcell.ColorDefault)

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
	Playlists, err := NewPlaylistNav(done)
	if err != nil {
		panic(err)
	}

	Playlists.Table.SetBackgroundColor(tcell.ColorDefault)
	PlaylistActions = map[string]*Action{
		"playEntry": NewAction(Playlists.PlaySelectEntry, nil),
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			Main.AfterContextClose(func() { App.SetFocus(mainS) })
			p := NewPlaylistView(Main)
			r, _ := Playlists.Table.GetSelection()
			p.CurrentPlaylist = &(*Playlists.Playlists)[r]
			mainS.AddPage("PLAYLIST", p.I.View, true, true)
			App.SetFocus(mainS)
			return nil
		}, nil),
	}
	Playlists.MapActions(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(Navbar, 6, 3, false).
		AddItem(Playlists.Table, 0, 6, false).
		AddItem(imagePreviewer, 9, 3, false)

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 17, 1, false).
		AddItem(mainS, 0, 4, false)

	searchBarFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchbar, 3, 1, false).
		AddItem(sNavExpViewFlex, 0, 1, false)

	MainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBarFlex, 0, 8, false).
		AddItem(pBar, 5, 1, false)

	rootPages := tview.NewPages()
	rootPages.AddPage("Main", MainFlex, true, true)

	Main.Primitive("root", rootPages)
	App.EnableMouse(true)
	App.SetRoot(Main.Root, true).SetFocus(Playlists.Table)

	return &Application{
		App:            App,
		MainS:          mainS,
		Navbar:         Navbar,
		SearchBar:      searchbar,
		ProgressBar:    pBar,
		Pages:          rootPages,
		ImagePreviewer: imagePreviewer,
	}

}
