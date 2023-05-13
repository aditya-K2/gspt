package config

import (
	"flag"
)

var (
	Flags = &Flag{}
)

type Flag struct {
	ConfigPath     string
	HideImage      bool
	RoundedCorners bool
	UseIcons       bool
}

func parseFlags() {
	flag.StringVar(&Flags.ConfigPath, "c", UserConfigPath,
		"Specify The Directory to check for config.yml file.")
	flag.BoolVar(&Flags.HideImage, "hide-image", Config.HideImage,
		"Do not display the cover art image.")
	flag.BoolVar(&Flags.RoundedCorners, "rounded-corners", Config.RoundedCorners,
		"Enable Rounded Corners")
	flag.BoolVar(&Flags.UseIcons, "icons", Config.UseIcons,
		"Use Icons")
	flag.Parse()
}
