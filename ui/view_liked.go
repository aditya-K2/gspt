package ui

import (
	"fmt"

	"github.com/aditya-K2/gspt/spt"
	"github.com/zmb3/spotify/v2"
)

type LikedSongsView struct {
	*DefaultView
	likedSongs *spt.LikedSongs
}

func NewLikedSongsView() *LikedSongsView {
	l := &LikedSongsView{
		&DefaultView{&defView{}},
		nil,
	}
	return l
}

func (p *LikedSongsView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if p.likedSongs == nil {
			msg := SendNotificationWithChan("Loading Liked Songs...")
			p.refreshState(func(err error) {
				if err != nil {
					msg <- err.Error()
					return
				}
				msg <- "Liked Songs Loaded Succesfully!"
			})
		}
		if p.likedSongs != nil {
			for _, v := range *p.likedSongs {
				c = append(c, []Content{
					{Content: v.Name, Style: TrackStyle},
					{Content: v.Artists[0].Name, Style: ArtistStyle},
					{Content: v.Album.Name, Style: AlbumStyle},
				})
			}
		}
		return c
	}
}

func (l *LikedSongsView) ContextHandler() func(start, end, sel int) {
	return func(start, end, sel int) {
		// Assuming that there are no external effects on the user's playlists
		// (i.e Any Creation or Deletion of Playlists while the context Menu is
		// open
		// TODO: Better Error Handler
		userPlaylists, ch := spt.CurrentUserPlaylists()
		err := <-ch
		if err != nil {
			SendNotification("Error Retrieving User Playlists")
			return
		}
		tracks := make([]spotify.ID, 0)
		for k := start; k <= end; k++ {
			tracks = append(tracks, (*l.likedSongs)[k].ID)
		}
		aerr := spt.AddTracksToPlaylist((*userPlaylists)[sel].ID, tracks...)
		if aerr != nil {
			SendNotification(aerr.Error())
			return
		} else {
			SendNotification(fmt.Sprintf("Added %d tracks to %s", len(tracks), (*userPlaylists)[sel].Name))
		}
	}
}

func (l *LikedSongsView) OpenEntry() {
	r, _ := Main.Table.GetSelection()
	if err := spt.PlaySong((*l.likedSongs)[r].URI); err != nil {
		SendNotification(err.Error())
	}
}

func (l *LikedSongsView) Name() string { return "LikedSongsView" }

func (p *LikedSongsView) refreshState(errHandler func(error)) {
	cl, ch := spt.CurrentUserSavedTracks()
	p.likedSongs = cl
	go func() {
		err := <-ch
		errHandler(err)
		if err == nil {
			p.likedSongs = cl
		}
	}()
}

func (l *LikedSongsView) RefreshState() {
	// TODO: Better Error Handler
	l.refreshState(func(err error) {
		if err != nil {
			SendNotification(err.Error())
		}
	})
}
