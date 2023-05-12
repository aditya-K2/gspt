package ui

import (
	"github.com/aditya-K2/gspt/spt"
)

type AlbumsView struct {
	*DefaultViewNone
	savedAlbums *spt.SavedAlbums
}

func NewAlbumsView() *AlbumsView {
	a := &AlbumsView{
		&DefaultViewNone{&defView{}},
		nil,
	}
	return a
}

func (a *AlbumsView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if a.savedAlbums == nil {
			msg := SendNotificationWithChan("Loading Albums from your Library...")
			sa, ch := spt.CurrentUserSavedAlbums()
			go func() {
				err := <-ch
				if err != nil {
					msg <- err.Error()
				} else {
					msg <- "Albums loaded Succesfully!"
				}
			}()
			a.savedAlbums = sa
		}
		if a.savedAlbums != nil {
			for _, v := range *a.savedAlbums {
				c = append(c, []Content{
					{Content: v.Name, Style: AlbumStyle},
					{Content: artistName(v.Artists), Style: ArtistStyle},
					{Content: v.ReleaseDate, Style: TimeStyle},
				})
			}
		}
		return c
	}
}

func (a *AlbumsView) OpenEntry() {
	r, _ := Main.GetSelection()
	albumView.SetAlbum((*a.savedAlbums)[r].Name, &(*a.savedAlbums)[r].ID)
	SetCurrentView(albumView)
}

func (a *AlbumsView) PlayEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlayContext((*a.savedAlbums)[r].URI); err != nil {
		SendNotification(err.Error())
	}
}

func (a *AlbumsView) QueueSelectEntry() {
	r, _ := Main.GetSelection()
	alb := (*a.savedAlbums)[r]
	msg := SendNotificationWithChan("Queueing " + alb.Name + "...")
	go func() {
		if err := spt.QueueAlbum(alb.ID); err != nil {
			msg <- err.Error()
		} else {
			msg <- (alb.Name) + " queued succesfully!"
		}
	}()
}

func (a *AlbumsView) Name() string { return "AlbumsView" }
