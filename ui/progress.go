package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/zmb3/spotify/v2"

	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/utils"
	"github.com/rivo/tview"
)

var (
	state     *spotify.PlayerState
	stateLock sync.Mutex
	ctrackId  spotify.ID
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

func GetProgressGlyph(width, percentage float64, btext string) string {
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
		GetProgressGlyph(float64(_width-OFFSET-1),
			percentage,
			self.BarText),
		x, y+2, _width-OFFSET, tview.AlignRight, tcell.ColorWhite)
}

func RefreshProgress() {
	s, err := spt.GetPlayerState()
	if err != nil {
		SendNotification(err.Error())
		return
	}
	stateLock.Lock()
	state = s
	stateLock.Unlock()
	if Ui != nil && Ui.CoverArt != nil {
		if s.Item != nil && s.Item.ID != ctrackId {
			ctrackId = s.Item.ID
			Ui.CoverArt.RefreshState()
		}
	}
}

func RefreshProgressLocal() {
	stateLock.Lock()
	if state != nil {
		if state.Item != nil && state.Playing {
			if state.Item.Duration-state.Progress >= 1000 {
				state.Progress += 1000
			} else {
				state.Progress += state.Item.Duration - state.Progress
				go func() {
					RefreshProgress()
				}()
			}
		}
	}
	stateLock.Unlock()
}

func updateRoutine() {
	RefreshProgress()
	go func() {
		localTicker := time.NewTicker(time.Second)
		spotifyTicker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-spotifyTicker.C:
				{
					RefreshProgress()
				}
			case <-localTicker.C:
				{
					RefreshProgressLocal()
				}
			}
		}
	}()
}

func progressFunc() (string, string, string, float64) {
	percentage := 0.0
	barTitle := " - "
	barText := "---:---"
	barTopTitle := "[]"
	if state != nil {
		barTopTitle = fmt.Sprintf("[Device: %s Shuffle: %t Repeat: %s]", state.Device.Name, state.ShuffleState, state.RepeatState)
		if state.Item != nil {
			barTitle = fmt.Sprintf("%s - %s", state.Item.Name, state.Item.Artists[0].Name)
			barText = utils.StrTime(float64(state.Progress/1000)) + "/" + utils.StrTime(float64(state.Item.Duration/1000))
			percentage = (float64(state.Progress) / float64(state.Item.Duration)) * 100
		}
	}
	return barTitle, barTopTitle, barText, percentage
}
