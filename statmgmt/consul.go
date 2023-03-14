package statmgmt

import (
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/consulapi"
	"etri-sfpoc-controller/model"
	"fmt"
)

func connectConsul(consulAddr, origin string) error {
	id, ok := config.Params["id"].(string)
	if !ok {
		return errors.New("id is invalid error")
	}

	key := fmt.Sprintf("agent/%s", id)
	agent := model.Agent{
		ID: key,
	}

	err := consulapi.Connect(consulAddr)
	if err != nil {
		return err
	}

	err = consulapi.RegisterAgent(agent, fmt.Sprintf("http://%s:4000", origin))
	if err != nil {
		return err
	}

	go consulapi.UpdateTTL(func() (bool, error) { return true, nil }, key)

	// go consul_api.Monitor(func(s string) { fmt.Println(s) }, context.Background())

	return nil
}
