package main

import (
	"bytes"
	"encoding/json"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanage"
	"etri-sfpoc-controller/notifier"
	"etri-sfpoc-controller/router"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func subscribe(token string) {
	fmt.Println("token: ", token)
	var addr = flag.String("addr", config.Params["serverAddr"].(string), "http service address")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/push/v1/" + token}
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

			notifier.Box.Publish(
				notifier.NewEvent(
					"title",
					obj,
					notifier.SubtokenRcvCtrlMsg,
				),
			)
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

func register() (string, error) {
	// Controller 이름을 읽어옴
	payload := map[string]interface{}{}
	payload["cname"] = config.Params["cname"]
	fmt.Println("cname: ", payload["cname"])
	fmt.Println(config.Params["cname"])
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Controller 등록 메시지 송신
	resp, err := http.Post(
		fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v1/ctrls"),
		"application/json",
		bytes.NewReader(b),
	)

	if err != nil {
		return "", err
	}

	// 응답 메시지 수신
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &payload)

	// 등록 후 생성된 Controller ID 저장
	config.Set("cid", payload["cid"].(string))

	return payload["cid"].(string), nil
}
func main() {
	cid, _ := config.Params["cid"]
	if cid == "blank" {
		var err error
		cid, err = register()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("cid: ", cid)
	}

	run, cancel := devmanage.NewManager()
	go run()
	go router.NewRouter().Run(config.Params["bind"].(string))
	subscribe(cid.(string))
	cancel()
}
