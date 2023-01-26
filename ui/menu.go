package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type menu struct {
	Menu     *tview.Table
	title    string
	content  []string
	sHandler func(s string)
}

func newMenu() *menu {
	c := &menu{}

	menu := tview.NewTable()
	menu.SetBorder(true)
	menu.SetSelectable(true, false)
	c.Menu = menu

	return c
}

func (c *menu) Size(mw, mh int) (int, int, int, int) {
	cheight := len(c.content) + 3
	cwidth := 30
	epx := 4

	return mw/2 - (cwidth/2 + epx), (mh/2 - (cheight/2 + epx)), cwidth, cheight
}

func (c *menu) ContentHandler() {
	if c.title != "" {
		c.Menu.SetCell(0, 0,
			GetCell(c.title, tcell.StyleDefault.
				Foreground(tcell.ColorWhite).
				Bold(true)).SetSelectable(false))
	}
	for k := range c.content {
		c.Menu.SetCell(k+1, 0,
			GetCell(c.content[k], defaultstyle))
	}
}

func (c *menu) SelectionHandler() func(s string) {
	return c.sHandler
}

func (c *menu) SetSelectionHandler(f func(s string)) {
	c.sHandler = f
}

func (c *menu) Primitive() *tview.Table { return c.Menu }

func (c *menu) Content(s []string) { c.content = s }

func (c *menu) Title(s string) { c.title = s }
