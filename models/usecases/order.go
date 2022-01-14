package usecases

import "Dp-218_GO_micro/models"

type OrderUseCases interface {
	CountTripDistance(order models.Order) (int, error)
	CountTripAmountMoney(order models.Order) (int, error)
	CompleteOrder(order *models.Order) error
}
