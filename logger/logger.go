package logger

import (
	"fmt"
	"os"
	"time"
)

func Start() {
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	dt := time.Now()
	f.WriteString(fmt.Sprintf("start at - %d.%d.%d - %d:%d\n", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute()))
	f.Close()
}

func Stop() int {
	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	dt := time.Now()
	f.WriteString(fmt.Sprintf("stop at - %d.%d.%d - %d:%d\n", dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute()))
	f.Close()

	return 1
}
