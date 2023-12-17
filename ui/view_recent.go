package ui

import (
	"fmt"

	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type RecentlyPlayedView struct {
	*DefaultView
	recentlyPlayed []spotify.RecentlyPlayedItem
}

func NewRecentlyPlayedView() *RecentlyPlayedView {
	r := &RecentlyPlayedView{
		&DefaultView{&defView{}},
		[]spotify.RecentlyPlayedItem{},
	}
	return r
}

func (r *RecentlyPlayedView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		for _, v := range r.recentlyPlayed {
			c = append(c, []Content{
				{Content: v.Track.Name, Style: TrackStyle},
				{Content: artistName(v.Track.Artists), Style: ArtistStyle},
				{Content: utils.StrTime(float64(v.Track.Duration / 1000)), Style: TimeStyle},
			})
		}
		return c
	}
}

func (r *RecentlyPlayedView) AddToPlaylist() {
	_r, _ := Main.GetSelection()
	addToPlaylist([]spotify.ID{r.recentlyPlayed[_r].Track.ID})
}

func (r *RecentlyPlayedView) AddToPlaylistVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	addToPlaylist(Map(r.recentlyPlayed[start:end+1],
		func(r spotify.RecentlyPlayedItem) spotify.ID {
			return r.Track.ID
		}))
	return nil
}

func (r *RecentlyPlayedView) QueueSongsVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	tracks := r.recentlyPlayed[start : end+1]
	msg := SendNotificationWithChan(fmt.Sprintf("Queueing %d tracks...", len(tracks)))
	go func() {
		err := spt.QueueTracks(Map(tracks,
			func(s spotify.RecentlyPlayedItem) spotify.ID {
				return s.Track.ID
			})...)
		if err != nil {
			msg <- err.Error()
		}
		msg <- fmt.Sprintf("Queued %d tracks!", len(tracks))
	}()
	return nil
}

func (r *RecentlyPlayedView) Name() string { return "RecentlyPlayedView" }

func (r *RecentlyPlayedView) RefreshState() {
	_r, err := spt.RecentlyPlayed()
	if err != nil {
		SendNotification(err.Error())
		return
	}
	r.recentlyPlayed = _r
}

func (re *RecentlyPlayedView) OpenEntry() {
	r, _ := Main.GetSelection()
	trackUri := re.recentlyPlayed[r].Track.URI
	contextUri := re.recentlyPlayed[r].PlaybackContext.URI
	if string(contextUri) != "" {
		if err := spt.PlaySongWithContextURI(contextUri, trackUri); err != nil {
			SendNotification(err.Error())
		}
	} else {
		if err := spt.PlaySong(trackUri); err != nil {
			SendNotification(err.Error())
		}
	}
}

func (re *RecentlyPlayedView) QueueEntry() {
	r, _ := Main.GetSelection()
	track := re.recentlyPlayed[r].Track
	msg := fmt.Sprintf("%s Queued Succesfully!", track.Name)
	if err := spt.QueueTracks(track.ID); err != nil {
		msg = err.Error()
	}
	SendNotification(msg)
}
