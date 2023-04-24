package spt

import (
	"errors"
	"strings"

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
		PlaybackOffset:  &spotify.PlaybackOffset{Position: &position},
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

func TogglePlayback() error {
	p, err := Client.PlayerCurrentlyPlaying(ctx())
	if err != nil {
		return err
	}
	if p.Playing {
		if err := Client.Pause(ctx()); err != nil {
			return err
		}
	} else {
		if err := Client.Play(ctx()); err != nil {
			return err
		}
	}
	return nil
}

func UriToID(uri spotify.URI) (spotify.ID, error) {
	a := strings.Split(string(uri), ":")
	if len(a) != 3 {
		return "", errors.New("Error Decoding the URI")
	}
	return spotify.ID(a[2]), nil
}
