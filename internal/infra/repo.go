package infra

import (
	"strings"

	"github.com/lusingander/birdfeeder/internal/domain"
)

func NewPostRepository() domain.PostRepository

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
	}
}

func parseCategories(category string) []string {
	return strings.Split(category, "/")
}
