package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/ui/preview"
	"github.com/lusingander/birdfeeder/internal/ui/tree"
	"github.com/lusingander/birdfeeder/internal/util"
)

type model struct {
	width, height int
	currentState  state

	tree    tree.Model
	preview preview.Model
	err     error
}

type state int

const (
	stateTree state = iota
	statePreview
)

func (model) Init() tea.Cmd {
	return readPosts
}

func readPosts() tea.Msg {
	posts, err := postRepository.ReadAllPosts()
	if err != nil {
		return errorMsg(err)
	}
	return initPostsMsg(posts)
}

type errorMsg error
type initPostsMsg []*domain.Post

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case initPostsMsg:
		return m.Update(tree.InitMsg(msg))
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tree, _ = m.tree.Update(msg)
		m.preview, _ = m.preview.Update(msg)
		return m, nil
	case errorMsg:
		m.err = msg
		return m, nil
	}

	switch m.currentState {
	case stateTree:
		m.tree, _ = m.tree.Update(msg)
		if m.tree.OpenPost {
			m.tree.OpenPost = false
			m.currentState = statePreview
			return m.Update(preview.InitMsg(m.tree))
		}
	case statePreview:
		m.preview, _ = m.preview.Update(msg)
		if m.preview.Close {
			m.preview.Close = false
			m.currentState = stateTree
			return m.Update(tree.ClosePreview(struct{}{}))
		}
	}

	return m, nil
}

func (m model) View() string {
	buf := util.NewBufferWrapper()
	m.internalView(buf)
	return buf.String()
}

func (m model) internalView(buf *util.BufferWrapper) {
	m.viewHeader(buf)
	if m.err != nil {
		buf.Writeln(m.err.Error())
		return
	}
	switch m.currentState {
	case stateTree:
		buf.Write(m.tree.View())
	case statePreview:
		buf.Write(m.preview.View())
	}
	m.viewFooter(buf)
}

func (m model) viewHeader(buf *util.BufferWrapper) {
	m.viewBreadcrumb(buf)
	m.viewHorizontalSeparator(buf)
}

func (m model) viewBreadcrumb(buf *util.BufferWrapper) {
	buf.Write("BIRDFEEDER")
	switch m.currentState {
	case stateTree:
		m.tree.ViewBreadcrumb(buf)
	case statePreview:
		m.preview.ViewBreadcrumb(buf)
	}
	buf.Writeln("")
}

func (m model) viewHorizontalSeparator(buf *util.BufferWrapper) {
	buf.Writeln(strings.Repeat("-", m.width))
}

func (model) viewFooter(buf *util.BufferWrapper) {
	buf.Write("")
}

// Start UI
func Start() error {
	initRepositories()

	m := model{
		currentState: stateTree,
	}
	m.tree = tree.New()
	m.preview = preview.New()

	p := tea.NewProgram(m)
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	return p.Start()
}
