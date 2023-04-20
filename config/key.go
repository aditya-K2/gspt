package config

import (
	"github.com/gdamore/tcell/v2"
)

type Key struct {
	R rune
	K tcell.Key
}

var (
	M = map[string]tcell.Key{
		"backspace":       tcell.KeyBackspace,
		"tab":             tcell.KeyTab,
		"esc":             tcell.KeyEsc,
		"escape":          tcell.KeyEscape,
		"enter":           tcell.KeyEnter,
		"backspace2":      tcell.KeyBackspace2,
		"ctrl-space":      tcell.KeyCtrlSpace,
		"ctrl-a":          tcell.KeyCtrlA,
		"ctrl-b":          tcell.KeyCtrlB,
		"ctrl-c":          tcell.KeyCtrlC,
		"ctrl-d":          tcell.KeyCtrlD,
		"ctrl-e":          tcell.KeyCtrlE,
		"ctrl-f":          tcell.KeyCtrlF,
		"ctrl-g":          tcell.KeyCtrlG,
		"ctrl-h":          tcell.KeyCtrlH,
		"ctrl-i":          tcell.KeyCtrlI,
		"ctrl-j":          tcell.KeyCtrlJ,
		"ctrl-k":          tcell.KeyCtrlK,
		"ctrl-l":          tcell.KeyCtrlL,
		"ctrl-m":          tcell.KeyCtrlM,
		"ctrl-n":          tcell.KeyCtrlN,
		"ctrl-o":          tcell.KeyCtrlO,
		"ctrl-p":          tcell.KeyCtrlP,
		"ctrl-q":          tcell.KeyCtrlQ,
		"ctrl-r":          tcell.KeyCtrlR,
		"ctrl-s":          tcell.KeyCtrlS,
		"ctrl-t":          tcell.KeyCtrlT,
		"ctrl-u":          tcell.KeyCtrlU,
		"ctrl-v":          tcell.KeyCtrlV,
		"ctrl-w":          tcell.KeyCtrlW,
		"ctrl-x":          tcell.KeyCtrlX,
		"ctrl-y":          tcell.KeyCtrlY,
		"ctrl-z":          tcell.KeyCtrlZ,
		"ctrl-leftsq":     tcell.KeyCtrlLeftSq,
		"ctrl-backslash":  tcell.KeyCtrlBackslash,
		"ctrl-rightsq":    tcell.KeyCtrlRightSq,
		"ctrl-carat":      tcell.KeyCtrlCarat,
		"ctrl-underscore": tcell.KeyCtrlUnderscore,
		"up":              tcell.KeyUp,
		"down":            tcell.KeyDown,
		"right":           tcell.KeyRight,
		"left":            tcell.KeyLeft,
		"up_left":         tcell.KeyUpLeft,
		"up_right":        tcell.KeyUpRight,
		"down_left":       tcell.KeyDownLeft,
		"down_right":      tcell.KeyDownRight,
		"center":          tcell.KeyCenter,
		"pgup":            tcell.KeyPgUp,
		"pgdn":            tcell.KeyPgDn,
		"home":            tcell.KeyHome,
		"end":             tcell.KeyEnd,
		"insert":          tcell.KeyInsert,
		"delete":          tcell.KeyDelete,
		"help":            tcell.KeyHelp,
		"exit":            tcell.KeyExit,
		"clear":           tcell.KeyClear,
		"cancel":          tcell.KeyCancel,
		"print":           tcell.KeyPrint,
		"pause":           tcell.KeyPause,
		"backtab":         tcell.KeyBacktab,
	}
	RuneKeys = map[rune]bool{
		'!':  true,
		'@':  true,
		'#':  true,
		'$':  true,
		'%':  true,
		'^':  true,
		'&':  true,
		'*':  true,
		'(':  true,
		')':  true,
		'-':  true,
		'=':  true,
		'_':  true,
		'+':  true,
		',':  true,
		'.':  true,
		'<':  true,
		'>':  true,
		'/':  true,
		'?':  true,
		'[':  true,
		']':  true,
		'{':  true,
		'}':  true,
		'|':  true,
		'\\': true,
		':':  true,
		';':  true,
		'\'': true,
		'"':  true,
		' ':  true,
	}
)

func (k *Key) Rune() rune {
	return k.R
}

func (k *Key) Key() tcell.Key {
	return k.K
}

func NewKey(s string) Key {
	if len(s) == 1 {
		a := []rune(s)
		if (a[0] >= 'A' && a[0] <= 'Z') ||
			(a[0] >= 'a' && a[0] <= 'z') ||
			(a[0] >= '0' && a[0] <= '9') ||
			(RuneKeys[a[0]]) {
			return Key{R: a[0]}
		}
	}
	if val, ok := M[s]; ok {
		return Key{K: val}
	}
	return Key{}
}
