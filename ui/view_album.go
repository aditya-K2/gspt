package ui

import (
	"fmt"

	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type AlbumView struct {
	*DefaultView
	currentAlbumID   *spotify.ID
	currentAlbumName string
	currentFullAlbum *spt.Album
}

func NewAlbumView() *AlbumView {
	t := &AlbumView{
		&DefaultView{&defView{}},
		nil, "", nil,
	}
	return t
}

func (a *AlbumView) SetAlbum(name string, al *spotify.ID) {
	a.currentAlbumID = al
	a.currentAlbumName = name
	a.currentFullAlbum = nil
}

func (a *AlbumView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)

		if a.currentAlbumID != nil {
			if a.currentFullAlbum == nil {
				msg := SendNotificationWithChan("Fetching %s....", a.currentAlbumName)
				al, ch := spt.GetAlbum(*a.currentAlbumID)
				go func() {
					err := <-ch
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Album Fetched Succesfully!"
					}
				}()
				a.currentFullAlbum = al
			}
			for _, v := range *(*a.currentFullAlbum).Tracks {
				ca := make([]Content, 0)
				ca = append(ca, Content{v.Name, TrackStyle})
				ca = append(ca, Content{artistName(v.Artists), ArtistStyle})
				ca = append(ca, Content{a.currentAlbumName, AlbumStyle})
				c = append(c, ca)
			}
		}

		return c
	}
}

func (a *AlbumView) AddToPlaylist() {
	r, _ := Main.GetSelection()
	track := (*(*a.currentFullAlbum).Tracks)[r]
	addToPlaylist([]spotify.ID{track.ID})
}

func (a *AlbumView) QueueEntry() {
	r, _ := Main.GetSelection()
	track := (*(*a.currentFullAlbum).Tracks)[r]
	msg := fmt.Sprintf("%s queued succesfully!", track.Name)
	if err := spt.QueueTracks(track.ID); err != nil {
		msg = err.Error()
	}
	SendNotification(msg)
}

func (a *AlbumView) AddToPlaylistVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	addToPlaylist(Map((*(*a.currentFullAlbum).Tracks)[start:end+1],
		func(s spotify.SimpleTrack) spotify.ID {
			return s.ID
		}))
	return nil
}

func (a *AlbumView) QueueSongsVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	tracks := (*(*a.currentFullAlbum).Tracks)[start : end+1]
	queueSongs(Map(tracks,
		func(s spotify.SimpleTrack) spotify.ID {
			return s.ID
		}))
	return nil
}

func (a *AlbumView) OpenEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlaySongWithContext(a.currentFullAlbum.URI, r); err != nil {
		SendNotification(err.Error())
	}
}

func (a *AlbumView) Name() string { return "AlbumView" }
