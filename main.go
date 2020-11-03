package main

import (
	"log"
	"os"

	"github.com/lusingander/birdfeeder/internal/infra"
	"github.com/lusingander/birdfeeder/internal/ui"
)

func run(args []string) error {
	cfg, err := infra.ReadConfig()
	if err != nil {
		return err
	}
	_ = cfg
	return ui.Start()
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
