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
	return &pb.Client{Id:uint64(i.ID),FiscalIdentity:i.FiscalIdentity,Name:i.Name,Surname:i.Surname,Balance:i.Balance,IsInvestor:i.IsInvestor}
}

func CastClient(i *pb.Client) *Client {
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
	*p = InvoiceState(value.(string))
	return nil
}

func (p InvoiceState) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *SellOrderState) Scan(value interface{}) error {
	*p = SellOrderState(value.(string))
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
	State				InvoiceState	`gorm:"type:invoice_state;default:'financing search'"`
}

func (i *Invoice) CastgRPC() *pb.Invoice {
	return &pb.Invoice{Id:uint64(i.ID),ClientId:uint64(i.ClientID),Amount:i.Amount,State:string(i.State)}
}

func  CastInvoice(i *pb.Invoice) *Invoice {
	return &Invoice{ClientID:uint(i.ClientId),Amount:i.Amount,State:InvoiceState(i.State)}
}

type SellOrder struct {
	gorm.Model
	InvoiceID			uint
	Size				string			`gorm:"type:amount_type"`
	Amount				string			`gorm:"type:amount_type"`
	Discount			string			`gorm:"->;type:discount_type generated always as (100 - (amount/size)*100) stored"`
	FinanSize			string			`gorm:"type:amount_type_calc;default:0"`
	FinanAmount			string			`gorm:"type:amount_type_calc;default:0"`
	State				SellOrderState	`gorm:"type:sell_order_state;default:'ongoing'"`
}

func (i *SellOrder) CastgRPC() *pb.SellOrder {
	return &pb.SellOrder{Id:uint64(i.ID),InvoiceId:uint64(i.InvoiceID),Size:i.Size,Amount:i.Amount,State:string(i.State)}
}

func  CastSellOrder(i *pb.SellOrder) *SellOrder {
	return &SellOrder{InvoiceID:uint(i.InvoiceId),Size:i.Size,Amount:i.Amount,State:SellOrderState(i.State)}
}

type Ledger struct {
	gorm.Model
	InvestorID			uint
	SellOrderID			uint
	Size				string			`gorm:"type:amount_type"`
	Amount				string			`gorm:"type:amount_type"`
	Balance				string			`gorm:"type:amount_type"`
	Discount			string			`gorm:"->;type:discount_type generated always as (100 - (amount/size)*100) stored"`
	ExpectedProfit		string			`gorm:"->;type:amount_type_calc generated always as (size - amount) stored"`
	IsAdjusted			bool			`gorm:"default:false"`				// adjusted to feet the size of an invoice.
}

func (i *Ledger) CastgRPC() *pb.Ledger {
	return &pb.Ledger{	Id:uint64(i.ID),InvestorId:uint64(i.InvestorID),SellOrderId:uint64(i.SellOrderID),
						Size:i.Size,Amount:i.Amount,Balance:i.Balance,CreatedAt:i.CreatedAt.String(),
						ExpectedProfit:i.ExpectedProfit,IsAdjusted:i.IsAdjusted	}
}

func  CastLedger(i *pb.Ledger) *Ledger {

	return &Ledger{	InvestorID:uint(i.InvestorId),SellOrderID:uint(i.SellOrderId),
						Size:i.Size,Amount:i.Amount,Balance:i.Balance,
						ExpectedProfit:i.ExpectedProfit,IsAdjusted:i.IsAdjusted	}
}


