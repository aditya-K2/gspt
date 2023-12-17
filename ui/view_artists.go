package ui

import (
	"github.com/aditya-K2/gspt/spt"
)

type ArtistsView struct {
	*DefaultViewNone
	followedArtists *spt.FollowedArtists
}

func NewArtistsView() *ArtistsView {
	a := &ArtistsView{
		&DefaultViewNone{&defView{}},
		nil,
	}
	return a
}

func (a *ArtistsView) Content() func() [][]Content {
	return func() [][]Content {
		c := make([][]Content, 0)
		if a.followedArtists == nil {
			msg := SendNotificationWithChan("Fetching Artists from your Library...")
			fa, ch := spt.CurrentUserFollowedArtists()
			go func() {
				err := <-ch
				if err != nil {
					msg <- err.Error()
				} else {
					msg <- "Artists loaded Succesfully!"
				}
			}()
			a.followedArtists = fa
		}
		if a.followedArtists != nil {
			for _, v := range *a.followedArtists {
				c = append(c, []Content{
					{Content: v.Name, Style: ArtistStyle},
					{Content: mergeGenres(v.Genres), Style: AlbumStyle},
				})
			}
		}
		return c
	}
}

func (a *ArtistsView) OpenEntry() {
	r, _ := Main.GetSelection()
	artistView.SetArtist(&(*a.followedArtists)[r].ID)
	SetCurrentView(artistView)
}

func (a *ArtistsView) PlayEntry() {
	r, _ := Main.GetSelection()
	if err := spt.PlayContext((*a.followedArtists)[r].URI); err != nil {
		SendNotification(err.Error())
	}
}

func (a *ArtistsView) Name() string { return "ArtistsView" }
