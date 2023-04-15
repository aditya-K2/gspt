package ui

import "github.com/gdamore/tcell/v2"

type Refreshable interface {
	RefreshState()
}

type ActionFunc func(e *tcell.EventKey) *tcell.EventKey

// Action represents the Event Handler to the events that are passed to the
// InputCapture of the tview Widgets. It will refresh the state of the
// Refreshable r upon every subsequent call to the function that is being
// returned by the GetFunc()
type Action struct {
	f           ActionFunc
	refreshable Refreshable
}

func NewAction(f ActionFunc, refreshes Refreshable) *Action {
	a := &Action{}
	a.SetFunc(f)
	a.SetRefreshable(refreshes)
	return a
}

func (a *Action) SetFunc(f ActionFunc) {
	a.f = f
}

func (a *Action) SetRefreshable(r Refreshable) {
	a.refreshable = r

}

func (a *Action) Func() ActionFunc {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if a != nil && a.f != nil {
			val := a.f(e)
			if a.refreshable != nil && val == nil {
				a.refreshable.RefreshState()
			}
			return val
		}
		return e
	}
}
