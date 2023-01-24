package ui

import (
	"github.com/rivo/tview"
)

type ContextMenu struct {
	Menu *tview.Table
	_s   []string
}

func NewContextMenu() *ContextMenu {
	c := &ContextMenu{}

	ctxMenu := tview.NewTable()
	ctxMenu.SetBorder(true)
	ctxMenu.SetSelectable(true, false)

	c.Menu = ctxMenu

	return c
}

func (c *ContextMenu) Size(mw, mh int) (int, int, int, int) {
	cslice := c.ContentHandler()
	cheight := len(cslice) + 3
	cwidth := 30
	epx := 4

	return mw/2 - (cwidth/2 + epx), (mh/2 - (cheight/2 + epx)), cwidth, cheight
}

func (c *ContextMenu) ContentHandler() []string {
	return []string{
		"Hello",
		"Bye",
	}
}

func (c *ContextMenu) SelectionHandler(s string) {
	c._s = append(c._s, s)
}

func (c ContextMenu) Title() string            { return "Add to Playlist: " }
func (c *ContextMenu) Primitive() *tview.Table { return c.Menu }
