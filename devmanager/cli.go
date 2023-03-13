package devmanager

import (
	"context"
	"log"
	"strings"

	"github.com/rjeczalik/notify"
)

var ctx context.Context
var cancel context.CancelFunc

var onConnected func(e DeviceControllerI) = nil

func AddOnConnected(h func(d DeviceControllerI)) {
	onConnected = h
}

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

func Close() {
	cancel()
}

func Watch() error {
	err := initDiscoverDevice()
	if err != nil {
		return err
	}

	filter := make(chan notify.EventInfo, 1)
	if err := notify.Watch("/dev", filter, notify.Create); err != nil {
		return err
	}

	defer notify.Stop(filter)

	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-filter:
			if strings.Contains(e.Path(), "/dev/ttyACM") || strings.Contains(e.Path(), "/dev/ttyUSB") {
				d, err := discover(e.Path())
				if err != nil {
					log.Println(err)
					continue
				}
				if onConnected != nil {
					go onConnected(d)
				}
			}
		}
	}
}
