package ui

import (
	"time"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    "github.com/raaican/raihand/internal/bot"
)

type tab int

type globalTickMsg time.Time

const (
    tabDashboard tab = iota
    tabLogs
    tabGuilds
    tabActions
)

var tabNames = []string{"Dashboard", "Logs", "Guilds", "Actions"}

type RootModel struct {
    bot       *bot.Bot
    activeTab tab
    width     int
    height    int
    dashboard DashboardModel
    logView   LogViewModel
    guilds    GuildsModel
    actions   ActionsModel
}

func NewRootModel(b *bot.Bot) RootModel {
    return RootModel{
        bot:       b,
        activeTab: tabDashboard,
        dashboard: NewDashboardModel(b),
        logView:   NewLogViewModel(b),
        guilds:    NewGuildsModel(b),
        actions:   NewActionsModel(b),
    }
}

func globalTickCmd() tea.Cmd {
	return tea.Every(500*time.Millisecond, func(t time.Time) tea.Msg {
		return globalTickMsg(t)
	})
}

func (m RootModel) Init() tea.Cmd {
    return tea.Batch(
        m.dashboard.Init(),
        m.logView.Init(),
        m.guilds.Init(),
        m.actions.Init(),
		globalTickCmd(),
    )
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

    switch msg := msg.(type) {
	case globalTickMsg:
		return m, globalTickCmd()

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        inner := m.height - 4
        m.dashboard = m.dashboard.SetSize(m.width, inner)
        m.logView   = m.logView.SetSize(m.width, inner)
        m.guilds    = m.guilds.SetSize(m.width, inner)
        m.actions   = m.actions.SetSize(m.width, inner)

    case tea.KeyPressMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "1":
            m.activeTab = tabDashboard
        case "2":
            m.activeTab = tabLogs
        case "3":
            m.activeTab = tabGuilds
        case "4":
            m.activeTab = tabActions
        case "tab":
            m.activeTab = (m.activeTab + 1) % tab(len(tabNames))
        }
    }

	var logCmd tea.Cmd
	m.logView, logCmd = m.logView.Update(msg)
	cmds = append(cmds, logCmd)

    var tabCmd tea.Cmd
    switch m.activeTab {
    case tabDashboard:
        m.dashboard, tabCmd = m.dashboard.Update(msg)
    case tabGuilds:
        m.guilds, tabCmd = m.guilds.Update(msg)
    case tabActions:
        m.actions, tabCmd = m.actions.Update(msg)
    }
	cmds = append(cmds, tabCmd)

    return m, tea.Batch(cmds...)
}

func (m RootModel) View() tea.View {
    tabs    := m.renderTabs()
    var content string
    switch m.activeTab {
    case tabDashboard:
        content = m.dashboard.View()
    case tabLogs:
        content = m.logView.View()
    case tabGuilds:
        content = m.guilds.View()
    case tabActions:
        content = m.actions.View()
    }
    help := styleHelp.Render("tab: next  1-4: jump  q: quit")
    out  := lipgloss.JoinVertical(lipgloss.Left, tabs, content, help)
    v    := tea.NewView(out)
    v.AltScreen = true
    return v
}

func (m RootModel) renderTabs() string {
    var tabs []string
    for i, name := range tabNames {
        if tab(i) == m.activeTab {
            tabs = append(tabs, styleTabActive.Render(name))
        } else {
            tabs = append(tabs, styleTabInactive.Render(name))
        }
    }
    return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}
