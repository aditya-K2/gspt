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
	fun ActionFunc
	r   Refreshable
}

func NewAction(f ActionFunc, r Refreshable) *Action {
	a := &Action{}
	a.SetFunc(f)
	a.SetRefreshable(r)
	return a
}

func (a *Action) SetFunc(f ActionFunc) {
	a.fun = f
}

func (a *Action) SetRefreshable(r Refreshable) {
	a.r = r

}

func (a *Action) Func() ActionFunc {
	return func(e *tcell.EventKey) *tcell.EventKey {
		val := a.fun(e)
		if a.r != nil && val == nil {
			a.r.RefreshState()
		}
		return val
	}
}
