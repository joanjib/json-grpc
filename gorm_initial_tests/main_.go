package main

import (

	"fmt"
//	"io/ioutil"
	"num/db"
	"num/models"
	"gorm.io/gorm"
)

func main () {


	if db.Db == nil  {
		panic("not initialitzed Db object")
	}
	db := db.Db

	var err error
//	typesDomains, err := ioutil.ReadFile("./sql/types-domains.sql")
    if err != nil {
        panic(err)
    }

	var r *gorm.DB

	ErrorTrack :
	if r != nil {
		fmt.Println("error ob")
		return
	}
	ErrorTrack_err :
	if err != nil {
		fmt.Println("error err")
		return
	}

//	r- = db.Exec  (string(typesDomains)	)

	// migrations:
	err = db.AutoMigrate(&models.Client{}	)
	err = db.AutoMigrate(&models.Invoice{}	)
	err = db.AutoMigrate(&models.SellOrder{})
	err = db.AutoMigrate(&models.Ledger{}	)

	// Insertions on all tables:
	var test *models.Client
	test = &models.Client{Balance: "101.2299",IsInvestor:true}
	r = db.Create(test)
	fmt.Println(test)
	idInvestor := test.ID

	var invoice *models.Invoice = &models.Invoice{ClientID:idInvestor,Amount:"100.22",State:models.InvoiceRejected}
	r = db.Create(invoice)
	id := invoice.ID
	fmt.Println(invoice)

	var sellOrder *models.SellOrder = &models.SellOrder{InvoiceID:id,Size:"1000.00",Amount:"900.00",State:models.SellOrderOngoing}
	r = db.Create(sellOrder)
	fmt.Println(sellOrder)
	sellOrderId := sellOrder.ID

	var ledger *models.Ledger = &models.Ledger{InvestorID:idInvestor,SellOrderID:sellOrderId,Size:"100.00",Amount:"90.00",Balance:"10000"}
	r = db.Create(ledger)
	fmt.Println(ledger)

	// retrieve of the ledger 2
	ledger = &models.Ledger{}
	r = db.First(ledger,2)
	// update of the mentioned ledger:
	ledger.Balance = "200"
	r = db.Save(ledger)

}

