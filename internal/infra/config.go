package infra

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/lusingander/birdfeeder/internal/domain"
)

func ReadConfig() (*domain.Config, error) {
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

	var cfg domain.Config
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}
