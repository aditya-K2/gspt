package ui

import (
	"github.com/aditya-K2/gspt/config"
	"github.com/gdamore/tcell/v2"
)

type defView struct {
	m             map[string]map[config.Key]string
	actions       map[string]*Action
	visualActions map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey
}

func (d *defView) SetMappings(m map[string]map[config.Key]string) {
	d.m = m
}

func (d *defView) SetActions(a map[string]*Action) {
	d.actions = a
}

func (d *defView) SetVisualActions(a map[string]func(start, end int, e *tcell.EventKey) *tcell.EventKey) {
	d.visualActions = a
}

func (d *defView) ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if d.m["normal"] != nil {
			var key config.Key
			if e.Key() == tcell.KeyRune {
				key = config.Key{R: e.Rune()}
			} else {
				key = config.Key{K: e.Key()}
			}
			if val, ok := d.m["normal"][key]; ok {
				return d.actions[val].Func()(e)
			}
		}
		return e
	}
}

func (d *defView) VisualCapture() func(start, end int, e *tcell.EventKey) *tcell.EventKey {
	return func(start, end int, e *tcell.EventKey) *tcell.EventKey {
		if d.m["visual"] != nil {
			var key config.Key
			if e.Key() == tcell.KeyRune {
				key = config.Key{R: e.Rune()}
			} else {
				key = config.Key{K: e.Key()}
			}
			if val, ok := d.m["visual"][key]; ok {
				return d.visualActions[val](start, end, e)
			}
		}
		return e
	}
}

type DefaultViewNone struct {
	*defView
}

func (a *DefaultViewNone) DisableVisualMode() bool { return true }
func (a *DefaultViewNone) VisualCapture() func(start, end int, e *tcell.EventKey) *tcell.EventKey {
	return nil
}

type DefaultView struct {
	*defView
}

func (d *DefaultView) DisableVisualMode() bool { return false }
