package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) SendMessage(channelID, content string) error {
	if b.Status() != StatusOnline {
		return fmt.Errorf("bot is offline")
	}
	_, err := b.session.ChannelMessageSend(channelID, content)
	if err != nil {
		b.addLog("ERROR", fmt.Sprintf("SendMessage to %s: %v", channelID, err))
		return err
	}
	b.addLog("INFO", fmt.Sprintf("Sent message to channel %s", channelID))
	return nil
}

func (b *Bot) SetStatus(activityType discordgo.ActivityType, text string) error {
	err := b.session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{Name: text, Type: activityType},
		},
		Status: "online",
	})
	if err != nil {
		b.addLog("ERROR", fmt.Sprintf("SetStatus: %v", err))
		return err
	}
	b.addLog("INFO", fmt.Sprintf("Status updated: %s", text))
	return nil
}

func (b *Bot) GetChannels(guildID string) ([]*discordgo.Channel, error) {
	channels, err := b.session.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("fetching channels: %w", err)
	}
	var text []*discordgo.Channel
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildText {
			text = append(text, ch)
		}
	}
	return text, nil
}

func (b *Bot) KickMember(guildID, userID, reason string) error {
	err := b.session.GuildMemberDeleteWithReason(guildID, userID, reason)
	if err != nil {
		b.addLog("ERROR", fmt.Sprintf("KickMember %s: %v", userID, err))
		return err
	}
	b.addLog("WARN", fmt.Sprintf("Kicked member %s from guild %s (reason: %s)", userID, guildID, reason))
	return nil
}

func (b *Bot) BanMember(guildID, userID string, deleteMessageDays int, reason string) error {
	err := b.session.GuildBanCreateWithReason(guildID, userID, reason, deleteMessageDays)
	if err != nil {
		b.addLog("ERROR", fmt.Sprintf("BanMember %s: %v", userID, err))
		return err
	}
	b.addLog("WARN", fmt.Sprintf("Banned member %s from guild %s (reason: %s)",userID, guildID, reason))
	return nil
}
