package ui

import (
    "fmt"
    "strings"

    tea "charm.land/bubbletea/v2"
    "charm.land/bubbles/v2/textinput"
    "charm.land/lipgloss/v2"
    "github.com/bwmarrin/discordgo"
    "github.com/raaican/raihand/internal/bot"
)

type actionResult struct {
    err error
    msg string
}

type action int

const (
    actionSendMessage action = iota
    actionSetStatus
    actionKick
    actionBan
    actionCount
)

var actionLabels = []string{
    "Send Message",
    "Set Status",
    "Kick Member",
    "Ban Member",
}

type ActionsModel struct {
    bot         *bot.Bot
    cursor      action
    inputs      []textinput.Model
    focused     int
    result      string
    resultIsErr bool
    width       int
    height      int
}

func NewActionsModel(b *bot.Bot) ActionsModel {
    return ActionsModel{
        bot:    b,
        inputs: makeInputs(),
    }
}

func makeInputs() []textinput.Model {
    placeholders := []string{
        "Channel ID",       // 0 - SendMessage
        "Message content",  // 1 - SendMessage
        "",                 // 2 - (unused gap)
        "Status text",      // 3 - SetStatus
        "Guild ID",         // 4 - Kick/Ban
        "User ID",          // 5 - Kick/Ban
        "Reason",           // 6 - Kick/Ban
    }
    inputs := make([]textinput.Model, len(placeholders))
    for i, ph := range placeholders {
        t := textinput.New()
        t.Placeholder = ph
        t.CharLimit = 100
        inputs[i] = t
    }
    inputs[0].Focus()
    return inputs
}

func (m ActionsModel) Init() tea.Cmd {
    return textinput.Blink
}

func (m ActionsModel) Update(msg tea.Msg) (ActionsModel, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case actionResult:
        if msg.err != nil {
            m.result = msg.err.Error()
            m.resultIsErr = true
        } else {
            m.result = msg.msg
            m.resultIsErr = false
        }

    case tea.KeyPressMsg:
        switch msg.String() {
        case "up", "shift+tab":
            if m.cursor > 0 {
                m.cursor--
                m.resetInputs()
            }
        case "down":
            if int(m.cursor) < int(actionCount)-1 {
                m.cursor++
                m.resetInputs()
            }
        case "enter":
            startIdx := m.inputStartFor(m.cursor)
            count    := m.inputCountFor(m.cursor)
            if m.focused < count-1 {
                m.inputs[startIdx+m.focused].Blur()
                m.focused++
                m.inputs[startIdx+m.focused].Focus()
            } else {
                return m, m.executeAction()
            }
        case "esc":
            m.resetInputs()
        }
    }

    for i := range m.inputs {
        var cmd tea.Cmd
        m.inputs[i], cmd = m.inputs[i].Update(msg)
        cmds = append(cmds, cmd)
    }
    return m, tea.Batch(cmds...)
}

func (m *ActionsModel) resetInputs() {
    startIdx := m.inputStartFor(m.cursor)
    for i := range m.inputs {
        m.inputs[i].Blur()
        m.inputs[i].SetValue("")
    }
    m.focused = 0
    m.inputs[startIdx].Focus()
    m.result = ""
}

func (m ActionsModel) inputCountFor(a action) int {
    switch a {
    case actionSendMessage:
        return 2
    case actionSetStatus:
        return 1
    case actionKick, actionBan:
        return 3
    }
    return 1
}

func (m ActionsModel) inputStartFor(a action) int {
    switch a {
    case actionSendMessage:
        return 0
    case actionSetStatus:
        return 3
    case actionKick, actionBan:
        return 4
    }
    return 0
}

func (m ActionsModel) executeAction() tea.Cmd {
    b := m.bot
    switch m.cursor {
    case actionSendMessage:
        channelID := strings.TrimSpace(m.inputs[0].Value())
        content   := strings.TrimSpace(m.inputs[1].Value())
        return func() tea.Msg {
            if err := b.SendMessage(channelID, content); err != nil {
                return actionResult{err: err}
            }
            return actionResult{msg: "Message sent!"}
        }
    case actionSetStatus:
        text := strings.TrimSpace(m.inputs[3].Value())
        return func() tea.Msg {
            if err := b.SetStatus(discordgo.ActivityTypeGame, text); err != nil {
                return actionResult{err: err}
            }
            return actionResult{msg: "Status updated!"}
        }
    case actionKick:
        guildID := strings.TrimSpace(m.inputs[4].Value())
        userID  := strings.TrimSpace(m.inputs[5].Value())
        reason  := strings.TrimSpace(m.inputs[6].Value())
        return func() tea.Msg {
            if err := b.KickMember(guildID, userID, reason); err != nil {
                return actionResult{err: err}
            }
            return actionResult{msg: "Member kicked."}
        }
    case actionBan:
        guildID := strings.TrimSpace(m.inputs[4].Value())
        userID  := strings.TrimSpace(m.inputs[5].Value())
        reason  := strings.TrimSpace(m.inputs[6].Value())
        return func() tea.Msg {
            if err := b.BanMember(guildID, userID, 1, reason); err != nil {
                return actionResult{err: err}
            }
            return actionResult{msg: "Member banned."}
        }
    }
    return nil
}

func (m ActionsModel) SetSize(w, h int) ActionsModel {
    m.width = w
    m.height = h
    return m
}

func (m ActionsModel) View() string {
    title := styleTitle.Render("actions")

    var actionList string
    for i, label := range actionLabels {
        if action(i) == m.cursor {
            actionList += styleSelected.Render(fmt.Sprintf("> %s", label)) + "\n"
        } else {
            actionList += styleItem.Render(fmt.Sprintf("  %s", label)) + "\n"
        }
    }

    listBox := styleBorder.Width(24).Render(actionList)

    startIdx := m.inputStartFor(m.cursor)
    count    := m.inputCountFor(m.cursor)
    var formFields string
    for i := 0; i < count; i++ {
        formFields += m.inputs[startIdx+i].View() + "\n"
    }

    hint   := styleHelp.Render("enter: next/submit  esc: reset  down/up: change action")
    formW  := m.width - 32
    if formW < 20 {
        formW = 20
    }
    formBox := styleBorder.Width(formW).Render(
        styleTitle.Render(actionLabels[m.cursor]) + "\n\n" +
        formFields + "\n" + hint,
    )

    row := lipgloss.JoinHorizontal(lipgloss.Top, listBox, formBox)

    var resultLine string
    if m.result != "" {
        if m.resultIsErr {
            resultLine = styleLogError.Render("x " + m.result)
        } else {
            resultLine = styleLogInfo.Render("Y " + m.result)
        }
    }

    return lipgloss.JoinVertical(lipgloss.Left, title, row, resultLine)
}
