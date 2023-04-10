package spt

import "github.com/zmb3/spotify/v2"

func AddTracksToPlaylist(playlistId spotify.ID, t ...spotify.ID) error {
	_, err := Client.AddTracksToPlaylist(ctx(), playlistId, t...)
	return err
}
