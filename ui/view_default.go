package ui

import "github.com/aditya-K2/gspt/spt"

type DefaultView struct {
}

func (d *DefaultView) ContextOpener() func(m *Root, s func(s int)) {
	return func(m *Root, s func(s int)) {
		c := NewMenu()
		cc := []string{}
		plist, err := spt.CurrentUserPlaylists(func(s bool, err error) {})
		if err != nil {
			SendNotification("Error Retrieving User Playlists")
			return
		}
		for _, v := range *(plist) {
			cc = append(cc, v.Name)
		}
		c.Content(cc)
		c.Title("Add to Playlist")
		c.SetSelectionHandler(s)
		m.AddCenteredWidget(c)
	}
}

func (d *DefaultView) ContextKey() rune { return 'a' }
