package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	application = "birdfeeder"
	configFile  = "config.json"
)

type config struct {
	Team  string
	Token string
}

func readConfig() (*config, error) {
	conf, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(conf, application, configFile)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}

func initAllPosts(cfg *config) error {
	posts, err := fetchAllPosts(cfg)
	if err != nil {
		return err
	}
	return savePosts(posts)
}

func run(args []string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	_ = cfg

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
