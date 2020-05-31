package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getCacheDirPath() (string, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cache, application), nil
}

func savePosts(posts []*postDetail) error {
	cacheDir, err := getCacheDirPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.Mkdir(cacheDir, os.ModePerm)
	}
	for _, p := range posts {
		err = savePost(p, cacheDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func savePost(post *postDetail, cacheDir string) error {
	name := fmt.Sprintf("%d", post.Number)
	path := filepath.Join(cacheDir, name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := json.Marshal(post)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func readPosts() ([]*postDetail, error) {
	cacheDir, err := getCacheDirPath()
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}
	posts := make([]*postDetail, 0)
	for _, f := range files {
		if shouldSkip(f) {
			continue
		}
		path := filepath.Join(cacheDir, f.Name())
		post, err := readPost(path)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func shouldSkip(f os.FileInfo) bool {
	return f.IsDir() || f.Name() == metadataFile
}

func readPost(path string) (*postDetail, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var post postDetail
	err = json.NewDecoder(f).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
