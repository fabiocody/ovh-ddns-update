package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Address struct {
	ID        uint8  `gorm:"primaryKey"`
	Value     string `gorm:"unique"`
	UpdatedAt time.Time
}

var DB *gorm.DB

func SetupDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(Args.Database))
	HandleError(err)
	err = DB.AutoMigrate(&Address{})
	HandleError(err)
}

func GetPreviousIP() string {
	var address Address
	if result := DB.First(&address); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warning("Previous IP not found")
			return "not_found"
		}
		HandleError(result.Error)
	}
	log.Debugf("Previous IP is %s", address.Value)
	return address.Value
}

func SaveCurrentIP(ip string) {
	address := Address{ID: 1, Value: ip}
	result := DB.Save(&address)
	HandleError(result.Error)
}
