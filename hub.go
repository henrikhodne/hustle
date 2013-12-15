package hustle

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/garyburd/redigo/redis"
)

type hub struct {
	addr     string
	pubRedis redis.Conn
	subRedis redis.Conn
	subs     map[string]chan bool
	subsLock *sync.Mutex
}

func newHub(addr string) (*hub, error) {
	pubConn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	subConn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &hub{
		addr:     addr,
		pubRedis: pubConn,
		subRedis: subConn,
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

	return h.pubRedis.Do("PUBLISH", ep.Channel, payloadBytes)
}

func (h *hub) Subscribe(channelID, subscriptionID string) (chan *wsMessage, chan bool) {
	if channelID == "" {
		log.Println("channel id not present")
		return nil, nil
	}

	if subscriptionID == "" {
		log.Println("subscription id not present")
		return nil, nil
	}

	h.subsLock.Lock()
	defer h.subsLock.Unlock()

	if _, ok := h.subs[subscriptionID]; ok {
		log.Printf("subscription %s already present\n", subscriptionID)
		return nil, nil
	}

	subQuitChan := make(chan bool)
	h.subs[subscriptionID] = subQuitChan
	log.Printf("subscription %s added\n", subscriptionID)

	subChan := make(chan *wsMessage)

	go func() {
		psc := &redis.PubSubConn{h.subRedis}
		psc.Subscribe(channelID)

		log.Printf("subscription %s waiting for messages on channel %s",
			subscriptionID, channelID)

		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Printf("subscription %s received redis message on channel %s\n",
					subscriptionID, channelID)
				msg := newWsMessage()
				err := json.Unmarshal(v.Data, msg)
				if err != nil {
					log.Printf("error unmarshaling message: %v\n", err)
				} else {
					log.Printf("subscription %s sending message to subChan: %#v\n", subscriptionID, msg)
					subChan <- msg
				}
			case error:
				log.Printf("subscription %s error on channel %s\n",
					subscriptionID, channelID, v)
			default:
				log.Printf("subscription %s received unknown on channel %s: %#v\n",
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
