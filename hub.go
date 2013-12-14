package hustle

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type hub struct {
	addr     string
	redis    redis.Conn
	subs     map[string]chan bool
	subsLock *sync.Mutex
}

func newHub(addr string) (*hub, error) {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &hub{
		addr:     addr,
		redis:    conn,
		subs:     make(map[string]chan bool),
		subsLock: &sync.Mutex{},
	}, nil
}

func (h *hub) PublishEvent(ep *eventPayload) (interface{}, error) {
	log.Printf("hub publishing event %#v to channel %q\n", ep, ep.Channel)
	payloadBytes, err := json.Marshal(ep)
	if err != nil {
		return nil, err
	}

	return h.redis.Do("PUBLISH", ep.Channel, payloadBytes)
}

func (h *hub) Subscribe(channelID, subscriptionID string) (chan redis.Message, chan bool) {
	h.subsLock.Lock()
	defer h.subsLock.Unlock()

	if _, ok := h.subs[subscriptionID]; ok {
		log.Printf("subscription %s already present\n", subscriptionID)
		return nil, nil
	}

	subQuitChan := make(chan bool)
	h.subs[subscriptionID] = subQuitChan

	subChan := make(chan redis.Message)

	go func() {
		psc := &redis.PubSubConn{h.redis}
		psc.Subscribe(channelID)

		for {
			select {
			case <-subQuitChan:
				return
			}

			switch v := psc.Receive().(type) {
			case redis.Message:
				subChan <- v
			case error:
				log.Printf("subscription %s error on channel %s\n",
					subscriptionID, channelID, v)
			}
		}
	}()

	return subChan, subQuitChan
}

func (h *hub) Unsubscribe(subscriptionID string) {
	h.subsLock.Lock()
	defer h.subsLock.Unlock()

	if subQuitChan, ok := h.subs[subscriptionID]; ok {
		subQuitChan <- true
		delete(h.subs, subscriptionID)
	}
}
