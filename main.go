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
	teal     = lipgloss.Color("#73daca")
	lavender = lipgloss.Color("#bb9af7")
	magenta  = lipgloss.Color("#bb9af7")
	blue     = lipgloss.Color("#7aa2f7")
	darkGray = lipgloss.Color("#24283b")
	slate    = lipgloss.Color("#565f89")
	orange   = lipgloss.Color("#ff9e64")
	green    = lipgloss.Color("#9ece6a")
	red      = lipgloss.Color("#f7768e")
	bg       = lipgloss.Color("#1a1b26")
	white    = lipgloss.Color("#c0caf5")
)

// --- Styles ---
var (
	logoStyle = lipgloss.NewStyle().
			Foreground(teal).
			Bold(true).
			MarginBottom(1)

	frameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blue).
			Padding(1, 2).
			Align(lipgloss.Center, lipgloss.Center)

	sidebarStyle = lipgloss.NewStyle().
			PaddingRight(2).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(slate)

	mainContentStyle = lipgloss.NewStyle().
				PaddingLeft(4)

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

// --- State Constants ---
const (
	stateDashboard = iota
	stateAlgorithmMenu
	stateVisualizer
	stateDocs
)

type algorithm struct {
	name       string
	desc       string
	complexity string
	stability  string
	docPath    string
}

type model struct {
	state      int
	width      int
	height     int
	cursor     int
	subCursor  int
	menu       []algorithm
	docsText   string

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
		state: stateDashboard,
		menu: []algorithm{
			{"BUBBLE SORT", "Classic O(n²) sorting algorithm. Ideal for visualizing the basic concept of swapping and iterations.", "O(n²)", "STABLE", "docs/bubble_sort/README.md"},
			{"QUICK SORT", "Highly efficient O(n log n) divide-and-conquer algorithm. Selects a pivot to partition data.", "O(n log n)", "UNSTABLE", "docs/quick_sort/README.md"},
			{"MERGE SORT", "Reliable O(n log n) stable sort. Recursively divides array into halves and merges them back.", "O(n log n)", "STABLE", "docs/merge_sort/README.md"},
			{"EXIT ENGINE", "Terminate the AlgoScope visualization engine and return to host shell.", "N/A", "N/A", ""},
		},
		cursor: 0,
		subCursor: 0,
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
			if m.state == stateDashboard {
				return m, tea.Quit
			}
			m.state = stateDashboard
		case "up", "k":
			if m.state == stateDashboard {
				if m.cursor > 0 { m.cursor-- }
			} else if m.state == stateAlgorithmMenu {
				if m.subCursor > 0 { m.subCursor-- }
			}
		case "down", "j":
			if m.state == stateDashboard {
				if m.cursor < len(m.menu)-1 { m.cursor++ }
			} else if m.state == stateAlgorithmMenu {
				if m.subCursor < 1 { m.subCursor++ }
			}
		case "enter", " ":
			if m.state == stateDashboard {
				if m.cursor == 3 {
					return m, tea.Quit
				}
				m.state = stateAlgorithmMenu
				m.subCursor = 0
			} else if m.state == stateAlgorithmMenu {
				if m.subCursor == 0 {
					m.state = stateVisualizer
					m.initSort()
					return m, tick()
				} else {
					// Load docs
					path := m.menu[m.cursor].docPath
					if path != "" {
						content, err := os.ReadFile(path)
						if err == nil {
							m.docsText = string(content)
							m.state = stateDocs
						} else {
							m.docsText = "Documentation not found for this algorithm."
							m.state = stateDocs
						}
					} else {
						m.docsText = "No documentation available for " + m.menu[m.cursor].name
						m.state = stateDocs
					}
				}
			}
		case "esc", "b":
			if m.state == stateAlgorithmMenu {
				m.state = stateDashboard
			} else if m.state == stateVisualizer || m.state == stateDocs {
				m.state = stateAlgorithmMenu
			}
		case "r":
			if m.state == stateVisualizer {
				m.initSort()
				return m, tick()
			}
		}
	case tickMsg:
		if m.state != stateVisualizer || m.done {
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
	if m.width < 90 || m.height < 24 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, 
			lipgloss.NewStyle().Foreground(red).Bold(true).Render("Terminal too small for AlgoScope CLI\nMinimum required: 90x24"))
	}

	var content string
	switch m.state {
	case stateDashboard:
		content = m.dashboardView()
	case stateAlgorithmMenu:
		content = m.algorithmMenuView()
	case stateVisualizer:
		content = m.sortView()
	case stateDocs:
		content = m.docsView()
	}

	// Calculate optimal frame size with margins
	fWidth := m.width - 4
	fHeight := m.height - 2

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		frameStyle.Width(fWidth).Height(fHeight).Render(content))
}

