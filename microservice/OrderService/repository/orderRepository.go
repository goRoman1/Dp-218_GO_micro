package repository

import (
	"OrderService/proto"
	"context"
	"database/sql"
	"fmt"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, info *proto.TripInfo) (*proto.Order, error)
}

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (or *OrderRepo) CreateOrder(ctx context.Context, info *proto.TripInfo) (*proto.Order, error) {
	fmt.Println("Create Order called on Order_micro")
	var order = &proto.Order{}
	order.UserID = info.UserID
	order.ScooterID = info.ScooterID
	order.StatusStartID = info.StatusStartID
	order.StatusEndID = info.StatusEndID

	querySQL := `INSERT INTO orders(user_id, scooter_id, status_start_id, status_end_id) 
					VALUES ($1, $2, $3, $4) RETURNING id`
	err := or.db.QueryRowContext(ctx, querySQL, order.UserID, order.ScooterID, order.StatusStartID, order.StatusEndID).Scan(&order.Id)
	if err != nil {
		return nil, err
	}

	fmt.Println("Order created on Order_service")
	return order, nil
}
