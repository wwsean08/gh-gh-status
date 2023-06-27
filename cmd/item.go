package cmd

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
)

type listDelegator struct {
	style *listItemsStyle
}

type listItemsStyle struct {
	// The Normal state.
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	// The selected item state.
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	// Charcters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

func NewItemDelegator() *listDelegator {
	return &listDelegator{}
}

func (d *listDelegator) Render(w io.Writer, m list.Model, index int, item list.Item) {

}

func (d *listDelegator) Height() int {
	return 2
}

func (d *listDelegator) Spacing() int {
	return 1
}

func (d *listDelegator) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
