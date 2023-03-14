package devmanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.bug.st/serial"
)

func getToken() uint8 {
	// return uint8(rand.Intn(253) + 1)
	return 100
}

// GetToken generates a random token by a given length
func getMessage(code, token uint8, payload []byte) ([]byte, error) {

	if payload == nil {
		return []byte{code, token, 255}, nil
	}

	b := make([]byte, 0, len(payload)+3)
	length := uint8(len(payload) + 3)
	b = append(b, code, token, length)
	b = append(b, payload...)
	b = append(b, 255)
	// Note that err == nil only if we read len(b) bytes.

	return b, nil
}

func readMessage(reader io.ReadCloser) ([]byte, error) {
	buf := bytes.Buffer{}
	b := make([]byte, 1)
	len := 0

	var err error
	for {
		_, err = reader.Read(b)
		if err != nil {
			return nil, err
		}

		if b[0] == 255 {
			return buf.Bytes(), nil
		}

		buf.Write(b)
		len++
	}
}

func initDiscoverDevice() error {
	fs, err := os.ReadDir("/dev")
	if err != nil {
		return err
	}

	for _, f := range fs {
		if strings.Contains(f.Name(), "ttyACM") || strings.Contains(f.Name(), "ttyUSB") {
			d, err := discover(filepath.Join("/dev", f.Name()))
			if err != nil {
				log.Println(err)
				continue
			}
			if onConnected != nil {
				onConnected(d)
			}
		}
	}

	return nil
}

func discover(iface string) (DeviceControllerI, error) {
	options := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(iface, options)
	if err != nil {
		port.Close()
		return nil, err
	}

	devCtrl := &deviceController{
		port:      port,
		recvMsgCh: make(chan []byte),
		status:    ControllerStatusReady,
	}
	go devCtrl.Run()

	code, b, err := devCtrl.Do(1, []byte("init"))
	if err != nil {
		devCtrl.Close()
		return nil, err
	}

	if code != 205 {
		devCtrl.Close()
		return nil, errors.New("invalid response error")
	}

	var initInformation map[string]interface{}
	err = json.Unmarshal(b, &initInformation)
	if err != nil {
		devCtrl.Close()
		return nil, err
	}

	var ok bool
	devCtrl.ctrlName, ok = initInformation["uuid"].(string)

	if !ok {
		devCtrl.Close()
		return nil, errors.New("not imported uuid error")
	}

	devCtrl.serviceName, ok = initInformation["sname"].(string)
	if !ok {
		devCtrl.Close()
		return nil, errors.New("not imported sname error")
	}

	return devCtrl, nil
}
