package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
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
	text     string
	position int
	msg      chan string
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
		text:  s,
		timer: t,
		msg:   nil,
	}
}

// Get A Pointer to A Notification Struct with a close channel
func newNotificationWithChan(s string, c chan string) *notification {
	return &notification{
		Box:   tview.NewBox(),
		text:  s,
		timer: time.Second,
		msg:   c,
	}
}

// Draw Function for the Notification Primitive
func (self *notification) Draw(screen tcell.Screen) {
	pos := (self.position*3 + self.position + 1)

	var (
		_, _, COL, _ int = Flex.GetRect()
		TEXTLENGTH   int = len(self.text)
		HEIGHT       int = 3
		TextPosition int = 1
	)

	self.Box.SetBackgroundColor(tcell.ColorBlack)
	self.SetRect(COL-(TEXTLENGTH+7), pos, TEXTLENGTH+4, HEIGHT)
	self.DrawForSubclass(screen, self.Box)
	tview.Print(screen, self.text,
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
		n.position = npos
		root.Root.AddPage(currentTime, n, false, true)
		App.Draw()
		App.SetFocus(Main)
		if n.msg != nil {
			n.text = <-n.msg
			App.Draw()
		}
		time.Sleep(n.timer)
		root.Root.RemovePage(currentTime)
		posArr.Free(npos)
		App.SetFocus(Main)
		App.Draw()
	}()
}

func SendNotification(format string, args ...any) {
	SendNotificationWithTimer(time.Second, fmt.Sprintf(format, args...))
}

func SendNotificationWithTimer(t time.Duration, format string, args ...any) {
	go func() {
		c <- newNotificationWithTimer(fmt.Sprintf(format, args...), t)
	}()
}

// SendNotificationWithChan sends a notification that won't be closed unless
// an update message isn't sent over the channel that it returns. The message
// string received over the channel is then displayed for half a second and the
// notification is then removed.
func SendNotificationWithChan(format string, args ...any) chan string {
	close := make(chan string)
	go func() {
		c <- newNotificationWithChan(fmt.Sprintf(format, args...), close)
	}()
	return close
}
