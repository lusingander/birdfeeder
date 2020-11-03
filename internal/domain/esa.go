package domain

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
}
