package main

import (
	"fmt"

	"github.com/aditya-K2/gspot/gspotify"
	"github.com/aditya-K2/gspot/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	i := ui.NewInteractiveView()
	m := ui.NewMain(i.View)
	var err error
	gspotify.Client, err = gspotify.NewClient()
	if err != nil {
		panic(err)
	}
	albs, err := gspotify.CurrentUserSavedAlbums(func(status bool, err error) {
		fmt.Println("Done")
	})
	if err != nil {
		panic(err)
	}
	content := func() [][]ui.Content {
		a := *albs
		c := make([][]ui.Content, 0)
		for _, v := range a {
			c = append(c, []ui.Content{
				{Content: v.Name, Style: ui.Defaultstyle.Foreground(tcell.ColorBlue)},
				{Content: v.Artists[0].Name, Style: ui.Defaultstyle.Foreground(tcell.ColorPink)},
				{Content: v.ReleaseDate, Style: ui.Defaultstyle.Foreground(tcell.ColorGreen)},
			})
		}
		return c
	}
	i.SetContentFunc(content)
	playlists, err := gspotify.CurrentUserPlaylists(func(status bool, err error) {
		fmt.Println("Done")
	})
	contextOpener := func() {
		c := ui.NewMenu()
		if err != nil {
			panic(err)
		}

		cc := []string{}
		for _, v := range *playlists {
			cc = append(cc, v.Name)
		}
		c.Content(cc)
		c.Title("Add to Playlist")
		c.SetSelectionHandler(i.SelectionHandler)
		m.AddCenteredWidget(c)
	}
	i.SetContextKey('a')
	i.SetContextOpener(contextOpener)
	contextHandler := func(start, end, sel int) {
		for k := start; k <= end; k++ {
			if err := gspotify.AddAlbumToPlaylist((*albs)[k].ID, (*playlists)[sel].ID); err != nil {
				panic(err)
			}
		}
	}
	i.SetContextHandler(contextHandler)
	i.SetExternalCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyEnter {
			r, _ := i.View.GetSelection()
			if err := gspotify.PlayContext(&(*albs)[r].URI); err != nil {
				panic(err)
			}
		}
		return e
	})
	if err := tview.NewApplication().SetRoot(m.Root, true).Run(); err != nil {
		panic(err)
	}
}
