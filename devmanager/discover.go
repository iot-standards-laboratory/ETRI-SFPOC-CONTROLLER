package devmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

func DiscoverController(iface string) (DeviceControllerI, error) {
	log.Println("start DiscoverController()")
	defer log.Println("exit DiscoverController()")
	err := changePermission(iface)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	options := serial.OpenOptions{
		PortName:        iface,
		BaudRate:        57600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 16,
	}

	port, err := serial.Open(options)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// encoder := json.NewEncoder(port)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	isDone := false
	timer := time.AfterFunc(time.Second*10, func() {
		isDone = true
	})
	defer timer.Stop()

	for !isDone {
		code, payload, err := readMessage(port)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if code == 199 {
			err = initDevice(port)
			if err != nil {
				return nil, err
			}
			rcvMsg := map[string]interface{}{}
			err = json.Unmarshal(payload, &rcvMsg)
			if err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			ctrlName := rcvMsg["uuid"].(string)
			serviceName := rcvMsg["sname"].(string)
			fmt.Println("service name: ", serviceName)
			return NewDeviceController(port, ctrlName, serviceName), nil
		} else if code == 201 {
			_, err = port.Write([]byte{10, byte(rand.Int() % 256), 3, 255})
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, errors.New("timeout error")
}

func initDevice(port io.ReadWriter) error {
	log.Println("start initDevice()")
	defer log.Println("exit initDevice()")
	var err error
	okchan := make(chan error)
	defer close(okchan)
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()
	_, err = port.Write([]byte{11, byte(rand.Int() % 256), 3, 255})
	if err != nil {
		return err
	}

	isTerminated := false
	defer func() { isTerminated = true }()

	go func() {
		log.Println("start initDevice gorouting end")
		defer log.Println("exit initDevice gorouting end")
		for {
			code, _, err := readMessage(port)
			if err != nil {
				log.Println("okchan : ", err)

				if isTerminated {
					return
				}

				okchan <- err
				return
			}

			if isTerminated {
				return
			}

			if code == 200 {
				okchan <- nil
				return
			}
		}
	}()

	for i := 0; i < 20; i++ {
		select {
		case <-ticker.C:
			fmt.Println("retransmission command to change mode as timeout")
			_, err = port.Write([]byte{11, byte(rand.Int() % 256), 3, 255})
			if err != nil {
				return err
			}
		case err := <-okchan:
			// err is nil or error
			return err
		}
	}

	return err
}
