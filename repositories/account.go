//go:generate mockgen -source=account.go -destination=../repositories/mock/mock_account.go -package=mock
package repositories

import (
	"Dp-218_GO_micro/models"
	"time"
)

// AccountRepo - interface for money account repository
type AccountRepo interface {
	GetAccountsByOwner(user models.User) (*models.AccountList, error)
	GetAccountByID(accountID int) (models.Account, error)
	GetAccountByNumber(number string) (models.Account, error)
	AddAccount(account *models.Account) error
	UpdateAccount(accountID int, accountData models.Account) (models.Account, error)
}

// AccountTransactionRepo - interface for money transaction repository
type AccountTransactionRepo interface {
	GetAccountTransactionByID(transID int) (models.AccountTransaction, error)
	AddAccountTransaction(accountTransaction *models.AccountTransaction) error
	GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error)
	GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error) //nolint:lll
	GetAccountTransactionsByOrder(order models.Order) (*models.AccountTransactionList, error)
	GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ...models.Account) (*models.AccountTransactionList, error) //nolint:lll
}

// PaymentTypeRepo - interface for payment type repository
type PaymentTypeRepo interface {
	GetPaymentTypeById(paymentTypeID int) (models.PaymentType, error)
}
