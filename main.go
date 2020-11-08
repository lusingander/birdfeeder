package main

import (
	"log"
	"os"

	"github.com/lusingander/birdfeeder/internal/ui"
)

func run(args []string) error {
	return ui.Start()
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
