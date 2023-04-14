package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/zmb3/spotify/v2"
)

type SearchView struct {
	*DefaultViewNone
	search  string
	results *spotify.SearchResult
}

func NewSearchView() *SearchView {
	s := &SearchView{
		&DefaultViewNone{&defView{}},
		"", nil,
	}
	return s
}

func (a *SearchView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if a.results != nil {
			c = append(c, []Content{{"Tracks", NotSelectableStyle}})
			if a.results.Tracks != nil {
				for _, v := range a.results.Tracks.Tracks {
					c = append(c, []Content{
						{Content: v.Name, Style: TrackStyle},
						{Content: v.Artists[0].Name, Style: ArtistStyle},
						{Content: v.Album.Name, Style: AlbumStyle},
					})
				}
			}
			if a.results.Albums != nil {
				c = append(c, []Content{{"Albums", NotSelectableStyle}})
				for _, v := range a.results.Albums.Albums {
					c = append(c, []Content{
						{Content: v.Name, Style: AlbumStyle},
						{Content: v.Artists[0].Name, Style: ArtistStyle},
						{Content: v.ReleaseDate, Style: TimeStyle},
					})
				}
			}
			if a.results.Artists != nil {
				c = append(c, []Content{{"Artists", NotSelectableStyle}})
				for _, v := range a.results.Artists.Artists {
					c = append(c, []Content{
						{Content: v.Name, Style: AlbumStyle},
					})
				}
			}
			if a.results.Playlists != nil {
				c = append(c, []Content{{"Playlists", NotSelectableStyle}})
				for _, v := range a.results.Playlists.Playlists {
					c = append(c, []Content{
						{Content: v.Name, Style: PlaylistNavStyle},
						{Content: v.Owner.DisplayName, Style: ArtistStyle},
					})
				}
			}
		}
		return c
	}
}

func (a *SearchView) RefreshState() {
	if a.search != "" {
		results, err := spt.Search(a.search)
		if err != nil {
			SendNotification("Error retrieving Artist Top Tracks: " + err.Error())
			return
		}
		a.results = results
	}
}

func (a *SearchView) SetSearch(s string) {
	a.search = s
	a.RefreshState()
}

func (a *SearchView) Name() string { return "SearchView" }
