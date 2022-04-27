package devmanager

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

var ctx context.Context
var cancel context.CancelFunc

var onDiscovered func(port io.ReadWriter, sname, dname string) error = nil

var onConnected func(e Event) = nil

func AddOnDiscovered(h func(port io.ReadWriter, sname, dname string) error) {
	onDiscovered = h
}

func AddOnConnected(h func(e Event)) {
	onConnected = h
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

func Close() {
	cancel()
}

func Watch() error {
	var err error
	iface, err := discoverDevice()
	if err == nil {
		// onDiscover
		if onConnected != nil {
			onConnected(NewEvent(map[string]interface{}{"port": iface}, "onConnected"))
		}
		go discover(iface)
	} else if err.Error() != "not found device" {
		return err
	}

	for {
		iface, err = WatchNewDevice(ctx)
		if err == nil {
			// onDiscover
			go discover(iface)
		} else if err.Error() != "not found device" {
			return err
		}
	}
}

func discover(iface string) error {

	defer log.Println("exit discover()")
	err := changePermission(iface)
	if err != nil {
		log.Println(err)
		return err
	}

	options := serial.OpenOptions{
		PortName:        iface,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 16,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Println(err)
		return err
	}

	reader := bufio.NewReader(port)
	// encoder := json.NewEncoder(port)

	sndMsg := map[string]interface{}{}
	sndMsg["code"] = 0
	sndMsg["token"], err = GetToken()
	if err != nil {
		log.Println(err)
		return err
	}

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println(err)
			return err
		}
		rcvMsg := map[string]interface{}{}
		err = json.Unmarshal(line, &rcvMsg)
		if err == nil {
			if rcvMsg["code"] == 1.0 {
				port.Write([]byte(`{"code": 255, "token": "initial", "mode": 0}`))
				time.Sleep(time.Second)
				continue
			}

			if onDiscovered != nil {
				return onDiscovered(port, rcvMsg["sname"].(string), rcvMsg["uuid"].(string))
			}

			return nil
		}
	}
	// return nil
}

// func Write(obj interface{}) error {
// 	if port == nil {
// 		return errors.New("device is not connected")
// 	}
// 	enc := json.NewEncoder(port)
// 	err := enc.Encode(obj)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
