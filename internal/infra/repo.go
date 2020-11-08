package infra

import (
	"strings"

	"github.com/lusingander/birdfeeder/internal/domain"
)

func NewConfigRepository() domain.ConfigRepository {
	return &configRepository{}
}

type configRepository struct{}

func (configRepository) ReadConfig() (*domain.Config, error) {
	return readConfig()
}

func NewPostRepository() domain.PostRepository {
	return &postRepository{}
}

type postRepository struct{}

func (postRepository) ReadAllPosts() ([]*domain.Post, error) {
	posts, err := readPosts()
	if err != nil {
		return nil, err
	}
	ret := make([]*domain.Post, len(posts))
	for i, p := range posts {
		ret[i] = toPost(p)
	}
	return ret, nil
}

func toPost(p *postDetail) *domain.Post {
	categories := parseCategories(p.Category)
	return &domain.Post{
		Number:     p.Number,
		Title:      p.Name,
		Body:       p.BodyMd,
		Wip:        p.Wip,
		Categories: categories,
		Tags:       p.Tags,
		UpdatedAt:  p.UpdatedAt,
	}
}

func parseCategories(category string) []string {
	return strings.Split(category, "/")
}

func NewMetaRepository() domain.MetaRepository {
	return &metaRepository{}
}

type metaRepository struct{}

func (metaRepository) ReadMeta() (*domain.Meta, error) {
	meta, err := readMetadata()
	if err != nil {
		return nil, err
	}
	return &domain.Meta{
		LastUpdate: meta.LastUpdate,
	}, nil
}
