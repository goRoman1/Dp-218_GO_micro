package usecases

import (
	"Dp218GO/models"
	"time"
)

//AccountUsecases - interface for user accounting usecases
type AccountUsecases interface {
	MoneyFromCents(cents int) models.Money
	CentsFromMoney(money models.Money) int
	CalculateMoneyAmountByDate(account models.Account, byTime time.Time) (models.Money, error)
	CalculateProfitForPeriod(account models.Account, start, end time.Time) (models.Money, error)
	CalculateLossForPeriod(account models.Account, start, end time.Time) (models.Money, error)
	AddMoneyToAccount(account models.Account, amountCents int) error
	TakeMoneyFromAccount(account models.Account, amountCents int) error
}
