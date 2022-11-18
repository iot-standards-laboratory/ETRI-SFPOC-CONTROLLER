package devmanager

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type Token []byte

func (t Token) String() string {
	return hex.EncodeToString(t)
}

func (t Token) Hash() uint64 {
	return crc64.Checksum(t, crc64.MakeTable(crc64.ISO))
}

// GetToken generates a random token by a given length
func GetToken() (Token, error) {
	b := make([]byte, 1)
	_, err := rand.Read(b)
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
