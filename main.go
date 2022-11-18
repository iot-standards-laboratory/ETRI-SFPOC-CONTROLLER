package main

import (
	"errors"
	"etri-sfpoc-controller/config"
	"etri-sfpoc-controller/devmanager"
	"etri-sfpoc-controller/model/cachestorage"
	"etri-sfpoc-controller/router"
	"etri-sfpoc-controller/statmgmt"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

// func makeDeviceController(port io.ReadWriter, did, dname, sname string) devmanager.DeviceControllerI {
// 	fmt.Println("makeDeviceController()")
// 	// model.AddDeviceController 에서 등록된 디바이스 목록에 해당 디바이스 추가할 것!!
// 	ctrl := devmanager.NewDeviceController(port, dname, did)

// }

var interrupt chan os.Signal

func waitInterrupt() {
	// waiting interrupt
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	log.Println("receive interrupt")
}

func main() {
	cfg := flag.Bool("init", false, "create initial config file")
	flag.Parse()

	if *cfg {
		config.CreateInitFile()
		return
	}

	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		fmt.Println("config file doesn't exist")
		fmt.Println("please add -init option to create config file")
		return
	}
	config.LoadConfig()

	statmgmt.Bootup()

	if statmgmt.Status() == statmgmt.STATUS_DISCONNECTED {
		err := statmgmt.Connect()
		if err != nil {
			panic(err)
		}
	}

	// var cancel context.CancelFunc = nil
	// defer func() {
	// 	if cancel != nil {
	// 		cancel()
	// 	}
	// }()
	// if strings.Compare(config.Params["mode"].(string), string(config.MANAGEDBYEDGE)) == 0 {
	// 	cancel = manageSubscribe()
	// }

	devmanager.AddOnConnected(func(port string) {
		ctrl, err := devmanager.DiscoverController(port)
		if err != nil {
			log.Println(err)
			return
		}
		ctrl.AddOnUpdate(func(e interface{}) {
			fmt.Println(e)
		})

		ctrl.AddOnError(func(e error) {
			if e.Error() == "EOF" {
				log.Println("EOF error!!")
				ctrl.Close()
			}
		})

		ctrl.AddOnClose(func(key uint64) {
			cachestorage.RemoveDeviceController(key)
		})
		go ctrl.Run()
		cachestorage.AddDeviceController(ctrl)
	})

	go devmanager.Watch()
	// go devManagerTest()

	go router.NewRouter().Run(config.Params["bind"].(string))

	waitInterrupt()
	// do something before program exit
	// websocket close

}
