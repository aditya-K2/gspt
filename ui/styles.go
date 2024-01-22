package ui

import (
	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/tview"
)

var (
	borders = map[string]map[string]rune{
		config.CornersRounded: {
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
		config.CornersDefault: {
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
	tview.Borders.TopLeft = borders[cfg.Corners]["TopLeft"]
	tview.Borders.TopRight = borders[cfg.Corners]["TopRight"]
	tview.Borders.BottomRight = borders[cfg.Corners]["BottomRight"]
	tview.Borders.BottomLeft = borders[cfg.Corners]["BottomLeft"]
	tview.Borders.Vertical = borders[cfg.Corners]["Vertical"]
	tview.Borders.Horizontal = borders[cfg.Corners]["Horizontal"]
	tview.Borders.TopLeftFocus = borders[cfg.Corners]["TopLeftFocus"]
	tview.Borders.TopRightFocus = borders[cfg.Corners]["TopRightFocus"]
	tview.Borders.BottomRightFocus = borders[cfg.Corners]["BottomRightFocus"]
	tview.Borders.BottomLeftFocus = borders[cfg.Corners]["BottomLeftFocus"]
	tview.Borders.VerticalFocus = borders[cfg.Corners]["VerticalFocus"]
	tview.Borders.HorizontalFocus = borders[cfg.Corners]["HorizontalFocus"]
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
