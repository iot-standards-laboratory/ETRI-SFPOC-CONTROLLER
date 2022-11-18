package devmanager

import (
	"hash/crc64"
	"io"
)

type DeviceControllerI interface {
	AddOnUpdate(func(e interface{}))
	// Sync(code uint8, cmd string) error
	Key() uint64
	Run()
	// close()
	// AddOnClose(func(dname string, did string, ctrl DeviceControllerI) error)
}

type deviceController struct {
	port        io.ReadWriter
	ctrlName    string
	latestToken Token
	ackCh       chan string
	onUpdate    func(e interface{})
	onClose     func(dname, did string, ctrl DeviceControllerI) error
}

func NewDeviceController(port io.ReadWriter, ctrlName string) DeviceControllerI {
	return &deviceController{
		port:     port,
		ctrlName: ctrlName,
		ackCh:    make(chan string),
	}
}

func (ctrl *deviceController) Key() uint64 {
	return crc64.Checksum([]byte(ctrl.ctrlName), crc64.MakeTable(crc64.ISO))
}

func (ctrl *deviceController) close() {
	close(ctrl.ackCh)
}

func (ctrl *deviceController) AddOnUpdate(h func(e interface{})) {
	ctrl.onUpdate = h
}

func (ctrl *deviceController) AddOnClose(h func(dname string, did string, ctrl DeviceControllerI) error) {
	ctrl.onClose = h
}

func (ctrl *deviceController) Run() {
	go func() {
		for {
			code, msg, err := readMessage(ctrl.port)
			if err != nil {
				return
			}

			if ctrl.onUpdate != nil {
				ctrl.onUpdate(
					map[string]interface{}{
						"code": code,
						"msg":  msg,
					},
				)
			}
		}
	}()
	// go ctrl.send()
}

// func (ctrl *deviceController) Sync(body map[string]interface{}) error {
// 	var err error
// 	ctrl.latestToken, err = GetToken()
// 	if err != nil {
// 		return err
// 	}

// 	ctrl.cmdCh <- NewEvent(map[string]interface{}{
// 		"code":  1,
// 		"body":  body,
// 		"token": ctrl.latestToken.String(),
// 	}, "command")

// 	return nil
// }

// func (ctrl *deviceController) send() {
// 	enc := json.NewEncoder(ctrl.port)
// 	// enc := json.NewEncoder(os.Stdout)
// 	ticker := time.NewTicker(time.Second)
// 	ticker.Stop()
// 	defer func() {
// 		ticker.Stop()
// 	}()
// 	var latestParams map[string]interface{}

// 	for {
// 		select {
// 		case e, ok := <-ctrl.cmdCh:
// 			if !ok {
// 				return
// 			}
// 			latestParams = e.Params()
// 			err := enc.Encode(latestParams)
// 			if err != nil {
// 				log.Println(err)
// 			}

// 			ticker.Reset(time.Second * 5)
// 		case tkn, ok := <-ctrl.ackCh:
// 			if !ok {
// 				return
// 			}
// 			if tkn == ctrl.latestToken.String() {
// 				fmt.Println("Acked!!")
// 				ticker.Stop()
// 			}
// 		case <-ticker.C:
// 			fmt.Println("retransmission as timeout : ", latestParams)
// 			err := enc.Encode(latestParams)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 		}
// 	}

// }

// func (ctrl *deviceController) recv() {
// 	defer ctrl.close()

// 	reader := bufio.NewReader(ctrl.port)

// 	for {
// 		b, _, err := reader.ReadLine()
// 		if err != nil {
// 			if err == io.EOF {
// 				log.Println("USB is disconnected")
// 				return
// 			}
// 		}
// 		recvObj := map[string]interface{}{}
// 		err = json.Unmarshal(b, &recvObj)

// 		if recvObj["code"] == 2.0 {
// 			ctrl.ackCh <- recvObj["token"].(string)
// 			continue
// 		}

// 		if err == nil && ctrl.onRecv != nil {
// 			ctrl.onRecv(NewEvent(recvObj, "recv"))
// 		}

// 	}
// }
