package gspotify

import (
	"context"
	"errors"

	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	SnapshotID string
	ID         spotify.ID
	Tracks     []spotify.PlaylistTrack
}

type Album struct {
	spotify.FullAlbum
	Tracks []spotify.SimpleTrack
}

type SavedAlbums []spotify.SavedAlbum
type UserPlaylists []spotify.SimplePlaylist

type Playable interface {
	Type() string
	Uri()
}

var (
	ctx           = func() context.Context { return context.Background() }
	Client        *spotify.Client
	playlistCache map[spotify.ID]*Playlist = make(map[spotify.ID]*Playlist)
	albumCache    map[spotify.ID]*Album    = make(map[spotify.ID]*Album)
	PageContinue                           = errors.New("CONTINUE")
)

func getPlaylistTracks(trackPage *spotify.PlaylistTrackPage) ([]spotify.PlaylistTrack, error) {
	tracks := make([]spotify.PlaylistTrack, 0)
	addTracks := func() {
		tracks = append(tracks, trackPage.Tracks...)
	}
	addTracks()
	for page := 1; ; page++ {
		if perr := Client.NextPage(ctx(), trackPage); perr == spotify.ErrNoMorePages {
			break
		} else if perr != nil {
			return nil, perr
		}
		addTracks()
	}
	return tracks, nil
}

func getAlbumTracks(trackPage *spotify.SimpleTrackPage) ([]spotify.SimpleTrack, error) {
	tracks := make([]spotify.SimpleTrack, 0)
	addTracks := func() {
		tracks = append(tracks, trackPage.Tracks...)
	}
	addTracks()
	for page := 1; ; page++ {
		if perr := Client.NextPage(ctx(), trackPage); perr == spotify.ErrNoMorePages {
			break
		} else if perr != nil {
			return nil, perr
		}
		addTracks()
	}
	return tracks, nil
}

func GetPlaylist(playlistId spotify.ID) (*Playlist, error) {
	if fp, err := Client.GetPlaylist(ctx(), playlistId); err != nil {
		return nil, err
	} else {
		if _, ok := playlistCache[fp.ID]; !ok || playlistCache[fp.ID].SnapshotID != fp.SnapshotID {
			if tracks, err := getPlaylistTracks(&fp.Tracks); err != nil {
				return nil, err
			} else {
				p := &Playlist{
					fp.SnapshotID,
					fp.ID,
					tracks,
				}
				playlistCache[fp.ID] = p
				return p, nil
			}
		} else {
			return playlistCache[fp.ID], nil
		}
	}
}

func GetAlbum(albumID spotify.ID) (*Album, error) {
	if fa, err := Client.GetAlbum(ctx(), albumID); err != nil {
		return nil, err
	} else {
		if _, ok := albumCache[fa.ID]; !ok {
			if tracks, err := getAlbumTracks(&fa.Tracks); err != nil {
				return nil, err
			} else {
				p := &Album{
					*fa,
					tracks,
				}
				albumCache[fa.ID] = p
				return p, nil
			}
		} else {
			return albumCache[fa.ID], nil
		}
	}
}

// CurrentUserSavedAlbums Returns the SavedAlbums in a very specific manner.
// It returns the first page and then starts a go routine in the background
// and keeps updating the SavedAlbums and sends nil to the provided channel
// if successful else sends the corresponding error.
func CurrentUserSavedAlbums(done func(status bool, err error)) (*SavedAlbums, error) {
	_a := make(SavedAlbums, 0)
	albums := &_a
	if sp, err := Client.CurrentUsersAlbums(ctx()); err != nil {
		return nil, err
	} else {
		addAlbums := func() {
			_a = append(_a, sp.Albums...)
		}
		addAlbums()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx(), sp); perr == spotify.ErrNoMorePages {
					done(true, nil)
					break
				} else if perr != nil {
					done(false, perr)
					break
				}
				addAlbums()
			}
		}()
		return albums, nil
	}
}

// CurrentUserPlaylists Returns the UserPlaylists in a very specific manner.
// It returns the first page and then starts a go routine in the background
// and keeps updating the UserPlaylists and sends nil to the provided channel
// if successful else sends the corresponding error.
func CurrentUserPlaylists(done func(status bool, err error)) (*UserPlaylists, error) {
	_p := make(UserPlaylists, 0)
	playlists := &_p
	if spp, err := Client.CurrentUsersPlaylists(ctx()); err != nil {
		return nil, err
	} else {
		addPlaylists := func() {
			_p = append(_p, spp.Playlists...)
		}
		addPlaylists()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx(), spp); perr == spotify.ErrNoMorePages {
					done(true, nil)
					break
				} else if perr != nil {
					done(false, perr)
					break
				}
				addPlaylists()
			}
		}()
		return playlists, nil
	}
}
