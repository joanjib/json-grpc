package models

import "database/sql/driver"
import "gorm.io/gorm"

type Client struct {
	gorm.Model
	FiscalIdentity		string		`gorm:"not null"`
	Name				string		`gorm:"not null"`
	Surname				string		`gorm:"not null"`
	Balance				string		`gorm:"type:numeric(11,2);not null;check:balance>0"`
	IsInvestor			bool		`gorm:"default:false"`
}

// enum InvoiceState definition		-		begin
type InvoiceState string

const (
	FinancingSearch	InvoiceState = "financing search"
	Rejected		InvoiceState = "rejected"
	Financed		InvoiceState = "financed"
)
func (p *InvoiceState) Scan(value interface{}) error {
	*p = InvoiceState(value.([]byte))
	return nil
}

func (p InvoiceState) Value() (driver.Value, error) {
	return string(p), nil
}

// enum InvoiceState definition		-		end

type Invoice struct {
	gorm.Model
	ClientID			uint
	Amount				string			`gorm:"type:amount_type"`
	State				InvoiceState	`gorm:"type:invoice_state"`
}
