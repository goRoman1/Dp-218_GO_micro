package postgres

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"context"
	"strconv"
	"strings"
	"time"
)

// AccountRepoDB - struct for implementing account repository
type AccountRepoDB struct {
	userRepo *UserRepoDB
	db       repositories.AnyDatabase
}

// NewAccountRepoDB - init of new Account repo
func NewAccountRepoDB(userRepo *UserRepoDB, db repositories.AnyDatabase) *AccountRepoDB {
	return &AccountRepoDB{userRepo, db}
}

// GetAccountsByOwner - gets list of accounts of given user from the DB
func (accdb *AccountRepoDB) GetAccountsByOwner(user models.User) (*models.AccountList, error) {
	list := &models.AccountList{}

	querySQL := `SELECT id, name, number FROM accounts WHERE owner_id = $1;`
	rows, err := accdb.db.QueryResult(context.Background(), querySQL, user.ID)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.ID, &account.Name, &account.Number)
		if err != nil {
			return list, err
		}
		account.User = user
		list.Accounts = append(list.Accounts, account)
	}

	return list, nil
}

// GetAccountByID - gets account entity by ID from the DB
func (accdb *AccountRepoDB) GetAccountByID(accountID int) (models.Account, error) {
	account := models.Account{}

	querySQL := `SELECT id, name, number, owner_id FROM accounts WHERE id = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, accountID)
	var userID int
	err := row.Scan(&account.ID, &account.Name, &account.Number, &userID)
	if err != nil {
		return account, err
	}
	account.User, err = accdb.userRepo.GetUserByID(userID)

	return account, err
}

// GetAccountByNumber - gets account entity by account number from the DB
func (accdb *AccountRepoDB) GetAccountByNumber(number string) (models.Account, error) {
	account := models.Account{}

	querySQL := `SELECT id, name, number, owner_id FROM accounts WHERE number = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, number)
	var userID int
	err := row.Scan(&account.ID, &account.Name, &account.Number, &userID)
	if err != nil {
		return account, err
	}
	account.User, err = accdb.userRepo.GetUserByID(userID)

	return account, err
}

// AddAccount - creates new account in the DB based on given entity
func (accdb *AccountRepoDB) AddAccount(account *models.Account) error {
	var id int
	querySQL := `INSERT INTO accounts(name, number, owner_id) VALUES($1, $2, $3) RETURNING id;`
	err := accdb.db.QueryResultRow(context.Background(), querySQL, account.Name, account.Number, account.User.ID).
		Scan(&id)
	if err != nil {
		return err
	}
	account.ID = id

	return nil
}

// UpdateAccount - updates account in the DB by ID and given entity
func (accdb *AccountRepoDB) UpdateAccount(accountID int, accountData models.Account) (models.Account, error) {
	account := models.Account{}
	querySQL := `UPDATE accounts 
		SET name=$1, number=$2, owner_id=$3 
		WHERE id=$4 RETURNING id, name, number, owner_id;`
	var userID int
	err := accdb.db.QueryResultRow(context.Background(), querySQL,
		accountData.Name, accountData.Number, accountData.User.ID, accountID).
		Scan(&account.ID, &account.Name, &account.Number, &userID)
	if err != nil {
		return account, err
	}
	account.User, err = accdb.userRepo.GetUserByID(userID)
	if err != nil {
		return account, err
	}

	return account, nil
}

