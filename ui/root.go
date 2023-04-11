package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Root struct {
	Root  *tview.Pages
	after func()
}

func NewRoot() *Root {
	m := &Root{}

	Root := tview.NewPages()

	m.Root = Root
	return m
}

func (m *Root) Primitive(name string, t tview.Primitive) {
	m.Root.AddPage(name, t, true, true)
}

func (m *Root) AfterContextClose(f func()) {
	m.after = f
}

type CenteredWidget interface {
	Primitive() *tview.Table
	ContentHandler()
	SelectionHandler() func(s int)
	Size(mw, mh int) (int, int, int, int)
}

func (m *Root) AddCenteredWidget(t CenteredWidget) {
	p := (t.Primitive())
	closec := make(chan bool)
	currentTime := time.Now().String()
	sHandler := t.SelectionHandler()
	_, _, w, h := m.Root.GetRect()

	closeCtx := func() {
		m.Root.RemovePage(currentTime)
		if m.after != nil {
			m.after()
		}
	}
	drawCtx := func() {
		m.Root.AddPage(currentTime, t.Primitive(), false, true)
		p.SetRect(t.Size(w, h))
	}
	redraw := func() {
		closeCtx()
		drawCtx()
	}
	deleteCtx := func() {
		closeCtx()
		closec <- true
	}

	capture := func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEscape {
			deleteCtx()
			return nil
		} else if e.Key() == tcell.KeyEnter {
			r, _ := p.GetSelection()
			sHandler(r)
			closeCtx()
			return nil
		}
		return e
	}
	p.SetInputCapture(capture)

	t.ContentHandler()

	resizeHandler := func() {
		dur := 500
		tck := time.NewTicker(time.Duration(dur) * time.Millisecond)
		go func() {
			for {
				select {
				case <-tck.C:
					{
						_, _, _w, _h := m.Root.GetRect()
						if _w != w || _h != h {
							w = _w
							h = _h
							redraw()
						}
					}
				case <-closec:
					{
						return
					}
				}
			}
		}()
	}
	resizeHandler()

	drawCtx()
}
