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


	db.AutoMigrate(&models.Client{}		)
	db.AutoMigrate(&models.Invoice{}	)
	db.AutoMigrate(&models.SellOrder{}	)
	db.AutoMigrate(&models.Ledger{}		)
	var test *models.Client

	test = &models.Client{Balance: "101.2299",IsInvestor:true}

	db.Create(test)
	fmt.Println(test)

	idInvestor := test.ID

	var invoice *models.Invoice = &models.Invoice{ClientID:idInvestor,Amount:"100.22",State:models.InvoiceRejected}
	db.Create(invoice)
	id := invoice.ID
	fmt.Println(invoice)

	var sellOrder *models.SellOrder = &models.SellOrder{InvoiceID:id,Size:"1000.00",Amount:"900.00",State:models.SellOrderOngoing}
	db.Create(sellOrder)
	fmt.Println(sellOrder)
	sellOrderId := sellOrder.ID

	var ledger *models.Ledger = &models.Ledger{InvestorID:idInvestor,SellOrderID:sellOrderId,Size:"100.00",Amount:"90.00",Balance:"10000"}
	db.Create(ledger)
	fmt.Println(ledger)

	ledger = &models.Ledger{}

	db.First(ledger,2)
	ledger.Balance = "200"
	db.Save(ledger)

}

