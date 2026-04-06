package tui

import (
	"fmt"
	"strings"
	"time"

	"AlgoScope-cli/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) SortView() string {
	if m.Engine == nil {
		return "Initializing Engine..."
	}

	header := lipgloss.NewStyle().
		Foreground(ui.Teal).
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.Blue).
		Render(" 󰓡 VISUALIZING: " + m.Menu[m.Cursor].Name + " ")

	vH := m.Height - 18
	if vH < 5 {
		vH = 5
	}

	var bars strings.Builder
	for h := vH - 1; h >= 0; h-- {
		var row strings.Builder
		for idx, val := range m.Engine.Array {
			char, style := " ", ui.BarStyle
			ch := val - (h * 8)
			if ch >= 8 {
				char = "█"
			} else if ch > 0 {
				char = ui.GetBarChar(ch)
			}
			if m.Engine.Done {
				style = ui.CompleteStyle
			} else if idx == m.Engine.J || idx == m.Engine.J+1 {
				style = ui.ActiveStyle
			}
			row.WriteString(style.Render(char))
		}
		bars.WriteString(row.String() + "\n")
	}

	visWidth := len(m.Engine.Array) + 4
	maxVisW := m.Width - 12
	if visWidth > maxVisW {
		visWidth = maxVisW
	}

	vis := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.Blue).
		Padding(1, 2).
		Width(visWidth).
		Render(bars.String())

	stats := lipgloss.JoinHorizontal(lipgloss.Center,
		ui.TagStyle.Background(ui.DarkGray).Foreground(ui.Lavender).Render(fmt.Sprintf("󱫐 ITER: %d", m.Engine.I)),
		ui.TagStyle.Background(ui.DarkGray).Foreground(ui.Teal).Render(fmt.Sprintf("󱎫 TIME: %s", m.Engine.Elapsed.Round(time.Second))),
	)

	if m.Engine.Done {
		stats = ui.CompleteStyle.Bold(true).Render(fmt.Sprintf("󰄬 SORTING COMPLETE! TOTAL TIME: %v", m.Engine.Elapsed.Round(time.Millisecond)))
	}

	help := lipgloss.NewStyle().Foreground(ui.Slate).MarginTop(1).Render("R Restart • ESC Back • Q Quit")

	return lipgloss.JoinVertical(lipgloss.Center, header, vis, stats, help)
}
