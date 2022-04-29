package commonutils

import (
	"context"
	"etri-sfpoc-controller/config"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func Subscribe(ctx context.Context, path, token string, handler func(parmas []byte), disconnectedHandler func()) {
	fmt.Println("token: ", token)
	var addr = config.Params["serverAddr"].(string)

	u := url.URL{Scheme: "ws", Host: addr, Path: path + token}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			if handler != nil {
				handler(message)
			}
		}
	}()

	for {
		select {
		case <-done:
			if disconnectedHandler != nil {
				disconnectedHandler()
			}
			return

		case <-ctx.Done():
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			return
		}
	}
}
