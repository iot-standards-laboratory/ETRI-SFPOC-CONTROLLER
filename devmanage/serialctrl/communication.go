package serialctrl

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

func recv(port io.Reader, r Receiver) {
	reader := bufio.NewReader(port)

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("USB is disconnected")
				_managerObj.onRemoved(port)
				return
			}
		}
		recvObj := map[string]interface{}{}
		err = json.Unmarshal(b, &recvObj)

		if err == nil {
			r.onRecv(port, recvObj)
		}
		// data = string(b)
	}
}

// for {
// 	fmt.Print("> ")
// 	cmd, _, _ := cmdReader.ReadLine()
// 	cmdTkns := strings.Split(string(cmd), " ")

// 	if cmdTkns[0] == "light" {
// 		command["code"] = 1
// 		if cmdTkns[1] == "on" {
// 			command["light"] = 100
// 		} else {
// 			command["light"] = 0
// 		}
// 	} else if cmdTkns[0] == "fan" {
// 		command["code"] = 2
// 		if cmdTkns[1] == "on" {
// 			command["status"] = 1
// 		} else {
// 			command["status"] = 0
// 		}
// 	} else if cmdTkns[0] == "servo" {
// 		command["code"] = 3
// 		angle, err := strconv.Atoi(cmdTkns[1])
// 		if err != nil {
// 			continue
// 		}
// 		command["angle"] = angle
// 	} else if cmdTkns[0] == "print" {
// 		fmt.Println(data)
// 	}
// 	err := encoder.Encode(command)
// 	if err != nil {
// 		panic(err)
// 	}
// }
