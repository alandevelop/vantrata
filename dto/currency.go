package dto

import (
	"github.com/google/uuid"
)

type Currency struct {
	Id   uuid.UUID
	Code CurrencyCode
}

type CurrencyCode string

const (
	CurrencyCodeEUR CurrencyCode = "EUR"
	CurrencyCodeUSD              = "USD"
)

func (s CurrencyCode) String() string {
	return string(s)
}

var CurrencyEUR = Currency{
	Id:   uuid.MustParse("0de138b8-a5e9-4944-baa9-2eb692de14ea"),
	Code: CurrencyCodeEUR,
}
