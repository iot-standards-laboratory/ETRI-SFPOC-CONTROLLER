package statmgmt

import (
	"bytes"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/mqtthandler"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

const (
	STATUS_INIT = iota
	STATUS_DISCONNECTED
	STATUS_CONNECTED
)

var status = STATUS_INIT
var bootupChannel chan interface{}

func init() {
	bootupChannel = make(chan interface{})
}
func Status() int {
	return status
}

func Bootup() error {
	id := config.Params["id"].(string)
	if strings.Compare(id, "blank") == 0 {
		if strings.Compare(config.Params["mode"].(string), string(config.MANAGEDBYEDGE)) == 0 {
			status = STATUS_INIT
			<-bootupChannel
			return nil
		}

		_uuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		id = _uuid.String()
		config.Set("id", id)
	}
	status = STATUS_DISCONNECTED
	return nil
}

func Register(accessToken string) error {
	// key is user's access token
	payload := map[string]interface{}{
		"name": config.Params["name"],
		"key":  accessToken,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("http://%s/%s", config.Params["edgeAddress"], "api/v2/agents"),
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// 응답 메시지 수신
	var body = map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&body)

	// 등록 후 생성된 Agent ID 저장
	config.Set("id", body["id"].(string))

	status = STATUS_DISCONNECTED
	bootupChannel <- struct{}{}
	return nil
}

func Connect() error {
	id, ok := config.Params["id"]
	if !ok {
		return errors.New("id is invalid error")
	}

	edgeAddress, ok := config.Params["edgeAddress"]
	if !ok {
		return errors.New("server address is invalid error")
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/api/v2/agents/%s", edgeAddress, id),
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}

	var body = map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&body)
	if err != nil {
		return err
	}

	// fmt.Println(body)
	// wsAddr, ok := body["wsAddr"].(string)
	// if !ok {
	// 	return errors.New("invalid params")
	// }

	consulAddr, ok := body["consulAddr"].(string)
	if !ok {
		return errors.New("invalid params")
	}

	mqttAddr, ok := body["mqttAddr"].(string)
	if !ok {
		return errors.New("invalid params")
	}

	origin, ok := body["origin"].(string)
	if !ok {
		return errors.New("invalid params")
	}

	var i = 0
	for i = 0; i < 10; i++ {
		err := connect(consulAddr, mqttAddr, origin)
		if err == nil {
			err = nil
			break
		}
		time.Sleep(time.Second * 3)
	}

	if err != nil {
		return err
	}

	return nil
}

func connect(consulAddr, mqttAddr, origin string) error {
	// err := connectCentrifuge(wsAddr)
	// if err != nil {
	// 	return err
	// }

	err := connectConsul(consulAddr, origin)
	if err != nil {
		return err
	}

	err = mqtthandler.ConnectMQTT(mqttAddr)
	if err != nil {
		return err
	}

	id, ok := config.Params["id"]
	if !ok {
		return errors.New("invalid id error")
	}
	err = mqtthandler.Subscribe(fmt.Sprintf("%s/#", id))
	if err != nil {
		return err
	}

	status = STATUS_CONNECTED
	return nil
}
