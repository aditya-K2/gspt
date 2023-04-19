package ui

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/tview"
	"github.com/aditya-K2/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/nfnt/resize"
	"github.com/zmb3/spotify/v2"
	"gitlab.com/diamondburned/ueberzug-go"
)

type CoverArt struct {
	*tview.Box
	image *ueberzug.Image
}

func newCoverArt() *CoverArt {
	return &CoverArt{
		tview.NewBox().SetBorder(true).SetBackgroundColor(tcell.ColorDefault),
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
	fw, fh, err := getFontWidth()
	if err != nil {
		return nil, err
	}
	img = resize.Resize(
		uint((ImgW*fw)+config.Config.ImageWidthExtraX),
		uint((ImgH*fh)+config.Config.ImageWidthExtraY),
		img,
		resize.Bilinear,
	)
	return img, nil
}

func fileName(a spotify.SimpleAlbum) string {
	return fmt.Sprintf(config.Config.CacheDir+"/%s.jpg", a.ID)
}

func getFontWidth() (int, int, error) {
	w, h, err := ueberzug.GetParentSize()
	if err != nil {
		return 0, 0, err
	}
	_, _, rw, rh := root.Root.GetRect()
	if rw == 0 || rh == 0 {
		return 0, 0, errors.New("Unable to get row width and height")
	}
	fw := w / rw
	fh := h / rh
	return fw, fh, nil
}

func (c *CoverArt) RefreshState() {
	if c.image != nil {
		c.image.Clear()
	}
	fw, fh, err := getFontWidth()
	if err != nil {
		SendNotification(err.Error())
		return
	}
	if state != nil {
		if state.Item != nil {
			if len(state.Item.Album.Images) > 0 {
				file := fileName(state.Item.Album)

				// Download Image if doesn't Exits
				if !utils.FileExists(file) {
					msg := SendNotificationWithChan("Downloading Cover Art...")
					f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
					defer f.Close()
					if err != nil {
						go func() {
							msg <- fmt.Sprintf("Error Downloading Image: %s", err.Error())
						}()
						return
					}
					err = state.Item.Album.Images[0].Download(f)
					if err != nil {
						go func() {
							msg <- fmt.Sprintf("Error Downloading Image: %s", err.Error())
						}()
						return
					}
					go func() {
						msg <- "Image Downloaded Succesfully!"
					}()
				}

				// Open the Image
				uimg, err := getImg(file)
				if err != nil {
					SendNotification(fmt.Sprintf("Error Rendering Image: %s", err.Error()))
					return
				}
				im, err := ueberzug.NewImage(uimg,
					int(ImgX*fw)+config.Config.AdditionalPaddingX,
					int(ImgY*fh)+config.Config.AdditionalPaddingY)
				if err != nil {
					SendNotification(fmt.Sprintf("Error Rendering Image: %s", err.Error()))
					return
				}
				c.image = im
			}
		}
	}
}
