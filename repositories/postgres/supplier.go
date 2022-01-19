package postgres

import (
	"Dp218GO/models"
	"Dp218GO/repositories"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jszwec/csvutil"
	"io"
	"os"
	"strings"
)

var userId = 9

// SupplierRepoDB - struct representing supplier repository
type SupplierRepoDB struct {
	db repositories.AnyDatabase
}

// NewSupplierRepoDB - supplier repo initialization
func NewSupplierRepoDB(db repositories.AnyDatabase) *SupplierRepoDB {
	return &SupplierRepoDB{db}
}

// GetModels - get list of all system scooter models with payment prices from the DB
func (s *SupplierRepoDB) GetModels() (*models.ScooterModelDTOList, error) {
	modelsOdtList := &models.ScooterModelDTOList{}
	pricesList := &models.SupplierPricesDTOList{}

	pricesList, err := s.getPrices()
	if err != nil {
		return modelsOdtList, err
	}

	querySQL := `SELECT * FROM scooter_models ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return modelsOdtList, err
	}

	for rows.Next() {
		var paymentTypeID int
		var model models.ScooterModelDTO
		err := rows.Scan(&model.ID, &paymentTypeID, &model.ModelName, &model.MaxWeight, &model.Speed)
		if err != nil {
			return modelsOdtList, err
		}

		model.Price, err = s.findSupplierPricesList(pricesList, paymentTypeID, userId)
		if err != nil {
			return modelsOdtList, err
		}

		model.SuppliersScooters, err = s.getSuppliersScootersByModelId(model.ID)
		if err != nil {
			return modelsOdtList, err
		}

		modelsOdtList.ScooterModelsDTO = append(modelsOdtList.ScooterModelsDTO, model)

	}
	return modelsOdtList, nil
}

// SelectModel - get scooter model from the DB by given ID
func (s *SupplierRepoDB) SelectModel(id int) (*models.ScooterModelDTO, error) {
	modelDTO := &models.ScooterModelDTO{}

	querySQL := `SELECT id, payment_type_id, model_name, max_weight, speed  FROM scooter_models WHERE id = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, id)

	var paymentTypeId int
	err := row.Scan(&modelDTO.ID, &paymentTypeId, &modelDTO.ModelName, &modelDTO.MaxWeight, &modelDTO.Speed)
	if err != nil {
		return modelDTO, err
	}

	modelDTO.Price, err = s.getPrice(paymentTypeId, userId)

	return modelDTO, err
}

// AddModel - create scooter model record in the DB based on given entity
func (s *SupplierRepoDB) AddModel(modelData *models.ScooterModelDTO) error {

	paymentTypeId, err := s.addPaymentTypeId(modelData.ModelName)
	if err != nil {
		return err
	}
	var modelId int
	querySQL := `INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed)
	   		VALUES($1, $2, $3, $4)
	   		RETURNING id;`
	err = s.db.QueryResultRow(context.Background(), querySQL, &paymentTypeId, modelData.ModelName, modelData.MaxWeight, modelData.Speed).Scan(&modelId)
	if err != nil {
		return err
	}

	var priceId int
	querySQL = `INSERT INTO supplier_prices(price, payment_type_id, user_id)
	   		VALUES($1, $2, $3)
	   		RETURNING id;`
	err = s.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&priceId)
	if err != nil {
		return err
	}
	return nil
}

//EditPrice - changes the price for the rental of a scooter which is associated with the model
func (s *SupplierRepoDB) EditPrice(modelData *models.ScooterModelDTO) error {
	price := &models.ScooterModelDTO{}
	paymentTypeId, err := s.getPaymentTypeByModelName(modelData.ModelName)
	if err != nil {
		return err
	}

	querySQL := `UPDATE supplier_prices SET price=$1 WHERE payment_type_id = $2 AND user_id = $3 RETURNING price;`
	err = s.db.QueryResultRow(context.Background(), querySQL, modelData.Price, paymentTypeId, userId).Scan(&price.Price)
	if err != nil {
		return err
	}

	return nil
}

// getPrices - reads all data from the table supplier_prices
func (s *SupplierRepoDB) getPrices() (*models.SupplierPricesDTOList, error) {
	list := &models.SupplierPricesDTOList{}

	querySQL := `SELECT * FROM supplier_prices ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var supplierPriceODT models.SupplierPricesDTO
		err := rows.Scan(&supplierPriceODT.ID, &supplierPriceODT.Price, &supplierPriceODT.PaymentTypeID, &supplierPriceODT.UserId)

		if err != nil {
			return list, err
		}

		list.SupplierPricesDTO = append(list.SupplierPricesDTO, supplierPriceODT)
	}
	return list, nil
}

//getPaymentTypeID - selects payment_type by scooter model id
func (s *SupplierRepoDB) getPaymentTypeID(modelId int) (int, error) {
	model := &models.ScooterModel{}

	querySQL := `SELECT payment_type_id FROM scooter_models WHERE id= $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, modelId)
	err := row.Scan(&model.PaymentType.ID)

	return model.PaymentType.ID, err
}

