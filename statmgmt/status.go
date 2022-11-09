package statmgmt

import (
	"bytes"
	"encoding/json"
	"etri-sfpoc-controller/config"
	"fmt"
	"io/ioutil"
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
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(b, &payload)

	// 등록 후 생성된 Controller ID 저장
	config.Set("cid", payload["cid"].(string))

	status = STATUS_DISCONNECTED
	return nil
}

func Connect(wsAddr string, consulAddr string) error {
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
