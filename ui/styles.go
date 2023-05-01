package ui

import (
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
	tview.Borders.TopLeft = borders[cfg.RoundedCorners]["TopLeft"]
	tview.Borders.TopRight = borders[cfg.RoundedCorners]["TopRight"]
	tview.Borders.BottomRight = borders[cfg.RoundedCorners]["BottomRight"]
	tview.Borders.BottomLeft = borders[cfg.RoundedCorners]["BottomLeft"]
	tview.Borders.Vertical = borders[cfg.RoundedCorners]["Vertical"]
	tview.Borders.Horizontal = borders[cfg.RoundedCorners]["Horizontal"]
	tview.Borders.TopLeftFocus = borders[cfg.RoundedCorners]["TopLeftFocus"]
	tview.Borders.TopRightFocus = borders[cfg.RoundedCorners]["TopRightFocus"]
	tview.Borders.BottomRightFocus = borders[cfg.RoundedCorners]["BottomRightFocus"]
	tview.Borders.BottomLeftFocus = borders[cfg.RoundedCorners]["BottomLeftFocus"]
	tview.Borders.VerticalFocus = borders[cfg.RoundedCorners]["VerticalFocus"]
	tview.Borders.HorizontalFocus = borders[cfg.RoundedCorners]["HorizontalFocus"]
}

func setStyles() {
	TrackStyle = cfg.Colors.Track.Style()
	AlbumStyle = cfg.Colors.Album.Style()
	ArtistStyle = cfg.Colors.Artist.Style()
	TimeStyle = cfg.Colors.Timestamp.Style()
	GenreStyle = cfg.Colors.Genre.Style()
	PlaylistNavStyle = cfg.Colors.PlaylistNav.Style()
	NavStyle = cfg.Colors.Nav.Style()
	ContextMenuStyle = cfg.Colors.ContextMenu.Style()
	NotSelectableStyle = cfg.Colors.Null.Style()
	tview.Styles.BorderColorFocus = cfg.Colors.BorderFocus.Foreground()
	tview.Styles.BorderColor = cfg.Colors.Border.Foreground()
}
