package main

import (

	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type Test struct {

	gorm.Model
	Price		string `gorm:"type:numeric(11,2)"`
}

func main () {

	dsn := "host=localhost user=joan password=joan123 dbname=num port=5432 sslmode=disable "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
	    panic("failed to connect database")
	}

	db.AutoMigrate(&Test{})
	var test *Test

	test = &Test{Price: "101.2299"}

	db.Create(test)
	fmt.Println(test)
	test = &Test{}


	db.First(test, 1)

	fmt.Println(test)
}

