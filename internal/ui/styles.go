package ui

import "github.com/charmbracelet/lipgloss"

// --- Colors (Tokyo Night Theme) ---
var (
	Teal     = lipgloss.Color("#73daca")
	Lavender = lipgloss.Color("#bb9af7")
	Magenta  = lipgloss.Color("#bb9af7")
	Blue     = lipgloss.Color("#7aa2f7")
	DarkGray = lipgloss.Color("#24283b")
	Slate    = lipgloss.Color("#565f89")
	Orange   = lipgloss.Color("#ff9e64")
	Green    = lipgloss.Color("#9ece6a")
	Red      = lipgloss.Color("#f7768e")
	Bg       = lipgloss.Color("#1a1b26")
	White    = lipgloss.Color("#c0caf5")
)

// --- Styles ---
var (
	LogoStyle = lipgloss.NewStyle().
			Foreground(Teal).
			Bold(true).
			MarginBottom(1)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Blue).
			Padding(1, 2).
			Align(lipgloss.Center, lipgloss.Center)

	SidebarStyle = lipgloss.NewStyle().
			PaddingRight(2).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(Slate)

	MainContentStyle = lipgloss.NewStyle().
				PaddingLeft(4)

	TagStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			MarginRight(1).
			Foreground(White)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Lavender).
			Bold(true).
			Underline(true).
			MarginBottom(1)

	SystemStatusStyle = lipgloss.NewStyle().
				Foreground(Slate).
				Italic(true)

	// Sort Visualization
	BarStyle      = lipgloss.NewStyle().Foreground(Blue)
	ActiveStyle   = lipgloss.NewStyle().Foreground(Magenta).Bold(true)
	CompleteStyle = lipgloss.NewStyle().Foreground(Green)
)

func GetBarChar(h int) string {
	chars := []string{" ", " ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	if h <= 0 {
		return " "
	}
	if h >= 8 {
		return "█"
	}
	return chars[h]
}
