package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/zmb3/spotify/v2"
)

type SearchContent struct {
	URI  spotify.URI
	ID   spotify.ID
	Name string
	Type string
}

type SearchView struct {
	*DefaultViewNone
	search        string
	results       *spotify.SearchResult
	searchContent []SearchContent
}

func NewSearchView() *SearchView {
	s := &SearchView{
		&DefaultViewNone{&defView{}},
		"", nil, []SearchContent{},
	}
	return s
}

func (a *SearchView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		generate := true
		if len(a.searchContent) != 0 {
			generate = false
		}
		if a.results != nil {
			if a.results.Tracks != nil {
				c = append(c, []Content{{"Tracks", NotSelectableStyle}})
				if generate {
					a.searchContent = append(a.searchContent, SearchContent{Type: "null"})
				}
				for _, v := range a.results.Tracks.Tracks {
					if generate {
						a.searchContent = append(a.searchContent, SearchContent{v.URI, v.ID, v.Name, "track"})
					}
					c = append(c, []Content{
						{Content: v.Name, Style: TrackStyle},
						{Content: artistName(v.Artists), Style: ArtistStyle},
						{Content: v.Album.Name, Style: AlbumStyle},
					})
				}
			}
			if a.results.Albums != nil {
				c = append(c, []Content{{"Albums", NotSelectableStyle}})
				if generate {
					a.searchContent = append(a.searchContent, SearchContent{Type: "null"})
				}
				for _, v := range a.results.Albums.Albums {
					if generate {
						a.searchContent = append(a.searchContent, SearchContent{v.URI, v.ID, v.Name, "album"})
					}
					c = append(c, []Content{
						{Content: v.Name, Style: AlbumStyle},
						{Content: artistName(v.Artists), Style: ArtistStyle},
						{Content: v.ReleaseDate, Style: TimeStyle},
					})
				}
			}
			if a.results.Artists != nil {
				c = append(c, []Content{{"Artists", NotSelectableStyle}})
				if generate {
					a.searchContent = append(a.searchContent, SearchContent{Type: "null"})
				}
				for _, v := range a.results.Artists.Artists {
					if generate {
						a.searchContent = append(a.searchContent, SearchContent{v.URI, v.ID, v.Name, "artist"})
					}
					c = append(c, []Content{
						{Content: v.Name, Style: AlbumStyle},
					})
				}
			}
			if a.results.Playlists != nil {
				c = append(c, []Content{{"Playlists", NotSelectableStyle}})
				if generate {
					a.searchContent = append(a.searchContent, SearchContent{Type: "null"})
				}
				for _, v := range a.results.Playlists.Playlists {
					if generate {
						a.searchContent = append(a.searchContent, SearchContent{v.URI, v.ID, v.Name, "playlist"})
					}
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

func (a *SearchView) OpenEntry() {
	r, _ := Main.GetSelection()
	switch a.searchContent[r].Type {
	case "track":
		{
			if err := spt.PlaySong(a.searchContent[r].URI); err != nil {
				SendNotification(err.Error())
			}
		}
	case "album":
		{
			albumView.SetAlbum(a.searchContent[r].Name, &a.searchContent[r].ID)
			SetCurrentView(albumView)
		}
	case "artist":
		{
			artistView.SetArtist(&a.searchContent[r].ID)
			SetCurrentView(artistView)
		}
	case "playlist":
		{
			if p, err := spt.GetFullPlaylist(a.searchContent[r].ID); err != nil {
				SendNotification("Error Opening the playlists: " + err.Error())
				return
			} else {
				playlistView.SetPlaylist(&p.SimplePlaylist)
				SetCurrentView(playlistView)
			}
		}
	}
}

func (a *SearchView) PlayEntry() {
	r, _ := Main.GetSelection()
	switch a.searchContent[r].Type {
	case "album", "artist", "playlist":
		{
			if err := spt.PlayContext(a.searchContent[r].URI); err != nil {
				SendNotification(err.Error())
			}
		}
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
	a.searchContent = make([]SearchContent, 0)
	a.RefreshState()
}

func (a *SearchView) Name() string { return "SearchView" }
