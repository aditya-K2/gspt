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

func (p *PlaylistView) SetPlaylist(pl *spotify.SimplePlaylist) {
	p.CurrentPlaylist = pl
	p.CurrentUserFullPlaylist = nil
}

func (p *PlaylistView) Content() [][]Content {
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

func (p *PlaylistView) ContextOpener(m *Main, s func(s int)) {
	c := NewMenu()
	cc := []string{}
	plist, err := spt.CurrentUserPlaylists(func(s bool, err error) {})
	if err != nil {
		panic(err)
	}
	for _, v := range *(plist) {
		cc = append(cc, v.Name)
	}
	c.Content(cc)
	c.Title("Add to Playlist")
	c.SetSelectionHandler(s)
	m.AddCenteredWidget(c)
}

func (p *PlaylistView) ContextHandler(start, end, sel int) {
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

func (p *PlaylistView) ExternalInputCapture(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyEnter {
		r, _ := Ui.MainS.View.GetSelection()
		if err := spt.PlaySongWithContext(&p.CurrentPlaylist.URI, r); err != nil {
			panic(err)
		}
	}
	return e
}

func (p *PlaylistView) ContextKey() rune {
	return 'a'
}

func (p *PlaylistView) Name() string { return "PlaylistView" }
