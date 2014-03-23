package hustle

// Event is all that I know about events so far.  dang.
type Event struct {
	Name     string      `json:"name"`
	Channels []string    `json:"channels"`
	Data     interface{} `json:"data"`
}

// Payloads transforms an event into a slice of payloads,
// one for each channel
func (evt *Event) Payloads(socketID string) []*eventPayload {
	ret := []*eventPayload{}
	for _, channel := range evt.Channels {
		ret = append(ret, &eventPayload{
			Event:    evt.Name,
			Channel:  channel,
			Data:     evt.Data,
			SocketID: socketID,
		})
	}
	return ret
}

type eventPayload struct {
	Event    string      `json:"event,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Channel  string      `json:"channel,omitempty"`
	SocketID string      `json:"socket_id,omitempty"`
}
