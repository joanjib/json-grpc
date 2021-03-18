package models

import "database/sql/driver"
import "gorm.io/gorm"

type Client struct {
	gorm.Model
	FiscalIdentity		string			`gorm:"not null"`
	Name				string			`gorm:"not null"`
	Surname				string			`gorm:"not null"`
	Balance				string			`gorm:"type:amount_type"`
	IsInvestor			bool			`gorm:"default:false"`
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

type SellOrder struct {
	gorm.Model
	InvoiceID			uint
	Size				string			`gorm:"type:amount_type"`
	Amount				string			`gorm:"type:amount_type"`
	FinanSize			string			`gorm:"type:amount_type_calc;default:0"`
	FinanAmount			string			`gorm:"type:amount_type_calc;default:0"`
	State				SellOrderState	`gorm:"type:sell_order_state"`
}
