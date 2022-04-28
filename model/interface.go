package model

import "gorm.io/gorm"

type DBHandlerI interface {
	// api for device
	GetDevices() ([]*Device, int, error)
	AddDevice(device *Device) error
	GetDevice(did string) (*Device, error)
	GetDeviceID(dname string) (string, error)
	// GetServiceForDevice(did string) (string, error)

	// api for service
	// GetSID(sname string) (string, error)
	// StatusCheck(did string, new map[string]interface{}) bool
}

type _DBHandler struct {
	db       *gorm.DB
	sidCache map[string]string
	states   map[string]map[string]interface{}
}

var DefaultDB DBHandlerI

func init() {
	var err error
	DefaultDB, err = NewSqliteHandler("dump.db")
	if err != nil {
		panic(err)
	}
}
