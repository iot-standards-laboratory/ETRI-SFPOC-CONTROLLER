package utils

import (
	"errors"
	"etri-sfpoc-controller/config"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func CMD_Init() error {
	id, ok := config.Params["id"].(string)
	if !ok {
		return errors.New("invalid agent id error")
	}
	edgeAddress, ok := config.Params["edgeAddress"]
	if !ok {
		return errors.New("invalid edge address error")
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/api/v2/agents", edgeAddress), nil)
	if err != nil {
		return err
	}

	req.Header.Add("agent_id", id)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(msg))
	}

	os.Remove("./config.properties")
	return nil
}

func CMD_Reboot() error {
	return nil
}
