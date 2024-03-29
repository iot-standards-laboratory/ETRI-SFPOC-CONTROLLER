package devmanager

import (
	"errors"
	"hash/crc64"
	"log"
	"sync"
	"time"

	"go.bug.st/serial"
)

type DeviceControllerI interface {
	AddOnError(func(err error))
	AddOnClose(func(key uint64))
	Do(code uint8, payload []byte) (int, []byte, error)
	ResetBuffer() error
	Name() string
	ServiceName() string
	ServiceID() string
	Key() uint64
	Run()
	Close()
	// AddOnClose(func(dname string, did string, ctrl DeviceControllerI) error)
}

const (
	ControllerStatusReady = iota
	ControllerStatusRunning
	ControllerStatusClosing
)

type deviceController struct {
	port        serial.Port
	ctrlName    string
	serviceName string
	serviceID   string
	mutex       sync.Mutex
	status      int
	recvMsgCh   chan []byte
	onError     func(err error)
	onClose     func(key uint64)
}

func (ctrl *deviceController) Name() string {
	return ctrl.ctrlName
}

func (ctrl *deviceController) Key() uint64 {
	return crc64.Checksum([]byte(ctrl.ctrlName), crc64.MakeTable(crc64.ISO))
}

func (ctrl *deviceController) ServiceName() string {
	return ctrl.serviceName
}

func (ctrl *deviceController) ServiceID() string {
	return ctrl.serviceID
}

func (ctrl *deviceController) Close() {
	ctrl.status = ControllerStatusClosing
	ctrl.port.Close()
	close(ctrl.recvMsgCh)

	if ctrl.onClose != nil {
		ctrl.onClose(ctrl.Key())
	}
}

func (ctrl *deviceController) AddOnError(h func(err error)) {
	ctrl.onError = h
}

func (ctrl *deviceController) AddOnClose(h func(key uint64)) {
	ctrl.onClose = h
}

func (ctrl *deviceController) Run() {
	ctrl.status = ControllerStatusRunning

	for ctrl.status != ControllerStatusClosing {
		b, err := readMessage(ctrl.port)
		if err != nil {
			ctrl.Close()
			return
		}

		ctrl.recvMsgCh <- b
	}
}

func (ctrl *deviceController) ResetBuffer() error {
	err := ctrl.port.ResetInputBuffer()
	if err != nil {
		return err
	}
	err = ctrl.port.ResetOutputBuffer()
	if err != nil {
		return err
	}

	return nil
}

func (ctrl *deviceController) Do(code uint8, payload []byte) (int, []byte, error) {
	msg, err := getMessage(code, getToken(), payload)
	if err != nil {
		return -1, nil, err
	}

	ctrl.mutex.Lock()
	defer ctrl.mutex.Unlock()
	ctrl.ResetBuffer()
	_, err = ctrl.port.Write(msg)
	if err != nil {
		return -1, nil, err
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			log.Println("retransmission command as timeout")
			ctrl.ResetBuffer()
			_, err = ctrl.port.Write(msg)
			if err != nil {
				return -1, nil, err
			}
		case recvMsg, ok := <-ctrl.recvMsgCh:
			if !ok {
				return -1, nil, errors.New("channel is closed error")
			}

			if recvMsg[1] != msg[1] {
				log.Println("retransmission command as invalid token:", recvMsg[1])
				ctrl.ResetBuffer()
				if err != nil {
					return -1, nil, err
				}
				continue
			}

			return int(recvMsg[0]), recvMsg[2:], nil
		}
	}

	return -1, nil, errors.New("timeout error")
}
