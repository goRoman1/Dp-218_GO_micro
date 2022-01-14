package models

import "time"

// PaymentType - entity for payment types
type PaymentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Money - struct to represent money amounts from cents into dollars-cents
type Money struct {
	Dollars int `json:"dollars"`
	Cents   int `json:"cents"`
}

// Account - entity for users banking Accounts
type Account struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number string `json:"number"`
	User   User   `json:"user"`
}

// AccountList - struct representing list of Accounts
type AccountList struct {
	Accounts []Account `json:"accounts"`
}

// AccountTransaction - entity representing single money transaction in the system
type AccountTransaction struct {
	ID          int         `json:"id"`
	DateTime    time.Time   `json:"date_time"`
	PaymentType PaymentType `json:"payment_type"`
	AccountFrom Account     `json:"account_from"`
	AccountTo   Account     `json:"account_to"`
	Order       Order       `json:"order"`
	AmountCents int         `json:"amount_cents"`
}

// AccountTransactionList - struct representing list of money transactions
type AccountTransactionList struct {
	AccountTransactions []AccountTransaction `json:"account_transactions"`
}

// GetAmountInMoney - converts money amount in cents into Money struct
func (accTrans *AccountTransaction) GetAmountInMoney() Money {
	coefCents := 1
	if accTrans.AmountCents < 0 {
		coefCents = -1
	}
	return Money{
		Dollars: accTrans.AmountCents / 100,
		Cents:   coefCents * accTrans.AmountCents % 100,
	}
}
