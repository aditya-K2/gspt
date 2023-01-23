package main

import (
	"errors"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	berr         = errors.New("Couldn't Get Base Selection in Interactive View")
	defaultfg    = tcell.ColorGreen
	defaultbg    = tcell.ColorDefault
	defaultstyle = tcell.StyleDefault.
			Foreground(defaultfg).
			Background(defaultbg)
)

type _range struct {
	Start int
	End   int
}

type interactiveView struct {
	visual            bool
	vrange            *_range
	baseSel           int
	Pages             *tview.Pages
	view              *tview.Table
	ctxContentHandler func() []string
	ctxSelectHandler  func(s string)
}

func NewInteractiveView() *interactiveView {
	view := tview.NewTable()
	view.SetSelectable(true, false)
	view.SetBackgroundColor(tcell.ColorDefault)
	pages := tview.NewPages()
	pages.AddPage("IView", view, true, true)

	i := &interactiveView{
		view:   view,
		Pages:  pages,
		vrange: &_range{},
		visual: false,
	}

	view.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		i.update()
		return i.view.GetInnerRect()
	})
	view.SetInputCapture(i.capture)
	return i
}

func (i *interactiveView) exitVisualMode() {
	if i.vrange.Start < i.baseSel {
		i.view.Select(i.vrange.Start, -1)
	} else if i.vrange.End > i.baseSel {
		i.view.Select(i.vrange.End, -1)
	}
	i.baseSel = -1
}

func (i *interactiveView) enterVisualMode() {
	row, _ := i.view.GetSelection()
	i.baseSel = row
	i.vrange.Start, i.vrange.End = row, row
}

func (i *interactiveView) toggleVisualMode() {
	if i.visual {
		i.exitVisualMode()
	} else if !i.visual {
		i.enterVisualMode()
	}
	i.visual = !i.visual
}

func (i *interactiveView) getHandler(s string) func(e *tcell.EventKey) *tcell.EventKey {
	vr := i.vrange
	check := func() {
		if vr.Start <= -1 {
			vr.Start = 0
		}
		if vr.End <= -1 {
			vr.End = 0
		}
		if vr.End >= i.view.GetRowCount() {
			vr.End = i.view.GetRowCount() - 1
		}
		if vr.Start >= i.view.GetRowCount() {
			vr.Start = i.view.GetRowCount() - 1
		}
	}
	funcMap := map[string]func(e *tcell.EventKey) *tcell.EventKey{
		"up": func(e *tcell.EventKey) *tcell.EventKey {
			if i.visual {
				check()
				if vr.End > i.baseSel {
					vr.End--
				} else if vr.Start <= i.baseSel {
					vr.Start--
				}
				if i.baseSel == -1 {
					panic(berr)
				}
				return nil
			}
			return e
		},
		"down": func(e *tcell.EventKey) *tcell.EventKey {
			if i.visual {
				check()
				if vr.Start < i.baseSel {
					vr.Start++
				} else if vr.Start == i.baseSel {
					vr.End++
				}
				if i.baseSel == -1 {
					panic(berr)
				}
				return nil
			}
			return e
		},
		"exitvisual": func(e *tcell.EventKey) *tcell.EventKey {
			if i.visual {
				i.exitVisualMode()
				i.visual = false
				return nil
			}
			return e
		},
		"top": func(e *tcell.EventKey) *tcell.EventKey {
			if i.visual {
				i.vrange.Start = 0
				i.vrange.End = i.baseSel
				i.view.ScrollToBeginning()
				return nil
			}
			return e
		},
		"bottom": func(e *tcell.EventKey) *tcell.EventKey {
			if i.visual {
				i.vrange.Start = i.baseSel
				i.vrange.End = i.view.GetRowCount() - 1
				i.view.ScrollToEnd()
				return nil
			}
			return e
		},
		"openContextMenu": func(e *tcell.EventKey) *tcell.EventKey {
			i.openContextMenu()
			return nil
		},
	}
	if val, ok := funcMap[s]; ok {
		return val
	} else {
		return nil
	}
}

func (i *interactiveView) capture(e *tcell.EventKey) *tcell.EventKey {
	switch e.Rune() {
	case 'j':
		{
			return i.getHandler("down")(e)
		}
	case 'k':
		{
			return i.getHandler("up")(e)
		}
	case 'v':
		{
			i.toggleVisualMode()
			return nil
		}
	case 'g':
		{
			return i.getHandler("top")(e)
		}
	case 'G':
		{
			return i.getHandler("bottom")(e)
		}
	case 'C':
		{
			return i.getHandler("openContextMenu")(e)
		}
	default:
		{
			if e.Key() == tcell.KeyEscape {
				return i.getHandler("exitvisual")(e)
			}
			return e
		}
	}
}

func GetCell(text string, st tcell.Style) *tview.TableCell {
	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetStyle(st)
}

func (i *interactiveView) setCtxContentHandler(f func() []string) {
	i.ctxContentHandler = f
}

func (i *interactiveView) setCtxSelectHandler(f func(s string)) {
	i.ctxSelectHandler = f
}

func (i *interactiveView) openContextMenu() {
	iv := i.view
	r, c := iv.GetSelection()
	cslice := i.ctxContentHandler()
	cwidth := 30
	cheight := len(cslice) + 2
	cpaddingx := 2
	cpaddingy := 1
	currentTime := time.Now().String()

	if i.visual {
		r, c = i.baseSel, 1
	}

	ctxMenu := tview.NewTable()
	ctxMenu.SetBorder(true)
	ctxMenu.SetSelectable(true, false)
	capture := func(e *tcell.EventKey) *tcell.EventKey {
		closeCtx := func() {
			i.Pages.RemovePage(currentTime)
		}
		if e.Key() == tcell.KeyEscape {
			closeCtx()
			return nil
		} else if e.Key() == tcell.KeyEnter {
			i.ctxSelectHandler(
				ctxMenu.GetCell(
					ctxMenu.GetSelection()).Text)
			closeCtx()
			return nil
		}
		return e
	}

	ctxMenu.SetInputCapture(capture)

	for k := range cslice {
		ctxMenu.SetCell(k, 0,
			GetCell(cslice[k], defaultstyle))
	}

	i.Pages.AddPage(currentTime, ctxMenu, false, true)
	ctxMenu.SetRect(c+cpaddingx, r+cpaddingy, cwidth, cheight)
}

func (i *interactiveView) update() {
	s := strings.Split("orem ipsum dolor sit amet, consectetur adipiscing elit. Nunc nec leo a tellus gravida convallis. Curabitur tempus purus nisi. Proin non enim convallis augue porta aliquet.", " ")
	i.view.Clear()
	for j := range s {
		b := ""
		if i.visual && (j >= i.vrange.Start && j <= i.vrange.End) {
			b = "[blue::]█[::]"
		}
		i.view.SetCell(j, 0,
			GetCell(b, defaultstyle))
		i.view.SetCell(j, 1,
			GetCell(s[j], defaultstyle))
		i.view.SetCell(j, 2,
			GetCell(s[j], defaultstyle.Foreground(tcell.ColorBlue)))
		i.view.SetCell(j, 3,
			GetCell(s[j], defaultstyle.Foreground(tcell.ColorYellow)))
	}
}
