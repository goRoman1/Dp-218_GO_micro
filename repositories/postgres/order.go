package postgres

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"context"
)

//OrderRepoDb is a repository for database connection.
type OrderRepoDb struct {
	db repositories.AnyDatabase
}

//NewOrderRepoDB creates new OrderRepoDb
func NewOrderRepoDB(db repositories.AnyDatabase) *OrderRepoDb {
	return &OrderRepoDb{db}
}

//CreateOrder creates a new order in the database table 'orders'.
func (ordb *OrderRepoDb) CreateOrder(user models.User, scooterID, startID, endID int, distance float64) (models.Order,
	error) {
	var order = models.Order{}
	order.UserID = user.ID
	order.ScooterID = scooterID
	order.StatusStartID = startID
	order.StatusEndID = endID
	order.Distance = distance

	querySQL := `INSERT INTO orders(user_id, scooter_id, status_start_id, status_end_id, distance, amount_cents) 
					VALUES ($1, $2, $3, $4, $5, 0) RETURNING id`
	err := ordb.db.QueryResultRow(context.Background(), querySQL, user.ID, scooterID, startID, endID,
		distance).Scan(&order.ID)
	if err != nil {
		return order, err
	}

	return order, nil
}

//UpdateOrder updates exist order and returns the new one.
func (ordb *OrderRepoDb) UpdateOrder(orderID int, orderData models.Order) (models.Order, error) {
	order := models.Order{}
	querySQL := `UPDATE orders 
					SET user_id=$1, scooter_id=$2, status_start_id=$3, status_end_id=$5, distance=$5, amount=$6
					WHERE id=$7 RETURNING id, user_id, scooter_id, status_start_id, status_end_id, distance, amount;`

	err := ordb.db.QueryResultRow(context.Background(), querySQL,
		orderData.UserID, orderData.ScooterID, orderData.StatusStartID, orderData.StatusEndID, orderData.Distance,
		orderData.Amount, orderID).
		Scan(&order.ID, &order.UserID, &order.ScooterID, &order.StatusStartID, &order.StatusEndID, &order.Distance,
			&order.Amount)
	if err != nil {
		return order, err
	}

	return order, nil
}

//DeleteOrder deletes the chosen order.
func (ordb *OrderRepoDb) DeleteOrder(orderID int) error {
	querySQL := `DELETE FROM orders WHERE id = $1;`
	_, err := ordb.db.QueryExec(context.Background(), querySQL, orderID)
	return err
}

//GetAllOrders returns all orders in the database.
func (ordb *OrderRepoDb) GetAllOrders() (*models.OrderList, error) {
	orderList := &models.OrderList{}

	querySQL := `SELECT * from orders`
	rows, err := ordb.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return orderList, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.ScooterID, &order.StatusStartID, &order.StatusEndID,
			&order.Distance, &order.Amount)
		if err != nil {
			return orderList, err
		}
		orderList.Orders = append(orderList.Orders, order)
	}
	return orderList, nil
}

//GetOrderByID returns exact order by it's ID.
func (ordb *OrderRepoDb) GetOrderByID(orderID int) (models.Order, error) {
	order := models.Order{}

	querySQL := `SELECT * 
					FROM orders
					WHERE id=$1`

	row := ordb.db.QueryResultRow(context.Background(), querySQL, orderID)
	err := row.Scan(&order.ID, &order.UserID, &order.ScooterID, &order.StatusStartID, &order.StatusEndID,
		&order.Distance, &order.Amount)
	if err != nil {
		return order, err
	}

	return order, nil
}

//GetOrdersByUserID returns a list of orders attached with user's ID.
func (ordb *OrderRepoDb) GetOrdersByUserID(userID int) (models.OrderList, error) {
	orderList := models.OrderList{}

	querySQL := `SELECT * 
					FROM orders 
					WHERE user_id=$1`

	rows, err := ordb.db.QueryResult(context.Background(), querySQL, userID)
	if err != nil {
		return orderList, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.ScooterID, &order.StatusStartID, &order.StatusEndID,
			&order.Distance, &order.Amount)
		if err != nil {
			return orderList, err
		}
		orderList.Orders = append(orderList.Orders, order)
	}
	return orderList, nil
}

//GetOrdersByScooterID returns a list of orders attached with scooter's ID.
func (ordb *OrderRepoDb) GetOrdersByScooterID(scooterID int) (models.OrderList, error) {
	orderList := models.OrderList{}
	querySQL := `SELECT * 
					FROM orders 
					WHERE scooter_id=$1`

	rows, err := ordb.db.QueryResult(context.Background(), querySQL, scooterID)
	if err != nil {
		return orderList, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.ScooterID, &order.StatusStartID, &order.StatusEndID,
			&order.Distance, &order.Amount)
		if err != nil {
			return orderList, err
		}
		orderList.Orders = append(orderList.Orders, order)
	}
	return orderList, nil
}

//GetScooterMileageByID returns total mileage for the chosen scooter.
func (ordb *OrderRepoDb) GetScooterMileageByID(scooterID int) (float64, error) {
	var mileageKm float64
	querySQL := `SELECT SUM(distance) 
					FROM orders 
					WHERE scooter_id=$1`

	row := ordb.db.QueryResultRow(context.Background(), querySQL, scooterID)
	err := row.Scan(&mileageKm)
	if err != nil {
		return 0, err
	}
	mileageKm = mileageKm / 1000

	return mileageKm, nil

}

//GetUserMileageByID returns total mileage for the chosen user.
func (ordb *OrderRepoDb) GetUserMileageByID(userID int) (float64, error) {
	var mileageKm float64
	querySQL := `SELECT SUM(distance) 
					FROM orders 
					WHERE user_id=$1`

	row := ordb.db.QueryResultRow(context.Background(), querySQL, userID)
	err := row.Scan(&mileageKm)
	if err != nil {
		return 0, err
	}
	mileageKm = mileageKm / 1000

	return mileageKm, nil
}
