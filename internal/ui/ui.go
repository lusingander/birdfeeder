package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
)

type model struct {
	posts []*domain.Post
	err   error
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
		m.posts = msg
		return m, nil
	case errorMsg:
		m.err = msg
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	buf := newBufferWrapper()
	m.internalView(buf)
	return buf.String()
}

func (m model) internalView(buf *bufferWrapper) {
	buf.writeln("- BIRDFEEDER -")
	if m.err != nil {
		buf.writeln(m.err.Error())
		return
	}
	for _, post := range m.posts {
		buf.writeln("[%d] %s", post.Number, post.Title)
	}
}

// Start UI
func Start() error {
	initRepositories()

	p := tea.NewProgram(model{})
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	return p.Start()
}
