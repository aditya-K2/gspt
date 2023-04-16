package ui

import (
	"time"

	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/utils"
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

	navMenu := newNavMenu([]navItem{
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

	navMenu.Table.SetBackgroundColor(tcell.ColorDefault)

	navMenu.Table.SetBorder(true)
	navMenu.Table.SetSelectable(true, false)

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
	globalMaps := map[string]*Action{
		"focus_search": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			Ui.App.SetFocus(searchbar)
			return nil
		}, nil),
		"choose_device": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			OpenDeviceMenu()
			return nil
		}, nil),
		"focus_nav": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			Ui.App.SetFocus(navMenu.Table)
			return nil
		}, nil),
		"focus_playlists": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			Ui.App.SetFocus(playlistNav.Table)
			return nil
		}, nil),
		"focus_main_view": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			Ui.App.SetFocus(Main.Table)
			return nil
		}, nil),
	}

	// Actions
	playlistNav.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"play_entry": NewAction(playlistNav.PlaySelectEntry, pBar),
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			r, _ := playlistNav.Table.GetSelection()
			playlistView.SetPlaylist(&(*playlistNav.Playlists)[r])
			SetCurrentView(playlistView)
			App.SetFocus(Main.Table)
			return nil
		}, nil),
	}))
	navMenu.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(navMenu.SelectEntry, nil),
	}))
	playlistView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(*tcell.EventKey) *tcell.EventKey {
			playlistView.PlaySelectEntry()
			return nil
		}, pBar),
	}))
	recentlyPlayedView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(recentlyPlayedView.SelectEntry, pBar),
	}))
	topTracksView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey { topTracksView.OpenSelectEntry(); return nil }, pBar),
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey { topTracksView.PlaySelectedEntry(); return nil }, pBar),
	}))
	likedSongsView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			likedSongsView.OpenEntry()
			return nil
		}, pBar),
	}))
	searchView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			searchView.SelectEntry()
			return nil
		}, nil),
	}))
	artistsView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistsView.OpenArtist()
			return nil
		}, nil),
	}))
	artistView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.OpenEntry()
			return nil
		}, pBar),
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.PlayEntry()
			return nil
		}, pBar),
	}))
	albumsView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.OpenAlbum()
			return nil
		}, nil),
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.PlaySelectEntry()
			return nil
		}, pBar),
	}))
	albumView.SetActions(utils.MergeMaps(globalMaps, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumView.PlaySelectEntry()
			return nil
		}, pBar),
	}))

	mappings := config.GenerateMappings()

	// Mappings
	playlistNav.SetMappings(mappings["playlist_nav"])
	playlistNav.Table.SetInputCapture(playlistNav.ExternalInputCapture())
	navMenu.SetMappings(mappings["nav_menu"])
	navMenu.Table.SetInputCapture(navMenu.ExternalInputCapture())
	playlistView.SetMappings(mappings["playlist_view"])
	recentlyPlayedView.SetMappings(mappings["recently_played_view"])
	topTracksView.SetMappings(mappings["top_tracks_view"])
	likedSongsView.SetMappings(mappings["liked_songs_view"])
	albumsView.SetMappings(mappings["albums_view"])
	albumView.SetMappings(mappings["album_view"])
	artistsView.SetMappings(mappings["artists_view"])
	artistView.SetMappings(mappings["artist_view"])
	searchView.SetMappings(mappings["search_view"])

	searchNavFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(navMenu.Table, 6, 3, false).
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

	loadStyles()
	config.OnConfigChange = loadStyles

	Ui = &Application{
		App:         App,
		Main:        Main,
		CoverArt:    coverArt,
		PlaylistNav: playlistNav,
		NavMenu:     navMenu,
		SearchBar:   searchbar,
		ProgressBar: pBar,
		Root:        Root,
	}

	return Ui
}
