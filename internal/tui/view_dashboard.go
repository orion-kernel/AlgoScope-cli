package tui

import (
	"AlgoScope-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) DashboardView() string {
	// Logo
	logo := ui.LogoStyle.Render(
		"                                                                             \n" +
			"     ▄▄    ▄▄               ▄▄▄▄▄                             ▄   ▄▄▄▄ ▄▄    \n" +
			"   ▄█▀▀█▄   ██             ██▀▀▀▀█▄                           ▀██████▀  ██   \n" +
			"   ██  ██   ██    ▄▄       ▀██▄  ▄▀                             ██      ██ ▀▀\n" +
			"   ██▀▀██   ██ ▄████ ▄███▄   ▀██▄▄  ▄███▀ ▄███▄ ████▄ ▄█▀█▄     ██      ██ ██\n" +
			" ▄ ██  ██   ██ ██ ██ ██ ██ ▄   ▀██▄ ██    ██ ██ ██ ██ ██▄█▀     ██      ██ ██\n" +
			" ▀██▀  ▀█▄█▄██▄▀████▄▀███▀ ▀██████▀▄▀███▄▄▀███▀▄████▀▄▀█▄▄▄     ▀█████ ▄██▄██\n" +
			"                  ██                            ██                           \n" +
			"                ▀▀▀                             ▀                            \n" +
			"                - ADVANCED ALGORITHM VISUALIZATION ENGINE -   ")

	availW := m.Width - 8
	
	// Sidebar
	var menuItems []string
	for i, item := range m.Menu {
		txt := item.Name
		if i == m.Cursor {
			txt = lipgloss.NewStyle().Foreground(ui.Teal).Bold(true).Background(ui.DarkGray).Padding(0, 1).Render(" " + txt)
		} else {
			txt = lipgloss.NewStyle().Foreground(ui.Slate).Padding(0, 1).Render("  " + txt)
		}
		menuItems = append(menuItems, txt)
	}

	sbWidth := 30
	if availW < 80 { sbWidth = 25 }
	
	sidebar := ui.SidebarStyle.Width(sbWidth).Render(lipgloss.JoinVertical(lipgloss.Left, menuItems...))

	// Main Panel
	sel := m.Menu[m.Cursor]

	compBadge := ui.TagStyle.Background(ui.Orange).Render(sel.Complexity)
	stabBadge := ui.TagStyle.Background(ui.Blue).Render(sel.Stability)
	if sel.Name == "EXIT ENGINE" {
		compBadge, stabBadge = "", ""
	}

	mcWidth := availW - sbWidth - 4
	
	content := ui.MainContentStyle.Width(mcWidth).Render(lipgloss.JoinVertical(lipgloss.Left,
		ui.TitleStyle.Render(sel.Name),
		lipgloss.JoinHorizontal(lipgloss.Top, compBadge, stabBadge),
		"",
		lipgloss.NewStyle().Foreground(ui.White).Width(mcWidth).Render(sel.Desc),
		"",
		ui.SystemStatusStyle.Render("System: Ready"),
		ui.SystemStatusStyle.Render("Core: Online"),
	))

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	footer := lipgloss.NewStyle().Foreground(ui.Slate).MarginTop(2).Render("▲▼ Navigate • Enter Select • Q Quit")

	return lipgloss.JoinVertical(lipgloss.Center, logo, body, footer)
}
