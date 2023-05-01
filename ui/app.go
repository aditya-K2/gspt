package ui

import (
	"time"

	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/tview"
	"github.com/aditya-K2/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	ImgX  int
	ImgY  int
	ImgW  int
	ImgH  int
	start = true
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

var (
	App      *tview.Application
	coverArt *CoverArt
	Main     *interactiveView
	root     *Root
	Flex     *tview.Flex
	cfg      = config.Config
)

func onConfigChange() {
	setStyles()
	setBorderRunes()
	if coverArt != nil && !cfg.HideImage {
		SendNotification("HERE")
		coverArt.RefreshState()
	}
}

func rectWatcher() {
	// Wait Until the ImagePreviewer is drawn
	// Ensures that cover art is not drawn before the UI is rendered.
	// Ref Issue: #39
	drawCh := make(chan bool)
	go func() {
		for ImgX == 0 && ImgY == 0 {
			ImgX, ImgY, ImgW, ImgH = coverArt.GetRect()
		}
		drawCh <- true
	}()

	go func() {
		// Waiting for the draw channel
		draw := <-drawCh
		if draw {
			go func() {
				for {
					_ImgX, _ImgY, _ImgW, _ImgH := coverArt.GetRect()
					if start {
						progressRoutine()
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
					time.Sleep(time.Millisecond * time.Duration(cfg.RedrawInterval))
				}
			}()
		}
	}()
}

func NewApplication() *tview.Application {
	onConfigChange()
	config.OnConfigChange = onConfigChange

	App = tview.NewApplication()
	root = NewRoot()
	coverArt = newCoverArt()
	Main = NewInteractiveView()
	Main.SetBorder(true)

	progressBar := NewProgressBar().SetProgressFunc(progressFunc)
	searchbar := NewSearchBar()
	navMenu := NewNavMenu([]navItem{
		{"Albums", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			SetCurrentView(albumsView)
			App.SetFocus(Main)
			return nil
		}, nil)},
		{"Artists", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			SetCurrentView(artistsView)
			App.SetFocus(Main)
			return nil
		}, nil)},
		{"Liked Songs", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			SetCurrentView(likedSongsView)
			App.SetFocus(Main)
			return nil
		}, nil)},
		{"Recently Played", NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			recentlyPlayedView.RefreshState()
			SetCurrentView(recentlyPlayedView)
			App.SetFocus(Main)
			return nil
		}, nil)},
	})
	playlistNav := NewPlaylistNav()

	root.AfterContextClose(func() { App.SetFocus(Main) })

	// Define Actions
	globalActions := map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			GetCurrentView().OpenEntry()
			return nil
		}, progressBar),
		"focus_search": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(searchbar)
			return nil
		}, nil),
		"toggle_playback": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			if err := spt.TogglePlayback(); err != nil {
				SendNotification(err.Error())
			}
			return nil
		}, progressBar),
		"choose_device": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			openDeviceMenu()
			return nil
		}, nil),
		"focus_nav": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(navMenu)
			return nil
		}, nil),
		"focus_playlists": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(playlistNav)
			return nil
		}, nil),
		"focus_main_view": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(Main)
			return nil
		}, nil),
		"open_current_track_album": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			openCurrentAlbum()
			return nil
		}, nil),
		"open_current_track_artist": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			openCurrentArtist()
			return nil
		}, nil),
		"open_current_context": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			if state != nil && state.Item != nil {
				switch state.PlaybackContext.Type {
				case "artist":
					{
						openCurrentArtist()
					}
				case "album":
					{
						openCurrentAlbum()
					}
				case "playlist":
					{
						id, err := spt.UriToID(state.PlaybackContext.URI)
						if err != nil {
							SendNotification("Error switching contexts: " + err.Error())
							return e
						}
						p, err := spt.GetFullPlaylist(id)
						if err != nil {
							SendNotification("Error switching contexts: " + err.Error())
							return e
						}
						playlistView.SetPlaylist(&p.SimplePlaylist)
						SetCurrentView(playlistView)
						App.SetFocus(Main)
					}
				default:
					{
						SendNotification("No Context Found!")
					}
				}
			}
			return nil
		}, nil),
		"next": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			if err := spt.Next(); err != nil {
				SendNotification(err.Error())
				return e
			}
			return nil
		}, progressBar),
		"previous": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			if err := spt.Previous(); err != nil {
				SendNotification(err.Error())
				return e
			}
			return nil
		}, progressBar),
	}
	playlistNav.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"play_entry": NewAction(playlistNav.PlayEntry,
			progressBar),
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			r, _ := playlistNav.GetSelection()
			playlistView.SetPlaylist(&(*playlistNav.Playlists)[r])
			SetCurrentView(playlistView)
			App.SetFocus(Main)
			return nil
		}, nil),
	}))
	navMenu.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(navMenu.OpenEntry,
			nil),
	}))
	playlistView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"add_to_playlist": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			playlistView.AddToPlaylist()
			return nil
		}, nil),
	}))
	recentlyPlayedView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"add_to_playlist": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			recentlyPlayedView.AddToPlaylist()
			return nil
		}, nil),
	}))
	topTracksView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			topTracksView.PlaySelectedEntry()
			return nil
		}, progressBar),
	}))
	likedSongsView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"add_to_playlist": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			likedSongsView.AddToPlaylist()
			return nil
		}, nil),
	}))
	searchView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			searchView.PlayEntry()
			return nil
		}, progressBar),
	}))
	artistsView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistsView.PlayEntry()
			return nil
		}, nil),
	}))
	artistView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.PlayEntry()
			return nil
		}, progressBar),
	}))
	albumsView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.PlayEntry()
			return nil
		}, progressBar),
		"queue_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.QueueSelectEntry()
			return nil
		}, progressBar),
	}))
	albumView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"add_to_playlist": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumView.AddToPlaylist()
			return nil
		}, nil),
	}))

	// Visual Actions
	albumView.SetVisualActions(map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey{
		"add_to_playlist": albumView.AddToPlaylistVisual,
	})
	recentlyPlayedView.SetVisualActions(map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey{
		"add_to_playlist": recentlyPlayedView.AddToPlaylistVisual,
	})
	playlistView.SetVisualActions(map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey{
		"add_to_playlist": playlistView.AddToPlaylistVisual,
	})
	likedSongsView.SetVisualActions(map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey{
		"add_to_playlist": likedSongsView.AddToPlaylistVisual,
	})

	mappings := config.GenerateMappings()

	// Map Actions
	playlistNav.SetMappings(mappings["playlist_nav"])
	playlistNav.SetInputCapture(playlistNav.ExternalInputCapture())
	navMenu.SetMappings(mappings["nav_menu"])
	navMenu.SetInputCapture(navMenu.ExternalInputCapture())
	playlistView.SetMappings(mappings["playlist_view"])
	recentlyPlayedView.SetMappings(mappings["recently_played_view"])
	topTracksView.SetMappings(mappings["top_tracks_view"])
	likedSongsView.SetMappings(mappings["liked_songs_view"])
	albumsView.SetMappings(mappings["albums_view"])
	albumView.SetMappings(mappings["album_view"])
	artistsView.SetMappings(mappings["artists_view"])
	artistView.SetMappings(mappings["artist_view"])
	searchView.SetMappings(mappings["search_view"])

	// Set up UI
	navFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(navMenu, 6, 3, false).
		AddItem(playlistNav, 0, 6, false)

	if !cfg.HideImage {
		navFlex.AddItem(coverArt, 9, 3, false)
	}

	// mid
	mFlex := tview.NewFlex().
		AddItem(navFlex, 17, 1, false).
		AddItem(Main, 0, 4, false)

	// mid + top
	tFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchbar, 3, 1, false).
		AddItem(mFlex, 0, 1, false)

	Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tFlex, 0, 8, false).
		AddItem(progressBar, 5, 1, false)

	root.Primitive("Main", Flex)
	App.SetRoot(root, true).SetFocus(Main)

	// Start Routines
	InitNotifier()
	if !cfg.HideImage {
		go rectWatcher()
	} else {
		// Start Progress Routine directly
		progressRoutine()
	}

	go func() {
		for {
			if App != nil {
				App.Draw()
				time.Sleep(time.Second)
			}
		}
	}()

	SetCurrentView(topTracksView)
	topTracksView.RefreshState()

	return App
}
