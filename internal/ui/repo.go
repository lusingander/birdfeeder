package ui

import (
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/infra"
)

var postRepository domain.PostRepository

func initRepositories() {
	postRepository = infra.PostRepositoryImpl{}
}
