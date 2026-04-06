package tui

import (
	"os"
	"time"

	"AlgoScope-cli/internal/algo"
	"AlgoScope-cli/internal/engine"
	"AlgoScope-cli/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- State Constants ---
const (
	StateDashboard = iota
	StateAlgorithmMenu
	StateVisualizer
	StateDocs
)

type tickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(10*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type Model struct {
	State      int
	Width      int
	Height     int
	Cursor     int
	SubCursor  int
	Menu       []algo.Algorithm
	DocsText   string
	Engine     *engine.Engine
}

func InitialModel() Model {
	return Model{
		State:  StateDashboard,
		Menu:   algo.GetAlgorithms(),
		Cursor: 0,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.State == StateDashboard {
				return m, tea.Quit
			}
			m.State = StateDashboard
		case "up", "k":
			if m.State == StateDashboard {
				if m.Cursor > 0 { m.Cursor-- }
			} else if m.State == StateAlgorithmMenu {
				if m.SubCursor > 0 { m.SubCursor-- }
			}
		case "down", "j":
			if m.State == StateDashboard {
				if m.Cursor < len(m.Menu)-1 { m.Cursor++ }
			} else if m.State == StateAlgorithmMenu {
				if m.SubCursor < 1 { m.SubCursor++ }
			}
		case "enter", " ":
			if m.State == StateDashboard {
				if m.Menu[m.Cursor].Name == "EXIT ENGINE" {
					return m, tea.Quit
				}
				m.State = StateAlgorithmMenu
				m.SubCursor = 0
			} else if m.State == StateAlgorithmMenu {
				if m.SubCursor == 0 {
					m.State = StateVisualizer
					m.Engine = engine.NewEngine(m.Width, m.Height, m.Menu[m.Cursor].ID)
					return m, Tick()
				} else {
					// Load docs
					path := m.Menu[m.Cursor].DocPath
					if path != "" {
						content, err := os.ReadFile(path)
						if err == nil {
							m.DocsText = string(content)
							m.State = StateDocs
						} else {
							m.DocsText = "Documentation not found for this algorithm."
							m.State = StateDocs
						}
					} else {
						m.DocsText = "No documentation available for " + m.Menu[m.Cursor].Name
						m.State = StateDocs
					}
				}
			}
		case "esc", "b":
			if m.State == StateAlgorithmMenu {
				m.State = StateDashboard
			} else if m.State == StateVisualizer || m.State == StateDocs {
				m.State = StateAlgorithmMenu
			}
		case "r":
			if m.State == StateVisualizer {
				m.Engine = engine.NewEngine(m.Width, m.Height, m.Menu[m.Cursor].ID)
				return m, Tick()
			}
		}
	case tickMsg:
		if m.State != StateVisualizer || m.Engine == nil || m.Engine.Done {
			return m, nil
		}
		m.Engine.Tick()
		return m, Tick()
	}
	return m, nil
}

func (m Model) View() string {
	if m.Width < 90 || m.Height < 24 {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(ui.Red).Bold(true).Render("Terminal too small for AlgoScope CLI\nMinimum required: 90x24"))
	}

	var content string
	switch m.State {
	case StateDashboard:
		content = m.DashboardView()
	case StateAlgorithmMenu:
		content = m.AlgorithmMenuView()
	case StateVisualizer:
		content = m.SortView()
	case StateDocs:
		content = m.DocsView()
	}

	fWidth := m.Width - 4
	fHeight := m.Height - 2

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center,
		ui.FrameStyle.Width(fWidth).Height(fHeight).Render(content))
}
