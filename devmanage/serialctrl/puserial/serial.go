package puserial

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
)

func InitDevice() ([]string, error) {
	fs, err := ioutil.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	var result []string = nil
	for _, f := range fs {
		if strings.Contains(f.Name(), "ttyACM") {
			result = append(result, filepath.Join("/dev", f.Name()))
		}
	}

	return result, nil
}

func WatchNewDevice(ctx context.Context, ch_discover chan<- notify.EventInfo) error {
	defer close(ch_discover)

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
			if strings.Contains(e.Path(), "/dev/ttyACM") {
				fmt.Println(e.Path())
				ch_discover <- e
			}
		}
	}
}

// func InitUSB() (string, error) {
// 	fs, err := ioutil.ReadDir("/dev")
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, f := range fs {
// 		if strings.Contains(f.Name(), "ttyACM") {
// 			return filepath.Join("/dev", f.Name()), nil
// 		}
// 	}

// 	return "", errors.New("USB Not found")
// }

// func Run() {
// 	dev, err := InitUSB()
// 	if err != nil {
// 		if err.Error() == "USB Not found" {
// 			dev = DiscoverUSB()
// 		} else {
// 			log.Fatalln(err.Error())
// 		}
// 	}

// 	fmt.Println("discover: ", dev)
// 	time.Sleep(time.Second)

// 	log.Println("changing the mod of file")
// 	cmd := exec.Command("chmod", "a+rw", dev)
// 	b, _ := cmd.CombinedOutput()
// 	fmt.Println(string(b))

// 	// Set up options.
// 	options := serial.OpenOptions{
// 		PortName:        dev,
// 		BaudRate:        9600,
// 		DataBits:        8,
// 		StopBits:        1,
// 		MinimumReadSize: 16,
// 	}

// 	// Open the port.
// 	port, err := serial.Open(options)
// 	if err != nil {
// 		log.Fatalf("serial.Open: %v", err)
// 	}

// 	// Make sure to close it later.
// 	defer port.Close()

// 	// Write 4 bytes to the port.
// 	// b := []byte{0x00, 0x01, 0x02, 0x03}
// 	// n, err := port.Write(b)
// 	// if err != nil {
// 	// 	log.Fatalf("port.Write: %v", err)
// 	// }

// 	// fmt.Println("Wrote", n, "bytes.")
// 	// var tokenByte string

// 	reader := bufio.NewReader(port)
// 	// decoder := json.NewDecoder(port)
// 	encoder := json.NewEncoder(port)
// 	command := map[string]interface{}{}
// 	command["code"] = 1
// 	command["light"] = 0

// 	var data string
// 	go func() {
// 		for {
// 			b, _, _ := reader.ReadLine()
// 			recvObj := map[string]interface{}{}
// 			err := json.Unmarshal(b, &recvObj)
// 			// err = readJsonFromSerial(recvObj, decoder)
// 			if err != nil {
// 				if err.Error() == "EOF" {
// 					return
// 				}

// 				fmt.Println("error: ", string(b))
// 			}

// 			data = string(b)
// 			// fmt.Println("line : ", recvObj)
// 		}
// 	}()

// 	cmdReader := bufio.NewReader(os.Stdin)
// 	for {
// 		fmt.Print("> ")
// 		cmd, _, _ := cmdReader.ReadLine()
// 		cmdTkns := strings.Split(string(cmd), " ")

// 		if cmdTkns[0] == "light" {
// 			command["code"] = 1
// 			if cmdTkns[1] == "on" {
// 				command["light"] = 100
// 			} else {
// 				command["light"] = 0
// 			}
// 		} else if cmdTkns[0] == "fan" {
// 			command["code"] = 2
// 			if cmdTkns[1] == "on" {
// 				command["status"] = 1
// 			} else {
// 				command["status"] = 0
// 			}
// 		} else if cmdTkns[0] == "servo" {
// 			command["code"] = 3
// 			angle, err := strconv.Atoi(cmdTkns[1])
// 			if err != nil {
// 				continue
// 			}
// 			command["angle"] = angle
// 		} else if cmdTkns[0] == "print" {
// 			fmt.Println(data)
// 		}
// 		err := encoder.Encode(command)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }
