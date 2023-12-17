package ui

import (
	"fmt"

	"github.com/aditya-K2/gspt/spt"
	"github.com/zmb3/spotify/v2"
)

type TopTracksView struct {
	*DefaultViewNone
	topTracks  []spotify.FullTrack
	topArtists []spotify.FullArtist
}

func NewTopTracksView() *TopTracksView {
	t := &TopTracksView{
		&DefaultViewNone{&defView{}},
		[]spotify.FullTrack{},
		[]spotify.FullArtist{},
	}
	return t
}

func (a *TopTracksView) RefreshState() {
	topTracks, err := spt.GetTopTracks()
	if err != nil {
		SendNotification("Error retrieving Top Tracks: " + err.Error())
		return
	}
	a.topTracks = topTracks
	artists, err := spt.GetTopArtists()
	if err != nil {
		SendNotification("Error retrieving Top Artists: " + err.Error())
		return
	}
	a.topArtists = artists
}

func (a *TopTracksView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		c = append(c, []Content{{"Top Artists:", NotSelectableStyle}})
		for _, v := range a.topArtists {
			c = append(c, []Content{
				{Content: v.Name, Style: ArtistStyle},
				{Content: mergeGenres(v.Genres), Style: GenreStyle},
			})
		}
		c = append(c, []Content{{"Top Tracks:", NotSelectableStyle}})
		for _, v := range a.topTracks {
			c = append(c, []Content{
				{Content: v.Name, Style: TrackStyle},
				{Content: artistName(v.Artists), Style: ArtistStyle},
				{Content: v.Album.Name, Style: AlbumStyle},
			})
		}
		return c
	}
}

func (a *TopTracksView) handle(trackHandler, artistHandler func(r int)) {
	r, _ := Main.GetSelection()
	if r > 0 {
		if r < (len(a.topArtists) + 1) {
			if artistHandler != nil {
				artistHandler(r - 1)
			}
		} else if r != len(a.topArtists)+1 {
			if trackHandler != nil {
				trackHandler(r - 2 - len(a.topArtists))
			}
		}
	}
}

func (a *TopTracksView) PlayEntry() {
	a.handle(nil, func(r int) {
		if err := spt.PlayContext(a.topArtists[r].URI); err != nil {
			SendNotification(err.Error())
		}
	})

}

func (a *TopTracksView) OpenEntry() {
	artistHandler := func(r int) {
		artistView.SetArtist(&(a.topArtists)[r].ID)
		SetCurrentView(artistView)
	}
	trackHandler := func(r int) {
		if err := spt.PlaySong(a.topTracks[r].URI); err != nil {
			SendNotification(err.Error())
		}
	}
	a.handle(trackHandler, artistHandler)
}

func (a *TopTracksView) QueueEntry() {
	a.handle(func(r int) {
		msg := fmt.Sprintf("%s queued succesfully!", a.topTracks[r].Name)
		if err := spt.QueueTracks(a.topTracks[r].ID); err != nil {
			msg = err.Error()
		}
		SendNotification(msg)
	}, func(int) { SendNotification("Artists can not be queued!") })

}

func (a *TopTracksView) Name() string { return "TopTracksView" }
