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
)

func setStyles() {
	TrackStyle = config.Config.Colors.Track.Style()
	AlbumStyle = config.Config.Colors.Album.Style()
	ArtistStyle = config.Config.Colors.Artist.Style()
	TimeStyle = config.Config.Colors.Timestamp.Style()
	GenreStyle = config.Config.Colors.Genre.Style()
	PlaylistNavStyle = config.Config.Colors.PlaylistNav.Style()
	NavStyle = config.Config.Colors.Nav.Style()
	ContextMenuStyle = config.Config.Colors.ContextMenu.Style()
	NotSelectableStyle = config.Config.Colors.Null.Style()
	if coverArt != nil {
		coverArt.RefreshState()
	}
	if config.Config.RoundedCorners {
		tview.Borders.TopLeft = '╭'
		tview.Borders.TopRight = '╮'
		tview.Borders.BottomRight = '╯'
		tview.Borders.BottomLeft = '╰'
		tview.Borders.Vertical = '│'
		tview.Borders.Horizontal = '─'
		tview.Borders.TopLeftFocus = '╭'
		tview.Borders.TopRightFocus = '╮'
		tview.Borders.BottomRightFocus = '╯'
		tview.Borders.BottomLeftFocus = '╰'
		tview.Borders.VerticalFocus = '│'
		tview.Borders.HorizontalFocus = '─'
		tview.Styles.BorderColorFocus = config.Config.Colors.BorderFocus.Foreground()
		tview.Styles.BorderColor = config.Config.Colors.Border.Foreground()
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
					time.Sleep(time.Millisecond * time.Duration(config.Config.RedrawInterval))
				}
			}()
		}
	}()
}

func NewApplication() *tview.Application {
	setStyles()
	config.OnConfigChange = setStyles

	App = tview.NewApplication()
	root = NewRoot()
	coverArt = newCoverArt()
	Main = NewInteractiveView()
	Main.Table.SetBorder(true)

	progressBar := NewProgressBar().SetProgressFunc(progressFunc)
	searchbar := NewSearchBar()
	navMenu := NewNavMenu([]navItem{
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
	playlistNav := NewPlaylistNav()

	root.AfterContextClose(func() { App.SetFocus(Main.Table) })

	// Define Actions
	openCurrentArtist := func() {
		if state != nil && state.Item != nil {
			if len(state.Item.Artists) != 0 {
				artistView.SetArtist(&state.Item.Artists[0].ID)
				SetCurrentView(artistView)
				App.SetFocus(Main.Table)
			} else {
				SendNotification("No Artist Found!")
			}
		}
	}
	openCurrentAlbum := func() {
		if state != nil && state.Item != nil {
			albumView.SetAlbum(state.Item.Album.Name, &state.Item.Album.ID)
			SetCurrentView(albumView)
			App.SetFocus(Main.Table)
		}
	}
	globalActions := map[string]*Action{
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
			OpenDeviceMenu()
			return nil
		}, nil),
		"focus_nav": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(navMenu.Table)
			return nil
		}, nil),
		"focus_playlists": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(playlistNav.Table)
			return nil
		}, nil),
		"focus_main_view": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			App.SetFocus(Main.Table)
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
						p, err := spt.GetFullPlaylist(&id)
						if err != nil {
							SendNotification("Error switching contexts: " + err.Error())
							return e
						}
						playlistView.SetPlaylist(&p.SimplePlaylist)
						SetCurrentView(playlistView)
						App.SetFocus(Main.Table)
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
		"play_entry": NewAction(playlistNav.PlaySelectEntry,
			progressBar),
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			r, _ := playlistNav.Table.GetSelection()
			playlistView.SetPlaylist(&(*playlistNav.Playlists)[r])
			SetCurrentView(playlistView)
			App.SetFocus(Main.Table)
			return nil
		}, nil),
	}))
	navMenu.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(navMenu.SelectEntry,
			nil),
	}))
	playlistView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(*tcell.EventKey) *tcell.EventKey {
			playlistView.PlaySelectEntry()
			return nil
		}, progressBar),
	}))
	recentlyPlayedView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(recentlyPlayedView.SelectEntry,
			progressBar),
	}))
	topTracksView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			topTracksView.OpenSelectEntry()
			return nil
		}, progressBar),
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			topTracksView.PlaySelectedEntry()
			return nil
		}, progressBar),
	}))
	likedSongsView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			likedSongsView.OpenEntry()
			return nil
		}, progressBar),
	}))
	searchView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			searchView.SelectEntry()
			return nil
		}, nil),
	}))
	artistsView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistsView.OpenArtist()
			return nil
		}, nil),
	}))
	artistView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.OpenEntry()
			return nil
		}, progressBar),
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			artistView.PlayEntry()
			return nil
		}, progressBar),
	}))
	albumsView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.OpenAlbum()
			return nil
		}, nil),
		"play_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.PlaySelectEntry()
			return nil
		}, progressBar),
		"queue_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumsView.QueueSelectEntry()
			return nil
		}, progressBar),
	}))
	albumView.SetActions(utils.MergeMaps(globalActions, map[string]*Action{
		"open_entry": NewAction(func(e *tcell.EventKey) *tcell.EventKey {
			albumView.PlaySelectEntry()
			return nil
		}, progressBar),
	}))

	mappings := config.GenerateMappings()

	// Map Actions
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

	// Set up UI
	navFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(navMenu.Table, 6, 3, false).
		AddItem(playlistNav.Table, 0, 6, false)

	if !config.Config.HideImage {
		navFlex.AddItem(coverArt, 9, 3, false)
	}

	// mid
	mFlex := tview.NewFlex().
		AddItem(navFlex, 17, 1, false).
		AddItem(Main.Table, 0, 4, false)

	// mid + top
	tFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchbar, 3, 1, false).
		AddItem(mFlex, 0, 1, false)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tFlex, 0, 8, false).
		AddItem(progressBar, 5, 1, false)

	root.Primitive("Main", mainFlex)
	App.SetRoot(root.Root, true).SetFocus(Main.Table)

	// Start Routines
	InitNotifier()
	if !config.Config.HideImage {
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
