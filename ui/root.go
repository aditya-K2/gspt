package ui

import (
	"time"

	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type Root struct {
	*tview.Pages
	after func()
}

func NewRoot() *Root {
	return &Root{tview.NewPages(), nil}
}

func (m *Root) Primitive(name string, t tview.Primitive) {
	m.AddPage(name, t, true, true)
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
	_, _, w, h := m.GetRect()

	closeCtx := func() {
		m.RemovePage(currentTime)
		if m.after != nil {
			m.after()
		}
	}
	drawCtx := func() {
		m.AddPage(currentTime, t.Primitive(), false, true)
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
		dur := cfg.RedrawInterval
		tck := time.NewTicker(time.Duration(dur) * time.Millisecond)
		go func() {
			for {
				select {
				case <-tck.C:
					{
						_, _, _w, _h := m.GetRect()
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
