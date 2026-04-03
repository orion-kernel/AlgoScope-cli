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

// --- Configuration & Constants ---
const (
	delay = 10 * time.Millisecond
)

type viewState int

const (
	viewDashboard viewState = iota
	viewBubbleSort
)

// --- Styles ---
var (
	// Colors
	cyan    = lipgloss.Color("#00F5FF")
	purple  = lipgloss.Color("#7D56F4")
	pink    = lipgloss.Color("#FF06B7")
	gray    = lipgloss.Color("#888888")
	black   = lipgloss.Color("#000000")

	// Dashboard Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			Padding(1, 4).
			MarginBottom(2)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			BorderForeground(purple)

	buttonStyle = lipgloss.NewStyle().
			Foreground(purple).
			Padding(0, 2).
			MarginTop(1)

	activeButtonStyle = buttonStyle.
				Foreground(black).
				Background(cyan).
				Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true).
			MarginTop(2)

	// Sort Visualization Styles
	barStyle      = lipgloss.NewStyle().Foreground(purple)
	activeStyle   = lipgloss.NewStyle().Foreground(pink)
	completeStyle = lipgloss.NewStyle().Foreground(cyan)
	statsStyle    = lipgloss.NewStyle().Foreground(gray).Italic(true)
)

// --- Messages & Commands ---
type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// --- Model ---
type model struct {
	state  viewState
	width  int
	height int

	// Dashboard State
	cursor int
	menu   []string

	// Bubble Sort State
	array   []int
	i, j    int
	swapped bool
	done    bool
	start   time.Time
	elapsed time.Duration
}

func initialModel() model {
	return model{
		state:  viewDashboard,
		menu:   []string{"Bubble Sort", "Quit"},
		cursor: 0,
	}
}

func (m *model) initSort() {
	rand.Seed(time.Now().UnixNano())
	// Use 80% of screen width for the visualization
	visWidth := int(float64(m.width) * 0.8)
	if visWidth > 80 {
		visWidth = 80
	}
	if visWidth < 10 {
		visWidth = 10
	}

	// Max height for bars
	maxBarHeight := (m.height - 10) * 8
	if maxBarHeight < 8 {
		maxBarHeight = 8
	}

	arr := make([]int, visWidth)
	for i := range arr {
		arr[i] = rand.Intn(maxBarHeight-2) + 1
	}
	m.array = arr
	m.i = 0
	m.j = 0
	m.swapped = false
	m.done = false
	m.start = time.Now()
}

func (m model) Init() tea.Cmd {
	return nil
}

// --- Update ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.state == viewDashboard {
				return m, tea.Quit
			}
			m.state = viewDashboard
			return m, nil

		case "up", "k":
			if m.state == viewDashboard && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == viewDashboard && m.cursor < len(m.menu)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.state == viewDashboard {
				switch m.cursor {
				case 0: // Bubble Sort
					m.state = viewBubbleSort
					m.initSort()
					return m, tick()
				case 1: // Quit
					return m, tea.Quit
				}
			}

		case "r":
			if m.state == viewBubbleSort {
				m.initSort()
				return m, tick()
			}
		case "esc":
			m.state = viewDashboard
			return m, nil
		}

	case tickMsg:
		if m.state != viewBubbleSort || m.done {
			return m, nil
		}

		m.elapsed = time.Since(m.start)

		// Multiple steps per tick for better performance
		steps := 4
		if len(m.array) > 60 {
			steps = 8
		}
		for k := 0; k < steps; k++ {
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

// --- View ---
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var content string
	switch m.state {
	case viewDashboard:
		content = m.dashboardView()
	case viewBubbleSort:
		content = m.bubbleSortView()
	}

	// Center everything on screen
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m model) dashboardView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("ALGOSCOPE DASHBOARD"))
	b.WriteString("\n\n")

	var menuItems []string
	for i, item := range m.menu {
		if m.cursor == i {
			menuItems = append(menuItems, activeButtonStyle.Render("> "+item+" <"))
		} else {
			menuItems = append(menuItems, buttonStyle.Render("  "+item+"  "))
		}
	}

	b.WriteString(menuStyle.Render(lipgloss.JoinVertical(lipgloss.Center, menuItems...)))
	b.WriteString("\n")
	b.WriteString(footerStyle.Render("Use arrow keys to navigate • Enter to select"))

	return b.String()
}

func (m model) bubbleSortView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("BUBBLE SORT VISUALIZER"))
	b.WriteString("\n")

	// Render bars
	visHeight := m.height - 15
	if visHeight < 5 {
		visHeight = 5
	}

	var bars strings.Builder
	for h := visHeight - 1; h >= 0; h-- {
		var row strings.Builder
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

			row.WriteString(style.Render(char))
		}
		bars.WriteString(row.String())
		if h > 0 {
			bars.WriteString("\n")
		}
	}

	b.WriteString(lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(purple).
		Padding(1, 2).
		Width(len(m.array) + 4). // +4 to account for padding
		Render(bars.String()))

	status := ""
	if m.done {
		status = completeStyle.Render(fmt.Sprintf("\nSorting Complete! Time: %v", m.elapsed.Round(time.Millisecond)))
		status += "\nPress 'r' to restart, 'esc' for dashboard."
	} else {
		status = fmt.Sprintf("\nSorting... Iteration: %d/%d", m.i, len(m.array))
		status += statsStyle.Render(fmt.Sprintf(" | Time: %v", m.elapsed.Round(time.Second)))
		status += "\nPress 'esc' to return to dashboard."
	}
	b.WriteString(status)

	return b.String()
}

func getBarChar(h int) string {
	chars := []string{" ", " ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	if h <= 0 {
		return " "
	}
	if h >= 8 {
		return "█"
	}
	return chars[h]
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
