package preview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/lusingander/birdfeeder/internal/domain"
	"github.com/lusingander/birdfeeder/internal/ui/tree"
	"github.com/lusingander/birdfeeder/internal/ui/viewport"
	"github.com/lusingander/birdfeeder/internal/util"
)

type Model struct {
	post *domain.Post
	base tree.Model

	viewport viewport.Model

	Close bool
}

func New() Model {
	return Model{
		viewport: viewport.Model{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type InitMsg tree.Model

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.Close = true
			return m, nil
		}
	case InitMsg:
		m.base = tree.Model(msg)
		m.post = m.base.CurrentPost()
		m.viewport.GotoTop()
		m.viewport.SetContent(m.viewMarkdown())
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3 // header + footer
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

func (m Model) viewMarkdown() string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(m.viewport.Width),
	)
	if err != nil {
		return err.Error()
	}
	md, err := r.Render(m.post.Body)
	if err != nil {
		return err.Error()
	}
	return md
}

func (m Model) ViewBreadcrumb(buf *util.BufferWrapper) {
	m.base.ViewBreadcrumb(buf)
}
