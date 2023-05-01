package ui

import (
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

func mergeGenres(g []string) string {
	s := ""
	for k, v := range g {
		sep := ","
		if k == 0 {
			sep = ""
		}
		s += sep + v
	}
	return s
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
				{Content: v.Artists[0].Name, Style: ArtistStyle},
				{Content: v.Album.Name, Style: AlbumStyle},
			})
		}
		return c
	}
}

func (a *TopTracksView) PlaySelectedEntry() {
	r, _ := Main.GetSelection()
	if r > 0 {
		if r < (len(a.topArtists) + 1) {
			if err := spt.PlayContext(a.topArtists[r-1].URI); err != nil {
				SendNotification(err.Error())
			}
		}
	}

}

func (a *TopTracksView) OpenEntry() {
	r, _ := Main.GetSelection()
	if r > 0 {
		if r < (len(a.topArtists) + 1) {
			artistView.SetArtist(&(a.topArtists)[r-1].ID)
			SetCurrentView(artistView)
		} else if r != len(a.topArtists)+1 {
			if err := spt.PlaySong(a.topTracks[r-2-len(a.topArtists)].URI); err != nil {
				SendNotification(err.Error())
			}
		}
	}
}

func (a *TopTracksView) Name() string { return "TopTracksView" }