// findSupplierPricesList - find price in given price list by paymentTypeId and user Id
func (s *SupplierRepoDB) findSupplierPricesList(supplierPrice *models.SupplierPricesDTOList, paymentTypeId, userId int) (int, error) {
	for _, v := range supplierPrice.SupplierPricesDTO {
		if v.PaymentTypeID == paymentTypeId && v.UserId == userId {
			return v.Price, nil
		}
	}
	return 0, fmt.Errorf("not found paymentType id=%d", paymentTypeId)
}

// getPrice - selects a specific price by payment-type id and user id
func (s *SupplierRepoDB) getPrice(paymentTypeId, userId int) (int, error) {
	price := models.ScooterModelDTO{}
	querySQL := `SELECT price FROM supplier_prices WHERE payment_type_id = $1 AND user_id = $2;`
	row := s.db.QueryResultRow(context.Background(), querySQL, paymentTypeId, userId)
	err := row.Scan(&price.Price)

	return price.Price, err
}

// addPaymentTypeId - add payment type by model name
func (s *SupplierRepoDB) addPaymentTypeId(modelName string) (int, error) {
	var paymentTypeId int
	querySQL := `INSERT INTO payment_types (name) VALUES ($1) RETURNING id;`
	err := s.db.QueryResultRow(context.Background(), querySQL, modelName).Scan(&paymentTypeId)
	if err != nil {
		return 0, err
	}
	return paymentTypeId, nil
}

// getPaymentTypeByModelName - select payment type by model name
func (s *SupplierRepoDB) getPaymentTypeByModelName(modelName string) (int, error) {
	paymentType := models.PaymentType{}
	querySQL := `SELECT * FROM payment_types WHERE name = $1;`
	row := s.db.QueryResultRow(context.Background(), querySQL, modelName)
	err := row.Scan(&paymentType.ID, &paymentType.Name)
	return paymentType.ID, err
}

//getSuppliersScootersByModelId - get scooters related to scooter model
func (s *SupplierRepoDB) getSuppliersScootersByModelId(modelId int) (models.SuppliersScooterList, error) {
	list := models.SuppliersScooterList{}

	querySQL := `SELECT id, serial_number FROM scooters WHERE model_id = $1 ORDER BY id DESC;`
	rows, err := s.db.QueryResult(context.Background(), querySQL, modelId)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var scooter models.SuppliersScooter
		err := rows.Scan(&scooter.ID, &scooter.SerialNumber)
		if err != nil {
			return list, err
		}

		list.Scooters = append(list.Scooters, scooter)
	}
	return list, nil
}

//AddSuppliersScooter - adds a scooter to the scooter table, its serial number will be displayed in the scooter list
func (s *SupplierRepoDB) AddSuppliersScooter(modelId int, scooterSerial string) error {

	// проверить если уже существует
	var id int
	querySQL := `INSERT INTO scooters(model_id, owner_id, serial_number) 
	   		VALUES($1, $2, $3) 
			ON CONFLICT (serial_number) DO UPDATE
				SET model_id = $4
				RETURNING id;`
	err := s.db.QueryResultRow(context.Background(), querySQL, modelId, userId, scooterSerial, modelId).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

//DeleteSuppliersScooter - removes the scooter and its status from the list of scooters and from the database
func (s *SupplierRepoDB) DeleteSuppliersScooter(id int) error {
	querySQL := `DELETE FROM scooters WHERE id = $1;`
	_, err := s.db.QueryExec(context.Background(), querySQL, id)

	querySQL = `DELETE FROM scooter_statuses WHERE scooter_id = $1;`
	_, err = s.db.QueryExec(context.Background(), querySQL, id)
	return err
}

//ConvertToStruct - converting the received data from .csv into a structure for further work
func (s *SupplierRepoDB) ConvertToStruct(path string) []models.UploadedScooters {

	csvFile, _ := os.Open(path)
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'

	scooterHeader, _ := csvutil.Header(models.UploadedScooters{}, "csv")
	dec, _ := csvutil.NewDecoder(reader, scooterHeader...)

	var fileData []models.UploadedScooters
	for {
		var scooter models.UploadedScooters
		if err := dec.Decode(&scooter); err == io.EOF {
			break
		}
		fileData = append(fileData, scooter)
	}
	return fileData
}

// InsertToDb - enter the data received from the file into the database
func (s *SupplierRepoDB) InsertToDb(modelId int, scooters []models.UploadedScooters) error {
	valueStrings := make([]string, 0, len(scooters))
	valueArgs := make([]interface{}, 0, len(scooters)*3)
	for i, scooter := range scooters {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, modelId)
		valueArgs = append(valueArgs, userId)
		valueArgs = append(valueArgs, scooter.SerialNumber)
	}

	stmt := fmt.Sprintf("INSERT INTO scooters(model_id, owner_id, serial_number) VALUES %s ON CONFLICT (serial_number) DO UPDATE SET model_id = excluded.model_id", strings.Join(valueStrings, ","))
	if _, err := s.db.QueryExec(context.Background(), stmt, valueArgs...); err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return err
	}
	return nil
}
