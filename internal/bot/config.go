package bot

type VoiceConfig struct {
	ModChannelID    string
	CreateChannelID string
	ExcludedUsers   map[string]bool
}
