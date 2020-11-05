package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/ui/tree"
	"github.com/lusingander/birdfeeder/internal/util"
)

type model struct {
	width, height int

	tree tree.Model
	err  error
}

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
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case initPostsMsg:
		return m.Update(tree.InitMsg(msg))
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case errorMsg:
		m.err = msg
		return m, nil
	}
	m.tree, _ = m.tree.Update(msg)
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
	buf.Write(m.tree.View())
	m.viewFooter(buf)
}

func (m model) viewHeader(buf *util.BufferWrapper) {
	m.viewBreadcrumb(buf)
	m.viewHorizontalSeparator(buf)
}

func (m model) viewBreadcrumb(buf *util.BufferWrapper) {
	buf.Write(" BIRDFEEDER")
	m.tree.ViewBreadcrumb(buf)
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

	m := model{}
	m.tree = tree.New()

	p := tea.NewProgram(m)
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	return p.Start()
}
