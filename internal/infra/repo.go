package infra

import "github.com/lusingander/birdfeeder/internal/domain"

type PostRepositoryImpl struct{}

func (PostRepositoryImpl) ReadAllPosts() ([]*domain.Post, error) {
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
	return &domain.Post{
		Number:   p.Number,
		Title:    p.Name,
		Body:     p.BodyMd,
		Wip:      p.Wip,
		Category: p.Category,
		Tags:     p.Tags,
	}
}
