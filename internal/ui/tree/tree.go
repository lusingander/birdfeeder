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

	cursor    int
	current   *node
	histories []*history
}

func New() Model {
	return Model{
		viewport: viewport.Model{},
	}
}

type history struct {
	cursor int
	*node
}

func (m *Model) createHistory() {
	h := &history{m.cursor, m.current}
	m.histories = append(m.histories, h)
}

func (m *Model) goBackHistory() *history {
	h := m.histories[len(m.histories)-1]
	m.histories = m.histories[:len(m.histories)-1]
	return h
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
			if m.cursor < len(m.current.children)-1 {
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
		case "f":
			// should reconsider..
			m.viewport.ViewDown()
			m.cursor = m.viewport.YOffset
			m.viewport.SetContent(m.viewTree())
			return m, nil
		case "b":
			// should reconsider..
			m.viewport.ViewUp()
			m.cursor = m.viewport.YOffset
			m.viewport.SetContent(m.viewTree())
			return m, nil
		case "g":
			if m.cursor > 0 {
				m.cursor = 0
				m.viewport.GotoTop()
				m.viewport.SetContent(m.viewTree())
			}
			return m, nil
		case "G":
			if m.cursor < len(m.current.children)-1 {
				m.cursor = len(m.current.children) - 1
				m.viewport.GotoBottom()
				m.viewport.SetContent(m.viewTree())
			}
			return m, nil
		case "l":
			target := m.current.children[m.cursor]
			if len(target.children) > 0 {
				m.createHistory()
				m.current = target
				m.cursor = 0
				m.viewport.GotoTop()
				m.viewport.SetContent(m.viewTree())
			}
			return m, nil
		case "h":
			if len(m.histories) > 0 {
				h := m.goBackHistory()
				m.current = h.node
				m.cursor = h.cursor
				m.viewport.SetContent(m.viewTree())
			}
			return m, nil
		}
	case InitMsg:
		m.root = buildRoot(msg)
		m.current = m.root
		m.viewport.SetContent(m.viewTree())
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3 // header + footer
	}
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
	if m.current == nil {
		return ""
	}
	buf := util.NewBufferWrapper()
	for i, node := range m.current.children {
		if i == m.cursor {
			buf.Write("> ")
		} else {
			buf.Write("  ")
		}
		if node.post != nil {
			buf.Writeln("%s", node.name)
		} else {
			buf.Writeln("%s (%d)", node.name, len(node.children))
		}
	}
	return buf.String()
}

func (m Model) ViewBreadcrumb(buf *util.BufferWrapper) {
	buf.Write(" > POSTS")
	if len(m.histories) > 0 {
		for _, h := range m.histories[1:] {
			buf.Write(" > %s", h.name)
		}
		buf.Write(" > %s", m.current.name)
	}
}
