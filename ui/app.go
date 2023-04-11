package ui

import (
	"fmt"

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

var (
	TrackStyle       = tcell.StyleDefault.Foreground(tcell.ColorBlue)
	AlbumStyle       = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	ArtistStyle      = tcell.StyleDefault.Foreground(tcell.ColorPink)
	TimeStyle        = tcell.StyleDefault.Foreground(tcell.ColorOrange)
	PlaylistNavStyle = tcell.StyleDefault.Foreground(tcell.ColorCoral)
	NavStyle         = tcell.StyleDefault.Foreground(tcell.ColorPapayaWhip).Background(tcell.ColorBlack)
	ContextMenuStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
)

type Application struct {
	App            *tview.Application
	Main           *interactiveView
	NavMenu        *NavMenu
	SearchBar      *tview.Box
	ProgressBar    *tview.Box
	Root           *Root
	ImagePreviewer *tview.Box
}

func NewApplication() *Application {

	App := tview.NewApplication()
	Root := NewRoot()
	pBar := tview.NewBox().SetBorder(true).SetTitle("PROGRESS").SetBackgroundColor(tcell.ColorDefault)
	searchbar := tview.NewBox().SetBorder(true).SetTitle("SEARCH").SetBackgroundColor(tcell.ColorDefault)
	SetCurrentView(playlistView)
	Main := NewInteractiveView()
	Main.Table.SetBorder(true)

	NavMenu := newNavMenu([]navItem{
		{"Albums", NewAction(func(e *tcell.EventKey) *tcell.EventKey { SetCurrentView(albumsView); return nil }, nil)},
		{"Artists", NewAction(func(e *tcell.EventKey) *tcell.EventKey { fmt.Println("Artists"); return nil }, nil)},
		{"Liked Songs", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			likedSongsView.RefreshState()
			SetCurrentView(likedSongsView)
			App.SetFocus(Main.Table)
			return nil
		}, nil)},
		{"Recently Played", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			recentlyPlayedView.RefreshState()
			SetCurrentView(recentlyPlayedView)
			App.SetFocus(Main.Table)
			return nil
		}, nil)},
	})
	imagePreviewer := tview.NewBox()

	imagePreviewer.SetBorder(true)

	NavMenu.Table.SetBackgroundColor(tcell.ColorDefault)
	imagePreviewer.SetBackgroundColor(tcell.ColorDefault)

	NavMenu.Table.SetBorder(true)
	NavMenu.Table.SetSelectable(true, false)

	done := func(s bool, err error) {
		if s {
			App.Draw()
		}
	}
	playlistNav, err := NewPlaylistNav(done)
	if err != nil {
		panic(err)
	}

	Root.AfterContextClose(func() { App.SetFocus(Main.Table) })
	playlistNav.Table.SetBackgroundColor(tcell.ColorDefault)
	PlaylistActions = map[string]*Action{
		"playEntry": NewAction(playlistNav.PlaySelectEntry, nil),
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			r, _ := playlistNav.Table.GetSelection()
			playlistView.SetPlaylist(&(*playlistNav.Playlists)[r])
			SetCurrentView(playlistView)
			App.SetFocus(Main.Table)
			return nil
		}, nil),
	}
	playlistNav.MapActions(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
		tcell.KeyCtrlP: "playEntry",
	})

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(NavMenu.Table, 6, 3, false).
		AddItem(playlistNav.Table, 0, 6, false).
		AddItem(imagePreviewer, 9, 3, false)

	sNavExpViewFlex := tview.NewFlex().
		AddItem(searchNavFlex, 17, 1, false).
		AddItem(Main.Table, 0, 4, false)

	searchBarFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchbar, 3, 1, false).
		AddItem(sNavExpViewFlex, 0, 1, false)

	MainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBarFlex, 0, 8, false).
		AddItem(pBar, 5, 1, false)

	Root.Primitive("Main", MainFlex)
	App.EnableMouse(true)
	App.SetRoot(Root.Root, true).SetFocus(playlistNav.Table)

	Ui = &Application{
		App:            App,
		Main:           Main,
		NavMenu:        NavMenu,
		SearchBar:      searchbar,
		ProgressBar:    pBar,
		Root:           Root,
		ImagePreviewer: imagePreviewer,
	}

	return Ui
}
