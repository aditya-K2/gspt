package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aditya-K2/gspt/config"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v2"
)

var (
	fileName = "extras/CONFIG.md"
	f, oerr  = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
)

func write(body string) {
	if oerr != nil {
		panic(oerr)
	}
	if _, err := f.Write([]byte(body)); err != nil {
		panic(err)
	}
}

func runes() string {
	_s := ""
	for k := range config.RuneKeys {
		if k == '|' {
			_s += "\\|"
		} else {
			_s += fmt.Sprintf("%c, ", k)
		}
	}
	return strings.TrimSuffix(_s, ", ")
}

func main() {
	// Ready
	_m := map[tcell.Key]string{}
	_ma := []string{"a-z", "A-Z", "ctrl-a - ctrl-z", "0-9", runes()}
	for k, v := range config.M {
		_m[v] = k
		if !strings.HasPrefix(k, "ctrl-") {
			_ma = append(_ma, k)
		}
	}

	config.Config.CacheDir = "$XDG_CACHE_HOME"
	c, err := yaml.Marshal(config.Config)
	if err != nil {
		panic(err)
	}

	s := ""
	for view, modes := range config.DefaultMappings {
		s += fmt.Sprintf("    %s:\n", view)
		for mode, mappings := range modes {
			s += fmt.Sprintf("        %s:\n", mode)
			for key, function := range mappings {
				if key.R != 0 {
					s += fmt.Sprintf("            %s: \"%s\"\n", function, string(key.R))
				} else {
					s += fmt.Sprintf("            %s: \"%s\"\n", function, _m[key.K])
				}
			}
		}
	}

	_s := "|||\n|--|--|\n"
	for k := 0; k < len(_ma)-1; k++ {
		_s += fmt.Sprintf("| <kbd>%s</kbd> ", _ma[k])
		k++
		_s += fmt.Sprintf("| <kbd>%s</kbd> |\n", _ma[k])
	}

	// Override the file
	f.Truncate(0)

	// Write
	write("***This file is auto generated. If you find any bugs please open an issue***\n\n")

	write("***`ctrl-m` in mappings refers to the `enter` key.***\n")

	write("## Default Configuration\n```yml\n" + string(c) + "mappings:\n" + s + "```\n")

	write("## Keys available to map\n" + _s)

	fmt.Println("GENERATED CONFIG.MD")
}
