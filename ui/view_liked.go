package ui

import (
	"fmt"

	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
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
			msg := SendNotificationWithChan("Fetching Liked Songs...")
			p.refreshState(func(err error) {
				if err != nil {
					msg <- err.Error()
					return
				}
				msg <- "Liked Songs Fetched Succesfully!"
			})
		}
		if p.likedSongs != nil {
			for _, v := range *p.likedSongs {
				c = append(c, []Content{
					{Content: v.Name, Style: TrackStyle},
					{Content: artistName(v.Artists), Style: ArtistStyle},
					{Content: v.Album.Name, Style: AlbumStyle},
				})
			}
		}
		return c
	}
}

func (l *LikedSongsView) AddToPlaylist() {
	r, _ := Main.GetSelection()
	addToPlaylist([]spotify.ID{(*l.likedSongs)[r].ID})
}

func (l *LikedSongsView) AddToPlaylistVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	addToPlaylist(Map((*l.likedSongs)[start:end+1],
		func(s spotify.SavedTrack) spotify.ID {
			return s.ID
		}))
	return nil
}

func (l *LikedSongsView) OpenEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlaySong((*l.likedSongs)[r].URI); err != nil {
		SendNotification(err.Error())
	}
}

func (l *LikedSongsView) QueueEntry() {
	r, _ := Main.GetSelection()
	track := (*l.likedSongs)[r]
	msg := fmt.Sprintf("%s Queued Succesfully!", track.Name)
	if err := spt.QueueTracks(track.ID); err != nil {
		msg = err.Error()
	}
	SendNotification(msg)
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
