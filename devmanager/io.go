package devmanager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func send(port io.Writer, eventCh <-chan Event) error {
	// writer := json.NewEncoder(port)

	for e := range eventCh {
		// e.Params()
		fmt.Println(e)
		// writer.Encode(e.Params())
	}

	fmt.Println("sender died")
	return nil
}

func recv(port io.Reader, h func(e Event)) error {
	reader := bufio.NewReader(port)

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("USB is disconnected")

				return err
			}
		}
		recvObj := map[string]interface{}{}
		err = json.Unmarshal(b, &recvObj)

		if err == nil && h != nil {
			h(NewEvent(recvObj, "recv"))
		}
	}
}
