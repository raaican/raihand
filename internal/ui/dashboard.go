package ui

import (
    "fmt"
    "time"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    "github.com/raaican/raihand/internal/bot"
)

type tickMsg time.Time

type DashboardModel struct {
    bot    *bot.Bot
    width  int
    height int
}

func NewDashboardModel(b *bot.Bot) DashboardModel {
    return DashboardModel{bot: b}
}

func (m DashboardModel) Init() tea.Cmd {
    return tickCmd()
}

func tickCmd() tea.Cmd {
    return tea.Every(2*time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
    if _, ok := msg.(tickMsg); ok {
        return m, tickCmd()
    }
    return m, nil
}

func (m DashboardModel) SetSize(w, h int) DashboardModel {
    m.width = w
    m.height = h
    return m
}

func (m DashboardModel) View() string {
    status := m.bot.Status()
    guilds := m.bot.Guilds()
    logs   := m.bot.Logs()

    var statusStr string
    if status == bot.StatusOnline {
        statusStr = styleStatusOnline.Render("* ONLINE")
    } else {
        statusStr = styleStatusOffline.Render("x OFFLINE")
    }

    cards := []string{
        styleCard.Render(fmt.Sprintf("%s\n%s", styleTitle.Render("Status"), statusStr)),
        styleCard.Render(fmt.Sprintf("%s\n%s", styleTitle.Render("Guilds"), stylePrimary.Render(fmt.Sprintf("%d", len(guilds))))),
        styleCard.Render(fmt.Sprintf("%s\n%s", styleTitle.Render("Log Entries"), stylePrimary.Render(fmt.Sprintf("%d", len(logs))))),
        styleCard.Render(fmt.Sprintf("%s\n%s", styleTitle.Render("Time"), stylePrimary.Render(time.Now().Format("15:04:05")))),
    }
    row := lipgloss.JoinHorizontal(lipgloss.Top, cards...)

    recentTitle := styleTitle.Render("Recent Logs")
    var recentLogs string
    start := len(logs) - 5
    if start < 0 {
        start = 0
    }
    for _, e := range logs[start:] {
        recentLogs += renderLogLine(e) + "\n"
    }
    if recentLogs == "" {
        recentLogs = styleHelp.Render("No logs yet...")
    }
    logsBox := styleBorder.Width(m.width - 4).Render(recentTitle + "\n" + recentLogs)

    return lipgloss.JoinVertical(lipgloss.Left,
        styleTitle.Render("raihand control panel"),
        row,
        logsBox,
    )
}
