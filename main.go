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

// --- Colors (Tokyo Night Theme) ---
var (
	teal      = lipgloss.Color("#73daca")
	lavender  = lipgloss.Color("#bb9af7")
	magenta   = lipgloss.Color("#bb9af7")
	blue      = lipgloss.Color("#7aa2f7")
	darkGray  = lipgloss.Color("#24283b")
	slate     = lipgloss.Color("#565f89")
	orange    = lipgloss.Color("#ff9e64")
	green     = lipgloss.Color("#9ece6a")
	red       = lipgloss.Color("#f7768e")
	bg        = lipgloss.Color("#1a1b26")
	white     = lipgloss.Color("#c0caf5")
)

// --- Styles ---
var (
	logoStyle = lipgloss.NewStyle().
			Foreground(teal).
			Bold(true).
			MarginBottom(1)

	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(blue).
			Padding(1, 2)

	sidebarStyle = lipgloss.NewStyle().
			Width(30).
			PaddingRight(2).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(slate)

	mainContentStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Width(60)

	tagStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			MarginRight(1).
			Foreground(white)

	titleStyle = lipgloss.NewStyle().
			Foreground(lavender).
			Bold(true).
			Underline(true).
			MarginBottom(1)

	systemStatusStyle = lipgloss.NewStyle().
				Foreground(slate).
				Italic(true)

	// Sort Visualization
	barStyle      = lipgloss.NewStyle().Foreground(blue)
	activeStyle   = lipgloss.NewStyle().Foreground(magenta).Bold(true)
	completeStyle = lipgloss.NewStyle().Foreground(green)
)

type algorithm struct {
	name       string
	desc       string
	complexity string
	stability  string
}

type model struct {
	state  int // 0: Dashboard, 1: Bubble Sort
	width  int
	height int
	cursor int
	menu   []algorithm

	// Sort State
	array   []int
	i, j    int
	swapped bool
	done    bool
	start   time.Time
	elapsed time.Duration
}

func initialModel() model {
	return model{
		state: 0,
		menu: []algorithm{
			{"BUBBLE SORT", "Classic O(n²) sorting algorithm. Ideal for visualizing the basic concept of swapping and iterations.", "O(n²)", "STABLE"},
			{"QUICK SORT", "Highly efficient O(n log n) divide-and-conquer algorithm. Selects a pivot to partition data.", "O(n log n)", "UNSTABLE"},
			{"MERGE SORT", "Reliable O(n log n) stable sort. Recursively divides array into halves and merges them back.", "O(n log n)", "STABLE"},
			{"EXIT ENGINE", "Terminate the AlgoScope visualization engine and return to host shell.", "N/A", "N/A"},
		},
		cursor: 0,
	}
}

func (m *model) initSort() {
	rand.Seed(time.Now().UnixNano())
	visWidth := int(float64(m.width) * 0.7)
	if visWidth > 100 {
		visWidth = 100
	}
	if visWidth < 20 {
		visWidth = 20
	}

	maxHeight := (m.height - 20) * 8
	if maxHeight < 8 {
		maxHeight = 8
	}

	m.array = make([]int, visWidth)
	for i := range m.array {
		m.array[i] = rand.Intn(maxHeight-2) + 1
	}
	m.i, m.j, m.swapped, m.done, m.start = 0, 0, false, false, time.Now()
}

