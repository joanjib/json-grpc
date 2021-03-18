package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {

	dsn := "host=localhost user=joan password=joan123 dbname=num port=5432 sslmode=disable "
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
	    panic("failed to connect database")
	}
	fmt.Println("alert init db")
}
