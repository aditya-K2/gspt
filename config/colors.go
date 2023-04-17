package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/aditya-K2/utils"
	"github.com/gdamore/tcell/v2"
)

var (
	ColorError = func(s string) {
		_s := fmt.Sprintf("Wrong Color Provided: %s", s)
		utils.Print("RED", _s)
		os.Exit(-1)
	}
)

type Color struct {
	Fg     string `mapstructure:"fg"`
	Bg     string `mapstructure:"bg"`
	Bold   bool   `mapstructure:"bold"`
	Italic bool   `mapstructure:"italic"`
}

type Colors struct {
	Artist      Color `mapstructure:"artist"`
	Album       Color `mapstructure:"album"`
	Track       Color `mapstructure:"track"`
	Genre       Color `mapstructure:"genre"`
	Folder      Color `mapstructure:"folder"`
	Timestamp   Color `mapstructure:"timestamp"`
	PBarArtist  Color `mapstructure:"pbar_artist"`
	PBarTrack   Color `mapstructure:"pbar_track"`
	PlaylistNav Color `mapstructure:"playlist_nav"`
	Nav         Color `mapstructure:"nav"`
	ContextMenu Color `mapstructure:"context_menu"`
	BorderFocus Color `mapstructure:"border_focus"`
	Border      Color `mapstructure:"border"`

	Null Color
}

func (c Color) Foreground() tcell.Color {
	if strings.HasPrefix(c.Fg, "#") && len(c.Fg) == 7 {
		return tcell.GetColor(c.Fg)
	} else if val, ok := tcell.ColorNames[c.Fg]; ok {
		return val
	} else {
		ColorError(c.Fg)
		return tcell.ColorBlack
	}
}

func (c Color) Background() tcell.Color {
	if c.Bg == "" {
		return tcell.ColorBlack
	}
	if strings.HasPrefix(c.Bg, "#") && len(c.Bg) == 7 {
		return tcell.GetColor(c.Bg)
	} else if val, ok := tcell.ColorNames[c.Bg]; ok {
		return val
	} else {
		ColorError(c.Bg)
		return tcell.ColorBlack
	}
}

func (c Color) Style() tcell.Style {
	return tcell.StyleDefault.
		Foreground(c.Foreground()).
		Background(c.Background()).
		Bold(c.Bold).
		Italic(c.Italic)
}

func (c Color) String() string {
	style := ""
	if c.Bold {
		style += "b"
	}
	if c.Italic {
		style += "i"
	}
	checkColor := func(s string) string {
		var res string
		if _, ok := tcell.ColorNames[s]; ok {
			res = strings.ToLower(s)
		} else if strings.HasPrefix(s, "#") && len(s) == 7 {
			res = s
		} else {
			ColorError(s)
		}
		return res
	}
	foreground := checkColor(c.Fg)
	return fmt.Sprintf("[%s::%s]", foreground, style)
}

func NewColors() *Colors {
	return &Colors{
		Artist: Color{
			Fg:     "pink",
			Bold:   false,
			Italic: false,
		},
		Album: Color{
			Fg:     "green",
			Bold:   false,
			Italic: false,
		},
		Track: Color{
			Fg:     "blue",
			Bold:   false,
			Italic: false,
		},
		Timestamp: Color{
			Fg:     "red",
			Bold:   false,
			Italic: true,
		},
		Genre: Color{
			Fg:     "darkcyan",
			Bold:   true,
			Italic: false,
		},
		Folder: Color{
			Fg:     "yellow",
			Bold:   true,
			Italic: false,
		},
		PBarArtist: Color{
			Fg:     "blue",
			Bold:   true,
			Italic: false,
		},
		PBarTrack: Color{
			Fg:     "green",
			Bold:   true,
			Italic: true,
		},
		PlaylistNav: Color{
			Fg:     "coral",
			Bold:   false,
			Italic: false,
		},
		Nav: Color{
			Fg:     "papayawhip",
			Bold:   false,
			Italic: false,
		},
		ContextMenu: Color{
			Fg:     "turquoise",
			Bold:   true,
			Italic: false,
		},
		BorderFocus: Color{
			Fg:     "white",
			Bold:   false,
			Italic: false,
		},
		Border: Color{
			Fg:     "grey",
			Bold:   false,
			Italic: false,
		},
		Null: Color{
			Fg:     "white",
			Bold:   true,
			Italic: false,
		},
	}
}
