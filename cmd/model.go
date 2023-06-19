package cmd

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wwsean08/gh-gh-status/status"
	"time"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)

type item struct {
	component, status string
}

func (i item) Title() string       { return i.component }
func (i item) Description() string { return i.status }
func (i item) FilterValue() string { return i.component }

type statusMsg struct {
	status    *status.SystemStatus
	err       error
	timestamp *time.Time
}

type model struct {
	systemStatus *status.SystemStatus
	components   list.Model
	outageStatus list.Model
	lastUpdated  *time.Time
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case statusMsg:
		status := msg.status
		m.lastUpdated = msg.timestamp
		if status != nil {
			m.systemStatus = status

			items := make([]list.Item, 0, len(status.Components)-1)
			for _, component := range status.Components {
				if component.ID == IGNORE_GHSTATUS_COMPONENTID {
					continue
				}
				items = append(items, item{
					component: component.Component,
					status:    string(component.Status),
				})
			}
			m.components.SetItems(items)
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.components.SetSize(msg.Width-h, msg.Height-v)
		m.components.Title = "Component Status"
	}

	m.components, cmd = m.components.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.components.View())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
