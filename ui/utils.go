package ui

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/aditya-K2/gspt/config"
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
	return fmt.Sprintf(filepath.Join(config.Config.CacheDir, "%s.jpg"), a.ID)
}

func getFontWidth() (int, int, error) {
	w, h, err := ueberzug.GetParentSize()
	if err != nil {
		return 0, 0, err
	}
	_, _, rw, rh := root.Root.GetRect()
	if rw == 0 || rh == 0 {
		return 0, 0, errors.New("Unable to get row width and height")
	}
	fw := w / rw
	fh := h / rh
	return fw, fh, nil
}
