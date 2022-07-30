package main

import (
	"os"

	"github.com/Buzz2d0/xssfinder/internal/app"
)

const version = "v0.1.2"

func main() {
	a := app.New(version)
	if err := a.Run(os.Args); err != nil {
		panic(err)
	}
}
