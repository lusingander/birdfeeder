package main

import (
	"encoding/json"
	"fmt"
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

func getCacheDirPath() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(cache, application)
	return path, nil
}

func run(args []string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	fmt.Println(cfg)
	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}
