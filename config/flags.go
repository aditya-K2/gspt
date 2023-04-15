package config

import (
	"flag"
)

func parseFlags() {
	flag.StringVar(&ConfigPath, "c", ConfigPath,
		"Specify The Directory where to check for config.yml file.")
	flag.Parse()
}
