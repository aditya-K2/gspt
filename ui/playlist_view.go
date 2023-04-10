package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type PlaylistView struct {
	CurrentPlaylist         *spotify.SimplePlaylist
	CurrentUserFullPlaylist *spt.Playlist
	I                       *interactiveView
}

func NewPlaylistView(m *Main) *PlaylistView {
	i := NewInteractiveView()

	i.View.SetBorder(true)
	p := &PlaylistView{nil, nil, i}
	content := func() [][]Content {
		c := make([][]Content, 0)
		if p.CurrentPlaylist != nil {
			if p.CurrentUserFullPlaylist == nil {
				pf, err := spt.GetPlaylist(p.CurrentPlaylist.ID, func(bool, error) {})
				if err != nil {
					panic(err)
				}
				p.CurrentUserFullPlaylist = pf
			}
			for _, v := range *(*p.CurrentUserFullPlaylist).Tracks {
				c = append(c, []Content{
					{Content: v.Track.Name, Style: Defaultstyle.Foreground(tcell.ColorBlue)},
					{Content: v.Track.Artists[0].Name, Style: Defaultstyle.Foreground(tcell.ColorPink)},
					{Content: v.Track.Album.Name, Style: Defaultstyle.Foreground(tcell.ColorGreen)},
				})
			}
		} else {
			return [][]Content{{
				{"HElll", tcell.StyleDefault},
				{"HElll", tcell.StyleDefault},
				{"HElll", tcell.StyleDefault},
				{"HElll", tcell.StyleDefault},
			}}
		}
		return c
	}
	i.SetContentFunc(content)
	contextOpener := func() {
		c := NewMenu()
		cc := []string{}
		p, err := spt.CurrentUserPlaylists(func(s bool, err error) {})
		if err != nil {
			panic(err)
		}
		for _, v := range *(p) {
			cc = append(cc, v.Name)
		}
		c.Content(cc)
		c.Title("Add to Playlist")
		c.SetSelectionHandler(i.SelectionHandler)
		m.AddCenteredWidget(c)
	}
	i.SetContextKey('a')
	i.SetContextOpener(contextOpener)

	contextHandler := func(start, end, sel int) {
		// Assuming that there are no external effects on the user's playlists
		// (i.e Any Creation or Deletion of Playlists while the context Menu is
		// open
		ap, err := spt.CurrentUserPlaylists(func(s bool, err error) {})
		if err != nil {
			panic(err)
		}
		p.CurrentPlaylist = &(*ap)[sel]
		tracks := make([]spotify.ID, 0)
		for k := start; k <= end; k++ {
			tracks = append(tracks, (*(*p.CurrentUserFullPlaylist).Tracks)[k].Track.ID)
		}
		if err := spt.AddTracksToPlaylist((*ap)[sel].ID, tracks...); err != nil {
			panic(err)
		}
	}
	i.SetContextHandler(contextHandler)
	i.SetExternalCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEnter {
			r, _ := i.View.GetSelection()
			if err := spt.PlaySongWithContext(&p.CurrentPlaylist.URI, r); err != nil {
				panic(err)
			}
		}
		return e
	})

	return p
}
