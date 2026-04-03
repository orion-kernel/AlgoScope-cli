package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	width  = 60
	height = 15
	delay  = 10 * time.Millisecond
)

var (
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00F5FF")).Bold(true).MarginBottom(1)
	barStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))
	activeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7"))
	completeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00F5FF"))
	statsStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type model struct {
	array    []int
	i, j     int
	swapped  bool
	done     bool
	start    time.Time
	elapsed  time.Duration
}

func initialModel() model {
	rand.Seed(time.Now().UnixNano())
	arr := make([]int, width)
	for i := range arr {
		arr[i] = rand.Intn(height*8-2) + 1
	}
	return model{
		array: arr,
		i:     0,
		j:     0,
		start: time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			return initialModel(), tick()
		}

	case tickMsg:
		if m.done {
			return m, nil
		}

		m.elapsed = time.Since(m.start)

		// Multiple steps per tick for better performance
		for k := 0; k < 4; k++ {
			if m.i < len(m.array)-1 {
				if m.j < len(m.array)-m.i-1 {
					if m.array[m.j] > m.array[m.j+1] {
						m.array[m.j], m.array[m.j+1] = m.array[m.j+1], m.array[m.j]
						m.swapped = true
					}
					m.j++
				} else {
					if !m.swapped {
						m.done = true
						break
					}
					m.j = 0
					m.i++
					m.swapped = false
				}
			} else {
				m.done = true
				break
			}
		}

		return m, tick()
	}

	return m, nil
}

func getBarChar(h int) string {
	// Unicode lower block characters for sub-cell height
	chars := []string{" ", " ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	if h <= 0 {
		return " "
	}
	if h >= 8 {
		return "█"
	}
	return chars[h]
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("AlgoScope CLI - Bubble Sort"))
	b.WriteString("\n")

	// Render bars from top to bottom
	for h := height - 1; h >= 0; h-- {
		row := ""
		for index, val := range m.array {
			char := " "
			style := barStyle

			cellHeight := val - (h * 8)
			if cellHeight >= 8 {
				char = "█"
			} else if cellHeight > 0 {
				char = getBarChar(cellHeight)
			} else {
				char = " "
			}

			if m.done {
				style = completeStyle
			} else if index == m.j || index == m.j+1 {
				style = activeStyle
			}

			row += style.Render(char)
		}
		b.WriteString(row + "\n")
	}

	status := ""
	if m.done {
		status = completeStyle.Render(fmt.Sprintf("\nSorting Complete! Time: %v", m.elapsed.Round(time.Millisecond)))
		status += "\nPress 'r' to restart, 'q' to quit."
	} else {
		status = fmt.Sprintf("\nSorting... Iteration: %d/%d", m.i, len(m.array))
		status += statsStyle.Render(fmt.Sprintf(" | Time: %v", m.elapsed.Round(time.Second)))
	}
	b.WriteString(status)

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
