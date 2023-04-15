package config

import "github.com/gdamore/tcell/v2"

type Key struct {
	r rune
	k tcell.Key
}

var (
	m = map[string]tcell.Key{
		"Backspace":      tcell.KeyBackspace,
		"Tab":            tcell.KeyTab,
		"Esc":            tcell.KeyEsc,
		"Escape":         tcell.KeyEscape,
		"Enter ":         tcell.KeyEnter,
		"Backspace2":     tcell.KeyBackspace2,
		"CtrlSpace":      tcell.KeyCtrlSpace,
		"CtrlA":          tcell.KeyCtrlA,
		"CtrlB":          tcell.KeyCtrlB,
		"CtrlC":          tcell.KeyCtrlC,
		"CtrlD":          tcell.KeyCtrlD,
		"CtrlE":          tcell.KeyCtrlE,
		"CtrlF":          tcell.KeyCtrlF,
		"CtrlG":          tcell.KeyCtrlG,
		"CtrlH":          tcell.KeyCtrlH,
		"CtrlI":          tcell.KeyCtrlI,
		"CtrlJ":          tcell.KeyCtrlJ,
		"CtrlK":          tcell.KeyCtrlK,
		"CtrlL":          tcell.KeyCtrlL,
		"CtrlM":          tcell.KeyCtrlM,
		"CtrlN":          tcell.KeyCtrlN,
		"CtrlO":          tcell.KeyCtrlO,
		"CtrlP":          tcell.KeyCtrlP,
		"CtrlQ":          tcell.KeyCtrlQ,
		"CtrlR":          tcell.KeyCtrlR,
		"CtrlS":          tcell.KeyCtrlS,
		"CtrlT":          tcell.KeyCtrlT,
		"CtrlU":          tcell.KeyCtrlU,
		"CtrlV":          tcell.KeyCtrlV,
		"CtrlW":          tcell.KeyCtrlW,
		"CtrlX":          tcell.KeyCtrlX,
		"CtrlY":          tcell.KeyCtrlY,
		"CtrlZ":          tcell.KeyCtrlZ,
		"CtrlLeftSq":     tcell.KeyCtrlLeftSq,
		"CtrlBackslash":  tcell.KeyCtrlBackslash,
		"CtrlRightSq":    tcell.KeyCtrlRightSq,
		"CtrlCarat":      tcell.KeyCtrlCarat,
		"CtrlUnderscore": tcell.KeyCtrlUnderscore,
		"Up":             tcell.KeyUp,
		"Down":           tcell.KeyDown,
		"Right":          tcell.KeyRight,
		"Left":           tcell.KeyLeft,
		"UpLeft":         tcell.KeyUpLeft,
		"UpRight":        tcell.KeyUpRight,
		"DownLeft":       tcell.KeyDownLeft,
		"DownRight":      tcell.KeyDownRight,
		"Center":         tcell.KeyCenter,
		"PgUp":           tcell.KeyPgUp,
		"PgDn":           tcell.KeyPgDn,
		"Home":           tcell.KeyHome,
		"End":            tcell.KeyEnd,
		"Insert":         tcell.KeyInsert,
		"Delete":         tcell.KeyDelete,
		"Help":           tcell.KeyHelp,
		"Exit":           tcell.KeyExit,
		"Clear":          tcell.KeyClear,
		"Cancel":         tcell.KeyCancel,
		"Print":          tcell.KeyPrint,
		"Pause":          tcell.KeyPause,
		"Backtab":        tcell.KeyBacktab,
	}
)

func (k *Key) IsRune() bool {
	if k.r == 0 {
		return true
	}
	return false
}

func (k *Key) Rune() rune {
	return k.r
}

func (k *Key) Key() tcell.Key {
	return k.k
}

func NewKey(s string) Key {
	if len(s) == 1 {
		a := []rune(s)
		if (a[0] >= 'A' && a[0] <= 'Z') || (a[0] >= 'a' && a[0] <= 'z') {
			return Key{r: a[0]}
		}
	}
	if val, ok := m[s]; ok {
		return Key{k: val}
	}
	return Key{}
}