// GetAccountTransactionByID - gets transaction information from the DB by its ID
func (accdb *AccountRepoDB) GetAccountTransactionByID(transID int) (models.AccountTransaction, error) {
	accountTransaction := models.AccountTransaction{}

	querySQL := `SELECT 
		id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents 
		FROM account_transactions 
		WHERE id = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, transID)
	var paymentID int
	var accFromID, accToId int
	var orderId int
	err := row.Scan(&accountTransaction.ID, &accountTransaction.DateTime, &paymentID, &accFromID, &accToId, orderId, &accountTransaction.AmountCents)
	if err != nil {
		return accountTransaction, err
	}

	err = addTransactionComplexFields(accdb, &accountTransaction, paymentID, accFromID, accToId, orderId)
	if err != nil {
		return accountTransaction, err
	}

	return accountTransaction, err
}

func addTransactionComplexFields(accdb *AccountRepoDB, accountTransaction *models.AccountTransaction, paymentID, accFromID, accToID, orderId int) error {
	var err error
	accountTransaction.PaymentType, err = accdb.GetPaymentTypeById(paymentID)
	if err != nil {
		return err
	}
	accountTransaction.AccountFrom, err = accdb.GetAccountByID(accFromID)
	if err != nil && accFromID != 0 {
		return err
	}
	accountTransaction.AccountTo, err = accdb.GetAccountByID(accToID)
	if err != nil && accToID != 0 {
		return err
	}
	accountTransaction.Order, err = models.Order{}, nil //TODO: refactor when Orders implemented
	if err != nil && orderId != 0 {
		return err
	}
	return nil
}

//AddAccountTransaction - creates transaction record in the DB based on given entity
func (accdb *AccountRepoDB) AddAccountTransaction(accountTransaction *models.AccountTransaction) error {
	var id int
	querySQL := `INSERT INTO 
		account_transactions(date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents) 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id;`
	err := accdb.db.QueryResultRow(context.Background(), querySQL, accountTransaction.DateTime,
		accountTransaction.PaymentType.ID, accountTransaction.AccountFrom.ID, accountTransaction.AccountTo.ID,
		accountTransaction.Order.ID, accountTransaction.AmountCents).Scan(&id)
	if err != nil {
		return err
	}
	accountTransaction.ID = id

	return nil
}

func getTransactionsBySomeQuery(accdb *AccountRepoDB, querySQL string, params ...interface{}) (*models.AccountTransactionList, error) {
	list := &models.AccountTransactionList{}
	rows, err := accdb.db.QueryResult(context.Background(), querySQL, params...)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	type additionalTransData struct {
		paymentID int
		accFromID int
		accToID   int
		orderID   int
	}
	var transAdditionalData = make(map[models.AccountTransaction]additionalTransData)

	for rows.Next() {
		var accountTransaction models.AccountTransaction
		var paymentID int
		var accFromID, accToId int
		var orderID int
		err := rows.Scan(&accountTransaction.ID, &accountTransaction.DateTime,
			&paymentID, &accFromID, &accToId, &orderID, &accountTransaction.AmountCents)
		if err != nil {
			return list, err
		}

		transAdditionalData[accountTransaction] = additionalTransData{paymentID, accFromID, accToId, orderID}
	}

	for key, value := range transAdditionalData {
		err = addTransactionComplexFields(accdb, &key, value.paymentID, value.accFromID, value.accToID, value.orderID)
		if err != nil {
			return list, err
		}

		list.AccountTransactions = append(list.AccountTransactions, key)
	}

	return list, nil
}

// GetAccountTransactions - gets list of money transactions for given accounts from the DB
func (accdb *AccountRepoDB) GetAccountTransactions(accounts ...models.Account) (*models.AccountTransactionList, error) {
	querySQL := `SELECT 
		id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents 
		FROM account_transactions`
	var params []interface{}
	for i, acc := range accounts {
		if i == 0 {
			querySQL += ` WHERE FALSE`
		}
		paramIndex := strconv.Itoa(i + 1)
		querySQL += ` OR account_from_id = $` + paramIndex + ` OR account_to_id = $` + paramIndex
		params = append(params, acc.ID)
	}
	querySQL += `;`

	return getTransactionsBySomeQuery(accdb, querySQL, params...)
}

// GetAccountTransactionsInTimePeriod - gets list of money transactions for given accounts from start to end time from the DB
func (accdb *AccountRepoDB) GetAccountTransactionsInTimePeriod(start time.Time, end time.Time, accounts ...models.Account) (*models.AccountTransactionList, error) {
	querySQL := `SELECT 
		id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents 
		FROM account_transactions
		WHERE date_time>=$1 AND date_time<=$2`
	var params []interface{}
	params = append(params, start)
	params = append(params, end)
	var accountCondition = make([]string, len(accounts))
	for i, acc := range accounts {
		accountCondition[i] = `$` + strconv.Itoa(i+3)
		params = append(params, acc.ID)
	}
	if len(accounts) > 0 {
		conditionStr := strings.Join(accountCondition, ", ")
		querySQL += ` AND (account_from_id IN (` + conditionStr + `) OR account_to_id IN (` + conditionStr + `))`
	}
	querySQL += `;`

	return getTransactionsBySomeQuery(accdb, querySQL, params...)
}

// GetAccountTransactionsByOrder - gets list of money transactions for given order from the DB
func (accdb *AccountRepoDB) GetAccountTransactionsByOrder(order models.Order) (*models.AccountTransactionList, error) {
	querySQL := `SELECT 
		id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents 
		FROM account_transactions
		WHERE order_id=$1;`

	return getTransactionsBySomeQuery(accdb, querySQL, order.ID)
}

// GetAccountTransactionsByPaymentType - gets list of money transactions for given accounts & payment type from the DB
func (accdb *AccountRepoDB) GetAccountTransactionsByPaymentType(paymentType models.PaymentType, accounts ...models.Account) (*models.AccountTransactionList, error) {
	querySQL := `SELECT 
		id, date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents 
		FROM account_transactions
		WHERE payment_type_id=$1`
	var params []interface{}
	params = append(params, paymentType.ID)
	var accountCondition = make([]string, len(accounts))
	for i, acc := range accounts {
		accountCondition[i] = `$` + strconv.Itoa(i+2)
		params = append(params, acc.ID)
	}
	if len(accounts) > 0 {
		conditionStr := strings.Join(accountCondition, ", ")
		querySQL += ` AND (account_from_id IN (` + conditionStr + `) OR account_to_id IN (` + conditionStr + `))`
	}
	querySQL += `;`

	return getTransactionsBySomeQuery(accdb, querySQL, params...)
}

// GetPaymentTypeById - gets payment type entity by ID from the DB
func (accdb *AccountRepoDB) GetPaymentTypeById(paymentTypeId int) (models.PaymentType, error) {
	paymentType := models.PaymentType{}

	querySQL := `SELECT id, name FROM payment_types WHERE id = $1;`
	row := accdb.db.QueryResultRow(context.Background(), querySQL, paymentTypeId)
	err := row.Scan(&paymentType.ID, &paymentType.Name)

	return paymentType, err
}
