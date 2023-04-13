package ui

import (
	"fmt"
	"image"
	"os"

	"github.com/aditya-K2/utils"
	"github.com/nfnt/resize"
	"github.com/rivo/tview"
	"github.com/zmb3/spotify/v2"
	"gitlab.com/diamondburned/ueberzug-go"
)

type CoverArt struct {
	*tview.Box
	image *ueberzug.Image
}

func newCoverArt() *CoverArt {
	return &CoverArt{
		tview.NewBox(),
		nil,
	}
}

/* Gets the Image Struct from the provided path */
func getImg(uri string) (image.Image, error) {
	f, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	fw, fh := utils.GetFontWidth()
	img = resize.Resize(
		uint(float32(ImgW)*(fw+float32(10))),
		uint(float32(ImgH)*(fh+float32(10))),
		img,
		resize.Bilinear,
	)
	return img, nil
}

func fileName(a spotify.SimpleAlbum) string {
	return fmt.Sprintf("%s-%s.jpg", a.Name, a.Artists[0].Name)
}

func (c *CoverArt) RefreshState() {
	fw, fh := utils.GetFontWidth()
	if c.image != nil {
		c.image.Clear()
	}
	if state != nil {
		if state.Item != nil {
			if len(state.Item.Album.Images) > 0 {
				file := fileName(state.Item.Album)

				// Download Image if doesn't Exits
				if !utils.FileExists(file) {
					f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
					defer f.Close()
					if err != nil {
						SendNotification(fmt.Sprintf("Error Downloading Image: %s", err.Error()))
						return
					}
					err = state.Item.Album.Images[0].Download(f)
					if err != nil {
						SendNotification(fmt.Sprintf("Error Downloading Image: %s", err.Error()))
						return
					}
				}

				// Open the Image
				uimg, err := getImg(file)
				if err != nil {
					SendNotification(fmt.Sprintf("Error Rendering Image: %s", err.Error()))
					return
				}
				im, err := ueberzug.NewImage(uimg,
					int(float32(ImgX)*fw),
					int(float32(ImgY)*fh))
				if err != nil {
					SendNotification(fmt.Sprintf("Error Rendering Image: %s", err.Error()))
					return
				}
				c.image = im
			}
		}
	}
}
