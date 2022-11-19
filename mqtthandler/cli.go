package mqtthandler

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/model/cachestorage"
	"fmt"
	"hash/crc64"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const user = "etri"
const passwd = "etrismartfarm"

var client mqtt.Client

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	obj := map[string]interface{}{}
	err := json.Unmarshal(msg.Payload(), &obj)
	if err != nil {
		return
	}

	code, ok := obj["code"].(float64)
	if !ok {
		return
	}

	fmt.Println(code)
	if code-2.0 < 0.001 {
		ctrlName := msg.Topic()[strings.Index(msg.Topic(), "/")+1:]
		fmt.Println(ctrlName)
		if !ok {
			return
		}

		cmd, ok := obj["cmd"].(string)
		if !ok {
			return
		}

		ctrl, err := cachestorage.GetDeviceController(crc64.Checksum([]byte(ctrlName), crc64.MakeTable(crc64.ISO)))
		if err != nil {
			return
		}
		fmt.Println(cmd)
		ctrl.Sync([]byte(cmd))
	}

}

func ConnectMQTT(mqttAddr string) error {
	cid, ok := config.Params["cid"].(string)
	if !ok {
		return errors.New("invalid cid error")
	}
	opts := mqtt.NewClientOptions().AddBroker(mqttAddr).SetClientID(cid)

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(user)
	opts.SetPassword(passwd)

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func Subscribe(topic string) error {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func Unsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}
