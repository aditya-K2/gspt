package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
)

type AlbumsView struct {
	savedAlbums *spt.SavedAlbums
}

func (a *AlbumsView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if a.savedAlbums == nil {
			msg := SendNotificationWithChan("Loading Albums from your Library...")
			sa, err := spt.CurrentUserSavedAlbums(func(s bool, err error) {
				if s {
					msg <- "Done"
				} else {
					msg <- err.Error()
				}
			})
			if err != nil {
				SendNotification(err.Error())
			}
			a.savedAlbums = sa
		}
		for _, v := range *a.savedAlbums {
			c = append(c, []Content{
				{Content: v.Name, Style: Defaultstyle.Foreground(tcell.ColorGreen)},
				{Content: v.Artists[0].Name, Style: Defaultstyle.Foreground(tcell.ColorPink)},
				{Content: v.ReleaseDate, Style: Defaultstyle.Foreground(tcell.ColorOrange)},
			})
		}
		return c
	}
}

func (a *AlbumsView) ContextOpener() func(m *Root, s func(s int)) { return nil }
func (a *AlbumsView) ContextHandler() func(int, int, int)         { return nil }
func (a *AlbumsView) ExternalInputCapture() func(e *tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEnter {
			r, _ := Ui.Main.Table.GetSelection()
			albumView.SetAlbum(&(*a.savedAlbums)[r])
			SetCurrentView(albumView)
		}
		return e
	}
}
func (a *AlbumsView) ContextKey() rune { return 'a' }
func (a *AlbumsView) Name() string     { return "AlbumsView" }
