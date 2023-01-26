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

type CenteredWidget interface {
	Primitive() *tview.Table
	ContentHandler()
	SelectionHandler() func(s string)
	Size(mw, mh int) (int, int, int, int)
}

func (m *Main) addCenteredWidget(t CenteredWidget) {
	p := *(t.Primitive())
	closec := make(chan bool)
	currentTime := time.Now().String()
	sHandler := t.SelectionHandler()
	_, _, w, h := m.Root.GetRect()

	closeCtx := func() {
		m.Root.RemovePage(currentTime)
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
			sHandler(
				p.GetCell(
					p.GetSelection()).Text)
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

func (m *Main) OpenContextMenu() {
	c := newMenu()
	content := []string{}
	c.Content([]string{
		"Hello",
		"Bitches",
		"whatisup"})
	c.Title("Add to Playlist")
	sHandler := func(s string) {
		content = append(content, s)
	}
	c.SetSelectionHandler(sHandler)
	m.addCenteredWidget(c)
}
