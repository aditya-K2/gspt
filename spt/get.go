package spt

import (
	"context"
	"errors"

	"github.com/zmb3/spotify/v2"
)

type Playlist struct {
	SnapshotID string
	ID         spotify.ID
	Tracks     *[]spotify.PlaylistTrack
}

type Album struct {
	spotify.FullAlbum
	Tracks *[]spotify.SimpleTrack
}

type SavedAlbums []spotify.SavedAlbum
type UserPlaylists []spotify.SimplePlaylist
type LikedSongs []spotify.SavedTrack

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

// GetPlaylist retrieves a Spotify playlist by its ID and returns a
// Playlist object. If the playlist is already cached, it returns cached object.
// If not cached or if the snapshotID for the playlist changes, it retrieves
// the playlist from Spotify API and caches it. It uses a background go routine
// to fetch all pages of tracks for the playlist and appends them to the tracks
// list in the Playlist object. When done, it sends a true value to the done
// callback if successful, otherwise an error.
func GetPlaylist(playlistId spotify.ID, done func(bool, error)) (*Playlist, error) {
	if fp, err := Client.GetPlaylist(ctx(), playlistId); err != nil {
		return nil, err
	} else {
		if _, ok := playlistCache[fp.ID]; !ok || playlistCache[fp.ID].SnapshotID != fp.SnapshotID {
			tracks := &[]spotify.PlaylistTrack{}
			addTracks := func() {
				*tracks = append(*tracks, fp.Tracks.Tracks...)
			}
			addTracks()
			go func() {
				for page := 1; ; page++ {
					if perr := Client.NextPage(ctx(), &fp.Tracks); perr == spotify.ErrNoMorePages {
						done(true, nil)
						break
					} else if perr != nil {
						done(false, perr)
						return
					}
					addTracks()
				}
			}()
			p := &Playlist{
				fp.SnapshotID,
				fp.ID,
				tracks,
			}
			playlistCache[fp.ID] = p
			return p, nil
		} else {
			done(true, nil)
			return playlistCache[fp.ID], nil
		}
	}
}

// GetAlbum retrieves a Spotify album by its ID and returns an Album object.
// If the album is already cached, it returns the cached object.
// If not, it retrieves the album from the Spotify API and caches it.
// It uses a background go routine to fetch all pages of tracks for the album
// and appends them to the tracks list in the Album object.
// When done, it sends a true value to the done callback if successful
func GetAlbum(albumID spotify.ID, done func(bool, error)) (*Album, error) {
	if _, ok := albumCache[albumID]; !ok {
		fa, err := Client.GetAlbum(ctx(), albumID)
		if err != nil {
			return nil, err
		}
		tracks := &[]spotify.SimpleTrack{}
		addTracks := func() {
			*tracks = append(*tracks, fa.Tracks.Tracks...)
		}
		addTracks()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx(), &fa.Tracks); perr == spotify.ErrNoMorePages {
					done(true, nil)
					break
				} else if perr != nil {
					done(false, perr)
					return
				}
				addTracks()
			}
		}()
		p := &Album{
			*fa,
			tracks,
		}
		albumCache[fa.ID] = p
		return p, nil
	} else {
		done(true, nil)
		return albumCache[albumID], nil
	}

}

// CurrentUserSavedAlbums returns the SavedAlbums of the current user in a
// specific manner. It returns the first page and then starts a go routine
// in the background and keeps updating the SavedAlbums and calls the done
// function with a status of true and nil error if successful else calls the
// done function with a status of false and the corresponding error.
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
					return
				}
				addAlbums()
			}
		}()
		return albums, nil
	}
}

// CurrentUserPlaylists returns the UserPlaylists of the current user in a
// specific manner. It returns the first page and then starts a go routine in
// the background and keeps updating the UserPlaylists and calls the done
// function with a status of true and nil error if successful else calls the
// done function with a status of false and the corresponding error.
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
					return
				}
				addPlaylists()
			}
		}()
		return playlists, nil
	}
}

// CurrentUserSavedTracks returns the LikedSongs of the current user in a
// specific manner. It returns the first page and then starts a go routine in
// the background and keeps updating the LikedSongs and calls the done
// function with a status of true and nil error if successful else calls the
// done function with a status of false and the corresponding error.
func CurrentUserSavedTracks(done func(status bool, err error)) (*LikedSongs, error) {
	_p := make(LikedSongs, 0)
	playlists := &_p
	if ls, err := Client.CurrentUsersTracks(ctx()); err != nil {
		return nil, err
	} else {
		addTracks := func() {
			_p = append(_p, ls.Tracks...)
		}
		addTracks()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx(), ls); perr == spotify.ErrNoMorePages {
					done(true, nil)
					break
				} else if perr != nil {
					done(false, perr)
					return
				}
				addTracks()
			}
		}()
		return playlists, nil
	}
}

func RecentlyPlayed() ([]spotify.RecentlyPlayedItem, error) {
	return Client.PlayerRecentlyPlayedOpt(ctx(), &spotify.RecentlyPlayedOptions{Limit: 50})
}
