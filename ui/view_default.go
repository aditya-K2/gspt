package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
)

type defView struct {
	m       map[tcell.Key]string
	actions map[string]*Action
}

func (d *defView) SetMappings(m map[tcell.Key]string) {
	d.m = m
}

func (d *defView) SetActions(a map[string]*Action) {
	d.actions = a
}

func (d *defView) ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if d.m != nil {
			if val, ok := d.m[e.Key()]; ok {
				if d.actions != nil {
					return d.actions[val].Func()(e)
				}
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

type DefaultView struct {
	*defView
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

func (d *DefaultView) DisableVisualMode() bool { return false }
