package devmanager

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

var ctx context.Context
var cancel context.CancelFunc

var onDiscovered func(port io.ReadWriter, sname, dname string) = nil
var onRecv func(e Event) = nil

var onConnected func(e Event) = nil
var onDisconnected func(e Event) = nil

func AddOnDiscovered(h func(port io.ReadWriter, sname, dname string)) {
	onDiscovered = h
}

func AddRecvHandleFunc(h func(e Event)) {
	onRecv = h
}

func AddOnConnected(h func(e Event)) {
	onConnected = h
}

func AddOnDisconnected(h func(e Event)) {
	onDisconnected = h
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())
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
			fmt.Println("discivered!! - ", iface)
		} else if err.Error() != "not found device" {
			return err
		}

		// if discoverHandleFunc != nil {
		// 	discoverHandleFunc(NewEvent(map[string]interface{}{"iface": iface}, "discovered"))
		// }

		// fmt.Println("initDevice")

		// err = initDevice()
		// if err != nil {
		// 	log.Println(err)
		// 	port.Close()
		// 	continue
		// }

		// fmt.Println("recv")
		// err := recv(port, recvHandleFunc)
		// if err.Error() == "EOF" {
		// 	if onDisconnected != nil {
		// 		onDisconnected(NewEvent(map[string]interface{}{"iface": iface}, "disconnected"))
		// 	}
		// 	iface = ""
		// } else {
		// 	panic(err)
		// }
	}
}

func discover(iface string) error {
	var err error
	// err = changePermission(iface)
	// if err != nil {
	// 	return err
	// }

	options := serial.OpenOptions{
		PortName:        iface,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 16,
	}

	port, err := serial.Open(options)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(port)
	// encoder := json.NewEncoder(port)

	sndMsg := map[string]interface{}{}
	sndMsg["code"] = 0
	sndMsg["token"], err = GetToken()
	if err != nil {
		return err
	}

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		rcvMsg := map[string]interface{}{}
		err = json.Unmarshal(line, &rcvMsg)
		if err == nil {
			if rcvMsg["code"] != 1.0 {
				fmt.Println("initial done")
				return nil
			}

			if onDiscovered != nil {
				onDiscovered(port, rcvMsg["sname"].(string), rcvMsg["uuid"].(string))
			}
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
