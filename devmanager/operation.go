package devmanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/model"
	"fmt"
	"io/ioutil"
	"net/http"
)

// send post message to notify new device is discovered
func RegisterDeviceToEdge(dev map[string]interface{}) (string, error) {

	dev["cid"] = config.Params["cid"]
	b, err := json.Marshal(dev)
	if err != nil {
		return "", err
	}

	fmt.Println(string(b))

	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/devs/discover", config.Params["serverAddr"].(string)), "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusCreated {
		return "", errors.New("it is failed to get permission from admin")
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(b, &dev)
	fmt.Println(dev)

	return dev["did"].(string), nil
}

// send post message to notify already registered device is reconnected
func PostDeviceToEdge(dev *model.Device) error {

	b, err := json.Marshal(dev)
	if err != nil {
		return err
	}

	fmt.Println("post:", string(b))

	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/devs", config.Params["serverAddr"].(string)), "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		return errors.New("it is failed to get permission from admin")
	}

	return nil
}
