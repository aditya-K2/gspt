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
			if err := p.refreshState(func(err error) {
				go func() {
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Liked Songs Loaded Succesfully!"
					}
				}()
			}); err != nil {
				SendNotification(err.Error())
				return c
			}
		}
		for _, v := range *p.likedSongs {
			c = append(c, []Content{
				{Content: v.Name, Style: TrackStyle},
				{Content: v.Artists[0].Name, Style: ArtistStyle},
				{Content: v.Album.Name, Style: AlbumStyle},
			})
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
		userPlaylists, err := spt.CurrentUserPlaylists(func(err error) {})
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
	r, _ := Ui.Main.Table.GetSelection()
	if err := spt.PlaySong((*l.likedSongs)[r].URI); err != nil {
		SendNotification(err.Error())
	}
}

func (l *LikedSongsView) refreshState(f func(error)) error {
	pf, err := spt.CurrentUserSavedTracks(f)
	if err == nil {
		l.likedSongs = pf
	}
	return err
}

func (l *LikedSongsView) Name() string { return "LikedSongsView" }

func (l *LikedSongsView) RefreshState() {
	// TODO: Better Error Handler
	if err := l.refreshState(func(error) {}); err != nil {
		SendNotification(err.Error())
	}
}
