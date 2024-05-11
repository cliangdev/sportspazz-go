package server

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySqlStorage(dsn string, config gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &config)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
