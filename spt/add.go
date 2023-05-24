package spt

import (
	"errors"
	"fmt"

	"github.com/zmb3/spotify/v2"
)

func AddTracksToPlaylist(playlistId spotify.ID, t ...spotify.ID) error {
	_, err := Client.AddTracksToPlaylist(ctx, playlistId, t...)
	return err
}

func QueueTracks(ids ...spotify.ID) error {
	count := 0
	_ctx := ctx
	for _, id := range ids {
		if err := Client.QueueSong(_ctx, id); err != nil {
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

	ids := []spotify.ID{}
	for _, v := range *album.Tracks {
		ids = append(ids, v.ID)
	}

	if err := QueueTracks(ids...); err != nil {
		return err
	}
	return nil
}

func QueuePlaylist(id spotify.ID) error {
	playlist, c := GetPlaylist(id)
	if c == nil {
		return (errors.New("hi!"))
	}
	err := <-c
	if err != nil {
		return err
	}

	ids := []spotify.ID{}
	for _, v := range *playlist.Tracks {
		ids = append(ids, v.Track.ID)
	}

	if err := QueueTracks(ids...); err != nil {
		return err
	}
	return nil
}
