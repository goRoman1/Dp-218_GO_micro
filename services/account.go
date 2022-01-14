package services

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"errors"
	"time"
)

// constants for income & outcome payment types
const (
	PayIncomeTypeID  = 2
	PayOutcomeTypeID = 3
)

var ErrNotEnoughMoneyToTake = errors.New("can't take more money than you have")

// AccountService - structure for implementing accounting service
type AccountService struct {
	repoAccount            repositories.AccountRepo
	repoAccountTransaction repositories.AccountTransactionRepo
	repoPaymentType        repositories.PaymentTypeRepo
	clock                  Clock
}

type transactionsWithIncome struct {
	Transaction models.AccountTransaction
	IsIncome    bool
}

// NewAccountService - initialization of AccountService
func NewAccountService(repoAccount repositories.AccountRepo,
	repoAccountTransaction repositories.AccountTransactionRepo, repoPaymentType repositories.PaymentTypeRepo, clock Clock) *AccountService {

	return &AccountService{repoAccount, repoAccountTransaction,
		repoPaymentType, clock}
}

// GetAccountsByOwner - get user accounts list by user
func (accserv *AccountService) GetAccountsByOwner(user models.User) (*models.AccountList, error) {
	return accserv.repoAccount.GetAccountsByOwner(user)
}

// GetAccountByID - get account by ID
func (accserv *AccountService) GetAccountByID(accountID int) (models.Account, error) {
	return accserv.repoAccount.GetAccountByID(accountID)
}

// GetAccountByNumber - get account by its number
func (accserv *AccountService) GetAccountByNumber(number string) (models.Account, error) {
	return accserv.repoAccount.GetAccountByNumber(number)
}

// AddAccount - add account record
func (accserv *AccountService) AddAccount(account *models.Account) error {
	return accserv.repoAccount.AddAccount(account)
}

// UpdateAccount - update account record
func (accserv *AccountService) UpdateAccount(accountID int, accountData models.Account) (models.Account, error) {
	return accserv.repoAccount.UpdateAccount(accountID, accountData)
}

// GetAccountTransactionByID - get money transaction info by its ID
func (accserv *AccountService) GetAccountTransactionByID(transId int) (models.AccountTransaction, error) {
	return accserv.repoAccountTransaction.GetAccountTransactionByID(transId)
}

// AddAccountTransaction - add money transaction record
func (accserv *AccountService) AddAccountTransaction(accountTransaction *models.AccountTransaction) error {
	return accserv.repoAccountTransaction.AddAccountTransaction(accountTransaction)
}

// GetAccountTransactions - get money transactions for given accounts
func (accserv *AccountService) GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error) { //nolint:lll
	return accserv.repoAccountTransaction.GetAccountTransactions(accounts...)
}

// GetAccountTransactionsInTimePeriod - get money transactions for given accounts from start to end time
func (accserv *AccountService) GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error) { //nolint:lll
	return accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(start, end, accounts...)
}

// GetAccountTransactionsByOrder - get money transactions for given order
func (accserv *AccountService) GetAccountTransactionsByOrder(order models.Order) (*models.AccountTransactionList, error) { //nolint:lll
	return accserv.repoAccountTransaction.GetAccountTransactionsByOrder(order)
}

// GetAccountTransactionsByPaymentType - get money transactions for given accounts & given payment type
func (accserv *AccountService) GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ...models.Account) (*models.AccountTransactionList, error) { //nolint:lll
	return accserv.repoAccountTransaction.GetAccountTransactionsByPaymentType(paymentType, accounts...)
}

// GetPaymentTypeByID - get payment type by its ID
func (accserv *AccountService) GetPaymentTypeByID(paymentTypeId int) (models.PaymentType, error) {
	return accserv.repoPaymentType.GetPaymentTypeById(paymentTypeId)
}

// CalculateMoneyAmountByDate - count money total for given account by given time
func (accserv *AccountService) CalculateMoneyAmountByDate(account models.Account, byTime time.Time) (models.Money, error) {
	transactionsUpToDate, err := accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(time.UnixMilli(0), byTime, account)
	if err != nil {
		return models.Money{}, err
	}
	var amountCalculated int
	for _, trans := range transactionsUpToDate.AccountTransactions {
		if trans.AccountFrom.ID == account.ID {
			amountCalculated -= trans.AmountCents
		}
		if trans.AccountTo.ID == account.ID {
			amountCalculated += trans.AmountCents
		}
	}
	return accserv.MoneyFromCents(amountCalculated), nil
}

// CalculateProfitForPeriod - count profit for given period from start to end time
func (accserv *AccountService) CalculateProfitForPeriod(account models.Account, start, end time.Time) (models.Money, error) { //nolint:lll
	transactionsUpToDate, err := accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(start, end, account)
	if err != nil {
		return models.Money{}, err
	}
	var amountCalculated int
	for _, trans := range transactionsUpToDate.AccountTransactions {
		if trans.AccountTo.ID == account.ID {
			amountCalculated += trans.AmountCents
		}
	}
	return accserv.MoneyFromCents(amountCalculated), nil
}

