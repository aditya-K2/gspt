package ui

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/aditya-K2/gspt/spt"
	"github.com/zmb3/spotify/v2"
	"gitlab.com/diamondburned/ueberzug-go"
)

func openPlaylistMenu(handler func(playlistId spotify.SimplePlaylist)) {
	c := NewMenu()
	cc := []string{}
	// TODO: Better Error Handling
	plist, ch := spt.CurrentUserPlaylists()
	err := <-ch
	if err != nil {
		SendNotification(err.Error())
		return
	}
	for _, v := range *(plist) {
		cc = append(cc, v.Name)
	}
	c.Content(cc)
	c.Title("Add to Playlist")
	c.SetSelectionHandler(func(sel int) {
		handler((*plist)[sel])
	})
	root.AddCenteredWidget(c)
}

func openDeviceMenu() {
	m := NewMenu()
	cc := []string{}
	// TODO: Better Error Handling
	devices, err := spt.UserDevices()
	if err != nil {
		SendNotification(err.Error())
		return
	}
	for _, v := range devices {
		cc = append(cc, v.Name)
	}
	m.Content(cc)
	m.Title("Choose A Device")
	m.SetSelectionHandler(func(s int) {
		if err := spt.TransferPlayback(devices[s].ID); err != nil {
			SendNotification(err.Error())
		} else {
			RefreshProgress(true)
		}
	})
	root.AddCenteredWidget(m)
}

func addToPlaylist(tracks []spotify.ID) {
	openPlaylistMenu(func(sp spotify.SimplePlaylist) {
		aerr := spt.AddTracksToPlaylist(sp.ID, tracks...)
		if aerr != nil {
			SendNotification(aerr.Error())
			return
		} else {
			s := ""
			if len(tracks) > 1 {
				s = "s"
			}
			SendNotification("Added %d track%s to %s", len(tracks), s, sp.Name)
		}
	})
}

func fileName(a spotify.SimpleAlbum) string {
	return fmt.Sprintf(filepath.Join(cfg.CacheDir, "%s.jpg"), a.ID)
}

func getFontWidth() (int, int, error) {
	w, h, err := ueberzug.GetParentSize()
	if err != nil {
		return 0, 0, err
	}
	_, _, rw, rh := root.GetRect()
	if rw == 0 || rh == 0 {
		return 0, 0, errors.New("Unable to get row width and height")
	}
	fw := w / rw
	fh := h / rh
	return fw, fh, nil
}

func openCurrentArtist() {
	if state != nil && state.Item != nil {
		if len(state.Item.Artists) != 0 {
			artistView.SetArtist(&state.Item.Artists[0].ID)
			SetCurrentView(artistView)
			App.SetFocus(Main)
		} else {
			SendNotification("No Artist Found!")
		}
	}
}

func openCurrentAlbum() {
	if state != nil && state.Item != nil {
		albumView.SetAlbum(state.Item.Album.Name, &state.Item.Album.ID)
		SetCurrentView(albumView)
		App.SetFocus(Main)
	}
}

func mergeGenres(g []string) string {
	s := ""
	for k, v := range g {
		sep := ","
		if k == 0 {
			sep = ""
		}
		s += sep + v
	}
	return s
}

func artistName(s []spotify.SimpleArtist) string {
	if len(s) != 0 {
		return s[0].Name
	}
	return ""
}
