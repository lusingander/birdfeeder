package ui

import (
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/infra"
)

var postRepository domain.PostRepository
var configRepository domain.ConfigRepository
var metaRepository domain.MetaRepository

func initRepositories() {
	postRepository = infra.NewPostRepository()
	configRepository = infra.NewConfigRepository()
	metaRepository = infra.NewMetaRepository()
}
