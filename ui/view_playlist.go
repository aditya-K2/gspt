package ui

import (
	"fmt"

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
				msg := SendNotificationWithChan("Fetching %s....", p.currentPlaylist.Name)
				pf, ch := spt.GetPlaylist(p.currentPlaylist.ID)
				go func() {
					err := <-ch
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Playlist Fetched Succesfully!"
					}
				}()
				p.currentUserFullPlaylist = pf
			}
			if p.currentUserFullPlaylist != nil {
				for _, v := range *(*p.currentUserFullPlaylist).Tracks {
					c = append(c, []Content{
						{Content: v.Track.Name, Style: TrackStyle},
						{Content: artistName(v.Track.Artists), Style: ArtistStyle},
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
	addToPlaylist(Map((*(*p.currentUserFullPlaylist).Tracks)[start:end+1],
		func(s spotify.PlaylistTrack) spotify.ID {
			return s.Track.ID
		}))
	return nil
}

func (p *PlaylistView) OpenEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlaySongWithContext(p.currentPlaylist.URI, r); err != nil {
		SendNotification(err.Error())
	}
}

func (p *PlaylistView) QueueEntry() {
	r, _ := Main.GetSelection()
	track := (*(*p.currentUserFullPlaylist).Tracks)[r].Track
	msg := fmt.Sprintf("%s Queued Succesfully!", track.Name)
	if err := spt.QueueTracks(track.ID); err != nil {
		msg = err.Error()
	}
	SendNotification(msg)
}

func (p *PlaylistView) Name() string { return "PlaylistView" }
