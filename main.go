package main

import (
	"encoding/json"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanage"
	"etri-sfpoc-controller/router"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func subscribe() {
	var addr = flag.String("addr", "localhost:3000", "http service address")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/push/v1"}
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
			log.Printf("recv: %s", message)
			obj := map[string]interface{}{}
			err = json.Unmarshal(message, &obj)
			if err != nil {
				log.Println("read:", err)
				return
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
				return
			}
		}
	}
}
func main() {
	run, cancel := devmanage.NewManager()

	go run()
	router.NewRouter().Run(config.Params["bind"].(string))

	// subscribe()

	cancel()
}
