package tree

import (
	"sort"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/util"
)

type node struct {
	name     string
	post     *domain.Post
	children []*node
}

func buildRoot(posts []*domain.Post) *node {
	root := &node{
		children: []*node{},
	}
	for _, post := range posts {
		add(post, post.Categories, root)
	}
	sorter := func(i, j int) bool {
		return root.children[i].name < root.children[j].name
	}
	sort.Slice(root.children, sorter)
	return root
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

type Model struct {
	posts []*domain.Post
	root  *node

	viewport viewport.Model

	cursor int
}

func New() Model {
	return Model{
		viewport: viewport.Model{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type InitMsg []*domain.Post

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.cursor < len(m.root.children)-1 {
				m.cursor++
				if m.viewport.YOffset+m.viewport.Height < m.cursor+1 {
					m.viewport.LineDown(1)
				}
				m.viewport.SetContent(m.viewTree())
			}
			return m, nil
		case "k":
			if m.cursor > 0 {
				m.cursor--
				if m.viewport.YOffset > m.cursor {
					m.viewport.LineUp(1)
				}
				m.viewport.SetContent(m.viewTree())
			}
			return m, nil
		}
	case InitMsg:
		m.root = buildRoot(msg)
		m.viewport.SetContent(m.viewTree())
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2 // header + footer
	}
	m.viewport, _ = viewport.Update(msg, m.viewport)
	return m, nil
}

func (m Model) View() string {
	buf := util.NewBufferWrapper()
	m.internalView(buf)
	return buf.String()
}

func (m Model) internalView(buf *util.BufferWrapper) {
	buf.Writeln(viewport.View(m.viewport))
}

func (m Model) viewTree() string {
	if m.root == nil {
		return ""
	}
	buf := util.NewBufferWrapper()
	for i, node := range m.root.children {
		if i == m.cursor {
			buf.Write("> ")
		} else {
			buf.Write("  ")
		}
		buf.Writeln("%s (%d)", node.name, len(node.children))
	}
	return buf.String()
}