func (m model) Init() tea.Cmd { return nil }

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(10*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.state == 0 {
				return m, tea.Quit
			}
			m.state = 0
		case "up", "k":
			if m.state == 0 && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == 0 && m.cursor < len(m.menu)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.state == 0 {
				if m.cursor == 3 {
					return m, tea.Quit
				}
				m.state = 1
				m.initSort()
				return m, tick()
			}
		case "esc":
			m.state = 0
		case "r":
			if m.state == 1 {
				m.initSort()
				return m, tick()
			}
		}
	case tickMsg:
		if m.state != 1 || m.done {
			return m, nil
		}
		m.elapsed = time.Since(m.start)
		for k := 0; k < 6; k++ {
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
					m.j, m.i, m.swapped = 0, m.i+1, false
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

func (m model) View() string {
	if m.width < 20 || m.height < 10 {
		return "Terminal too small..."
	}

	var view string
	if m.state == 0 {
		view = m.dashboardView()
	} else {
		view = m.sortView()
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

func (m model) dashboardView() string {
	// Logo
	logo := logoStyle.Render(
		"▄▄▄▄· ▄▄▄·  ▐ ▄  ▄▄▄· .▄▄ ·  ▄▄·        ▄▄▄·▄▄▄ . \n" +
			"▐█ ▀█▪▐█ ▀█ •█▌▐█▐█ ▀█ ▐█ ▀. ▐█ ▌▪      ▐█ ▄█▀▄.▀· \n" +
			"▐█▀▀█▄▄█▀▀█ ▐█▐▐▌▄█▀▀█ ▄▀▀▀█▄██ ▄▄      ██▀·▐▀▀▪▄ \n" +
			"██▄▪▐█▐█ ▪▐▌██▐█▌▐█ ▪▐▌▐█▄▪▐█▐███▌      ▐█ ▪·▐█▄▄▌ \n" +
			"·▀▀▀▀  ▀  ▀ ▀▀ █▪ ▀  ▀  ▀▀▀▀ ·▀▀▀       ▀    ▀▀▀  \n" +
			"   - ADVANCED ALGORITHM VISUALIZATION ENGINE -   ")

	// Sidebar
	var menuItems []string
	for i, item := range m.menu {
		txt := item.name
		if i == m.cursor {
			txt = lipgloss.NewStyle().Foreground(teal).Bold(true).Background(darkGray).Padding(0, 1).Render(" "+txt)
		} else {
			txt = lipgloss.NewStyle().Foreground(slate).Padding(0, 1).Render("  "+txt)
		}
		menuItems = append(menuItems, txt)
	}
	sidebar := sidebarStyle.Render(lipgloss.JoinVertical(lipgloss.Left, menuItems...))

	// Main Panel
	sel := m.menu[m.cursor]

	compBadge := tagStyle.Background(orange).Render(sel.complexity)
	stabBadge := tagStyle.Background(blue).Render(sel.stability)
	if sel.name == "EXIT ENGINE" {
		compBadge, stabBadge = "", ""
	}

	content := mainContentStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(sel.name),
		lipgloss.JoinHorizontal(lipgloss.Top, compBadge, stabBadge),
		"",
		lipgloss.NewStyle().Foreground(white).Width(40).Render(sel.desc),
		"",
		systemStatusStyle.Render("System: Ready"),
		systemStatusStyle.Render("Core: Online"),
	))

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	footer := lipgloss.NewStyle().Foreground(slate).MarginTop(2).Render("▲▼ Navigate • Enter Select • Q Quit")

	return frameStyle.Render(lipgloss.JoinVertical(lipgloss.Center, logo, body, footer))
}

func (m model) sortView() string {
	header := lipgloss.NewStyle().Foreground(teal).Bold(true).Padding(0, 1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(blue).Render(" 󰓡 VISUALIZING: BUBBLE SORT ")

	vH := m.height - 20
	if vH < 5 {
		vH = 5
	}

	var bars strings.Builder
	for h := vH - 1; h >= 0; h-- {
		var row strings.Builder
		for idx, val := range m.array {
			char, style := " ", barStyle
			ch := val - (h * 8)
			if ch >= 8 {
				char = "█"
			} else if ch > 0 {
				char = getBarChar(ch)
			}
			if m.done {
				style = completeStyle
			} else if idx == m.j || idx == m.j+1 {
				style = activeStyle
			}
			row.WriteString(style.Render(char))
		}
		bars.WriteString(row.String() + "\n")
	}

	vis := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Padding(1, 2).
		Width(len(m.array) + 4).
		Render(bars.String())

	stats := lipgloss.JoinHorizontal(lipgloss.Center,
		tagStyle.Background(darkGray).Foreground(lavender).Render(fmt.Sprintf("󱫐 ITER: %d", m.i)),
		tagStyle.Background(darkGray).Foreground(teal).Render(fmt.Sprintf("󱎫 TIME: %s", m.elapsed.Round(time.Second))),
	)

	if m.done {
		stats = completeStyle.Bold(true).Render(fmt.Sprintf("󰄬 SORTING COMPLETE! TOTAL TIME: %v", m.elapsed.Round(time.Millisecond)))
	}

	help := lipgloss.NewStyle().Foreground(slate).MarginTop(1).Render("R Restart • ESC Back • Q Quit")

	return lipgloss.JoinVertical(lipgloss.Center, header, vis, stats, help)
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
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
