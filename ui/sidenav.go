package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	SideNavActions map[string]*Action
)

func NewSideNav(m map[tcell.Key]string) *tview.Table {
	T := tview.NewTable()

	mapSideNavActions(T, m)

	return T
}

func mapSideNavActions(T *tview.Table, f map[tcell.Key]string) {
	T.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if action, ok := f[e.Key()]; ok {
			SideNavActions[action].Func()(e)
		}
		return e
	})
}
