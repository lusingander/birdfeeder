package tree

import (
	"sort"
	"time"

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

func (n *node) postCount() int {
	if n.post != nil {
		return 1
	}
	c := 0
	for _, child := range n.children {
		c += child.postCount()
	}
	return c
}

func (n *node) updatedAt() time.Time {
	if n.post != nil {
		return n.post.UpdatedAt
	}
	t := time.Unix(0, 0)
	for _, child := range n.children {
		u := child.updatedAt()
		if u.After(t) {
			t = u
		}
	}
	return t
}

func buildRoot(posts []*domain.Post) *node {
	root := &node{
		children: []*node{},
	}
	for _, post := range posts {
		add(post, post.Categories, root)
	}
	root.sortNodesRecursive(byNameAsc)
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

type sortKey int

const (
	byNameAsc sortKey = iota
	byNameDesc
	byPostsCountAsc
	byPostsCountDesc
	byUpdatedAtAsc
	byUpdatedAtDesc
)

func (k sortKey) next() sortKey {
	return (k + 1) % 6
}

func (k sortKey) str() string {
	switch k {
	case byNameAsc:
		return "Name ↓"
	case byNameDesc:
		return "Name ↑"
	case byPostsCountAsc:
		return "Posts ↓"
	case byPostsCountDesc:
		return "Posts ↑"
	case byUpdatedAtAsc:
		return "Updated ↓"
	case byUpdatedAtDesc:
		return "Updated ↑"
	}
	return ""
}

func (n *node) sortNodesRecursive(key sortKey) {
	createSorter := func(nodes []*node) func(int, int) bool {
		switch key {
		case byNameAsc:
			return func(i, j int) bool { return nodes[i].name < nodes[j].name }
		case byNameDesc:
			return func(i, j int) bool { return nodes[i].name > nodes[j].name }
		case byPostsCountAsc:
			return func(i, j int) bool { return nodes[i].postCount() < nodes[j].postCount() }
		case byPostsCountDesc:
			return func(i, j int) bool { return nodes[i].postCount() > nodes[j].postCount() }
		case byUpdatedAtAsc:
			return func(i, j int) bool { return nodes[i].updatedAt().Before(nodes[j].updatedAt()) }
		case byUpdatedAtDesc:
			return func(i, j int) bool { return nodes[i].updatedAt().After(nodes[j].updatedAt()) }
		default:
			panic("Invalid key type")
		}
	}
	sort.Slice(n.children, createSorter(n.children))
	for _, child := range n.children {
		child.sortNodesRecursive(key)
	}
}

type Model struct {
	posts         []*domain.Post
	root          *node
	lastUpdateStr string

	viewport viewport.Model

	cursor    int
	current   *node
	histories []*history
	sortKey

	OpenPost bool
}

func New() Model {
	return Model{
		viewport: viewport.Model{},
		sortKey:  byNameAsc,
	}
}

func (m Model) CurrentPost() *domain.Post {
	return m.current.post
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

type InitMsg struct {
	Posts []*domain.Post
	Meta  *domain.Meta
}
type ClosePreview struct{}

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
			if target.post != nil {
				m.createHistory()
				m.current = target
				m.OpenPost = true
				return m, nil
			}
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
		case "H":
			m.cursor = m.viewport.YOffset
			m.viewport.SetContent(m.viewTree())
			return m, nil
		case "M":
			m.cursor = m.viewport.YOffset + (m.viewport.Height / 2)
			m.viewport.SetContent(m.viewTree())
			return m, nil
		case "L":
			m.cursor = m.viewport.YOffset + m.viewport.Height - 1
			m.viewport.SetContent(m.viewTree())
			return m, nil
		case "s":
			m.sortKey = m.sortKey.next()
			m.root.sortNodesRecursive(m.sortKey)
			m.viewport.SetContent(m.viewTree())
			return m, nil
		}
	case InitMsg:
		m.root = buildRoot(msg.Posts)
		m.current = m.root
		m.lastUpdateStr = msg.Meta.FormattedLastUpdate()
		m.viewport.SetContent(m.viewTree())
	case ClosePreview:
		h := m.goBackHistory()
		m.current = h.node
		m.viewport.SetContent(m.viewTree())
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4 // header + footer
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
			buf.Writeln("%s (%d)", node.name, node.postCount())
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

func (m Model) ViewFooter(buf *util.BufferWrapper) {
	lastUpdate := "last update: " + m.lastUpdateStr
	sortKey := "sort key: " + m.sortKey.str()
	buf.Writeln(util.Faint("%s, %s"), lastUpdate, sortKey)
	buf.Write(util.Faint("j/k: move cursor, h/l: change directory, Ctrl+C: quit"))
}
