package ui

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/gdamore/tcell/v2"
)

type AlbumsView struct {
	*DefaultView
	savedAlbums *spt.SavedAlbums
}

func (a *AlbumsView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if a.savedAlbums == nil {
			msg := SendNotificationWithChan("Loading Albums from your Library...")
			sa, err := spt.CurrentUserSavedAlbums(func(err error) {
				go func() {
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Albums loaded Succesfully!"
					}
				}()
			})
			if err != nil {
				SendNotification(err.Error())
			}
			a.savedAlbums = sa
		}
		for _, v := range *a.savedAlbums {
			c = append(c, []Content{
				{Content: v.Name, Style: AlbumStyle},
				{Content: v.Artists[0].Name, Style: ArtistStyle},
				{Content: v.ReleaseDate, Style: TimeStyle},
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
func (a *AlbumsView) DisableVisualMode() bool { return true }
func (a *AlbumsView) Name() string            { return "AlbumsView" }
