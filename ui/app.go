package ui

import (
	"time"

	"github.com/aditya-K2/gspt/config"
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
	TrackStyle         tcell.Style
	AlbumStyle         tcell.Style
	ArtistStyle        tcell.Style
	TimeStyle          tcell.Style
	GenreStyle         tcell.Style
	PlaylistNavStyle   tcell.Style
	NavStyle           tcell.Style
	ContextMenuStyle   tcell.Style
	NotSelectableStyle tcell.Style
)

func loadStyles() {
	TrackStyle = config.Config.Colors.Track.Style()
	AlbumStyle = config.Config.Colors.Album.Style()
	ArtistStyle = config.Config.Colors.Artist.Style()
	TimeStyle = config.Config.Colors.Timestamp.Style()
	GenreStyle = config.Config.Colors.Genre.Style()
	PlaylistNavStyle = config.Config.Colors.PlaylistNav.Style()
	NavStyle = config.Config.Colors.Nav.Style()
	ContextMenuStyle = config.Config.Colors.ContextMenu.Style()
	NotSelectableStyle = config.Config.Colors.Null.Style()
	if Ui != nil {
		Ui.CoverArt.RefreshState()
	}
}

type Application struct {
	App            *tview.Application
	CoverArt       *CoverArt
	Main           *interactiveView
	NavMenu        *NavMenu
	PlaylistNav    *PlaylistNav
	SearchBar      *tview.InputField
	ProgressBar    *ProgressBar
	Root           *Root
	ImagePreviewer *tview.Box
}

func NewApplication() *Application {

	App := tview.NewApplication()
	Root := NewRoot()
	pBar := NewProgressBar().SetProgressFunc(progressFunc)
	coverArt := newCoverArt()
	searchbar := NewSearchBar()
	SetCurrentView(topTracksView)
	topTracksView.RefreshState()
	Main := NewInteractiveView()
	Main.Table.SetBorder(true)

	NavMenu := newNavMenu([]navItem{
		{"Albums", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			SetCurrentView(albumsView)
			App.SetFocus(Main.Table)
			return nil
		}, nil)},
		{"Artists", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			SetCurrentView(artistsView)
			App.SetFocus(Main.Table)
			return nil
		}, nil)},
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

	// Actions
	playlistNav.SetActions(map[string]*Action{
		"playEntry": NewAction(playlistNav.PlaySelectEntry, pBar),
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			r, _ := playlistNav.Table.GetSelection()
			playlistView.SetPlaylist(&(*playlistNav.Playlists)[r])
			SetCurrentView(playlistView)
			App.SetFocus(Main.Table)
			return nil
		}, nil),
	})
	playlistView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(*tcell.EventKey) *tcell.EventKey {
			playlistView.PlaySelectEntry()
			return nil
		}, pBar),
	})
	recentlyPlayedView.SetActions(map[string]*Action{
		"openEntry": NewAction(recentlyPlayedView.SelectEntry, pBar),
	})
	topTracksView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey { topTracksView.OpenSelectEntry(); return nil }, pBar),
		"playEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey { topTracksView.PlaySelectedEntry(); return nil }, pBar),
	})
	likedSongsView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			likedSongsView.OpenEntry()
			return nil
		}, pBar),
	})
	searchView.SetActions(map[string]*Action{})
	artistsView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistsView.OpenArtist()
			return nil
		}, nil),
	})
	artistView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.OpenEntry()
			return nil
		}, pBar),
		"playEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.PlayEntry()
			return nil
		}, pBar),
	})
	albumsView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.OpenAlbum()
			return nil
		}, nil),
		"playEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.PlaySelectEntry()
			return nil
		}, pBar),
	})
	albumView.SetActions(map[string]*Action{
		"openEntry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumView.PlaySelectEntry()
			return nil
		}, pBar),
	})

	// Mappings
	playlistNav.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
		tcell.KeyCtrlP: "playEntry",
	})
	playlistNav.Table.SetInputCapture(playlistNav.ExternalInputCapture())
	playlistView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})
	recentlyPlayedView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})
	topTracksView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
		tcell.KeyCtrlP: "playEntry",
	})
	likedSongsView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})
	albumsView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
		tcell.KeyCtrlP: "playEntry",
	})
	albumView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})
	artistsView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
	})
	artistView.SetMappings(map[tcell.Key]string{
		tcell.KeyEnter: "openEntry",
		tcell.KeyCtrlP: "playEntry",
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
	App.SetRoot(Root.Root, true).SetFocus(Main.Table)

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
							RefreshProgress(false)
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
		if e.Rune() == '?' {
			Ui.App.SetFocus(searchbar)
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
		if e.Rune() == 'd' {
			OpenDeviceMenu()
			return nil
		}
		return e
	})

	loadStyles()
	config.OnConfigChange = loadStyles

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
