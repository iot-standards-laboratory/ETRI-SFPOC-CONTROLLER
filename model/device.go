package model

import "etrisfpocdatamodel"

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

	tx := s.db.Create(device)
	if tx.Error != nil {
		return tx.Error
	}

	tx.First(device, "did=?", device.DID)
	return nil

}

func (s *_DBHandler) IsExistDevice(dname string) bool {
	var device = Device{}

	result := s.db.First(&device, "dname=?", dname)

	return result.Error != nil
}

func (s *_DBHandler) GetDeviceID(dname string) (string, error) {
	var device Device
	tx := s.db.Select("did", "sname").First(&device, "dname=?", dname)

	if tx.Error != nil {
		return "", tx.Error
	}

	return device.DID, nil
}

func (s *_DBHandler) GetServiceForDevice(did string) (string, error) {
	var device Device
	tx := s.db.Select("sname").First(&device, "did=?", did)
	if tx.Error != nil {
		return "", tx.Error
	}

	return s.GetSID(device.SName)
}
