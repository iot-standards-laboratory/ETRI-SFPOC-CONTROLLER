package model

import (
	"etrisfpocdatamodel"
	"fmt"
)

type Device etrisfpocdatamodel.Device

func (s *_DBHandler) GetDevices() ([]*Device, int, error) {
	var devices []*Device

	result := s.db.Find(&devices)

	if result.Error != nil {
		return nil, -1, result.Error
	}
	return devices, int(result.RowsAffected), nil
}

func (s *_DBHandler) AddDevice(device *Device) error {

	fmt.Println("device name : ", device.DName)
	tx := s.db.Create(device)
	if tx.Error != nil {
		return tx.Error
	}

	tx.First(device, "did=?", device.DID)
	return nil

}

func (s *_DBHandler) GetDevice(did string) (*Device, error) {
	var device Device
	tx := s.db.First(&device, "did=?", did)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &device, nil
}

func (s *_DBHandler) GetDeviceID(dname string) (string, error) {
	var device Device
	tx := s.db.Select("did", "sname").First(&device, "dname=?", dname)

	if tx.Error != nil {
		return "", tx.Error
	}

	return device.DID, nil
}

// func (s *_DBHandler) GetServiceForDevice(did string) (string, error) {
// 	var device Device
// 	tx := s.db.Select("sname").First(&device, "did=?", did)
// 	if tx.Error != nil {
// 		return "", tx.Error
// 	}

// 	return cache.GetSID(device.SName)
// }
