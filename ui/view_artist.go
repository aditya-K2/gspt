package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/zmb3/spotify/v2"
)

type ArtistView struct {
	*DefaultViewNone
	artistID  *spotify.ID
	topTracks []spotify.FullTrack
	albums    []spotify.SimpleAlbum
}

func NewArtistView() *ArtistView {
	a := &ArtistView{
		&DefaultViewNone{&defView{}},
		nil,
		[]spotify.FullTrack{},
		[]spotify.SimpleAlbum{},
	}
	return a
}

func (a *ArtistView) SetArtist(id *spotify.ID) {
	a.artistID = id
	a.topTracks = []spotify.FullTrack{}
	a.albums = []spotify.SimpleAlbum{}
	go func() {
		a.RefreshState()
	}()
}

func (a *ArtistView) RefreshState() {
	msg := SendNotificationWithChan("Fetching Artist....")
	topTracks, err := spt.GetArtistTopTracks(*a.artistID)
	if err != nil {
		msg <- ("Error retrieving Artist Top Tracks: " + err.Error())
		return
	}
	a.topTracks = topTracks
	albums, err := spt.GetArtistAlbums(*a.artistID)
	if err != nil {
		msg <- ("Error retrieving Artist Albums: " + err.Error())
		return
	}
	a.albums = albums
	msg <- "Artist Fetched Succesfully!"
}

func (a *ArtistView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		c = append(c, []Content{{"Artist Albums: ", NotSelectableStyle}})
		for _, v := range a.albums {
			c = append(c, []Content{
				{Content: v.Name, Style: AlbumStyle},
				{Content: artistName(v.Artists), Style: ArtistStyle},
				{Content: v.ReleaseDate, Style: TimeStyle},
			})
		}
		c = append(c, []Content{{"Artist Top Tracks:", NotSelectableStyle}})
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

func (a *ArtistView) handle(albumHandler func(int), trackHandler func(int)) {
	r, _ := Main.GetSelection()
	if r > 0 {
		if r < (len(a.albums)+1) && len(a.albums) > 0 {
			albumHandler(r - 1)
		} else if r != len(a.albums)+1 && len(a.topTracks) > 0 {
			trackHandler(r - 2 - len(a.albums))
		}
	}
}

func (a *ArtistView) PlayEntry() {
	a.handle(func(r int) {
		if err := spt.PlayContext(a.albums[r].URI); err != nil {
			SendNotification(err.Error())
		}
	}, func(int) {})
}

func (a *ArtistView) OpenEntry() {
	a.handle(func(r int) {
		albumView.SetAlbum(a.albums[r].Name, &a.albums[r].ID)
		SetCurrentView(albumView)
	}, func(r int) {
		if err := spt.PlaySong(a.topTracks[r].URI); err != nil {
			SendNotification(err.Error())
		}
	})
}

func (a *ArtistsView) QueueEntry() {
}

func (a *ArtistView) Name() string { return "AlbumsView" }
