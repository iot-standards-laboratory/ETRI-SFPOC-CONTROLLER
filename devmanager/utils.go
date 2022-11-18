package devmanager

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetToken generates a random token by a given length
func GetMessage(payload []byte) ([]byte, error) {
	length := len(payload) + 2
	b := make([]byte, 1, length)
	_, err := rand.Read(b)
	// length
	b = append(b, byte(length+1))
	b = append(b, payload...)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func initDiscoverDevice() error {
	fs, err := ioutil.ReadDir("/dev")
	if err != nil {
		return err
	}

	for _, f := range fs {
		if strings.Contains(f.Name(), "ttyACM") || strings.Contains(f.Name(), "ttyUSB") {
			if onConnected != nil {
				go onConnected(filepath.Join("/dev", f.Name()))
			}
		}
	}

	return nil
}

func changePermission(iface string) error {
	log.Println("changing the mod of file")

	cmd := exec.Command("sudo", "chmod", "a+rw", iface)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
