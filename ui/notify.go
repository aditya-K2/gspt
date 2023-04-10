package ui

import (
	"sync"
	"time"

	"github.com/aditya-K2/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	maxNotifications = 3
	pm               sync.Mutex
	notAvailable     = -1
	posArr           = positionArray{}
	c                chan *notification
)

// Start Notification Service
func InitNotifier() {
	for _m := maxNotifications; _m != 0; _m-- {
		posArr = append(posArr, true)
	}
	c = make(chan *notification, maxNotifications)
	routine()
}

// notification Primitive
type notification struct {
	*tview.Box
	Text     string
	Position int
	close    chan bool
	timer    time.Duration
}

// Array for all available positions where the notification can be displayed.
type positionArray []bool

// Check If there is a position available.
func (p *positionArray) Available() bool {
	var t = false
	pm.Lock()
	for _, v := range *p {
		t = t || v
	}
	pm.Unlock()
	return t
}

func (p *positionArray) GetNextPosition() int {
	pm.Lock()
	v := *p
	for k := range v {
		if v[k] {
			v[k] = false
			pm.Unlock()
			return k
		}
	}
	pm.Unlock()
	return notAvailable
}

// Free a position
func (p *positionArray) Free(i int) {
	pm.Lock()
	v := *p
	v[i] = true
	pm.Unlock()
}

// Get A Pointer to A Notification Struct
func newNotificationWithTimer(s string, t time.Duration) *notification {
	return &notification{
		Box:   tview.NewBox(),
		Text:  s,
		timer: t,
		close: nil,
	}
}

// Get A Pointer to A Notification Struct with a close channel
func newNotificationWithChan(s string, c chan bool) *notification {
	return &notification{
		Box:   tview.NewBox(),
		Text:  s,
		close: c,
	}
}

// Draw Function for the Notification Primitive
func (self *notification) Draw(screen tcell.Screen) {
	termDetails := utils.GetWidth()
	pos := (self.Position*3 + self.Position + 1)

	var (
		COL          int = int(termDetails.Col)
		TEXTLENGTH   int = len(self.Text)
		HEIGHT       int = 3
		TextPosition int = 1
	)

	self.Box.SetBackgroundColor(tcell.ColorBlack)
	self.SetRect(COL-(TEXTLENGTH+7), pos, TEXTLENGTH+4, HEIGHT)
	self.DrawForSubclass(screen, self.Box)
	tview.Print(screen, self.Text,
		COL-(TEXTLENGTH+5), pos+TextPosition, TEXTLENGTH,
		tview.AlignCenter, tcell.ColorWhite)
}

// this routine checks for available position and sends notification if
// position is available.
func routine() {
	go func() {
		for {
			val := <-c
			// Wait until a new position isn't available
			for !posArr.Available() {
				continue
			}
			notify(val)
		}
	}()
}

func notify(n *notification) {
	go func() {
		currentTime := time.Now().String()
		npos := posArr.GetNextPosition()
		// Ensure a position is available.
		if npos == notAvailable {
			for !posArr.Available() {
			}
			npos = posArr.GetNextPosition()
		}
		n.Position = npos
		Ui.Root.Root.AddPage(currentTime, n, false, true)
		Ui.App.SetFocus(Ui.Main.Table)
		if n.close != nil {
			<-n.close
		} else {
			time.Sleep(n.timer)
		}
		Ui.Root.Root.RemovePage(currentTime)
		posArr.Free(npos)
		Ui.App.SetFocus(Ui.Main.Table)
	}()
}

func SendNotification(text string) {
	SendNotificationWithTimer(text, time.Second)
}

func SendNotificationWithTimer(text string, t time.Duration) {
	go func() {
		c <- newNotificationWithTimer(text, t)
	}()
}

func SendNotificationWithChan(text string) chan bool {
	close := make(chan bool)
	go func() {
		c <- newNotificationWithChan(text, close)
	}()
	return close
}
