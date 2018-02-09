package pusher_api_client

import (
	"encoding/json"
)

type Event string

const (
	ConnectionEstablishedEvent Event = "pusher:connection_established"
	ErrorEvent                 Event = "pusher:error"

	SubscribeEvent Event = "pusher:subscribe"
)

func (e Event) String() string {
	return string(e)
}

type Message struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

func (msg *Message) GetEvent() Event {
	return Event(msg.Event)
}

func (msg *Message) IsEvent(e Event) bool {
	return msg.Event == e.String()
}

func (msg *Message) UnmarshalData(v interface{}) error {
	return json.Unmarshal([]byte(msg.Data.(string)), v)
}

type ConnectionEstablishedMessage struct {
	SocketID        string `json:"socket_id"`
	ActivityTimeout int    `json:"activity_timeout"`
}

type ErrorMessage struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