// CalculateLossForPeriod - count loss for given period from start to end time
func (accserv *AccountService) CalculateLossForPeriod(account models.Account, start, end time.Time) (models.Money, error) { //nolint:lll
	transactionsUpToDate, err := accserv.repoAccountTransaction.GetAccountTransactionsInTimePeriod(start, end, account)
	if err != nil {
		return models.Money{}, err
	}
	var amountCalculated int
	for _, trans := range transactionsUpToDate.AccountTransactions {
		if trans.AccountFrom.ID == account.ID {
			amountCalculated += trans.AmountCents
		}
	}
	return accserv.MoneyFromCents(amountCalculated), nil
}

// AddMoneyToAccount - new transaction record to add money to given account
func (accserv *AccountService) AddMoneyToAccount(account models.Account, amountCents int) error {
	paymentType, err := accserv.repoPaymentType.GetPaymentTypeById(PayIncomeTypeID)
	if err != nil {
		return err
	}

	currentTime := accserv.clock.Now()

	accTransaction := &models.AccountTransaction{
		DateTime:    currentTime,
		PaymentType: paymentType,
		AccountFrom: models.Account{},
		AccountTo:   account,
		Order:       models.Order{},
		AmountCents: amountCents}

	return accserv.repoAccountTransaction.AddAccountTransaction(accTransaction)
}

// TakeMoneyFromAccount - new transaction record to get money from given account
func (accserv *AccountService) TakeMoneyFromAccount(account models.Account, amountCents int) error {
	paymentType, err := accserv.repoPaymentType.GetPaymentTypeById(PayOutcomeTypeID)
	if err != nil {
		return err
	}
	currentTime := accserv.clock.Now()
	totalMoney, err := accserv.CalculateMoneyAmountByDate(account, currentTime)
	if err != nil {
		return err
	}
	if accserv.CentsFromMoney(totalMoney) < amountCents {
		return ErrNotEnoughMoneyToTake
	}

	accTransaction := &models.AccountTransaction{
		DateTime:    currentTime,
		PaymentType: paymentType,
		AccountFrom: account,
		AccountTo:   models.Account{},
		Order:       models.Order{},
		AmountCents: amountCents}

	return accserv.repoAccountTransaction.AddAccountTransaction(accTransaction)
}

// MoneyFromCents - convert cents to Money struct (dollars, cents)
func (accserv *AccountService) MoneyFromCents(cents int) models.Money {
	coefCents := 1
	if cents < 0 {
		coefCents = -1
	}
	return models.Money{
		Dollars: cents / 100,
		Cents:   coefCents * cents % 100,
	}
}

// CentsFromMoney - convert Money struct (dollars, cents) to cents
func (accserv *AccountService) CentsFromMoney(money models.Money) int {
	return money.Dollars*100 + money.Cents
}

// GetAccountOutputStructByID - get more convenient structure for given account by its ID
func (accserv *AccountService) GetAccountOutputStructByID(accId int) (interface{}, error) {
	account, err := accserv.GetAccountByID(accId)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	moneyTotal, err := accserv.CalculateMoneyAmountByDate(account, now)
	if err != nil {
		return nil, err
	}

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthIncome, err := accserv.CalculateProfitForPeriod(account, monthStart, now)
	if err != nil {
		return nil, err
	}
	monthOutcome, err := accserv.CalculateLossForPeriod(account, monthStart, now)
	if err != nil {
		return nil, err
	}
	monthTransactions, err := accserv.GetAccountTransactionsInTimePeriod(monthStart, now, account)
	if err != nil {
		return nil, err
	}
	totalMonth := accserv.CentsFromMoney(monthIncome) - accserv.CentsFromMoney(monthOutcome)

	return struct {
		ID                  int
		Number              string
		Name                string
		TotalAmount         models.Money
		MonthlyIncome       models.Money
		MonthlyOutcome      models.Money
		MonthlyTransactions []transactionsWithIncome
		TotalMonthAmount    models.Money
	}{
		ID:                  account.ID,
		Number:              account.Number,
		Name:                account.Name,
		TotalAmount:         moneyTotal,
		MonthlyIncome:       monthIncome,
		MonthlyOutcome:      monthOutcome,
		MonthlyTransactions: addIncomeToTransactions(monthTransactions.AccountTransactions, account),
		TotalMonthAmount:    accserv.MoneyFromCents(totalMonth),
	}, nil
}

func addIncomeToTransactions(transactions []models.AccountTransaction, account models.Account) []transactionsWithIncome { //nolint:lll
	result := make([]transactionsWithIncome, len(transactions))
	for i := 0; i < len(transactions); i++ {
		result[i].Transaction = transactions[i]
		result[i].IsIncome = account.ID == transactions[i].AccountTo.ID
	}
	return result
}
