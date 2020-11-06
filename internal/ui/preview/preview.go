package preview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/ui/tree"
	"github.com/lusingander/birdfeeder/internal/util"
)

type Model struct {
	post *domain.Post
	base tree.Model

	Close bool
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type InitMsg struct {
	tree.Model
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.Close = true
			return m, nil
		}
	case InitMsg:
		m.base = msg.Model
		m.post = msg.Model.CurrentPost()
	}
	return m, nil
}

func (m Model) View() string {
	return ""
}

func (m Model) ViewBreadcrumb(buf *util.BufferWrapper) {
	m.base.ViewBreadcrumb(buf)
}
