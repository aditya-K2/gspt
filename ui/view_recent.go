package ui

import (
	"fmt"

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

func (r *RecentlyPlayedView) ContextHandler() func(start, end, sel int) {
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
			tracks = append(tracks, r.recentlyPlayed[k].Track.ID)
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
