package main

import (

	"fmt"
	"num/db"
	"num/models"
	"unsafe"
)

func main () {


	if db.Db == nil  {
		fmt.Println("nul")
	}
	db := db.Db


	db.AutoMigrate(&models.Client{})
	var test *models.Client

	test = &models.Client{Balance: "-101.2299"}

	db.Create(test)
	fmt.Println(test)
	test = &models.Client{}


	db.First(test, 1)

	fmt.Println(test)

	var a uint  = 22

	fmt.Println("size of uint",unsafe.Sizeof(a))
}

