package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PlaylistNav struct {
	Table     *tview.Table
	Playlists *spt.UserPlaylists
	c         chan bool
	done      func(bool, error)
}

var PlaylistActions map[string]*Action

func NewPlaylistNav(done func(s bool, e error)) (*PlaylistNav, error) {
	T := tview.NewTable()
	T.SetSelectable(true, false).SetBorder(true)
	p, err := spt.CurrentUserPlaylists(done)

	if err != nil {
		return nil, err
	}

	v := &PlaylistNav{T, p, make(chan bool), done}
	v.listen()

	T.SetDrawFunc(func(s tcell.Screen, x, y, w, h int) (int, int, int, int) {
		v.Draw()
		return T.GetInnerRect()
	})

	return v, nil
}

func (v *PlaylistNav) Draw() {
	for k, p := range *v.Playlists {
		v.Table.SetCell(k, 0,
			GetCell(p.Name, tcell.StyleDefault.Foreground(tcell.ColorRed)))
	}
}

func (v *PlaylistNav) MapActions(f map[tcell.Key]string) {
	v.Table.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if action, ok := f[e.Key()]; ok {
			PlaylistActions[action].Func()(e)
		}
		return e
	})
}

// func (v *PlaylistsView) SelectEntry(e *tcell.EventKey) *tcell.EventKey {
// 	r, _ := v.Table.GetSelection()
// 	if p, err := spt.GetPlaylist((*v.Playlists)[r].ID); err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println(p)
// 	}
// 	return nil
// }

func (v *PlaylistNav) PlaySelectEntry(e *tcell.EventKey) *tcell.EventKey {
	r, _ := v.Table.GetSelection()
	if err := spt.PlayContext(&(*v.Playlists)[r].URI); err != nil {
		panic(err)
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
