package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) RegisterVoiceMonitor(cfg VoiceConfig) {
	b.voiceCfg = cfg
	b.session.AddHandler(b.onVoiceStateUpdate)
}

func (b *Bot) onVoiceStateUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	cfg := b.voiceCfg

	if cfg.ExcludedUsers[e.UserID] {
		return
	}

	member, err := s.GuildMember(e.GuildID, e.UserID)
	if err != nil {
		b.addLog("ERROR", fmt.Sprintf("VoiceState: could not fetch member %s: %v", e.UserID, err))
		return
	}

	mention := "<@" + e.UserID + ">"
	_ = member // Nickname placeholder

	before := e.BeforeUpdate
	after := e.VoiceState

	send := func(msg string) {
		if cfg.ModChannelID == "" {
			b.addLog("WARN", "VoiceState: ModChannelID not configured")
			return
		}
		if _, err := s.ChannelMessageSend(cfg.ModChannelID, msg); err != nil {
			b.addLog("ERROR", fmt.Sprintf("VoiceState send: %v", err))
		}
		b.addLog("INFO", fmt.Sprintf("[VOICE] %s", msg))
	}

	beforeCh := ""
	if before != nil {
		beforeCh = before.ChannelID
	}
	afterCh := after.ChannelID

	switch {
	case beforeCh == "" && afterCh != "":
		if afterCh == cfg.CreateChannelID {
			return // VoiceMaster
		}
		send(fmt.Sprintf("%s joined <#%s>", mention, afterCh))

	// left a channel
	case beforeCh != "" && afterCh == "":
		if beforeCh == cfg.CreateChannelID {
			return //VoiceMaster
		}
		send(fmt.Sprintf("%s left <#%s>", mention, beforeCh))

	// moved or changed state
	case beforeCh != "" && afterCh != "":
		switch {
		case beforeCh == cfg.CreateChannelID:
			// created a channel via VoiceMaster
			send(fmt.Sprintf("%s create a channel > <#%s>", mention, afterCh))
		case afterCh != beforeCh && afterCh != cfg.CreateChannelID:
			// switched channels
			send(fmt.Sprintf("%s switched from <#%s> to <#%s>", mention, beforeCh, afterCh))

		case before != nil && !before.SelfMute && after.SelfMute:
			// muted
			send(fmt.Sprintf("%s muted", mention))

		case before != nil && before.SelfMute && !after.SelfMute:
			// unmuted
			send(fmt.Sprintf("%s unmuted", mention))
		}
	}
}
