package statmgmt

import (
	"bytes"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
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

func Status() int {
	return status
}

func Bootup() error {
	cid := config.Params["cid"].(string)
	if strings.Compare(cid, "blank") == 0 {
		if strings.Compare(config.Params["mode"].(string), string(config.MANAGEDBYEDGE)) == 0 {
			status = STATUS_INIT
			return nil
		}

		_uuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		cid = _uuid.String()
		config.Set("cid", cid)
	}
	status = STATUS_DISCONNECTED
	return nil
}

func Register(accessToken string) error {
	// key is user's access token
	payload := map[string]interface{}{
		"cname": config.Params["cname"],
		"key":   accessToken,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("http://%s/%s", config.Params["serverAddr"], "api/v2/ctrls"),
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

	// 등록 후 생성된 Controller ID 저장
	config.Set("cid", body["cid"].(string))

	status = STATUS_DISCONNECTED
	return nil
}

func Connect() error {
	cid, ok := config.Params["cid"]
	if !ok {
		return errors.New("cid is invalid error")
	}

	serverAddr, ok := config.Params["serverAddr"]
	if !ok {
		return errors.New("server address is invalid error")
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s/api/v2/ctrls/%s", serverAddr, cid),
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

	fmt.Println(body)
	wsAddr, ok := body["wsAddr"].(string)
	if !ok {
		return errors.New("invalid params")
	}

	consulAddr, ok := body["consulAddr"].(string)
	if !ok {
		return errors.New("invalid params")
	}

	var i = 0
	for i = 0; i < 10; i++ {
		err := connect(wsAddr, consulAddr)
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

func connect(wsAddr string, consulAddr string) error {
	err := connectCentrifuge(wsAddr)
	if err != nil {
		return err
	}
	err = connectConsul(consulAddr)
	if err != nil {
		return err
	}

	status = STATUS_CONNECTED
	return nil
}
