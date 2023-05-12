package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"

	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/tview"
	"github.com/aditya-K2/utils"
)

var (
	state    *spotify.PlayerState
	ctrackId spotify.ID
	devices  = map[string]string{
		"computer":     "󰍹",
		"tablet":       "",
		"smartphone":   "󰄜",
		"speaker":      "󰓃",
		"tv":           "",
		"avr":          "󰤽",
		"stb":          "󰤽",
		"audio_dongle": "󱡬",
		"game_console": "󰺵",
		"cast_video":   "󰄙",
		"cast_audio":   "󰄙",
		"automobile":   "",
	}
	playIcons = map[bool]string{
		true:  "",
		false: "",
	}
	shuffleIcons = map[bool]string{
		true:  "󰒟",
		false: "󰒞",
	}
	repeatIcons = map[string]string{
		"track":   "󰑘",
		"off":     "󰑗",
		"context": "󰑖",
	}
)

// ProgressBar is a two-lined Box. First line is the BarTitle
// Second being the actual progress done.
// Use SetProgressFunc to provide the callback which provides the Fields each time the ProgressBar will be Drawn.
// The progressFunc must return (BarTitle, BarTopTitle, BarText, percentage) respectively
type ProgressBar struct {
	*tview.Box
	BarTitle     string
	BarText      string
	BarTopTitle  string
	progressFunc func() (BarTitle string,
		BarTopTitle string,
		BarText string,
		percentage float64)
}

func (self *ProgressBar) SetProgressFunc(pfunc func() (string, string, string, float64)) *ProgressBar {
	self.progressFunc = pfunc
	return self
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{
		Box: tview.NewBox(),
	}
}

func getProgressGlyph(width, percentage float64, btext string) string {
	q := "[black:white:b]"
	var a string
	a += strings.Repeat(" ", int(width)-len(btext))
	a = utils.InsertAt(a, btext, int(width/2)-10)
	a = utils.InsertAt(a, "[-:-:-]", int(width*percentage/100))
	q += a
	return q
}

func (self *ProgressBar) Draw(screen tcell.Screen) {
	var (
		OFFSET int = 1
	)
	self.Box.SetBorder(true)
	self.Box.SetBackgroundColor(tcell.ColorDefault)
	var percentage float64
	self.BarTitle, self.BarTopTitle, self.BarText, percentage = self.progressFunc()
	self.DrawForSubclass(screen, self.Box)
	self.Box.SetTitle(self.BarTopTitle)
	self.Box.SetTitleAlign(tview.AlignRight)
	x, y, _width, _ := self.Box.GetInnerRect()
	tview.Print(screen, self.BarTitle, x+OFFSET, y, _width, tview.AlignLeft, tcell.ColorWhite)
	tview.Print(screen,
		getProgressGlyph(float64(_width-OFFSET-1),
			percentage,
			self.BarText),
		x, y+2, _width-OFFSET, tview.AlignRight, tcell.ColorWhite)
}
func (self *ProgressBar) RefreshState() {
	RefreshProgress(false)
}

func RefreshProgress(force bool) {
	s, err := spt.GetPlayerState()
	if err != nil {
		SendNotification(err.Error())
		return
	}
	state = s

	// Reset the "cached" if nothing is playing (TODO: better name)[
	if s.Item == nil {
		ctrackId = ""
	}

	if coverArt != nil {
		// If No Item is playing
		if (state.Item == nil) ||
			// An Item is Playing but doesn't match the cached Track ID
			(state.Item != nil && state.Item.ID != ctrackId) ||
			// Forced Redrawing
			force {
			if state.Item != nil {
				ctrackId = state.Item.ID
			}
			if !cfg.HideImage {
				coverArt.RefreshState()
			}
		}
	}
}

func refreshProgressLocal() {
	if state != nil {
		if state.Item != nil && state.Playing {
			if state.Item.Duration-state.Progress >= 1000 {
				state.Progress += 1000
			} else {
				state.Progress += state.Item.Duration - state.Progress
				go func() {
					RefreshProgress(false)
				}()
			}
		}
	}
}

func progressRoutine() {
	RefreshProgress(false)
	go func() {
		localTicker := time.NewTicker(time.Second)
		spotifyTicker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-spotifyTicker.C:
				{
					RefreshProgress(false)
				}
			case <-localTicker.C:
				{
					refreshProgressLocal()
				}
			}
		}
	}()
}

func deviceIcon(d spotify.PlayerDevice) string {
	if val, ok := devices[strings.ToLower(d.Type)]; cfg.UseNerdIcons && ok {
		return val
	}
	return "Device:"
}

func topTitle(playing, shuffle bool, repeat string, device spotify.PlayerDevice) (
	playState, deviceIcn, deviceName, shuffleState, repeatState string) {

	icon := ""
	playState = "Paused"
	shuffleState = fmt.Sprintf("%t", shuffle)

	if playing {
		playState = "Playing"
	}

	if cfg.UseNerdIcons {
		icon = playIcons[playing]
		icon += " "
		repeat = repeatIcons[repeat]
		shuffleState = shuffleIcons[shuffle]
	}

	playState = icon + playState
	deviceIcn = deviceIcon(device)
	deviceName = device.Name
	repeatState = repeat
	return
}

func progressFunc() (barTitle, barTopTitle, barText string, percentage float64) {
	percentage = 0
	barTitle = " - "
	barText = "---:---"
	barTopTitle = "[]"
	if state != nil {
		barTopTitle = fmt.Sprintf("[ %s %s %s Shuffle: %s Repeat: %s ]",
			wrap(
				topTitle(
					state.Playing,
					state.ShuffleState,
					state.RepeatState,
					state.Device))...)
		if state.Item != nil {
			barTitle = fmt.Sprintf("%s%s[-:-:-] - %s%s",
				cfg.Colors.PBarTrack.String(), state.Item.Name,
				cfg.Colors.PBarArtist.String(), artistName(state.Item.Artists))
			barText = utils.StrTime(float64(state.Progress/1000)) +
				"/" +
				utils.StrTime(float64(state.Item.Duration/1000))
			percentage = (float64(state.Progress) / float64(state.Item.Duration)) * 100
		}
	}
	return
}
