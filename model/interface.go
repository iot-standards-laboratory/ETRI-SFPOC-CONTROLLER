package model

import "gorm.io/gorm"

type DBHandlerI interface {
	GetDevices() ([]*Device, int, error)
	AddDevice(device *Device) error
	GetSID(sname string) (string, error)
	GetServiceForDevice(did string) (string, error)
	GetDevice(did string) (*Device, error)
	GetDeviceID(dname string) (string, error)
	// StatusCheck(did string, new map[string]interface{}) bool
}

type _DBHandler struct {
	db       *gorm.DB
	sidCache map[string]string
	states   map[string]map[string]interface{}
}
