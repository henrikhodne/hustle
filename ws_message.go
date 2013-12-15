package hustle

type wsMessage struct {
	Event       string      `json:"event,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Channel     string      `json:"channel,omitempty"`
	ChannelData interface{} `json:"channel_data,omitempty"`
	SocketID    string      `json:"socket_id,omitempty"`
	Auth        string      `json:"auth,omitempty"`
}

func newWsMessage() *wsMessage {
	return &wsMessage{}
}
