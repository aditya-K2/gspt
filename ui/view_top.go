package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type TopTracksView struct {
	topTracks  []spotify.FullTrack
	topArtists []spotify.FullArtist
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
				{Content: v.Genres[0], Style: AlbumStyle},
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

func (a *TopTracksView) ContextOpener() func(m *Root, s func(s int)) { return nil }
func (a *TopTracksView) ContextHandler() func(int, int, int)         { return nil }
func (a *TopTracksView) ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEnter {
			r, _ := Ui.Main.Table.GetSelection()
			if r > 0 {
				if r < (len(a.topArtists) + 1) {
					if err := spt.PlayContext(&a.topArtists[r-1].URI); err != nil {
						SendNotification(err.Error())
					}
				} else if r != len(a.topArtists)+1 {
					if err := spt.PlaySong(a.topTracks[r-2-len(a.topArtists)].URI); err != nil {
						SendNotification(err.Error())
					}
				}
			}
		}
		return e
	}
}

func (a *TopTracksView) ContextKey() rune        { return 'a' }
func (a *TopTracksView) DisableVisualMode() bool { return true }
func (a *TopTracksView) Name() string            { return "AlbumsView" }
