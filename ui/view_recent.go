package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

var (
	RecentlyPlayedViewActions = map[string]*Action{
		"selectEntry": NewAction(recentlyPlayedView.SelectEntry, nil),
	}
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
				{Content: v.Track.Artists[0].Name, Style: ArtistStyle},
				{Content: utils.StrTime(float64(v.Track.Duration / 1000)), Style: TimeStyle},
			})
		}
		return c
	}
}

func (r *RecentlyPlayedView) AddToPlaylist() {
	_r, _ := Main.Table.GetSelection()
	addToPlaylist([]spotify.ID{r.recentlyPlayed[_r].Track.ID})
}

func (r *RecentlyPlayedView) AddToPlaylistVisual(start, end int, e *tcell.EventKey) *tcell.EventKey {
	tracks := make([]spotify.ID, 0)
	for k := start; k <= end; k++ {
		tracks = append(tracks, r.recentlyPlayed[k].Track.ID)
	}
	addToPlaylist(tracks)
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

func (re *RecentlyPlayedView) SelectEntry(e *tcell.EventKey) *tcell.EventKey {
	r, _ := Main.Table.GetSelection()
	contextUri := re.recentlyPlayed[r].PlaybackContext.URI
	if string(contextUri) != "" {
		if err := spt.PlaySongWithContextURI(&re.recentlyPlayed[r].PlaybackContext.URI, &re.recentlyPlayed[r].Track.URI); err != nil {
			SendNotification(err.Error())
		}
	} else {
		if err := spt.PlaySong(re.recentlyPlayed[r].Track.URI); err != nil {
			SendNotification(err.Error())
		}
	}
	return nil
}
