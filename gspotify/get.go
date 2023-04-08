package gspotify

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	SnapshotID string
	Tracks     []spotify.PlaylistTrack
}

type Playable struct {
	ID   spotify.ID
	Name string
	Type string
}

var (
	playlistCache map[spotify.ID]*Playlist = make(map[spotify.ID]*Playlist)
	Client        *spotify.Client
	ctx           = context.Background()
)

func getPlaylistTracks(trackPage *spotify.PlaylistTrackPage) ([]spotify.PlaylistTrack, error) {
	tracks := make([]spotify.PlaylistTrack, 0)
	addTracks := func() {
		tracks = append(tracks, trackPage.Tracks...)
	}
	addTracks()
	for page := 1; ; page++ {
		if perr := Client.NextPage(ctx, trackPage); perr == spotify.ErrNoMorePages {
			break
		} else if perr != nil {
			return nil, perr
		}
		addTracks()
	}
	return tracks, nil
}

func GetPlaylist(playlistId spotify.ID) (*Playlist, error) {
	if fp, err := Client.GetPlaylist(ctx, playlistId); err != nil {
		return nil, err
	} else {
		if _, ok := playlistCache[fp.ID]; !ok || playlistCache[fp.ID].SnapshotID != fp.SnapshotID {
			if tracks, err := getPlaylistTracks(&fp.Tracks); err != nil {
				return nil, err
			} else {
				p := &Playlist{
					SnapshotID: fp.SnapshotID,
					Tracks:     tracks,
				}
				playlistCache[fp.ID] = p
				return p, nil
			}
		} else {
			return playlistCache[fp.ID], nil
		}
	}
}

func CurrentUserAllPlaylists() ([]spotify.SimplePlaylist, error) {
	playlists := make([]spotify.SimplePlaylist, 0)
	if spage, err := Client.CurrentUsersPlaylists(ctx); err != nil {
		return nil, err
	} else {
		addPlaylists := func() {
			playlists = append(playlists, spage.Playlists...)
		}
		addPlaylists()
		for page := 1; ; page++ {
			if perr := Client.NextPage(ctx, spage); perr == spotify.ErrNoMorePages {
				break
			} else if perr != nil {
				return nil, perr
			}
			addPlaylists()
		}
		return playlists, nil
	}
}
