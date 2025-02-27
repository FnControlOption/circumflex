package bubble

import (
	"fmt"
	"os"

	"clx/bfavorites"
	"clx/bubble/list"
	"clx/cli"
	"clx/settings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle()

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Run(config *settings.Config) {
	cli.ClearScreen()

	m := model{list: list.New(list.NewDefaultDelegate(), config, bfavorites.New(), 0, 0)}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
