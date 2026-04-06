package tui

import (
	"strings"

	"AlgoScope-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) DocsView() string {
	sel := m.Menu[m.Cursor]
	header := lipgloss.NewStyle().
		Foreground(ui.Teal).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.Blue).
		Padding(0, 3).
		MarginBottom(1).
		Render("󰈙 DOCUMENTATION: " + sel.Name)

	lines := strings.Split(m.DocsText, "\n")
	var docBody strings.Builder
	
	// Helper to parse bolding in any string
	parseBold := func(txt string) string {
		if !strings.Contains(txt, "**") {
			return txt
		}
		parts := strings.Split(txt, "**")
		var res strings.Builder
		for i, p := range parts {
			if i%2 == 1 {
				res.WriteString(lipgloss.NewStyle().Foreground(ui.Teal).Bold(true).Render(p))
			} else {
				res.WriteString(p)
			}
		}
		return res.String()
	}

	h1Skipped := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			docBody.WriteString("\n")
			continue
		}

		if !h1Skipped && strings.HasPrefix(trimmed, "# ") {
			h1Skipped = true
			continue
		}

		if strings.HasPrefix(trimmed, "## ") {
			title := strings.TrimPrefix(trimmed, "## ")
			docBody.WriteString("\n" + lipgloss.NewStyle().Foreground(ui.Blue).Bold(true).Render(title) + "\n")
		} else if strings.HasPrefix(trimmed, "> ") {
			content := strings.TrimPrefix(trimmed, "> ")
			docBody.WriteString(lipgloss.NewStyle().Foreground(ui.Slate).Italic(true).PaddingLeft(2).Render(content) + "\n")
		} else if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			content := strings.TrimPrefix(trimmed, "- ")
			content = strings.TrimPrefix(content, "* ")
			
			if strings.Contains(content, ":") {
				parts := strings.SplitN(content, ":", 2)
				key := lipgloss.NewStyle().Foreground(ui.Orange).Bold(true).Render(parseBold(parts[0]) + ":")
				val := lipgloss.NewStyle().Foreground(ui.White).Render(parseBold(parts[1]))
				docBody.WriteString("  • " + key + val + "\n")
			} else {
				docBody.WriteString("  • " + parseBold(content) + "\n")
			}
		} else if strings.HasPrefix(trimmed, "---") {
			docBody.WriteString(lipgloss.NewStyle().Foreground(ui.Slate).Render(strings.Repeat("─", 50)) + "\n")
		} else {
			docBody.WriteString(parseBold(line) + "\n")
		}
	}

	scrollArea := lipgloss.NewStyle().
		Padding(0, 4).
		Width(m.Width - 12).
		Height(m.Height - 10).
		Render(docBody.String())

	footer := lipgloss.NewStyle().Foreground(ui.Slate).MarginTop(1).Render("B/ESC Back • Q Quit")

	return lipgloss.JoinVertical(lipgloss.Center, header, scrollArea, footer)
}
