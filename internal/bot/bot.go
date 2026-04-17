package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Status string

const (
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

type Bot struct {
	mu      sync.RWMutex
	session *discordgo.Session
	status  Status
	logs    []LogEntry
	guilds  []*discordgo.Guild
	token   string
	logSubs []chan LogEntry
}

func New (token string) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}
	dg.Identify.Intents = discordgo.IntentsAll

	b := &Bot{
		token:   token,
		session: dg,
		status:  StatusOffline,
	}

	dg.AddHandler(b.onReady)
	dg.AddHandler(b.onMessage)
	dg.AddHandler(b.onGuildCreate)

	return b, nil
}

func (b *Bot) Connect() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("opening connection: %w", err)
	}
	return nil
}

func (b *Bot) Disconnect() error {
	b.mu.Lock()
	b.status = StatusOffline
	b.mu.Unlock()
	return b.session.Close()
}

func (b *Bot) Status() Status {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.status
}

func (b *Bot) Guilds() []*discordgo.Guild {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.guilds
}

func (b *Bot) Logs() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	cp := make([]LogEntry, len(b.logs))
	copy(cp, b.logs)
	return cp
}

func (b *Bot) SubscribeLogs() chan LogEntry {
	ch := make(chan LogEntry, 64)
	b.mu.Lock()
	b.logSubs = append(b.logSubs, ch)
	b.mu.Unlock()
	return ch
}

func (b *Bot) UnsubscribeLogs(ch chan LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, sub := range b.logSubs {
		if sub == ch {
			b.logSubs = append(b.logSubs[:i], b.logSubs[i+1:]...)
			close(ch)
			return
		}
	}
}

func (b *Bot) addLog(level, msg string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
	}
	b.mu.Lock()
	b.logs = append(b.logs, entry)
	if len(b.logs) > 500 {
		b.logs = b.logs[len(b.logs)-500:]
	}
	subs := make([]chan LogEntry, len(b.logSubs))
	copy(subs, b.logSubs)
	b.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- entry:
		default:
		}
	}
}

func (b *Bot) onReady(s *discordgo.Session, e *discordgo.Ready) {
	b.mu.Lock()
	b.status = StatusOnline
	b.guilds = e.Guilds
	b.mu.Unlock()
	b.addLog("INFO", fmt.Sprintf("Logged in as %s#%s", e.User.Username, e.User.Discriminator))
}

func (b *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	b.addLog("INFO", fmt.Sprintf("[MSG] %s: %s", m.Author.Username, m.Content))
}

func (b *Bot) onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	b.mu.Lock()
	found := false
	for i, guild := range b.guilds {
		if guild.ID == g.ID {
			b.guilds[i] = g.Guild
			found = true
			break
		}
	}
	if !found {
		b.guilds = append(b.guilds, g.Guild)
	}
	b.mu.Unlock()
	b.addLog("INFO", fmt.Sprintf("Guild available: %s (%s)", g.Name, g.ID))
}
