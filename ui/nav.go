package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type navItem struct {
	name   string
	action *Action
}

type NavMenu struct {
	*defView
	*tview.Table
	m []navItem
}

func NewNavMenu(m []navItem) *NavMenu {
	T := tview.NewTable()
	n := &NavMenu{&defView{}, T, m}

	T.SetDrawFunc(func(tcell.Screen, int, int, int, int) (int, int, int, int) {
		for k := range n.m {
			T.SetCell(k, 0,
				GetCell(n.m[k].name, NavStyle))
		}
		return T.GetInnerRect()
	})
	T.SetTitle("Library").SetTitleAlign(tview.AlignLeft)
	T.SetBackgroundColor(tcell.ColorDefault)
	T.SetBorder(true)
	T.SetSelectable(true, false)

	return n
}

func (n *NavMenu) OpenEntry(e *tcell.EventKey) *tcell.EventKey {
	r, _ := n.Table.GetSelection()
	if r < len(n.m) {
		return (*n.m[r].action).Func()(e)
	}
	return e
}

type PlaylistNav struct {
	*defView
	*tview.Table
	Playlists *spt.UserPlaylists
	c         chan bool
	done      func(error)
}

func NewPlaylistNav() *PlaylistNav {
	T := tview.NewTable()
	T.SetSelectable(true, false).SetBorder(true)
	T.SetTitle("Playlists").SetTitleAlign(tview.AlignLeft)
	T.SetBackgroundColor(tcell.ColorDefault)
	v := &PlaylistNav{&defView{}, T, nil, make(chan bool), func(err error) {
		if err != nil {
			SendNotification(err.Error())
		}
	}}
	T.SetDrawFunc(func(s tcell.Screen, x, y, w, h int) (int, int, int, int) {
		if v.Playlists == nil {
			v.RefreshState()
		}
		for k, p := range *v.Playlists {
			v.Table.SetCell(k, 0,
				GetCell(p.Name, PlaylistNavStyle))
		}
		return T.GetInnerRect()
	})
	return v
}

func (v *PlaylistNav) PlayEntry(e *tcell.EventKey) *tcell.EventKey {
	r, _ := v.Table.GetSelection()
	if err := spt.PlayContext(&(*v.Playlists)[r].URI); err != nil {
		SendNotification(err.Error())
	}
	return nil
}

func (v *PlaylistNav) RefreshState() {
	p, ch := spt.CurrentUserPlaylists()
	go func() {
		err := <-ch
		if err != nil {
			SendNotification(err.Error())
		}
	}()
	v.Playlists = p
}
