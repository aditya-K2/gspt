package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func OpenContextMenu(root *tview.Pages,
	ctxContentHandler func() []string, ctxSelectHandler func(s string)) {
	ctxMenu := tview.NewTable()

	_, _, w, h := root.GetRect()
	cslice := ctxContentHandler()
	cwidth := 30
	cheight := len(cslice) + 2
	currentTime := time.Now().String()
	epx := 4
	closec := make(chan bool)

	closeCtx := func() {
		root.RemovePage(currentTime)
	}

	drawCtx := func() {
		root.AddPage(currentTime, ctxMenu, false, true)
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
						_, _, _w, _h := root.GetRect()
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

	for k := range cslice {
		ctxMenu.SetCell(k, 0,
			GetCell(cslice[k], defaultstyle))
	}

	drawCtx()
}
