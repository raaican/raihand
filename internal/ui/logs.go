package ui

import (
    "fmt"
    "strings"

    tea "charm.land/bubbletea/v2"
    "charm.land/bubbles/v2/viewport"
    "github.com/raaican/raihand/internal/bot"
)

type logEntryMsg bot.LogEntry

type LogViewModel struct {
    bot      *bot.Bot
    viewport viewport.Model
    sub      chan bot.LogEntry
    width    int
    height   int
    lines    []string
}

func NewLogViewModel(b *bot.Bot) LogViewModel {
    vp := viewport.New()
    return LogViewModel{
        bot:      b,
        viewport: vp,
        sub:      b.SubscribeLogs(),
    }
}

func (m LogViewModel) Init() tea.Cmd {
    for _, e := range m.bot.Logs() {
        m.lines = append(m.lines, renderLogLine(e))
    }
    m.viewport.SetContent(strings.Join(m.lines, "\n"))
    m.viewport.GotoBottom()
    return waitForLog(m.sub)
}

func waitForLog(sub chan bot.LogEntry) tea.Cmd {
    return func() tea.Msg {
        return logEntryMsg(<-sub)
    }
}

func (m LogViewModel) Update(msg tea.Msg) (LogViewModel, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case logEntryMsg:
        m.lines = append(m.lines, renderLogLine(bot.LogEntry(msg)))
        if len(m.lines) > 1000 {
            m.lines = m.lines[len(m.lines)-1000:]
        }
        m.viewport.SetContent(strings.Join(m.lines, "\n"))
        m.viewport.GotoBottom()
        cmds = append(cmds, waitForLog(m.sub))

    case tea.KeyPressMsg:
        switch msg.String() {
        case "c":
            m.lines = nil
            m.viewport.SetContent("")
        }
    }

    var vpCmd tea.Cmd
    m.viewport, vpCmd = m.viewport.Update(msg)
    cmds = append(cmds, vpCmd)
    return m, tea.Batch(cmds...)
}

func (m LogViewModel) SetSize(w, h int) LogViewModel {
    m.width = w
    m.height = h
    // v2 viewport: SetWidth / SetHeight methods
    m.viewport.SetWidth(w - 4)
    m.viewport.SetHeight(h - 4)
    return m
}

func (m LogViewModel) View() string {
    header := styleTitle.Render("log monitor")
    help   := styleHelp.Render("↑/↓: scroll  c: clear")
    count  := styleLogTime.Render(fmt.Sprintf("%d entries", len(m.lines)))
    top    := fmt.Sprintf("%s  %s", header, count)
    return fmt.Sprintf("%s\n%s\n%s",
        top,
        styleBorder.Width(m.width-4).Render(m.viewport.View()),
        help,
    )
}

func renderLogLine(e bot.LogEntry) string {
    ts := styleLogTime.Render(e.Timestamp.Format("15:04:05"))
    var level string
    switch e.Level {
    case "INFO":
        level = styleLogInfo.Render("[INFO]")
    case "WARN":
        level = styleLogWarn.Render("[WARN]")
    case "ERROR":
        level = styleLogError.Render("[ERROR]")
    default:
        level = e.Level
    }
    return fmt.Sprintf("%s %s %s", ts, level, e.Message)
}
