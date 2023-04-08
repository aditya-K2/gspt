package gspotify

import "github.com/zmb3/spotify/v2"

func AddAlbumToPlaylist(albumId, playlistId spotify.ID) error {
	a, err := GetAlbum(albumId)
	if err != nil {
		return err
	}
	t := []spotify.ID{}
	for _, v := range (*a).Tracks {
		t = append(t, v.ID)
	}
	_, err = Client.AddTracksToPlaylist(ctx(), playlistId, t...)
	return err
}
