package hustle

import (
	"github.com/garyburd/redigo/redis"
)

type hub struct {
	addr  string
	redis redis.Conn
}

func newHub(addr string) (*hub, error) {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &hub{
		addr:  addr,
		redis: conn,
	}, nil
}

func (h *hub) PublishEvent(channel, name, socketID string, data []byte) (interface{}, error) {
	payload, err := newEventPayload(channel, name, socketID, data)
	if err != nil {
		return nil, err
	}

	return h.redis.Do("PUBLISH", channel, payload)
}
