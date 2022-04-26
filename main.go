package main

import (
	"bytes"
	"encoding/json"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"etri-sfpoc-controller/logger"
	"etri-sfpoc-controller/notifier"
	"etri-sfpoc-controller/router"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

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

			obj := map[string]interface{}{}
			err = json.Unmarshal(message, &obj)
			if err != nil {
				log.Println("read:", err)
				return
			}

			// log.Printf("recv: %s", obj["value"].(map[string]interface{})["dname"].(string))

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

func deviceManagerSetup() {
	devmanager.AddOnDiscovered(func(port io.ReadWriter, sname, dname string) {
		// do register procedure
		fmt.Println(sname)
		port.Write([]byte(`{"code": 1, "token": "initial", "mode": 1}`))
	})

	
	go devmanager.Watch()
	// serialctrl.AddRecvListener(serialctrl.NewEventHandler(func(e serialctrl.Event) {
	// 	param := e.Params()
	// 	fmt.Println("RECV: ", e.Params())
	// 	sid, err := model.DB.GetSID(param["sname"].(string))
	// 	if err != nil {
	// 		log.Println(err)
	// 	}

	// 	if sid == "not installed service" || sid == "not exist service" {
	// 		return
	// 	}

	// 	b, _ := json.Marshal(param)
	// 	req, err := http.NewRequest("PUT", "http://"+config.Params["serverAddr"].(string)+"/svc/"+sid+"/api/v1", bytes.NewReader(b))
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	resp, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	dec := json.NewDecoder(resp.Body)
	// 	var respObj map[string]interface{}
	// 	err = dec.Decode(&respObj)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	fmt.Println(respObj)

	// 	// serialctrl.Sync("DEVICE-A-UUID", map[string]interface{}{"ctrlValue": 100})
	// }))

	// serialctrl.AddRegisterHandleFunc(func(e serialctrl.Event) {
	// 	param := e.Params()
	// 	payload := map[string]interface{}{"sname": param["sname"], "dname": param["uuid"], "type": "sensor"}
	// 	fmt.Println("payload : ", payload)
	// 	respCh := make(chan bool)
	// 	go devmanage.RegisterDevice(payload, respCh)
	// 	<-respCh
	// })

	// serialctrl.AddRemoveHandleFunc(func(e serialctrl.Event) {
	// 	param := e.Params()
	// 	fmt.Println(param["uuid"].(string), " is removed!!")
	// })
}

func main() {

	logger.Start()

	exitCh := make(chan interface{})
	go func() {
		time.Sleep(10 * time.Hour)
		exitCh <- true
	}()

	cid := config.Params["cid"].(string)
	if cid == "blank" {
		var err error
		cid, err = register()
		if err != nil {
			panic(err)
		}
	}

	// go subscribe(cid)
	deviceManagerSetup()
	router.NewRouter().Run(config.Params["bind"].(string))
}
