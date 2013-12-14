package hustle

type wsMessage struct {
	Event    string         `json:"event"`
	SocketID string         `json:"socket_id"`
	Data     *wsMessageData `json:"data"`
}

type wsMessageData struct {
	Channel     string      `json:"channel"`
	Auth        string      `json:"auth"`
	ChannelData interface{} `json:"channel_data"`
}

func newWsMessage() *wsMessage {
	return &wsMessage{
		Data: &wsMessageData{},
	}
}
