package config

import (
	"flag"
)

func parseFlags() {
	flag.StringVar(&ConfigPath, "c", ConfigPath,
		"Specify The Directory to check for config.yml file.")
	flag.BoolVar(&Config.HideImage, "hide-image", Config.HideImage,
		"Do not display the cover art image.")
	flag.Parse()
}
