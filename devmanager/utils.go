package devmanager

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/notify"
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
	b := make([]byte, 4)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func discoverDevice() (string, error) {
	fs, err := ioutil.ReadDir("/dev")
	if err != nil {
		return "", err
	}

	for _, f := range fs {
		if strings.Contains(f.Name(), "ttyACM") || strings.Contains(f.Name(), "ttyUSB") {
			return filepath.Join("/dev", f.Name()), nil
		}
	}

	return "", errors.New("not found device")
}

func WatchNewDevice(ctx context.Context) (string, error) {
	log.Println("Watching device")
	filter := make(chan notify.EventInfo, 1)
	if err := notify.Watch("/dev", filter, notify.Create); err != nil {
		return "", err
	}
	defer notify.Stop(filter)

	for {
		select {
		case <-ctx.Done():
			return "", nil
		case e := <-filter:
			if strings.Contains(e.Path(), "/dev/ttyACM") || strings.Contains(e.Path(), "/dev/ttyUSB") {
				return e.Path(), nil

			}
		}
	}
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
