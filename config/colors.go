// # Colors and Style
//
// You can change `colors` and `styles` for some of the aspects of `gomp`
//
// #### Let's say to you want to change Color of Artist from the default Purple to Red
//
// In your `config.yml`
// ```yml
// COLORS:
//
// artist:
//
//	foreground: Red
//
// # Another Example
// pbar_artist:
//
//	foreground: "#ff0000" # For Hex Values
//	bold: True # Changes the Style
//	italic: False
//
// ```
//
// ![Dec30(Fri)012241PM](https://user-images.githubusercontent.com/51816057/210048064-b2816095-10f2-4f0b-83ed-0e87d636b894.png)
// ![Dec30(Fri)012315PM](https://user-images.githubusercontent.com/51816057/210048069-8e91509a-17a5-46da-a65e-ff8f427dde17.png)
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
	DColors = map[string]tcell.Color{
		"Black":       tcell.ColorBlack,
		"Maroon":      tcell.ColorMaroon,
		"Green":       tcell.ColorGreen,
		"Olive":       tcell.ColorOlive,
		"Navy":        tcell.ColorNavy,
		"Purple":      tcell.ColorPurple,
		"Teal":        tcell.ColorTeal,
		"Silver":      tcell.ColorSilver,
		"Gray":        tcell.ColorGray,
		"Red":         tcell.ColorRed,
		"Lime":        tcell.ColorLime,
		"Yellow":      tcell.ColorYellow,
		"Blue":        tcell.ColorBlue,
		"Fuchsia":     tcell.ColorFuchsia,
		"Aqua":        tcell.ColorAqua,
		"White":       tcell.ColorWhite,
		"16":          tcell.Color16,
		"17":          tcell.Color17,
		"18":          tcell.Color18,
		"19":          tcell.Color19,
		"20":          tcell.Color20,
		"21":          tcell.Color21,
		"22":          tcell.Color22,
		"23":          tcell.Color23,
		"24":          tcell.Color24,
		"25":          tcell.Color25,
		"26":          tcell.Color26,
		"27":          tcell.Color27,
		"28":          tcell.Color28,
		"29":          tcell.Color29,
		"30":          tcell.Color30,
		"31":          tcell.Color31,
		"32":          tcell.Color32,
		"33":          tcell.Color33,
		"34":          tcell.Color34,
		"35":          tcell.Color35,
		"36":          tcell.Color36,
		"37":          tcell.Color37,
		"38":          tcell.Color38,
		"39":          tcell.Color39,
		"40":          tcell.Color40,
		"41":          tcell.Color41,
		"42":          tcell.Color42,
		"43":          tcell.Color43,
		"44":          tcell.Color44,
		"45":          tcell.Color45,
		"46":          tcell.Color46,
		"47":          tcell.Color47,
		"48":          tcell.Color48,
		"49":          tcell.Color49,
		"50":          tcell.Color50,
		"51":          tcell.Color51,
		"52":          tcell.Color52,
		"53":          tcell.Color53,
		"54":          tcell.Color54,
		"55":          tcell.Color55,
		"56":          tcell.Color56,
		"57":          tcell.Color57,
		"58":          tcell.Color58,
		"59":          tcell.Color59,
		"60":          tcell.Color60,
		"61":          tcell.Color61,
		"62":          tcell.Color62,
		"63":          tcell.Color63,
		"64":          tcell.Color64,
		"65":          tcell.Color65,
		"66":          tcell.Color66,
		"67":          tcell.Color67,
		"68":          tcell.Color68,
		"69":          tcell.Color69,
		"70":          tcell.Color70,
		"71":          tcell.Color71,
		"72":          tcell.Color72,
		"73":          tcell.Color73,
		"74":          tcell.Color74,
		"75":          tcell.Color75,
		"76":          tcell.Color76,
		"77":          tcell.Color77,
		"78":          tcell.Color78,
		"79":          tcell.Color79,
		"80":          tcell.Color80,
		"81":          tcell.Color81,
		"82":          tcell.Color82,
		"83":          tcell.Color83,
		"84":          tcell.Color84,
		"85":          tcell.Color85,
		"86":          tcell.Color86,
		"87":          tcell.Color87,
		"88":          tcell.Color88,
		"89":          tcell.Color89,
		"90":          tcell.Color90,
		"91":          tcell.Color91,
		"92":          tcell.Color92,
		"93":          tcell.Color93,
		"94":          tcell.Color94,
		"95":          tcell.Color95,
		"96":          tcell.Color96,
		"97":          tcell.Color97,
		"98":          tcell.Color98,
		"99":          tcell.Color99,
		"100":         tcell.Color100,
		"101":         tcell.Color101,
		"102":         tcell.Color102,
		"103":         tcell.Color103,
		"104":         tcell.Color104,
		"105":         tcell.Color105,
		"106":         tcell.Color106,
		"107":         tcell.Color107,
		"108":         tcell.Color108,
		"109":         tcell.Color109,
		"110":         tcell.Color110,
		"111":         tcell.Color111,
		"112":         tcell.Color112,
		"113":         tcell.Color113,
		"114":         tcell.Color114,
		"115":         tcell.Color115,
		"116":         tcell.Color116,
		"117":         tcell.Color117,
		"118":         tcell.Color118,
		"119":         tcell.Color119,
		"120":         tcell.Color120,
		"121":         tcell.Color121,
		"122":         tcell.Color122,
		"123":         tcell.Color123,
		"124":         tcell.Color124,
		"125":         tcell.Color125,
		"126":         tcell.Color126,
		"127":         tcell.Color127,
		"128":         tcell.Color128,
		"129":         tcell.Color129,
		"130":         tcell.Color130,
		"131":         tcell.Color131,
		"132":         tcell.Color132,
		"133":         tcell.Color133,
		"134":         tcell.Color134,
		"135":         tcell.Color135,
		"136":         tcell.Color136,
		"137":         tcell.Color137,
		"138":         tcell.Color138,
		"139":         tcell.Color139,
		"140":         tcell.Color140,
		"141":         tcell.Color141,
		"142":         tcell.Color142,
		"143":         tcell.Color143,
		"144":         tcell.Color144,
		"145":         tcell.Color145,
		"146":         tcell.Color146,
		"147":         tcell.Color147,
		"148":         tcell.Color148,
		"149":         tcell.Color149,
		"150":         tcell.Color150,
		"151":         tcell.Color151,
		"152":         tcell.Color152,
		"153":         tcell.Color153,
		"154":         tcell.Color154,
		"155":         tcell.Color155,
		"156":         tcell.Color156,
		"157":         tcell.Color157,
		"158":         tcell.Color158,
		"159":         tcell.Color159,
		"160":         tcell.Color160,
		"161":         tcell.Color161,
		"162":         tcell.Color162,
		"163":         tcell.Color163,
		"164":         tcell.Color164,
		"165":         tcell.Color165,
		"166":         tcell.Color166,
		"167":         tcell.Color167,
		"168":         tcell.Color168,
		"169":         tcell.Color169,
		"170":         tcell.Color170,
		"171":         tcell.Color171,
		"172":         tcell.Color172,
		"173":         tcell.Color173,
		"174":         tcell.Color174,
		"175":         tcell.Color175,
		"176":         tcell.Color176,
		"177":         tcell.Color177,
		"178":         tcell.Color178,
		"179":         tcell.Color179,
		"180":         tcell.Color180,
		"181":         tcell.Color181,
		"182":         tcell.Color182,
		"183":         tcell.Color183,
		"184":         tcell.Color184,
		"185":         tcell.Color185,
		"186":         tcell.Color186,
		"187":         tcell.Color187,
		"188":         tcell.Color188,
		"189":         tcell.Color189,
		"190":         tcell.Color190,
		"191":         tcell.Color191,
		"192":         tcell.Color192,
		"193":         tcell.Color193,
		"194":         tcell.Color194,
		"195":         tcell.Color195,
		"196":         tcell.Color196,
		"197":         tcell.Color197,
		"198":         tcell.Color198,
		"199":         tcell.Color199,
		"200":         tcell.Color200,
		"201":         tcell.Color201,
		"202":         tcell.Color202,
		"203":         tcell.Color203,
		"204":         tcell.Color204,
		"205":         tcell.Color205,
		"206":         tcell.Color206,
		"207":         tcell.Color207,
		"208":         tcell.Color208,
		"209":         tcell.Color209,
		"210":         tcell.Color210,
		"211":         tcell.Color211,
		"212":         tcell.Color212,
		"213":         tcell.Color213,
		"214":         tcell.Color214,
		"215":         tcell.Color215,
		"216":         tcell.Color216,
		"217":         tcell.Color217,
		"218":         tcell.Color218,
		"219":         tcell.Color219,
		"220":         tcell.Color220,
		"221":         tcell.Color221,
		"222":         tcell.Color222,
		"223":         tcell.Color223,
		"224":         tcell.Color224,
		"225":         tcell.Color225,
		"226":         tcell.Color226,
		"227":         tcell.Color227,
		"228":         tcell.Color228,
		"229":         tcell.Color229,
		"230":         tcell.Color230,
		"231":         tcell.Color231,
		"232":         tcell.Color232,
		"233":         tcell.Color233,
		"234":         tcell.Color234,
		"235":         tcell.Color235,
		"236":         tcell.Color236,
		"237":         tcell.Color237,
		"238":         tcell.Color238,
		"239":         tcell.Color239,
		"240":         tcell.Color240,
		"241":         tcell.Color241,
		"242":         tcell.Color242,
		"243":         tcell.Color243,
		"244":         tcell.Color244,
		"245":         tcell.Color245,
		"246":         tcell.Color246,
		"247":         tcell.Color247,
		"248":         tcell.Color248,
		"249":         tcell.Color249,
		"250":         tcell.Color250,
		"251":         tcell.Color251,
		"252":         tcell.Color252,
		"253":         tcell.Color253,
		"254":         tcell.Color254,
		"255":         tcell.Color255,
		"AliceBlue":   tcell.ColorAliceBlue,
		"AquaMarine":  tcell.ColorAquaMarine,
		"Azure":       tcell.ColorAzure,
		"Beige":       tcell.ColorBeige,
		"Bisque":      tcell.ColorBisque,
		"BlueViolet":  tcell.ColorBlueViolet,
		"Brown":       tcell.ColorBrown,
		"BurlyWood":   tcell.ColorBurlyWood,
		"CadetBlue":   tcell.ColorCadetBlue,
		"Chartreuse":  tcell.ColorChartreuse,
		"Chocolate":   tcell.ColorChocolate,
		"Coral":       tcell.ColorCoral,
		"Cornsilk":    tcell.ColorCornsilk,
		"Crimson":     tcell.ColorCrimson,
		"DarkBlue":    tcell.ColorDarkBlue,
		"DarkCyan":    tcell.ColorDarkCyan,
		"DarkGray":    tcell.ColorDarkGray,
		"DarkGreen":   tcell.ColorDarkGreen,
		"DarkKhaki":   tcell.ColorDarkKhaki,
		"DarkMagenta": tcell.ColorDarkMagenta,
		"DarkOrange":  tcell.ColorDarkOrange,
		"DarkOrchid":  tcell.ColorDarkOrchid,
		"DarkRed":     tcell.ColorDarkRed,
		"DarkSalmon":  tcell.ColorDarkSalmon,
		"DarkViolet":  tcell.ColorDarkViolet,
		"DeepPink":    tcell.ColorDeepPink,
		"Pink":        tcell.ColorPink,
		"DeepSkyBlue": tcell.ColorDeepSkyBlue,
		"DimGray":     tcell.ColorDimGray,
		"DodgerBlue":  tcell.ColorDodgerBlue,
		"FireBrick":   tcell.ColorFireBrick,
		"FloralWhite": tcell.ColorFloralWhite,
		"ForestGreen": tcell.ColorForestGreen,
		"Gainsboro":   tcell.ColorGainsboro,
		"GhostWhite":  tcell.ColorGhostWhite,
		"Gold":        tcell.ColorGold,
		"Goldenrod":   tcell.ColorGoldenrod,
		"GreenYellow": tcell.ColorGreenYellow,
		"Honeydew":    tcell.ColorHoneydew,
		"HotPink":     tcell.ColorHotPink,
		"IndianRed":   tcell.ColorIndianRed,
		"Indigo":      tcell.ColorIndigo,
		"Ivory":       tcell.ColorIvory,
		"Khaki":       tcell.ColorKhaki,
		"Lavender":    tcell.ColorLavender,
		"LawnGreen":   tcell.ColorLawnGreen,
		"LightBlue":   tcell.ColorLightBlue,
		"LightCoral":  tcell.ColorLightCoral,
		"LightCyan":   tcell.ColorLightCyan,
		"LightGray":   tcell.ColorLightGray,
		"LightGreen":  tcell.ColorLightGreen,
		"LightPink":   tcell.ColorLightPink,
		"LightSalmon": tcell.ColorLightSalmon,
		"Linen":       tcell.ColorLinen,
		"MediumBlue":  tcell.ColorMediumBlue,
		"MintCream":   tcell.ColorMintCream,
		"OldLace":     tcell.ColorOldLace,
		"Orchid":      tcell.ColorOrchid,
		"PaleGreen":   tcell.ColorPaleGreen,
		"PapayaWhip":  tcell.ColorPapayaWhip,
		"PeachPuff":   tcell.ColorPeachPuff,
		"Peru":        tcell.ColorPeru,
		"Plum":        tcell.ColorPlum,
		"PowderBlue":  tcell.ColorPowderBlue,
		"RosyBrown":   tcell.ColorRosyBrown,
		"SaddleBrown": tcell.ColorSaddleBrown,
		"Salmon":      tcell.ColorSalmon,
		"SandyBrown":  tcell.ColorSandyBrown,
		"SeaGreen":    tcell.ColorSeaGreen,
		"Seashell":    tcell.ColorSeashell,
		"Sienna":      tcell.ColorSienna,
		"Skyblue":     tcell.ColorSkyblue,
		"SlateBlue":   tcell.ColorSlateBlue,
		"SlateGray":   tcell.ColorSlateGray,
		"Snow":        tcell.ColorSnow,
		"SpringGreen": tcell.ColorSpringGreen,
		"SteelBlue":   tcell.ColorSteelBlue,
		"Tan":         tcell.ColorTan,
		"Thistle":     tcell.ColorThistle,
		"Tomato":      tcell.ColorTomato,
		"Turquoise":   tcell.ColorTurquoise,
		"Violet":      tcell.ColorViolet,
		"Wheat":       tcell.ColorWheat,
		"WhiteSmoke":  tcell.ColorWhiteSmoke,
		"YellowGreen": tcell.ColorYellowGreen,
	}
)

