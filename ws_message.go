package hustle

type wsMessage struct {
	Event       string      `json:"event,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Channel     string      `json:"channel,omitempty"`
	SocketID    string      `json:"socket_id,omitempty"`
	ChannelData interface{} `json:"channel_data,omitempty"`
	Auth        string      `json:"auth,omitempty"`
}

func newWsMessage() *wsMessage {
	return &wsMessage{}
}
