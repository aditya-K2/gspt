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
			msg := SendNotificationWithChan("Loading Artists from your Library...")
			fa, err := spt.CurrentUserFollowedArtists(func(err error) {
				go func() {
					if err != nil {
						msg <- err.Error()
					} else {
						msg <- "Artists loaded Succesfully!"
					}
				}()
			})
			if err != nil {
				SendNotification(err.Error())
			}
			a.followedArtists = fa
		}
		for _, v := range *a.followedArtists {
			c = append(c, []Content{
				{Content: v.Name, Style: ArtistStyle},
				// {Content: v.Genres[0], Style: AlbumStyle},
			})
		}
		return c
	}
}

func (a *ArtistsView) OpenArtist() {
	r, _ := Main.Table.GetSelection()
	artistView.SetArtist(&(*a.followedArtists)[r].ID)
	artistView.RefreshState()
	SetCurrentView(artistView)
}

func (a *ArtistsView) Name() string { return "ArtistsView" }
