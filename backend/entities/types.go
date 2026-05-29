package entities

type ChannelType string

const (
	ChannelTypeTelegramBot ChannelType = "telegram"
	ChannelTypeSlackBot    ChannelType = "slack"
	ChannelTypeAudio       ChannelType = "audio"
	ChannelTypeVideo       ChannelType = "video"
)
