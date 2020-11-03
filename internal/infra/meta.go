package infra

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	metadataFile = "meta.json"
)

type metadata struct {
	LastUpdate time.Time `json:"last_update"`
}

func newMetadata() *metadata {
	return &metadata{
		LastUpdate: time.Now(),
	}
}

func saveMetadata() error {
	cacheDir, err := getCacheDirPath()
	if err != nil {
		return err
	}
	path := filepath.Join(cacheDir, metadataFile)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	meta := newMetadata()
	bytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
