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
	ctx           = context.Background()
	Client        *spotify.Client
	playlistCache map[spotify.ID]*Playlist = make(map[spotify.ID]*Playlist)
	albumCache    map[spotify.ID]*Album    = make(map[spotify.ID]*Album)
	PageContinue                           = errors.New("CONTINUE")
)

func GetPlaylist(playlistId spotify.ID) (*Playlist, chan error) {
	c := make(chan error)
	if fp, err := Client.GetPlaylist(ctx, playlistId); err != nil {
		go func() { c <- err }()
		return nil, c
	} else {
		if _, ok := playlistCache[fp.ID]; !ok || playlistCache[fp.ID].SnapshotID != fp.SnapshotID {
			tracks := &[]spotify.PlaylistTrack{}
			addTracks := func() {
				*tracks = append(*tracks, fp.Tracks.Tracks...)
			}
			addTracks()
			go func() {
				for page := 1; ; page++ {
					if perr := Client.NextPage(ctx, &fp.Tracks); perr == spotify.ErrNoMorePages {
						c <- nil
						break
					} else if perr != nil {
						c <- perr
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
			return p, c
		}
		go func() { c <- nil }()
		return playlistCache[fp.ID], c
	}
}

func GetAlbum(albumID spotify.ID) (*Album, chan error) {
	c := make(chan error)
	if _, ok := albumCache[albumID]; !ok {
		fa, err := Client.GetAlbum(ctx, albumID)
		if err != nil {
			go func() { c <- err }()
			return nil, c
		}
		tracks := &[]spotify.SimpleTrack{}
		addTracks := func() {
			*tracks = append(*tracks, fa.Tracks.Tracks...)
		}
		addTracks()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx, &fa.Tracks); perr == spotify.ErrNoMorePages {
					c <- nil
					break
				} else if perr != nil {
					c <- perr
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
		return p, c
	} else {
		go func() { c <- nil }()
		return albumCache[albumID], c
	}

}

func CurrentUserSavedAlbums() (*SavedAlbums, chan error) {
	c := make(chan error)
	_a := make(SavedAlbums, 0)
	albums := &_a
	if sp, err := Client.CurrentUsersAlbums(ctx); err != nil {
		go func() { c <- err }()
		return nil, c
	} else {
		addAlbums := func() {
			_a = append(_a, sp.Albums...)
		}
		addAlbums()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx, sp); perr == spotify.ErrNoMorePages {
					c <- nil
					break
				} else if perr != nil {
					c <- perr
					return
				}
				addAlbums()
			}
		}()
		return albums, c
	}
}

func CurrentUserPlaylists() (*UserPlaylists, chan error) {
	c := make(chan error)
	_p := make(UserPlaylists, 0)
	playlists := &_p
	if spp, err := Client.CurrentUsersPlaylists(ctx); err != nil {
		go func() { c <- err }()
		return nil, c
	} else {
		addPlaylists := func() {
			_p = append(_p, spp.Playlists...)
		}
		addPlaylists()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx, spp); perr == spotify.ErrNoMorePages {
					c <- nil
					break
				} else if perr != nil {
					c <- perr
					return
				}
				addPlaylists()
			}
		}()
		return playlists, c
	}
}

func CurrentUserSavedTracks() (*LikedSongs, chan error) {
	c := make(chan error)
	_p := make(LikedSongs, 0)
	playlists := &_p
	if ls, err := Client.CurrentUsersTracks(ctx); err != nil {
		go func() { c <- err }()
		return nil, c
	} else {
		addTracks := func() {
			_p = append(_p, ls.Tracks...)
		}
		addTracks()
		go func() {
			for page := 1; ; page++ {
				if perr := Client.NextPage(ctx, ls); perr == spotify.ErrNoMorePages {
					c <- nil
					break
				} else if perr != nil {
					c <- perr
					return
				}
				addTracks()
			}
		}()
		return playlists, c
	}
}

func CurrentUserFollowedArtists() (*FollowedArtists, chan error) {
	c := make(chan error)
	// TODO: Check if this is the proper implementation
	_a := make(FollowedArtists, 0)
	artists := &_a
	if ar, err := Client.CurrentUsersFollowedArtists(ctx); err != nil {
		go func() { c <- err }()
		return nil, c
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
					c <- nil
				}
				if ar, err = Client.CurrentUsersFollowedArtists(ctx, spotify.After(string(ap))); err != nil {
					if err == spotify.ErrNoMorePages {
						c <- nil
						break
					} else {
						c <- err
						break
					}
				} else {
					addArtists()
				}
			}
		}()
		return artists, c
	}
}

func RecentlyPlayed() ([]spotify.RecentlyPlayedItem, error) {
	return Client.PlayerRecentlyPlayedOpt(ctx, &spotify.RecentlyPlayedOptions{Limit: 50})
}

func GetPlayerState() (*spotify.PlayerState, error) {
	return Client.PlayerState(ctx)
}

func GetTopTracks() ([]spotify.FullTrack, error) {
	c, err := Client.CurrentUsersTopTracks(ctx, spotify.Limit(topTracksLimit))
	if c != nil {
		return c.Tracks, err
	} else {
		return []spotify.FullTrack{}, errors.New("No Top Tracks Found!")
	}
}

func GetTopArtists() ([]spotify.FullArtist, error) {
	c, err := Client.CurrentUsersTopArtists(ctx, spotify.Limit(topTracksLimit))
	if c != nil {
		return c.Artists, err
	} else {
		return []spotify.FullArtist{}, errors.New("No Top Artists Found!")
	}
}

func GetArtistTopTracks(artistID spotify.ID) ([]spotify.FullTrack, error) {
	c, err := Client.CurrentUser(ctx)
	if err != nil {
		return []spotify.FullTrack{}, err
	}
	return Client.GetArtistsTopTracks(ctx, artistID, c.Country)
}

func GetArtistAlbums(artistID spotify.ID) ([]spotify.SimpleAlbum, error) {
	c, err := Client.GetArtistAlbums(ctx, artistID, albumtypes)
	if err != nil {
		return []spotify.SimpleAlbum{}, err
	}
	return c.Albums, nil
}

func Search(s string) (*spotify.SearchResult, error) {
	return Client.Search(ctx, s,
		spotify.SearchTypePlaylist|
			spotify.SearchTypeAlbum|
			spotify.SearchTypeTrack|
			spotify.SearchTypeArtist)
}

func UserDevices() ([]spotify.PlayerDevice, error) {
	return Client.PlayerDevices(ctx)
}

func TransferPlayback(deviceId spotify.ID) error {
	s, err := GetPlayerState()
	if err != nil {
		return errors.New("Unable to get Current Player State!")
	}
	err = Client.PauseOpt(ctx, &spotify.PlayOptions{DeviceID: &s.Device.ID})
	return Client.TransferPlayback(ctx, deviceId, true)
}

func GetFullPlaylist(id spotify.ID) (*spotify.FullPlaylist, error) {
	return Client.GetPlaylist(ctx, id)
}
