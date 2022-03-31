package devmanage

import (
	"bytes"
	"encoding/json"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/notifier"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func StatusReport(e notifier.IEvent) {
	fmt.Println(e.Body())
}

func StatusChanged(e notifier.IEvent) {
	fmt.Println("status chagned : ", e.Body())
}

func HandleCtrlMsg(e notifier.IEvent) {
	fmt.Println("HandleCtrlMsg: ", e.Body())
}

func RegisterDevice(dev map[string]interface{}, waitCh chan bool) {
	dev["cid"] = config.Params["cid"]
	b, err := json.Marshal(dev)
	if err != nil {
		log.Println(err)
		waitCh <- false
	}

	resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/devs/discover", config.Params["serverAddr"].(string)), "application/json", bytes.NewReader(b))
	if err != nil || resp.StatusCode != http.StatusCreated {
		log.Println(err)
		waitCh <- false
		return
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusCreated {
		log.Println(err)
		waitCh <- false
		return
	}

	json.Unmarshal(b, &dev)
	log.Println("registered] ", dev)
	waitCh <- true
}

// func RegisterDevice(payload map[string]interface{}, cancelCh <-chan struct{}) ([]byte, error) {
// 	payload["cid"] = config.Params["cid"].(string)
// 	b, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	req, err := http.NewRequestWithContext(
// 		ctx,
// 		"POST",
// 		"http://"+config.Params["serverAddr"].(string)+"/api/v1/devs",
// 		bytes.NewReader(b),
// 	)
// 	req.Header.Set("Content-Type", "application/json")
// 	var resp *http.Response
// 	done := make(chan bool)
// 	go func() {
// 		resp, err = http.DefaultClient.Do(req)
// 		if resp == nil {
// 			return
// 		}
// 		done <- true
// 	}()

// 	select {
// 	case <-done:
// 		if err != nil {
// 			return nil, err
// 		}
// 	case <-cancelCh:
// 		return nil, errors.New("cancel error")
// 	}

// 	b, err = ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode == http.StatusCreated {
// 		db, err := model.GetDBHandler("sqlite", "./dump.db")
// 		if err != nil {
// 			return nil, err
// 		}

// 		var device model.Device
// 		json.Unmarshal(b, &device)
// 		err = db.AddDevice(&device)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// db := model.GetDBHandler("sqlite", "./dump.db")
// 	return b, nil
// }
