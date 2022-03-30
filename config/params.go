package config

import (
	"errors"
	"os"

	"github.com/magiconair/properties"
)

const Mode = "debug"

var Params = map[string]interface{}{}

func init() {

	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		createInitFile()
	}

	p := properties.MustLoadFile("./config.properties", properties.UTF8)
	Params["serverAddr"] = p.GetString("serverAddr", "localhost:3000")
	Params["bind"] = p.GetString("bind", ":3000")
	Params["bind"] = p.GetString("bind", ":3000")
}

func createInitFile() {
	f, err := os.Create("./config.properties")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := properties.NewProperties()
	p.SetValue("serverAddr", "localhost:3000")
	p.SetValue("bind", ":3000")
	p.Write(f, properties.UTF8)

}
