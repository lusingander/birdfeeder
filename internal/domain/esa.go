package domain

import "time"

type PostRepository interface {
	ReadAllPosts() ([]*Post, error)
}

type Post struct {
	Number     int
	Title      string
	Body       string
	Wip        bool
	Categories []string
	Tags       []string
	UpdatedAt  time.Time
}
