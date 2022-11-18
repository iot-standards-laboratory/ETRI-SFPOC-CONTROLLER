package devmanager

import (
	"context"
	"strings"

	"io"

	"github.com/rjeczalik/notify"
)

var ctx context.Context
var cancel context.CancelFunc

var onConnected func(e string) = nil

func AddOnConnected(h func(e string)) {
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
				go onConnected(e.Path())
			}
		}
	}
}

func readMessage(reader io.Reader) (int, []byte, error) {
	buf := make([]byte, 255)
	b := make([]byte, 1)
	len := 0

	var err error
	for {
		_, err = reader.Read(b)
		if err != nil {
			return 254, nil, err
		}

		if b[0] == 255 {
			return int(buf[0]), buf[1:len], nil
		}

		buf[len] = b[0]
		len++
	}
}
