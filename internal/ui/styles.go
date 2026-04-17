package ui

import "charm.land/lipgloss/v2"

var (
	colorPrimary = lipgloss.Color("#7289DA")
	colorSuccess = lipgloss.Color("#43B581")
	colorWarning = lipgloss.Color("#FAA61A")
	colorError   = lipgloss.Color("#F04747")
	colorMuted   = lipgloss.Color("#72767D")
	colorBg      = lipgloss.Color("#2C2F33")
	colorBgLight = lipgloss.Color("#36393F")
	colorText    = lipgloss.Color("#DCDDDE")
	colorTextDim = lipgloss.Color("#8E9297")

	styleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorPrimary).
			Padding(0, 2)

	styleTabInactive = lipgloss.NewStyle().
				Foreground(colorTextDim).
				Background(colorBgLight).
				Padding(0, 2)

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary)

	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	styleHelp = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true).
			MarginTop(1)

	styleStatusOnline = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorSuccess)

	styleStatusOffline = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorError)

	styleLogInfo = lipgloss.NewStyle().
			Foreground(colorSuccess)

	styleLogWarn = lipgloss.NewStyle().
			Foreground(colorWarning)

	styleLogError = lipgloss.NewStyle().
			Foreground(colorError)

	styleLogTime = lipgloss.NewStyle().
			Foreground(colorTextDim)

	styleCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBgLight).
			Padding(0, 1)

	stylePrimary = lipgloss.NewStyle().
			Foreground(colorPrimary)

	styleSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorPrimary).
			Padding(0, 1)

	styleItem = lipgloss.NewStyle().
			Foreground(colorText).
			Padding(0, 1)
)
