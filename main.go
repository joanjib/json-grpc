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
	db.AutoMigrate(&models.SellOrder{})
	var test *models.Client

	test = &models.Client{Balance: "101.2299"}

	db.Create(test)
	fmt.Println(test)

	id := test.ID

	var invoice *models.Invoice = &models.Invoice{ClientID:id,Amount:"100.22",State:models.InvoiceRejected}
	db.Create(invoice)
	id = invoice.ID
	fmt.Println(invoice)

	var sellOrder *models.SellOrder = &models.SellOrder{InvoiceID:id,Size:"1000.00",Amount:"900.00",State:models.SellOrderOngoing}
	db.Create(sellOrder)
	fmt.Println(sellOrder)

}

