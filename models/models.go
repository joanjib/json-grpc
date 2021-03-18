package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	FiscalIdentity		string		`gorm:"not null"`
	Name				string		`gorm:"not null"`
	Surname				string		`gorm:"not null"`
	Balance				string		`gorm:"type:numeric(11,2);not null;check:balance>0"`
	IsInvestor			bool		`gorm:"default:false"`
}

type Invoice struct {
	gorm.Model
	ClientID			uint
	amount				string		`gorm:"type:numeric(11,2)"`
}