type Color struct {
	Fg     string `mapstructure:"foreground"`
	Bg     string `mapstructure:"background"`
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

	Null Color
}

func (c Color) Foreground() tcell.Color {
	if strings.HasPrefix(c.Fg, "#") && len(c.Fg) == 7 {
		return tcell.GetColor(c.Fg)
	} else if val, ok := DColors[c.Fg]; ok {
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
	} else if val, ok := DColors[c.Bg]; ok {
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
		if _, ok := DColors[s]; ok {
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
			Fg:     "Pink",
			Bold:   false,
			Italic: false,
		},
		Album: Color{
			Fg:     "Green",
			Bold:   false,
			Italic: false,
		},
		Track: Color{
			Fg:     "Blue",
			Bold:   false,
			Italic: false,
		},
		Timestamp: Color{
			Fg:     "Red",
			Bold:   false,
			Italic: true,
		},
		Genre: Color{
			Fg:     "DarkCyan",
			Bold:   true,
			Italic: false,
		},
		Folder: Color{
			Fg:     "Yellow",
			Bold:   true,
			Italic: false,
		},
		PBarArtist: Color{
			Fg:     "Blue",
			Bold:   true,
			Italic: false,
		},
		PBarTrack: Color{
			Fg:     "Green",
			Bold:   true,
			Italic: true,
		},
		PlaylistNav: Color{
			Fg:     "Coral",
			Bold:   false,
			Italic: false,
		},
		Nav: Color{
			Fg:     "PapayaWhip",
			Bold:   false,
			Italic: false,
		},
		ContextMenu: Color{
			Fg:     "Turquoise",
			Bold:   true,
			Italic: false,
		},
		Null: Color{
			Fg:     "White",
			Bold:   true,
			Italic: false,
		},
	}
}
