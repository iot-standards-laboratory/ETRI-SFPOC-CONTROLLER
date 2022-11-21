package devmanager

import (
	"bytes"
	"errors"
	"hash/crc64"
	"io"
	"log"
	"sync"
	"time"
)

type DeviceControllerI interface {
	AddOnUpdate(func(e interface{}))
	AddOnError(func(err error))
	AddOnClose(func(key uint64))
	Sync(cmd []byte) error
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
	port     io.ReadWriter
	ctrlName string
	status   int

	syncMutex sync.Mutex
	ackCh     chan uint8

	onUpdate func(e interface{})
	onError  func(err error)
	onClose  func(key uint64)
	done     sync.WaitGroup
}

func NewDeviceController(port io.ReadWriter, ctrlName string) DeviceControllerI {
	return &deviceController{
		port:     port,
		ctrlName: ctrlName,
		status:   ControllerStatusReady,
		ackCh:    make(chan uint8),
	}
}

func (ctrl *deviceController) Key() uint64 {
	return crc64.Checksum([]byte(ctrl.ctrlName), crc64.MakeTable(crc64.ISO))
}

func (ctrl *deviceController) Close() {
	ctrl.status = ControllerStatusClosing

	ctrl.done.Wait()
	if ctrl.onClose != nil {
		ctrl.onClose(ctrl.Key())
	}
}

func (ctrl *deviceController) AddOnUpdate(h func(e interface{})) {
	ctrl.onUpdate = h
}

func (ctrl *deviceController) AddOnError(h func(err error)) {
	ctrl.onError = h
}

func (ctrl *deviceController) AddOnClose(h func(key uint64)) {
	ctrl.onClose = h
}

func (ctrl *deviceController) Run() {
	defer log.Println("ctrl", ctrl, "is stoped")
	ctrl.done.Add(1)
	defer ctrl.done.Done()

	ctrl.status = ControllerStatusRunning
	for {
		if ctrl.status == ControllerStatusClosing {
			return
		}

		code, msg, err := readMessage(ctrl.port)
		if err != nil {
			if ctrl.onError != nil {
				go ctrl.onError(err)
			}
			time.Sleep(time.Millisecond * 300)
			continue
		}

		if code == 201 {
			if ctrl.onUpdate != nil {
				ctrl.onUpdate(
					map[string]interface{}{
						"code": code,
						"msg":  string(msg),
					},
				)
			}
		} else if code == 200 {
			ctrl.ackCh <- msg[0]
		}
	}
}

func (ctrl *deviceController) Sync(cmd []byte) error {
	ctrl.syncMutex.Lock()
	defer ctrl.syncMutex.Unlock()

	msg, err := GetMessage(cmd)
	if err != nil {
		return err
	}
	tkn := uint8(msg[0])
	err = ctrl.writeMessage(msg)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 20; i++ {
		select {
		case <-ticker.C:
			log.Println("retransmission command to change mode as timeout")
			err = ctrl.writeMessage(msg)
			if err != nil {
				return err
			}
		case ackNum := <-ctrl.ackCh:
			if ackNum == tkn {
				return nil
			}
		}
	}

	return errors.New("timeout error")
}

func (ctrl *deviceController) writeMessage(payload []byte) error {
	buf := bytes.Buffer{}

	buf.WriteByte(byte(2)) // write code
	buf.Write(payload)
	buf.WriteByte(byte(255))

	_, err := ctrl.port.Write(buf.Bytes())
	return err
}