func (m model) dashboardView() string {
	// Logo
	logo := logoStyle.Render(
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

	// Available space inside frame (accounting for padding)
	availW := m.width - 8 // 4 (margins) + 4 (frame padding)
	
	// Sidebar
	var menuItems []string
	for i, item := range m.menu {
		txt := item.name
		if i == m.cursor {
			txt = lipgloss.NewStyle().Foreground(teal).Bold(true).Background(darkGray).Padding(0, 1).Render(" " + txt)
		} else {
			txt = lipgloss.NewStyle().Foreground(slate).Padding(0, 1).Render("  " + txt)
		}
		menuItems = append(menuItems, txt)
	}

	// Sidebar width: fixed for consistency, or slightly dynamic
	sbWidth := 30
	if availW < 80 { sbWidth = 25 }
	
	sidebar := sidebarStyle.Width(sbWidth).Render(lipgloss.JoinVertical(lipgloss.Left, menuItems...))

	// Main Panel
	sel := m.menu[m.cursor]

	compBadge := tagStyle.Background(orange).Render(sel.complexity)
	stabBadge := tagStyle.Background(blue).Render(sel.stability)
	if sel.name == "EXIT ENGINE" {
		compBadge, stabBadge = "", ""
	}

	// Content width: remaining available space
	mcWidth := availW - sbWidth - 4 // extra gap
	
	content := mainContentStyle.Width(mcWidth).Render(lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(sel.name),
		lipgloss.JoinHorizontal(lipgloss.Top, compBadge, stabBadge),
		"",
		lipgloss.NewStyle().Foreground(white).Width(mcWidth).Render(sel.desc),
		"",
		systemStatusStyle.Render("System: Ready"),
		systemStatusStyle.Render("Core: Online"),
	))

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	footer := lipgloss.NewStyle().Foreground(slate).MarginTop(2).Render("▲▼ Navigate • Enter Select • Q Quit")

	return lipgloss.JoinVertical(lipgloss.Center, logo, body, footer)
}

func (m model) algorithmMenuView() string {
	sel := m.menu[m.cursor]
	
	title := lipgloss.NewStyle().
		Foreground(lavender).
		Bold(true).
		MarginBottom(2).
		Render("ALGORITHM: " + sel.name)

	options := []string{"󰓡 START VISUALIZER", "󰈙 READ DOCUMENTATION"}
	var renderedOpts []string
	for i, opt := range options {
		style := lipgloss.NewStyle().Padding(1, 2).Width(30).Align(lipgloss.Center).Margin(1)
		if i == m.subCursor {
			style = style.Foreground(bg).Background(teal).Bold(true)
		} else {
			style = style.Foreground(teal).Border(lipgloss.NormalBorder()).BorderForeground(teal)
		}
		renderedOpts = append(renderedOpts, style.Render(opt))
	}

	menu := lipgloss.JoinHorizontal(lipgloss.Top, renderedOpts...)
	
	help := lipgloss.NewStyle().Foreground(slate).MarginTop(4).Render("▲▼ Navigate • Enter Select • B/ESC Back")

	return lipgloss.Place(m.width-8, m.height-6, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, menu, help))
}

func (m model) docsView() string {
	sel := m.menu[m.cursor]
	header := lipgloss.NewStyle().
		Foreground(teal).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Padding(0, 3).
		MarginBottom(1).
		Render("󰈙 DOCUMENTATION: " + sel.name)

	lines := strings.Split(m.docsText, "\n")
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
				res.WriteString(lipgloss.NewStyle().Foreground(teal).Bold(true).Render(p))
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

		// Skip the first H1 as it's already in the header box
		if !h1Skipped && strings.HasPrefix(trimmed, "# ") {
			h1Skipped = true
			continue
		}

		if strings.HasPrefix(trimmed, "## ") {
			title := strings.TrimPrefix(trimmed, "## ")
			docBody.WriteString("\n" + lipgloss.NewStyle().Foreground(blue).Bold(true).Render(title) + "\n")
		} else if strings.HasPrefix(trimmed, "> ") {
			content := strings.TrimPrefix(trimmed, "> ")
			docBody.WriteString(lipgloss.NewStyle().Foreground(slate).Italic(true).PaddingLeft(2).Render(content) + "\n")
		} else if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			content := strings.TrimPrefix(trimmed, "- ")
			content = strings.TrimPrefix(content, "* ")
			
			if strings.Contains(content, ":") {
				parts := strings.SplitN(content, ":", 2)
				key := lipgloss.NewStyle().Foreground(orange).Bold(true).Render(parseBold(parts[0]) + ":")
				val := lipgloss.NewStyle().Foreground(white).Render(parseBold(parts[1]))
				docBody.WriteString("  • " + key + val + "\n")
			} else {
				docBody.WriteString("  • " + parseBold(content) + "\n")
			}
		} else if strings.HasPrefix(trimmed, "---") {
			docBody.WriteString(lipgloss.NewStyle().Foreground(slate).Render(strings.Repeat("─", 50)) + "\n")
		} else {
			// Regular text with bolding support
			docBody.WriteString(parseBold(line) + "\n")
		}
	}

	scrollArea := lipgloss.NewStyle().
		Padding(0, 4).
		Width(m.width - 12).
		Height(m.height - 10).
		Render(docBody.String())

	footer := lipgloss.NewStyle().Foreground(slate).MarginTop(1).Render("B/ESC Back • Q Quit")

	return lipgloss.JoinVertical(lipgloss.Center, header, scrollArea, footer)
}

func (m model) sortView() string {
	header := lipgloss.NewStyle().Foreground(teal).Bold(true).Padding(0, 1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(blue).Render(" 󰓡 VISUALIZING: BUBBLE SORT ")

	vH := m.height - 18 // Adjusted for frame
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

	visWidth := len(m.array) + 4
	maxVisW := m.width - 12
	if visWidth > maxVisW {
		visWidth = maxVisW
	}

	vis := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Padding(1, 2).
		Width(visWidth).
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
