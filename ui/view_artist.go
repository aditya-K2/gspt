package ui

import "github.com/zmb3/spotify/v2"

type ArtistView struct {
	topTracks []spotify.FullTrack
	albums    []spotify.SimpleAlbum
}
