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

const user = "etrimqtt"
const passwd = "fainal2311"

var client mqtt.Client

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	prefixLen := len(config.Params["id"].(string))
	if strings.HasSuffix(msg.Topic(), "/init") {
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
		err := utils.CMD_Reboot()
		if err != nil {
			log.Fatalln(err)
		}

		go func() {
			time.Sleep(time.Second * 2)
			os.Exit(0)
		}()
	}

	if strings.HasSuffix(msg.Topic(), "post") {
		totalLen := len(msg.Topic())
		ctrlKey := msg.Topic()[prefixLen+1 : totalLen-5]

		nKey, err := strconv.ParseUint(ctrlKey, 0, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		ctrl, err := cachestorage.GetDeviceController(uint64(nKey))
		if err != nil {
			log.Println(err)
			return
		}
		var m_payload map[string]interface{}
		err = json.Unmarshal(msg.Payload(), &m_payload)
		if err != nil {
			log.Println(err)
			return
		}
		payload, err := json.Marshal(map[string]interface{}{
			m_payload["path"].(string): m_payload["value"],
		})
		if err != nil {
			log.Println(err)
			return
		}

		Publish(msg.Topic()[:totalLen-4]+"content/actuator", payload)

		if len(msg.Payload()) == 0 {
			return
		}
		ctrl.Do(2, msg.Payload())
	}

	if strings.HasSuffix(msg.Topic(), "get") {
		totalLen := len(msg.Topic())
		ctrlKey := msg.Topic()[prefixLen+1 : totalLen-4]
		nKey, err := strconv.ParseUint(ctrlKey, 0, 64)
		if err != nil {
			return
		}
		ctrl, err := cachestorage.GetDeviceController(uint64(nKey))
		if err != nil {
			fmt.Println(err)
			return
		}

		code, payload, err := ctrl.Do(1, msg.Payload())
		if err != nil {
			fmt.Println(err)
			return
		}
		if code != 205 {
			fmt.Println("invalid response code error:", code)
			return
		}

		if len(payload) == 0 {
			fmt.Println("^^")
			return
		}
		Publish(msg.Topic()[:totalLen-3]+"content/"+string(msg.Payload()), payload)
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
		id, ok := config.Params["id"]
		if !ok {
			os.Exit(-1)
		}
		err := Subscribe(fmt.Sprintf("%s/init", id))
		if err != nil {
			os.Exit(-1)
		}
		err = Subscribe(fmt.Sprintf("%s/reboot", id))
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
