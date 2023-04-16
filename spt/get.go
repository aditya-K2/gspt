package spt

import (
	"context"
	"errors"

	"github.com/zmb3/spotify/v2"
)

var (
	topTracksLimit = 15
	albumtypes     = []spotify.AlbumType{
		spotify.AlbumTypeAlbum,
		spotify.AlbumTypeSingle,
		spotify.AlbumTypeAppearsOn,
		spotify.AlbumTypeCompilation,
	}
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
type FollowedArtists []spotify.FullArtist
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
// list in the Playlist object. When done, it passes the error value to the
// error handler.
func GetPlaylist(playlistId spotify.ID, errHandler func(error)) (*Playlist, error) {
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
						errHandler(nil)
						break
					} else if perr != nil {
						errHandler(perr)
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
			errHandler(nil)
			return playlistCache[fp.ID], nil
		}
	}
}

// GetAlbum retrieves a Spotify album by its ID and returns an Album object.
// If the album is already cached, it returns the cached object.
// If not, it retrieves the album from the Spotify API and caches it.
// It uses a background go routine to fetch all pages of tracks for the album
// and appends them to the tracks list in the Album object.
// When done, it passes the error value to the error handler.
func GetAlbum(albumID spotify.ID, errHandler func(error)) (*Album, error) {
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
					errHandler(nil)
					break
				} else if perr != nil {
					errHandler(perr)
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
		errHandler(nil)
		return albumCache[albumID], nil
	}

}

// CurrentUserSavedAlbums returns the SavedAlbums of the current user in a
// specific manner. It returns the first page and then starts a go routine
// in the background and keeps updating the SavedAlbums
// When done, it passes the error value to the error handler.
func CurrentUserSavedAlbums(errHandler func(err error)) (*SavedAlbums, error) {
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
					errHandler(nil)
					break
				} else if perr != nil {
					errHandler(perr)
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
// the background and keeps updating the UserPlaylists
// When done, it passes the error value to the error handler.
func CurrentUserPlaylists(errHandler func(err error)) (*UserPlaylists, error) {
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
					errHandler(nil)
					break
				} else if perr != nil {
					errHandler(perr)
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
// the background and keeps updating the LikedSongs
// When done, it passes the error value to the error handler.
func CurrentUserSavedTracks(errHandler func(err error)) (*LikedSongs, error) {
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
					errHandler(nil)
					break
				} else if perr != nil {
					errHandler(perr)
					return
				}
				addTracks()
			}
		}()
		return playlists, nil
	}
}

func CurrentUserFollowedArtists(errHandler func(error)) (*FollowedArtists, error) {
	// TODO: Check if this is the proper implementation
	_a := make(FollowedArtists, 0)
	artists := &_a
	if ar, err := Client.CurrentUsersFollowedArtists(ctx()); err != nil {
		return artists, err
	} else {
		ap := spotify.ID("")
		addArtists := func() {
			_a = append(_a, ar.Artists...)
			ap = _a[len(_a)-1].ID
		}
		addArtists()
		go func() {
			for {
				if len(ar.Artists) == 0 {
					errHandler(nil)
				}
				if ar, err = Client.CurrentUsersFollowedArtists(ctx(), spotify.After(string(ap))); err != nil {
					if err == spotify.ErrNoMorePages {
						errHandler(nil)
						break
					} else {
						errHandler(err)
						break
					}
				} else {
					addArtists()
				}
			}
		}()
		return artists, nil
	}
}

func RecentlyPlayed() ([]spotify.RecentlyPlayedItem, error) {
	return Client.PlayerRecentlyPlayedOpt(ctx(), &spotify.RecentlyPlayedOptions{Limit: 50})
}

func GetPlayerState() (*spotify.PlayerState, error) {
	return Client.PlayerState(ctx())
}

func GetTopTracks() ([]spotify.FullTrack, error) {
	c, err := Client.CurrentUsersTopTracks(ctx(), spotify.Limit(topTracksLimit))
	if c != nil {
		return c.Tracks, err
	} else {
		return []spotify.FullTrack{}, errors.New("No Top Tracks Found!")
	}
}

func GetTopArtists() ([]spotify.FullArtist, error) {
	c, err := Client.CurrentUsersTopArtists(ctx(), spotify.Limit(topTracksLimit))
	if c != nil {
		return c.Artists, err
	} else {
		return []spotify.FullArtist{}, errors.New("No Top Artists Found!")
	}
}

func GetArtistTopTracks(artistID spotify.ID) ([]spotify.FullTrack, error) {
	c, err := Client.CurrentUser(ctx())
	if err != nil {
		return []spotify.FullTrack{}, err
	}
	return Client.GetArtistsTopTracks(ctx(), artistID, c.Country)
}

func GetArtistAlbums(artistID spotify.ID) ([]spotify.SimpleAlbum, error) {
	c, err := Client.GetArtistAlbums(ctx(), artistID, albumtypes)
	if err != nil {
		return []spotify.SimpleAlbum{}, err
	}
	return c.Albums, nil
}

func Search(s string) (*spotify.SearchResult, error) {
	return Client.Search(ctx(), s,
		spotify.SearchTypePlaylist|
			spotify.SearchTypeAlbum|
			spotify.SearchTypeTrack|
			spotify.SearchTypeArtist)
}

func UserDevices() ([]spotify.PlayerDevice, error) {
	return Client.PlayerDevices(ctx())
}

func TransferPlayback(deviceId spotify.ID) error {
	s, err := GetPlayerState()
	if err != nil {
		return errors.New("Unable to get Current Player State!")
	}
	err = Client.PauseOpt(ctx(), &spotify.PlayOptions{DeviceID: &s.Device.ID})
	return Client.TransferPlayback(ctx(), deviceId, true)
}

func GetSimplePlaylist(id *spotify.ID) (*spotify.FullPlaylist, error) {
	return Client.GetPlaylist(ctx(), *id)
}
