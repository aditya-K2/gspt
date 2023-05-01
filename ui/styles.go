package ui

import (
	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/tview"
)

var (
	borders = map[bool]map[string]rune{
		true: {
			"TopLeft":          '╭',
			"TopRight":         '╮',
			"BottomRight":      '╯',
			"BottomLeft":       '╰',
			"Vertical":         '│',
			"Horizontal":       '─',
			"TopLeftFocus":     '╭',
			"TopRightFocus":    '╮',
			"BottomRightFocus": '╯',
			"BottomLeftFocus":  '╰',
			"VerticalFocus":    '│',
			"HorizontalFocus":  '─',
		},
		false: {
			"TopLeft":          tview.Borders.TopLeft,
			"TopRight":         tview.Borders.TopRight,
			"BottomRight":      tview.Borders.BottomRight,
			"BottomLeft":       tview.Borders.BottomLeft,
			"Vertical":         tview.Borders.Vertical,
			"Horizontal":       tview.Borders.Horizontal,
			"TopLeftFocus":     tview.Borders.TopLeftFocus,
			"TopRightFocus":    tview.Borders.TopRightFocus,
			"BottomRightFocus": tview.Borders.BottomRightFocus,
			"BottomLeftFocus":  tview.Borders.BottomLeftFocus,
			"VerticalFocus":    tview.Borders.VerticalFocus,
			"HorizontalFocus":  tview.Borders.HorizontalFocus,
		},
	}
)

func setBorderRunes() {
	tview.Borders.TopLeft = borders[config.Config.RoundedCorners]["TopLeft"]
	tview.Borders.TopRight = borders[config.Config.RoundedCorners]["TopRight"]
	tview.Borders.BottomRight = borders[config.Config.RoundedCorners]["BottomRight"]
	tview.Borders.BottomLeft = borders[config.Config.RoundedCorners]["BottomLeft"]
	tview.Borders.Vertical = borders[config.Config.RoundedCorners]["Vertical"]
	tview.Borders.Horizontal = borders[config.Config.RoundedCorners]["Horizontal"]
	tview.Borders.TopLeftFocus = borders[config.Config.RoundedCorners]["TopLeftFocus"]
	tview.Borders.TopRightFocus = borders[config.Config.RoundedCorners]["TopRightFocus"]
	tview.Borders.BottomRightFocus = borders[config.Config.RoundedCorners]["BottomRightFocus"]
	tview.Borders.BottomLeftFocus = borders[config.Config.RoundedCorners]["BottomLeftFocus"]
	tview.Borders.VerticalFocus = borders[config.Config.RoundedCorners]["VerticalFocus"]
	tview.Borders.HorizontalFocus = borders[config.Config.RoundedCorners]["HorizontalFocus"]
}

func setStyles() {
	TrackStyle = config.Config.Colors.Track.Style()
	AlbumStyle = config.Config.Colors.Album.Style()
	ArtistStyle = config.Config.Colors.Artist.Style()
	TimeStyle = config.Config.Colors.Timestamp.Style()
	GenreStyle = config.Config.Colors.Genre.Style()
	PlaylistNavStyle = config.Config.Colors.PlaylistNav.Style()
	NavStyle = config.Config.Colors.Nav.Style()
	ContextMenuStyle = config.Config.Colors.ContextMenu.Style()
	NotSelectableStyle = config.Config.Colors.Null.Style()
	tview.Styles.BorderColorFocus = config.Config.Colors.BorderFocus.Foreground()
	tview.Styles.BorderColor = config.Config.Colors.Border.Foreground()
}
