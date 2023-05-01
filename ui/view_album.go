package ui

import (
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
				msg := SendNotificationWithChan("Loading %s....", a.currentAlbumName)
				al, ch := spt.GetAlbum(*a.currentAlbumID)
				go func() {
					err := <-ch
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Album Loaded Succesfully!"
					}
				}()
				a.currentFullAlbum = al
			}
			for _, v := range *(*a.currentFullAlbum).Tracks {
				ca := make([]Content, 0)
				ca = append(ca, Content{v.Name, TrackStyle})
				ca = append(ca, Content{v.Artists[0].Name, ArtistStyle})
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

func (a *AlbumView) AddToPlaylistVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	tracks := make([]spotify.ID, 0)
	sTracks := (*(*a.currentFullAlbum).Tracks)
	for k := start; k <= end; k++ {
		tracks = append(tracks, sTracks[k].ID)
	}
	addToPlaylist(tracks)
	return nil
}

func (a *AlbumView) PlayEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlaySongWithContext(&a.currentFullAlbum.URI, r); err != nil {
		SendNotification(err.Error())
	}
}

func (a *AlbumView) Name() string { return "AlbumView" }
