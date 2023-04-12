package main

import (
	"github.com/aditya-K2/gspt/spt"
	"github.com/aditya-K2/gspt/ui"
)

func main() {
	if err := spt.InitClient(); err != nil {
		panic(err)
	}
	if err := ui.NewApplication().App.Run(); err != nil {
		panic(err)
	}
}
