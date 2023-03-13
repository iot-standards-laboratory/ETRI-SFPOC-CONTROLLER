package mqtthandler

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/model/cachestorage"
	"etri-sfpoc-controller/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const user = "etri"
const passwd = "etrismartfarm"

var client mqtt.Client

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if strings.HasSuffix(msg.Topic(), "/init") {
		log.Println("init!!")
		err := utils.CMD_Init()
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			time.Sleep(time.Second * 2)
			os.Exit(0)
		}()
	}

	if strings.HasSuffix(msg.Topic(), "/reboot") {
		log.Println("reboot!!")
		err := utils.CMD_Reboot()
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			time.Sleep(time.Second * 2)
			os.Exit(0)
		}()
	}

	ctrlKey := msg.Topic()[strings.Index(msg.Topic(), "/")+1:]
	if ctrlKey[len(ctrlKey)-1] != 'c' {
		return
	}

	obj := map[string]interface{}{}
	err := json.Unmarshal(msg.Payload(), &obj)
	if err != nil {
		return
	}

	code, ok := obj["code"].(float64)
	if !ok {
		return
	}

	if code-2.0 < 0.001 {

		ctrlKey = ctrlKey[:len(ctrlKey)-1]

		cmd, ok := obj["cmd"].(string)
		if !ok {
			return
		}

		nKey, err := strconv.ParseInt(ctrlKey, 0, 64)
		if err != nil {
			return
		}
		ctrl, err := cachestorage.GetDeviceController(uint64(nKey))
		if err != nil {
			fmt.Println(err)
			return
		}

		ctrl.Do(2, []byte(cmd))
	}

}

func ConnectMQTT(mqttAddr string) error {
	id, ok := config.Params["id"].(string)
	if !ok {
		return errors.New("invalid id error")
	}
	opts := mqtt.NewClientOptions().AddBroker(mqttAddr).SetClientID(id)

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername(user)
	opts.SetPassword(passwd)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(time.Duration(time.Second * 5))
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("connect!!")
		id, ok := config.Params["id"]
		if !ok {
			os.Exit(-1)
		}
		err := Subscribe(fmt.Sprintf("%s/#", id))
		if err != nil {
			os.Exit(-1)
		}
	})

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func Publish(topic string, payload []byte) error {
	tkn := client.Publish(topic, 0, false, payload)
	return tkn.Error()
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
