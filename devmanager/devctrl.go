package devmanager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
)

type DeviceControllerI interface {
	AddOnRecv(func(e Event))
	Sync(map[string]interface{}) error
	Run()
	close()
}

type deviceController struct {
	port        io.Reader
	dname       string
	did         string
	latestToken Token
	wg          sync.WaitGroup
	cmdCh       chan Event
	ackCh       chan Token
	onRecv      func(e Event)
}

func NewDeviceController(port io.Reader, dname, did string) DeviceControllerI {
	return &deviceController{
		port:  port,
		dname: dname,
		did:   did,
		cmdCh: make(chan Event),
		ackCh: make(chan Token),
	}
}

func (ctrl *deviceController) close() {
	close(ctrl.cmdCh)
	close(ctrl.ackCh)
}

func (ctrl *deviceController) AddOnRecv(h func(e Event)) {
	ctrl.onRecv = h
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
		"token": ctrl.latestToken,
	}, "command")

	return nil
}

func (ctrl *deviceController) Run() {
	ctrl.wg.Add(2)
	go ctrl.recv()
	go ctrl.send()
	ctrl.wg.Wait()
}

func (ctrl *deviceController) send() {
	// writer := json.NewEncoder(port)
	defer ctrl.wg.Done()

	for e := range ctrl.cmdCh {
		// e.Params()
		fmt.Println(e)
		// writer.Encode(e.Params())
	}

	fmt.Println("sender died")
}

func (ctrl *deviceController) recv() {
	defer ctrl.wg.Done()
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
	}
}
