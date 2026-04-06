package tui

import (
	"AlgoScope-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) AlgorithmMenuView() string {
	sel := m.Menu[m.Cursor]
	
	title := lipgloss.NewStyle().
		Foreground(ui.Lavender).
		Bold(true).
		MarginBottom(2).
		Render("ALGORITHM: " + sel.Name)

	options := []string{"󰓡 START VISUALIZER", "󰈙 READ DOCUMENTATION"}
	var renderedOpts []string
	for i, opt := range options {
		style := lipgloss.NewStyle().Padding(1, 2).Width(30).Align(lipgloss.Center).Margin(1)
		if i == m.SubCursor {
			style = style.Foreground(ui.Bg).Background(ui.Teal).Bold(true)
		} else {
			style = style.Foreground(ui.Teal).Border(lipgloss.NormalBorder()).BorderForeground(ui.Teal)
		}
		renderedOpts = append(renderedOpts, style.Render(opt))
	}

	menu := lipgloss.JoinHorizontal(lipgloss.Top, renderedOpts...)
	
	help := lipgloss.NewStyle().Foreground(ui.Slate).MarginTop(4).Render("▲▼ Navigate • Enter Select • B/ESC Back")

	return lipgloss.Place(m.Width-8, m.Height-6, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, menu, help))
}
