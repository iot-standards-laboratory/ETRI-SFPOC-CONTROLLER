package statmgmt

import (
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/centrifuge_api"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/model/cachestorage"
	"fmt"
	"hash/crc64"

	"github.com/centrifugal/centrifuge-go"
)

func connectCentrifuge(wsAddr string) error {
	cid, ok := config.Params["cid"]
	if !ok {
		return errors.New("cid is invalid error")
	}

	cname, ok := config.Params["cname"]
	if !ok {
		return errors.New("cname is invalid error")
	}

	token, _ := centrifuge_api.IssueJWT(cid.(string), cname.(string), "ctrl", "/", nil)

	connectedHandler := func(e centrifuge.ConnectedEvent) {
	}

	disconnectedHandler := func(e centrifuge.DisconnectedEvent) {
	}

	err := centrifuge_api.NewClient(
		wsAddr,
		token,
		connectedHandler,
		disconnectedHandler,
	)
	if err != nil {
		centrifuge_api.ResetClient()
		return err
	}

	sub, err := centrifuge_api.AddSubscription(fmt.Sprintf("public:%s", cid))
	if err != nil {
		return err
	}

	sub.OnPublication(func(e centrifuge.PublicationEvent) {
		// control device
		msg := map[string]interface{}{}
		err := json.Unmarshal(e.Data, &msg)
		if err != nil {
			return
		}

		code, ok := msg["code"].(float64)
		if !ok {
			return
		}

		fmt.Println(code)
		if code-2.0 < 0.001 {
			fmt.Println(msg)
			ctrlName, ok := msg["ctrlName"].(string)
			if !ok {
				return
			}

			fmt.Println(ctrlName)
			cmd, ok := msg["cmd"].(string)
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
	})

	err = sub.Subscribe()
	if err != nil {
		return err
	}

	return nil
}
