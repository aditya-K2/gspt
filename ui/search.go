package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/aditya-K2/tview"
)

func NewSearchBar() *tview.InputField {
	T := tview.NewInputField()

	T.SetFieldBackgroundColor(tcell.ColorDefault)
	T.SetBackgroundColor(tcell.ColorDefault)
	T.SetTitle("Search").SetTitleAlign(tview.AlignLeft)
	T.SetBorder(true)
	T.SetAutocompleteStyles(
		tcell.ColorBlack,
		tcell.StyleDefault,
		tcell.StyleDefault.Reverse(true),
	)
	// T.SetAutocompleteMatchFieldWidth(true)
	T.SetDoneFunc(func(k tcell.Key) {
		switch k {
		case tcell.KeyEscape:
			{
				Ui.App.SetFocus(Ui.Main.Table)
				Ui.SearchBar.SetText("")
			}
		case tcell.KeyEnter:
			{
				searchView.SetSearch(T.GetText())
				SetCurrentView(searchView)
				Ui.App.SetFocus(Ui.Main.Table)
				Ui.SearchBar.SetText("")
			}
		}
	})

	T.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyCtrlP {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		}
		if e.Key() == tcell.KeyCtrlN {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		}
		return e
	})

	return T
}
