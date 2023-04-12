package spt

import (
	"github.com/zmb3/spotify/v2"
)

func play(options *spotify.PlayOptions) error {
	return Client.PlayOpt(ctx(), options)
}

func PlaySong(uri spotify.URI) error {
	return play(&spotify.PlayOptions{
		URIs: []spotify.URI{uri},
	})
}

func PlaySongWithContext(context *spotify.URI, position int) error {
	return play(&spotify.PlayOptions{
		PlaybackContext: context,
		PlaybackOffset:  &spotify.PlaybackOffset{Position: position},
	})
}

func PlaySongWithContextURI(context, uri *spotify.URI) error {
	return play(&spotify.PlayOptions{
		PlaybackContext: context,
		PlaybackOffset:  &spotify.PlaybackOffset{URI: *uri},
	})
}

func PlayContext(context *spotify.URI) error {
	return play(&spotify.PlayOptions{
		PlaybackContext: context,
	})
}
