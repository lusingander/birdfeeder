package tree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/util"
)

type node struct {
	name     string
	post     *domain.Post
	children []*node
}

type Model struct {
	posts []*domain.Post
	root  *node
}

func New(posts []*domain.Post) Model {
	root := &node{
		children: []*node{},
	}
	for _, post := range posts {
		add(post, post.Categories, root)
	}
	return Model{
		posts: posts,
		root:  root,
	}
}

func add(post *domain.Post, categories []string, target *node) {
	if len(categories) == 0 || categories[0] == "" {
		newNode := &node{
			name:     post.Title,
			post:     post,
			children: []*node{},
		}
		target.children = append(target.children, newNode)
		return
	}
	for _, child := range target.children {
		if child.name == categories[0] {
			add(post, categories[1:], child)
			return
		}
	}
	newNode := &node{
		name:     categories[0],
		children: []*node{},
	}
	target.children = append(target.children, newNode)
	add(post, categories[1:], newNode)
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	buf := util.NewBufferWrapper()
	m.internalView(buf)
	return buf.String()
}

func (m Model) internalView(buf *util.BufferWrapper) {
	if m.root == nil {
		return
	}
	for _, node := range m.root.children {
		buf.Writeln("%s (%d)", node.name, len(node.children))
	}
}
