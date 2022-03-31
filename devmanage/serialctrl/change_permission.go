package serialctrl

import (
	"fmt"
	"log"
	"os/exec"
)

func ChangePermission(iface string) error {
	log.Println("changing the mod of file")

	cmd := exec.Command("chmod", "a+rw", iface)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
