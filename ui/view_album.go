package ui

import (
	"fmt"

	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type AlbumView struct {
	*PlaylistView
	currentAlbum     *spotify.SavedAlbum
	currentFullAlbum *spt.Album
}

func (a *AlbumView) SetAlbum(al *spotify.SavedAlbum) {
	a.currentAlbum = al
	a.currentFullAlbum = nil
}

func (a *AlbumView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)

		if a.currentAlbum != nil {
			if a.currentFullAlbum == nil {
				al, err := spt.GetAlbum(a.currentAlbum.ID, func(bool, error) {})
				if err != nil {
					panic(err)
				}
				a.currentFullAlbum = al
			}
			for _, v := range *(*a.currentFullAlbum).Tracks {
				ca := make([]Content, 0)
				ca = append(ca, Content{v.Name, TrackStyle})
				ca = append(ca, Content{v.Artists[0].Name, ArtistStyle})
				ca = append(ca, Content{a.currentAlbum.Name, AlbumStyle})
				c = append(c, ca)
			}
		}

		return c
	}
}

func (a *AlbumView) ContextHandler() func(start, end, sel int) {
	return func(start, end, sel int) {
		// Assuming that there are no external effects on the user's playlists
		// (i.e Any Creation or Deletion of Playlists while the context Menu is
		// open
		userPlaylists, err := spt.CurrentUserPlaylists(func(s bool, err error) {})
		if err != nil {
			SendNotification("Error Retrieving User Playlists")
			return
		}
		sp := (*userPlaylists)[sel]
		tracks := make([]spotify.ID, 0)
		for k := start; k <= end; k++ {
			tracks = append(tracks, (*(*a.currentFullAlbum).Tracks)[k].ID)
		}
		aerr := spt.AddTracksToPlaylist(sp.ID, tracks...)
		if aerr != nil {
			SendNotification(aerr.Error())
			return
		} else {
			SendNotification(fmt.Sprintf("Added %d tracks to %s", len(tracks), sp.Name))
		}
	}
}

func (a *AlbumView) ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEnter {
			r, _ := Ui.Main.Table.GetSelection()
			if err := spt.PlaySongWithContext(&a.currentAlbum.URI, r); err != nil {
				SendNotification(err.Error())
			}
		}
		return e
	}
}

func (a *AlbumView) Name() string { return "AlbumView" }
