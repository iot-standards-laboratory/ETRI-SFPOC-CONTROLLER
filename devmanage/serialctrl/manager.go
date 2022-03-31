package serialctrl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/jacobsa/go-serial/serial"
)

type _manager struct {
	devicesWithUUID  map[string]*_device
	devicesWithIface map[interface{}]*_device
	chanForSync      map[string]chan map[string]interface{}
	RegisterListener EventHandler
	SyncListener     EventHandler
	RecvListener     EventHandler
}

func (m *_manager) onRegistered(dev *_device) {
	// if m.RegisterListener != nil {
	// 	m.RegisterListener.Handle(&EventStruct{key: })
	// }
	if registerHandleFunc != nil {
		param := map[string]interface{}{}
		param["uuid"] = dev.UUID
		param["sname"] = dev.Sname

		registerHandleFunc(&EventStruct{key: dev.UUID, params: param})
	}
	go recv(dev.Iface, m)
}

func (m *_manager) onRemoved(port io.Reader) {
	dev := _managerObj.devicesWithIface[port]
	delete(m.devicesWithUUID, dev.UUID)
	delete(m.devicesWithIface, port)
	ch, ok := m.chanForSync[dev.UUID]
	if ok {
		close(ch)
		delete(m.chanForSync, dev.UUID)
	}

	if removeHandleFunc != nil {
		param := map[string]interface{}{}
		param["uuid"] = dev.UUID

		removeHandleFunc(&EventStruct{key: dev.UUID, params: param})
	}
	// log.Println(m.devicesWithUUID)
	// log.Println(m.devicesWithIface)
}

func (m *_manager) onAdded(iface string) {
	err := ChangePermission(iface)
	if err != nil {
		panic(err)
	}

	// Set up options.
	options := serial.OpenOptions{
		PortName:        iface,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 16,
	}

	// Open the port.
	go func() {
		port, err := serial.Open(options)
		if err != nil {
			log.Printf("serial.Open: %v", err)
			return
		}

		reader := bufio.NewReader(port)
		encoder := json.NewEncoder(port)

		sndMsg := map[string]interface{}{}
		sndMsg["code"] = 100

		for {
			b, _, _ := reader.ReadLine()
			rcvMsg := map[string]interface{}{}
			err := json.Unmarshal(b, &rcvMsg)

			if err != nil {
				continue
			}

			code, ok := rcvMsg["code"].(float64)
			if ok && code == 100 {
				newDevice := &_device{
					UUID:      rcvMsg["uuid"].(string),
					IfaceName: iface,
					Iface:     port,
					Sname:     rcvMsg["sname"].(string),
					states:    map[string]interface{}{},
				}
				m.devicesWithUUID[rcvMsg["uuid"].(string)] = newDevice
				m.devicesWithIface[port] = newDevice

				m.onRegistered(newDevice)
				fmt.Println("onAdded sub-routine is died")
				return
			}
			encoder.Encode(sndMsg)
		}
	}()

	fmt.Println("onAdded main-routine is died")
}

func (m *_manager) onSync(key interface{}, params map[string]interface{}) {
	if m.SyncListener != nil {
		m.SyncListener.Handle(&EventStruct{key: key, params: params})
	}
}

func (m *_manager) onRecv(key interface{}, params map[string]interface{}) {
	if m.RecvListener != nil {
		m.RecvListener.Handle(&EventStruct{key: key, params: params})
	}
}

func (m *_manager) addRecvListener(h EventHandler) {
	m.RecvListener = &RecvHandler{next: h, devices: m.devicesWithIface, chanForSync: m.chanForSync}
}
