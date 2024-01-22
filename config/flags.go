package config

import (
	"flag"
)

var (
	Flags = &Flag{}
)

type Flag struct {
	ConfigPath string
	Version    bool
	Image      string
	Corners    string
	UseIcons   bool
}

func parseFlags() {
	flag.BoolVar(&Flags.Version, "v", false,
		"Do not display the cover art image.")
	flag.BoolVar(&Flags.Version, "version", false,
		"Do not display the cover art image.")
	flag.BoolVar(&Flags.UseIcons, "icons", Config.UseIcons,
		"Use Icons")
	flag.StringVar(&Flags.ConfigPath, "c", userConfigPath,
		"Specify The Directory to check for config.yml file.")
	flag.StringVar(&Flags.Image, "image", "",
		"Show or Hide Image ( 'show' | 'hidden' )")
	flag.StringVar(&Flags.Corners, "corners", "",
		"Enable or disable Rounded Corners ( 'rounded' | 'default' )")
	flag.Parse()
}
