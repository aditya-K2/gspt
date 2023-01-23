package main

import (
	"errors"
	"strings"

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

type InteractiveView struct {
	visual  bool
	vrange  *_range
	baseSel int
	View    *tview.Table
}

func NewInteractiveView() *InteractiveView {
	view := tview.NewTable()
	view.SetSelectable(true, false)
	i := &InteractiveView{
		View:   view,
		vrange: &_range{},
		visual: false,
	}

	view.SetInputCapture(i.capture)
	return i
}

func (i *InteractiveView) exitVisualMode() {
	if i.vrange.Start < i.baseSel {
		i.View.Select(i.vrange.Start, -1)
	} else if i.vrange.End > i.baseSel {
		i.View.Select(i.vrange.End, -1)
	}
	i.baseSel = -1
}

func (i *InteractiveView) enterVisualMode() {
	row, _ := i.View.GetSelection()
	i.baseSel = row
	i.vrange.Start, i.vrange.End = row, row
}

func (i *InteractiveView) toggleVisualMode() {
	if i.visual {
		i.exitVisualMode()
	} else if !i.visual {
		i.enterVisualMode()
	}
	i.visual = !i.visual
}

func (i *InteractiveView) GetHandler(s string) func(e *tcell.EventKey) *tcell.EventKey {
	vr := i.vrange
	check := func() {
		if vr.Start <= -1 {
			vr.Start = 0
		}
		if vr.End <= -1 {
			vr.End = 0
		}
		if vr.End >= i.View.GetRowCount() {
			vr.End = i.View.GetRowCount() - 1
		}
		if vr.Start >= i.View.GetRowCount() {
			vr.Start = i.View.GetRowCount() - 1
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
				i.View.ScrollToBeginning()
				return nil
			}
			return e
		},
		"bottom": func(e *tcell.EventKey) *tcell.EventKey {
			if i.visual {
				i.vrange.Start = i.baseSel
				i.vrange.End = i.View.GetRowCount() - 1
				i.View.ScrollToEnd()
				return nil
			}
			return e
		},
	}
	if val, ok := funcMap[s]; ok {
		return val
	} else {
		return nil
	}
}

func (i *InteractiveView) capture(e *tcell.EventKey) *tcell.EventKey {
	switch e.Rune() {
	case 'j':
		{
			return i.GetHandler("down")(e)
		}
	case 'k':
		{
			return i.GetHandler("up")(e)
		}
	case 'v':
		{
			i.toggleVisualMode()
			return nil
		}
	case 'g':
		{
			return i.GetHandler("top")(e)
		}
	case 'G':
		{
			return i.GetHandler("bottom")(e)
		}
	default:
		{
			if e.Key() == tcell.KeyEscape {
				return i.GetHandler("exitvisual")(e)
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

func (i *InteractiveView) Update() {
	s := strings.Split("orem ipsum dolor sit amet, consectetur adipiscing elit. Nunc nec leo a tellus gravida convallis. Curabitur tempus purus nisi. Proin non enim convallis augue porta aliquet.", " ")
	i.View.Clear()
	for j := range s {
		b := ""
		if i.visual && (j >= i.vrange.Start && j <= i.vrange.End) {
			b = "[blue::]█[::]"
		}
		i.View.SetCell(j, 0,
			GetCell(b, defaultstyle))
		i.View.SetCell(j, 1,
			GetCell(s[j], defaultstyle))
		i.View.SetCell(j, 2,
			GetCell(s[j], defaultstyle.Foreground(tcell.ColorBlue)))
		i.View.SetCell(j, 3,
			GetCell(s[j], defaultstyle.Foreground(tcell.ColorYellow)))
	}
}
