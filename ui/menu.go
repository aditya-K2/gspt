package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/aditya-K2/tview"
)

var (
	minheight = 7
)

// Menu implements the CenteredWidget interface for the Main Screen.
type menu struct {
	Menu     *tview.Table
	title    string
	content  []string
	sHandler func(s int)
}

func NewMenu() *menu {
	c := &menu{}

	menu := tview.NewTable()
	menu.SetBorder(true)
	menu.SetBorderPadding(1, 1, 1, 1)
	menu.SetBorderStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite))
	menu.SetBackgroundColor(tcell.ColorDefault)
	menu.SetSelectable(true, false)
	c.Menu = menu

	return c
}

func (c *menu) Size(mw, mh int) (int, int, int, int) {
	cheight := mh / 5
	if cheight < minheight {
		cheight = minheight
	}
	cwidth := 30
	epx := 4

	return mw/2 - (cwidth/2 + epx), (mh/2 - (cheight/2 + epx)), cwidth, cheight
}

func (c *menu) ContentHandler() {
	for k := range c.content {
		c.Menu.SetCell(k, 0,
			GetCell(c.content[k], ContextMenuStyle))
	}
}

func (c *menu) SelectionHandler() func(s int) {
	return c.sHandler
}

func (c *menu) SetSelectionHandler(f func(s int)) {
	c.sHandler = f
}

func (c *menu) Primitive() *tview.Table { return c.Menu }

func (c *menu) Content(s []string) { c.content = s }

func (c *menu) Title(s string) { c.Menu.SetTitle(s) }
