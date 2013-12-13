package hustle

import (
	"encoding/json"
)

// Event is all that I know about events so far.  dang.
type Event struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     []byte   `json:"data"`
}

type eventPayload struct {
	Event    string `json:"event"`
	Data     []byte `json:"data"`
	Channel  string `json:"channel"`
	SocketID string `json:"socket_id"`
}

func newEventPayload(channel, name, socketID string, data []byte) (string, error) {
	jsonBytes, err := json.Marshal(&eventPayload{
		Event:    name,
		Channel:  channel,
		Data:     data,
		SocketID: socketID,
	})
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
