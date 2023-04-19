package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type NavMenu struct {
	*defView
	Table *tview.Table
	m     []navItem
}

type navItem struct {
	name   string
	action *Action
}

func newNavMenu(m []navItem) *NavMenu {
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

func (n *NavMenu) SelectEntry(e *tcell.EventKey) *tcell.EventKey {
	r, _ := n.Table.GetSelection()
	if r < len(n.m) {
		return (*n.m[r].action).Func()(e)
	}
	return e
}

type PlaylistNav struct {
	*defView
	Table     *tview.Table
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
	v.listen()
	T.SetDrawFunc(func(s tcell.Screen, x, y, w, h int) (int, int, int, int) {
		v.Draw()
		return T.GetInnerRect()
	})
	return v
}

func (v *PlaylistNav) Draw() {
	if v.Playlists == nil {
		done := func(err error) {
			if err != nil {
				SendNotification(err.Error())
				return
			}
			App.Draw()
		}
		p, err := spt.CurrentUserPlaylists(done)
		if err != nil {
			SendNotification(err.Error())
			return
		}
		v.Playlists = p
	}
	for k, p := range *v.Playlists {
		v.Table.SetCell(k, 0,
			GetCell(p.Name, PlaylistNavStyle))
	}
}

func (v *PlaylistNav) PlaySelectEntry(e *tcell.EventKey) *tcell.EventKey {
	r, _ := v.Table.GetSelection()
	if err := spt.PlayContext(&(*v.Playlists)[r].URI); err != nil {
		SendNotification(err.Error())
	}
	return nil
}

func (v *PlaylistNav) listen() {
	go func() {
		for {
			if <-v.c {
				p, err := spt.CurrentUserPlaylists(v.done)
				if err != nil {
					panic(err)
				}
				v.Playlists = p
			}
		}
	}()
}

func (v *PlaylistNav) RefreshState() {
	go func() {
		v.c <- true
	}()
}
