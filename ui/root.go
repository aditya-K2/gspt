package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Main struct {
	Root *tview.Pages
}

func NewMain() *Main {
	m := &Main{}

	iv := NewInteractiveView()
	Root := tview.NewPages()
	Root.AddPage("iview", iv.View, true, true)

	m.Root = Root
	return m
}

func (m *Main) OpenContextMenu(title string,
	ctxContentHandler func() []string, ctxSelectHandler func(s string)) {
	ctxMenu := tview.NewTable()

	_, _, w, h := m.Root.GetRect()
	cslice := ctxContentHandler()
	cwidth := 30
	cheight := len(cslice) + 3
	currentTime := time.Now().String()
	epx := 4
	closec := make(chan bool)

	closeCtx := func() {
		m.Root.RemovePage(currentTime)
	}

	drawCtx := func() {
		m.Root.AddPage(currentTime, ctxMenu, false, true)
		ctxMenu.SetRect(w/2-(cwidth/2+epx), (h/2 - (cheight/2 + epx)), cwidth, cheight)
	}

	redraw := func() {
		closeCtx()
		drawCtx()
	}

	deleteCtx := func() {
		closeCtx()
		closec <- true
	}

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

	ctxMenu.SetBorder(true)
	ctxMenu.SetSelectable(true, false)
	capture := func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEscape {
			deleteCtx()
			return nil
		} else if e.Key() == tcell.KeyEnter {
			ctxSelectHandler(
				ctxMenu.GetCell(
					ctxMenu.GetSelection()).Text)
			deleteCtx()
			return nil
		}
		return e
	}

	ctxMenu.SetInputCapture(capture)

	ctxMenu.SetCell(0, 0,
		GetCell(title, tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Bold(true)).SetSelectable(false))
	for k := range cslice {
		ctxMenu.SetCell(k+1, 0,
			GetCell(cslice[k], defaultstyle))
	}

	drawCtx()
}
