package ui

import (
    "fmt"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    "github.com/raaican/raihand/internal/bot"
)

type GuildsModel struct {
    bot    *bot.Bot
    cursor int
    width  int
    height int
}

func NewGuildsModel(b *bot.Bot) GuildsModel {
    return GuildsModel{bot: b}
}

func (m GuildsModel) Init() tea.Cmd { return nil }

func (m GuildsModel) Update(msg tea.Msg) (GuildsModel, tea.Cmd) {
    guilds := m.bot.Guilds()
    if len(guilds) == 0 {
        return m, nil
    }
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(guilds)-1 {
                m.cursor++
            }
        }
    }
    return m, nil
}

func (m GuildsModel) SetSize(w, h int) GuildsModel {
    m.width = w
    m.height = h
    return m
}

func (m GuildsModel) View() string {
    guilds := m.bot.Guilds()
    title  := styleTitle.Render("guilds")

    if len(guilds) == 0 {
        return lipgloss.JoinVertical(lipgloss.Left, title, styleHelp.Render("No guilds available. Is the bot online?"))
    }

    var list string
    for i, g := range guilds {
        line := fmt.Sprintf("%-30s  ID: %s  Members: %d", g.Name, g.ID, g.MemberCount)
        if i == m.cursor {
            list += styleSelected.Render(line) + "\n"
        } else {
            list += styleItem.Render(line) + "\n"
        }
    }

    var detail string
    if m.cursor < len(guilds) {
        g := guilds[m.cursor]
        detail = styleBorder.Width(m.width/2 - 4).Render(
            styleTitle.Render("Selected Guild") + "\n" +
            fmt.Sprintf("Name:    %s\n", g.Name) +
            fmt.Sprintf("ID:      %s\n", g.ID) +
            fmt.Sprintf("Members: %d\n", g.MemberCount) +
            fmt.Sprintf("Region:  %s\n", g.PreferredLocale),
        )
    }

    listBox := styleBorder.Width(m.width/2 - 4).Render(list)
    row     := lipgloss.JoinHorizontal(lipgloss.Top, listBox, detail)
    help    := styleHelp.Render("↑/↓ or j/k: navigate")
    return lipgloss.JoinVertical(lipgloss.Left, title, row, help)
}
