package devmanager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

type DeviceControllerI interface {
	AddOnRecv(func(e Event))
	Sync(map[string]interface{}) error
	Run()
	close()
	AddOnClose(func(dname string, did string, ctrl DeviceControllerI) error)
}

type deviceController struct {
	port        io.ReadWriter
	dname       string
	did         string
	latestToken Token
	cmdCh       chan Event
	ackCh       chan string
	onRecv      func(e Event)
	onClose     func(dname, did string, ctrl DeviceControllerI) error
}

func NewDeviceController(port io.ReadWriter, dname, did string) DeviceControllerI {
	return &deviceController{
		port:  port,
		dname: dname,
		did:   did,
		cmdCh: make(chan Event),
		ackCh: make(chan string),
	}
}

func (ctrl *deviceController) close() {
	close(ctrl.cmdCh)
	close(ctrl.ackCh)
	if ctrl.onClose != nil {
		ctrl.onClose(ctrl.dname, ctrl.did, ctrl)
	}
}

func (ctrl *deviceController) AddOnRecv(h func(e Event)) {
	ctrl.onRecv = h
}

func (ctrl *deviceController) AddOnClose(h func(dname string, did string, ctrl DeviceControllerI) error) {
	ctrl.onClose = h
}

func (ctrl *deviceController) Sync(body map[string]interface{}) error {
	var err error
	ctrl.latestToken, err = GetToken()
	if err != nil {
		return err
	}

	ctrl.cmdCh <- NewEvent(map[string]interface{}{
		"code":  1,
		"body":  body,
		"token": ctrl.latestToken.String(),
	}, "command")

	return nil
}

func (ctrl *deviceController) Run() {
	go ctrl.recv()
	go ctrl.send()
}

func (ctrl *deviceController) send() {
	enc := json.NewEncoder(ctrl.port)
	// enc := json.NewEncoder(os.Stdout)
	ticker := time.NewTicker(time.Second)
	ticker.Stop()
	defer func() {
		ticker.Stop()
	}()
	var latestParams map[string]interface{}

	for {
		select {
		case e, ok := <-ctrl.cmdCh:
			if !ok {
				return
			}
			latestParams = e.Params()
			err := enc.Encode(latestParams)
			if err != nil {
				log.Println(err)
			}

			ticker.Reset(time.Second * 5)
		case tkn, ok := <-ctrl.ackCh:
			if !ok {
				return
			}
			if tkn == ctrl.latestToken.String() {
				fmt.Println("Acked!!")
				ticker.Stop()
			}
		case <-ticker.C:
			fmt.Println("retransmission as timeout : ", latestParams)
			err := enc.Encode(latestParams)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func (ctrl *deviceController) recv() {
	defer ctrl.close()

	reader := bufio.NewReader(ctrl.port)

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("USB is disconnected")
				return
			}
		}
		recvObj := map[string]interface{}{}
		err = json.Unmarshal(b, &recvObj)

		if err == nil && ctrl.onRecv != nil {
			ctrl.onRecv(NewEvent(recvObj, "recv"))
		}

		if recvObj["code"] == 2.0 {
			ctrl.ackCh <- recvObj["token"].(string)
		}
	}
}
