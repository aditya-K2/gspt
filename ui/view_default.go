package ui

import "github.com/aditya-K2/gspt/spt"

type DefaultView struct {
}

func (d *DefaultView) ContextOpener() func(m *Root, s func(s int)) {
	return func(m *Root, s func(s int)) {
		c := NewMenu()
		cc := []string{}
		// TODO: Better Error Handling
		plist, err := spt.CurrentUserPlaylists(func(err error) {})
		if err != nil {
			SendNotification(err.Error())
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