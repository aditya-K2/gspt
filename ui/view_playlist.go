package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type PlaylistView struct {
	*DefaultView
	currentPlaylist         *spotify.SimplePlaylist
	currentUserFullPlaylist *spt.Playlist
}

func NewPlaylistView() *PlaylistView {
	p := &PlaylistView{
		&DefaultView{&defView{}},
		nil,
		nil,
	}
	return p
}

func (p *PlaylistView) SetPlaylist(pl *spotify.SimplePlaylist) {
	p.currentPlaylist = pl
	p.currentUserFullPlaylist = nil
}

func (p *PlaylistView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if p.currentPlaylist != nil {
			if p.currentUserFullPlaylist == nil {
				msg := SendNotificationWithChan("Loading %s....", p.currentPlaylist.Name)
				pf, ch := spt.GetPlaylist(p.currentPlaylist.ID)
				go func() {
					err := <-ch
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Playlist Loaded Succesfully!"
					}
				}()
				p.currentUserFullPlaylist = pf
			}
			if p.currentUserFullPlaylist != nil {
				for _, v := range *(*p.currentUserFullPlaylist).Tracks {
					c = append(c, []Content{
						{Content: v.Track.Name, Style: TrackStyle},
						{Content: v.Track.Artists[0].Name, Style: ArtistStyle},
						{Content: v.Track.Album.Name, Style: AlbumStyle},
					})
				}
			}
		}
		return c
	}
}

func (p *PlaylistView) AddToPlaylist() {
	r, _ := Main.GetSelection()
	addToPlaylist([]spotify.ID{(*(*p.currentUserFullPlaylist).Tracks)[r].Track.ID})
}

func (p *PlaylistView) AddToPlaylistVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	tracks := make([]spotify.ID, 0)
	sTracks := (*(*p.currentUserFullPlaylist).Tracks)
	for k := start; k <= end; k++ {
		tracks = append(tracks, sTracks[k].Track.ID)
	}
	addToPlaylist(tracks)
	return nil
}

func (p *PlaylistView) OpenEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlaySongWithContext(&p.currentPlaylist.URI, r); err != nil {
		SendNotification(err.Error())
	}
}

func (p *PlaylistView) Name() string { return "PlaylistView" }
