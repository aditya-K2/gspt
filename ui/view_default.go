package ui

import (
	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"
)

type defView struct {
	m             map[config.Key]string
	vm            map[config.Key]string
	actions       map[string]*Action
	visualActions map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey
}

func (d *defView) SetMappings(m map[config.Key]string) {
	d.m = m
}

func (d *defView) SetVisualMappings(vm map[config.Key]string) {
	d.vm = vm
}

func (d *defView) SetActions(a map[string]*Action) {
	d.actions = a
}

func (d *defView) ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if d.m != nil {
			var key config.Key
			if e.Key() == tcell.KeyRune {
				key = config.Key{R: e.Rune()}
			} else {
				key = config.Key{K: e.Key()}
			}
			if val, ok := d.m[key]; ok {
				return d.actions[val].Func()(e)
			}
		}
		return e
	}
}

func (d *defView) VisualCapture() func(start, end int, e *tcell.EventKey) *tcell.EventKey {
	return func(start, end int, e *tcell.EventKey) *tcell.EventKey {
		if d.vm != nil {
			var key config.Key
			if e.Key() == tcell.KeyRune {
				key = config.Key{R: e.Rune()}
			} else {
				key = config.Key{K: e.Key()}
			}
			if val, ok := d.vm[key]; ok {
				return d.visualActions[val](start, end, e)
			}
		}
		return e
	}
}

type DefaultViewNone struct {
	*defView
}

func (a *DefaultViewNone) ContextOpener() func(m *Root, s func(s int)) { return nil }
func (a *DefaultViewNone) ContextHandler() func(int, int, int)         { return nil }
func (a *DefaultViewNone) ContextKey() rune                            { return 'a' }
func (a *DefaultViewNone) DisableVisualMode() bool                     { return true }
func (a *DefaultViewNone) VisualCapture() func(start, end int, e *tcell.EventKey) *tcell.EventKey {
	return nil
}

type DefaultView struct {
	*defView
}

func (d *DefaultView) DisableVisualMode() bool { return false }

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
