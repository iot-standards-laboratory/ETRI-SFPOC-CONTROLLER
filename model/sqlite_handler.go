package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteHandler(path string) (DBHandlerI, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Device{})

	return &_DBHandler{
		db:       db,
		sidCache: map[string]string{},
		states:   map[string]map[string]interface{}{},
	}, nil
}
