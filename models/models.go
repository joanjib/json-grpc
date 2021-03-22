package models

import (
	"database/sql/driver"
	"gorm.io/gorm"
	pb "num/arexservices"
)

type Client struct {
	gorm.Model
	FiscalIdentity		string			`gorm:"not null"`
	Name				string			`gorm:"not null"`
	Surname				string			`gorm:"not null"`
	Balance				string			`gorm:"type:amount_type"`
	IsInvestor			bool			`gorm:"default:false"`
}

func (i *Client) CastgRPC() *pb.Client {
	return &pb.Client{FiscalIdentity:i.FiscalIdentity,Name:i.Name,Surname:i.Surname,Balance:i.Balance,IsInvestor:i.IsInvestor}
}

func CastGorm(i *pb.Client) *Client {
	return &Client{FiscalIdentity:i.FiscalIdentity,Name:i.Name,Surname:i.Surname,Balance:i.Balance,IsInvestor:i.IsInvestor}
}

// enum InvoiceState and SellOrder definition		-		begin
type InvoiceState string
type SellOrderState string

const (
	InvoiceFinancingSearch	InvoiceState = "financing search"
	InvoiceRejected			InvoiceState = "rejected"
	InvoiceFinanced			InvoiceState = "financed"

	SellOrderOngoing		SellOrderState = "ongoing"
	SellOrderReversed		SellOrderState = "reversed"
	SellOrderLocked			SellOrderState = "locked"
	SellOrderCommitted		SellOrderState = "committed"
)
func (p *InvoiceState) Scan(value interface{}) error {
	*p = InvoiceState(value.([]byte))
	return nil
}

func (p InvoiceState) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *SellOrderState) Scan(value interface{}) error {
	*p = SellOrderState(value.([]byte))
	return nil
}

func (p SellOrderState) Value() (driver.Value, error) {
	return string(p), nil
}
// enum InvoiceState and SellOrder definition		-		end

type Invoice struct {
	gorm.Model
	ClientID			uint
	Amount				string			`gorm:"type:amount_type"`
	State				InvoiceState	`gorm:"type:invoice_state"`
}

func (i *Invoice) CastgRPC() *pb.Invoice {
	return &pb.Invoice{ClientId:uint64(i.ClientID),Amount:i.Amount,State:string(i.State)}
}


type SellOrder struct {
	gorm.Model
	InvoiceID			uint
	Size				string			`gorm:"type:amount_type"`
	Amount				string			`gorm:"type:amount_type"`
	FinanSize			string			`gorm:"type:amount_type_calc;default:0"`
	FinanAmount			string			`gorm:"type:amount_type_calc;default:0"`
	State				SellOrderState	`gorm:"type:sell_order_state"`
}


type Ledger struct {
	gorm.Model
	InvestorID			uint
	SellOrderID			uint
	Size				string			`gorm:"type:amount_type"`
	Amount				string			`gorm:"type:amount_type"`
	Balance				string			`gorm:"type:amount_type"`
	discount			string			`gorm:"type:discount_type"`
	expectedProfit		string			`gorm:"type:amount_type_calc"`
	IsAdjusted			bool			`gorm:"default:false"`				// adjusted to feet the size of an invoice.
}
