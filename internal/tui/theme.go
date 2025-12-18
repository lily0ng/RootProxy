package tui

import "github.com/charmbracelet/lipgloss"

const (
	ColorHTBBlack  = "#0C0C0C"
	ColorHTBPurple = "#9D4EDD"
	ColorHTBGreen  = "#00FF00"
	ColorHTBGray   = "#2D2D2D"
	ColorHTBRed    = "#FF0055"
)

type Theme struct {
	Bg      lipgloss.Color
	Fg      lipgloss.Color
	Accent  lipgloss.Color
	Success lipgloss.Color
	Danger  lipgloss.Color
	Muted   lipgloss.Color
	PanelBg lipgloss.Color
	Border  lipgloss.Color
}

func HTBDark() Theme {
	return Theme{
		Bg:      lipgloss.Color(ColorHTBBlack),
		Fg:      lipgloss.Color("#EDEDED"),
		Accent:  lipgloss.Color(ColorHTBPurple),
		Success: lipgloss.Color(ColorHTBGreen),
		Danger:  lipgloss.Color(ColorHTBRed),
		Muted:   lipgloss.Color(ColorHTBGray),
		PanelBg: lipgloss.Color("#141414"),
		Border:  lipgloss.Color(ColorHTBPurple),
	}
}
