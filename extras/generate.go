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
		_s += fmt.Sprintf("%c, ", k)
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
	for k, v := range config.DefaultMappings {
		s += fmt.Sprintf("    %s:\n", k)
		for x, y := range v {
			if x.R != 0 {
				s += fmt.Sprintf("        %s: \"%s\"\n", y, string(x.R))
			} else {
				s += fmt.Sprintf("        %s: \"%s\"\n", y, _m[x.K])
			}
		}
	}

	_s := "|||\n|--|--|\n"
	for k := 0; k < len(_ma)-1; k++ {
		if k != len(_ma)-1 {
			_s += "|" + _ma[k]
			k++
			_s += "|" + _ma[k] + "|\n"
		}
	}

	// Write
	write("***Auto generated*** (If you find any bugs please open an issue)\n")

	write("# Default Configuration\n```yml\n" + string(c) + "mappings:\n" + s + "```\n")

	write("# Available Keys\n" + _s)

	fmt.Println("GENERATED CONFIG.MD")
}
