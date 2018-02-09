package pusher_api_client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	PusherAPIURL = "ws://ws-%s.pusher.com:80/app/%s?client=%s&version=%s&protocol=7"
)

type (
	Client struct {
		Debug bool

		AppID         string
		Cluster       string
		ClientName    string
		ClientVersion string

		connection  *websocket.Conn
		subscribers map[Event][]chan Message
	}
)

func (cl *Client) Connect() error {
	u := fmt.Sprintf(PusherAPIURL, cl.Cluster, cl.AppID, cl.ClientName, cl.ClientVersion)

	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return err
	}

	cl.connection = conn

	// go client.pingPong()
	go cl.pollMessages()

	return nil
}

func (client *Client) Close() {
	client.connection.Close()
}

func (client *Client) write(msg Message) error {
	if client.Debug {
		log.Printf("send: %+v\n", msg)
	}

	return client.connection.WriteJSON(msg)
}

func (c *Client) Subscribe(event Event) chan Message {
	if c.subscribers == nil {
		c.subscribers = make(map[Event][]chan Message)
	}

	ch := make(chan Message)

	c.subscribers[event] = append(c.subscribers[event], ch)

	return ch
}

// func (client *PusherAPIClient) pingPong() {
// 	connectionEstablishedMessage := <-client.createChannel("pusher:connection_established")

// 	var msg PusherAPIMessageConnectionEstablished
// 	err := connectionEstablishedMessage.UnmarshalData(&msg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	interval := time.Minute
// 	if msg.ActivityTimeout > 0 {
// 		interval = time.Duration(msg.ActivityTimeout) * time.Second
// 	}

// 	pongMessage := &PusherAPIMessage{"pusher:pong", "", nil}

// 	for range time.Tick(interval) {
// 		err := client.write(pongMessage)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

func (client *Client) pollMessages() {
	for {
		var msg Message
		err := client.connection.ReadJSON(&msg)
		if err != nil {
			log.Fatal(err)
		}

		if client.Debug {
			log.Printf("received: %+v\n", msg)
		}

		if msg.IsEvent(ErrorEvent) {
			var errorMessage ErrorMessage
			err = msg.UnmarshalData(errorMessage)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("got error %d: %s", errorMessage.Code, errorMessage.Message)
		}

		if subscribers, found := client.subscribers[msg.GetEvent()]; found {
			for i, subscriber := range subscribers {
				subscriber <- msg

				if client.Debug {
					log.Printf("pushed to %s #%d", msg.GetEvent(), i)
				}
			}
		}
	}
}
