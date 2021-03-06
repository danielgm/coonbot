package slack

type IMCreatedEvent struct {
	Type    string             `json:"type"`
	User    string             `json:"user"`
	Channel ChannelCreatedInfo `json:"channel"`
}

type IMHistoryChangedEvent ChannelHistoryChangedEvent
type IMOpenEvent ChannelInfoEvent
type IMCloseEvent ChannelInfoEvent
type IMMarkedEvent ChannelInfoEvent
type IMMarkedHistoryChanged ChannelInfoEvent
