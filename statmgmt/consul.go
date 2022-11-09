package statmgmt

import (
	"context"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/consul_api"
	"etri-sfpoc-controller/model"
	"fmt"
)

func connectConsul(consulAddr string) error {
	cid, ok := config.Params["cid"].(string)
	if !ok {
		return errors.New("cid is invalid error")
	}

	key := fmt.Sprintf("ctrl/%s", cid)
	ctrl := model.Controller{
		CID: key,
	}

	err := consul_api.Connect(consulAddr)
	if err != nil {
		return err
	}

	err = consul_api.RegisterCtrl(ctrl, "http://localhost:8080")
	if err != nil {
		return err
	}

	go consul_api.UpdateTTL(func() (bool, error) { return true, nil }, key)

	go consul_api.Monitor(func(s string) { fmt.Println(s) }, context.Background())

	return nil
}
