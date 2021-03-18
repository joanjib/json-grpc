package main

import (

	"fmt"
	"num/db"
	"num/models"
)

func main () {


	if db.Db == nil  {
		fmt.Println("nul")
	}
	db := db.Db


	db.AutoMigrate(&models.Client{})
	db.AutoMigrate(&models.Invoice{})
	var test *models.Client

	test = &models.Client{Balance: "101.2299"}

	db.Create(test)
	fmt.Println(test)

	id := test.ID

	var invoice *models.Invoice = &models.Invoice{ClientID:id,Amount:"100.22",State:models.Rejected}
	db.Create(invoice)
	fmt.Println(invoice)

}

