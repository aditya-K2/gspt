package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	ImgX  int
	ImgY  int
	ImgW  int
	ImgH  int
	start = true
	Ui    *Application
)

var (
	TrackStyle       = tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorBlack)
	AlbumStyle       = tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	ArtistStyle      = tcell.StyleDefault.Foreground(tcell.ColorPink).Background(tcell.ColorBlack)
	TimeStyle        = tcell.StyleDefault.Foreground(tcell.ColorOrange).Background(tcell.ColorBlack)
	PlaylistNavStyle = tcell.StyleDefault.Foreground(tcell.ColorCoral).Background(tcell.ColorBlack)
	NavStyle         = tcell.StyleDefault.Foreground(tcell.ColorPapayaWhip).Background(tcell.ColorBlack)
	ContextMenuStyle = tcell.StyleDefault.Foreground(tcell.ColorPink).Background(tcell.ColorDefault).Bold(true)
)

type Application struct {
	App            *tview.Application
	CoverArt       *CoverArt
	Main           *interactiveView
	NavMenu        *NavMenu
	PlaylistNav    *PlaylistNav
	SearchBar      *tview.Box
	ProgressBar    *ProgressBar
	Root           *Root
	ImagePreviewer *tview.Box
}

func NewApplication() *Application {

	App := tview.NewApplication()
	Root := NewRoot()
	pBar := NewProgressBar().SetProgressFunc(progressFunc)
	coverArt := newCoverArt()
	searchbar := tview.NewBox().SetBorder(true).SetTitle("SEARCH").SetBackgroundColor(tcell.ColorDefault)
	SetCurrentView(playlistView)
	Main := NewInteractiveView()
	Main.Table.SetBorder(true)

	NavMenu := newNavMenu([]navItem{
		{"Albums", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			SetCurrentView(albumsView)
			App.SetFocus(Main.Table)
			return nil
		}, nil)},
		{"Artists", NewAction(func(e *tcell.EventKey) *tcell.EventKey { fmt.Println("Artists"); return nil }, nil)},
		{"Liked Songs", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
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

	NavMenu.Table.SetBackgroundColor(tcell.ColorDefault)

	NavMenu.Table.SetBorder(true)
	NavMenu.Table.SetSelectable(true, false)

	playlistNav, err := NewPlaylistNav(func(err error) {
		if err != nil {
			panic(err)
		}
		// Draw the App again after all the user playlists are retrieved.
		App.Draw()
	})
	if err != nil {
		panic(err)
	}

	Root.AfterContextClose(func() { App.SetFocus(Main.Table) })
	playlistNav.Table.SetBackgroundColor(tcell.ColorDefault)
	PlaylistActions = map[string]*Action{
		"playEntry": NewAction(playlistNav.PlaySelectEntry, pBar),
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
	recentlyPlayedView.MapActions(map[tcell.Key]string{
		tcell.KeyEnter: "selectEntry",
	})
	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(NavMenu.Table, 6, 3, false).
		AddItem(playlistNav.Table, 0, 6, false).
		AddItem(coverArt, 9, 3, false)

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
	App.SetRoot(Root.Root, true).SetFocus(playlistNav.Table)

	InitNotifier()
	updateRoutine()

	rectWatcher := func() {
		redrawInterval := 300

		// Wait Until the ImagePreviewer is drawn
		// Ensures that cover art is not drawn before the UI is rendered.
		// Ref Issue: #39
		drawCh := make(chan bool)
		go func() {
			for ImgX == 0 && ImgY == 0 {
				ImgX, ImgY, ImgW, ImgH = Ui.CoverArt.GetRect()
			}
			drawCh <- true

		}()

		go func() {
			// Waiting for the draw channel
			draw := <-drawCh
			if draw {
				go func() {
					for {
						_ImgX, _ImgY, _ImgW, _ImgH := Ui.CoverArt.GetRect()
						if start {
							RefreshProgress()
							start = false
						}
						if _ImgX != ImgX || _ImgY != ImgY ||
							_ImgW != ImgW || _ImgH != ImgH {
							ImgX = _ImgX
							ImgY = _ImgY
							ImgW = _ImgW
							ImgH = _ImgH
							coverArt.RefreshState()
						}
						time.Sleep(time.Millisecond * time.Duration(redrawInterval))
					}
				}()
			}
		}()
	}

	go rectWatcher()

	go func() {
		for {
			if Ui != nil && Ui.App != nil {
				Ui.App.Draw()
				time.Sleep(time.Second)
			}
		}
	}()

	App.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Rune() == '1' {
			Ui.App.SetFocus(NavMenu.Table)
			return nil
		}
		if e.Rune() == '2' {
			Ui.App.SetFocus(playlistNav.Table)
			return nil
		}
		if e.Rune() == '3' {
			Ui.App.SetFocus(Main.Table)
			return nil
		}
		return e
	})
	Ui = &Application{
		App:         App,
		Main:        Main,
		CoverArt:    coverArt,
		PlaylistNav: playlistNav,
		NavMenu:     NavMenu,
		SearchBar:   searchbar,
		ProgressBar: pBar,
		Root:        Root,
	}

	return Ui
}
