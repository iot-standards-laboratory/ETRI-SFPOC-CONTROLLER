package common

import (
	"encoding/json"
	"etri-sfpoc-controller/config"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func Subscribe(path, token string, handler func(parmas []byte)) {
	fmt.Println("token: ", token)
	var addr = config.Params["serverAddr"].(string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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

			obj := map[string]interface{}{}
			err = json.Unmarshal(message, &obj)
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("recv: %s", obj["value"].(map[string]interface{})["did"].(string))
			if handler != nil {
				handler(message)
			}
		}
	}()

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			os.Exit(0)
		}
	}
}
