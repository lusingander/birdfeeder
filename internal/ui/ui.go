package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/ui/tree"
	"github.com/lusingander/birdfeeder/internal/util"
)

type model struct {
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
		m.tree = tree.New(msg)
		return m, nil
	case errorMsg:
		m.err = msg
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	buf := util.NewBufferWrapper()
	m.internalView(buf)
	return buf.String()
}

func (m model) internalView(buf *util.BufferWrapper) {
	buf.Writeln("- BIRDFEEDER -")
	if m.err != nil {
		buf.Writeln(m.err.Error())
		return
	}
	buf.Writeln(m.tree.View())
}

// Start UI
func Start() error {
	initRepositories()

	p := tea.NewProgram(model{})
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	return p.Start()
}
