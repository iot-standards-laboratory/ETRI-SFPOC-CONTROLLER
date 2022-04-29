package config

import (
	"errors"
	"os"

	"github.com/magiconair/properties"
)

const Mode = "debug"

var Params = map[string]interface{}{}

func LoadConfig() {
	p := properties.MustLoadFile("./config.properties", properties.UTF8)
	Params["serverAddr"] = p.GetString("serverAddr", "localhost:3000")
	Params["mode"] = p.GetString("mode", "standalone")
	Params["bind"] = p.GetString("bind", ":4000")
	Params["cname"] = p.GetString("cname", "controllerA")
	Params["cid"] = p.GetString("cid", "blank")
	Params["sid"] = p.GetString("sid", "blank")
}

func CreateInitFile() {
	f, err := os.Create("./config.properties")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := properties.NewProperties()
	p.SetValue("serverAddr", "localhost:3000")
	p.SetValue("mode", STANDALONE)
	p.SetValue("bind", ":4000")
	p.SetValue("cname", "controllerA")
	p.SetValue("cid", "blank")
	p.SetValue("sid", "blank")
	p.Write(f, properties.UTF8)

}

func Set(key, value string) {

	var p *properties.Properties
	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		p = properties.NewProperties()
	} else {
		p = properties.MustLoadFile("./config.properties", properties.UTF8)
		os.Remove("./config.properties")
	}

	f, err := os.Create("./config.properties")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p.SetValue(key, value)
	p.Write(f, properties.UTF8)
	Params[key] = value
}
