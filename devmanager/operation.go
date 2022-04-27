package devmanager

import (
	"fmt"
	"time"
)

func RegisterDevice(dev map[string]interface{}, waitCh chan bool) error {
	fmt.Println("register device()")
	time.Sleep(time.Second * 2)
	return nil
	// dev["cid"] = config.Params["cid"]
	// b, err := json.Marshal(dev)
	// if err != nil {
	// 	log.Println(err)
	// 	waitCh <- false
	// }

	// resp, err := http.Post(fmt.Sprintf("http://%s/api/v1/devs/discover", config.Params["serverAddr"].(string)), "application/json", bytes.NewReader(b))
	// if err != nil || resp.StatusCode != http.StatusCreated {
	// 	log.Println(err)
	// 	waitCh <- false
	// 	return
	// }

	// b, err = ioutil.ReadAll(resp.Body)
	// if err != nil || resp.StatusCode != http.StatusCreated {
	// 	log.Println(err)
	// 	waitCh <- false
	// 	return
	// }

	// json.Unmarshal(b, &dev)
	// fmt.Println(dev)
	// d1 := serialctrl.DevWithUUID[dev["dname"].(string)]
	// d1.Did = dev["did"].(string)

	// d2 := serialctrl.DevWithIface[d1.Iface]
	// d2.Did = dev["did"].(string)

	// fmt.Println("d1: ", d1)
	// fmt.Println("d2: ", d2)
	// log.Println("registered] ", dev)
	// waitCh <- true
}
