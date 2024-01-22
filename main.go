package main

import (
	"github.com/aditya-K2/gspt/config"
	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/gspt/ui"
)

func main() {
	if err := config.ReadConfig(); err != nil {
		panic(err)
	}
	if err := spt.InitClient(); err != nil {
		panic(err)
	}
	if err := ui.NewApplication().Run(); err != nil {
		panic(err)
	}
}
