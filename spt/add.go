package spt

import (
	"errors"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func AddTracksToPlaylist(playlistId spotify.ID, t ...spotify.ID) error {
	_, err := Client.AddTracksToPlaylist(ctx(), playlistId, t...)
	return err
}

func QueueTracks(ids ...spotify.SimpleTrack) error {
	count := 0
	_ctx := ctx()
	for _, v := range ids {
		if err := Client.QueueSong(_ctx, v.ID); err != nil {
			return errors.New(fmt.Sprintf("%s | Tracks Queued: %d", err.Error(), count))
		}
	}
	return nil
}

func QueueAlbum(id spotify.ID) error {
	album, c := GetAlbum(id)
	if c == nil {
		return (errors.New("hi!"))
	}
	err := <-c
	if err != nil {
		return err
	}
	if err := QueueTracks(*album.Tracks...); err != nil {
		return err
	}
	return nil
}
