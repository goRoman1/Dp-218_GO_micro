package postgres

import (
	"Dp218GO/models"
	"Dp218GO/repositories"
	"context"
	"fmt"
)

//ScooterRepoDB is a repository for database connection.
type ScooterRepoDB struct {
	db repositories.AnyDatabase
}

//NewScooterRepoDB creates new ScooterRepoDB
func NewScooterRepoDB(db repositories.AnyDatabase) *ScooterRepoDB {
	return &ScooterRepoDB{db}
}

//GetAllScooters returns the list of all scooters in the database in ScooterDTO view.
func (scdb *ScooterRepoDB) GetAllScooters() (*models.ScooterListDTO, error) {
	scooterList := &models.ScooterListDTO{}

	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.can_be_rent
					FROM scooters as s 
					JOIN scooter_models as sm 
					ON s.model_id=sm.id 
					JOIN scooter_statuses as ss 
					ON s.id=ss.scooter_id 
					ORDER BY s.id`

	rows, err := scdb.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return scooterList, err
	}
	defer rows.Close()

	for rows.Next() {
		var scooter models.ScooterDTO
		err := rows.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain, &scooter.CanBeRent)
		if err != nil {
			return scooterList, err
		}
		scooterList.Scooters = append(scooterList.Scooters, scooter)
	}
	return scooterList, nil
}

func (scdb *ScooterRepoDB) GetAllScootersByStationID(stationID int) (*models.ScooterListDTO, error) {
	scooterList := &models.ScooterListDTO{}

	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.can_be_rent
					FROM scooters as s 
					JOIN scooter_models as sm 
					ON s.model_id=sm.id 
					JOIN scooter_statuses as ss 
					ON s.id=ss.scooter_id 
					WHERE ss.station_id=$1
					ORDER BY s.id`

	rows, err := scdb.db.QueryResult(context.Background(), querySQL, stationID)
	if err != nil {
		return scooterList, err
	}
	defer rows.Close()

	for rows.Next() {
		var scooter models.ScooterDTO
		err := rows.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain, &scooter.CanBeRent)
		if err != nil {
			return scooterList, err
		}
		scooterList.Scooters = append(scooterList.Scooters, scooter)
	}
	return scooterList, nil
}

//GetScooterById returns exact scooter by it's ID.
func (scdb *ScooterRepoDB) GetScooterById(scooterId int) (models.ScooterDTO, error) {
	scooter := models.ScooterDTO{}
	querySQL := `SELECT s.id, sm.max_weight, sm.model_name, ss.battery_remain, ss.can_be_rent
					FROM scooters as s 
					JOIN scooter_models as sm 
					ON s.model_id=sm.id 
					JOIN scooter_statuses as ss 
					ON s.id=ss.scooter_id 
					WHERE s.id=$1`

	row := scdb.db.QueryResultRow(context.Background(), querySQL, scooterId)
	err := row.Scan(&scooter.ID, &scooter.MaxWeight, &scooter.ScooterModel, &scooter.BatteryRemain, &scooter.CanBeRent)
	if err != nil {
		return scooter, err
	}

	return scooter, nil
}

//GetScooterStatus returns the ScooterStatus model of the chosen scooter by its ID.
func (scdb *ScooterRepoDB) GetScooterStatus(scooterID int) (models.ScooterStatus, error) {
	var scooterStatus = models.ScooterStatus{}
	scooter, err := scdb.GetScooterById(scooterID)
	if err != nil {
		fmt.Println(err)
		return models.ScooterStatus{}, err
	}
	scooterStatus.Scooter = scooter

	querySQL := `SELECT battery_remain, latitude, longitude 
					FROM scooter_statuses
					WHERE scooter_id=$1`

	row := scdb.db.QueryResultRow(context.Background(), querySQL, scooterID)
	err = row.Scan(&scooterStatus.BatteryRemain,
		&scooterStatus.Location.Latitude, &scooterStatus.Location.Longitude)
	if err != nil {
		return scooterStatus, err
	}

	return scooterStatus, nil
}

//CreateScooterStatusInRent creates a new record in ScooterStatusesInRent by scooter's ID and returns the
//ScooterStatusInRent model.
func (scdb *ScooterRepoDB) CreateScooterStatusInRent(scooterID int) (models.ScooterStatusInRent, error) {
	var scooterStatusInRent models.ScooterStatusInRent
	scooterStatus, err := scdb.GetScooterStatus(scooterID)
	if err != nil {
		fmt.Println(err)
		return scooterStatusInRent, err
	}

	scooterStatusInRent.Location = scooterStatus.Location

	querySQL := `INSERT INTO scooter_statuses_in_rent(date_time, latitude, longitude) 
					VALUES(now(), $1, $2) RETURNING id, date_time`

	err = scdb.db.QueryResultRow(context.Background(), querySQL, scooterStatus.Location.Latitude,
		scooterStatus.Location.Longitude).Scan(&scooterStatusInRent.ID, &scooterStatusInRent.DateTime)
	if err != nil {
		fmt.Println(err)
		return scooterStatusInRent, err
	}

	return scooterStatusInRent, nil

}

//SendCurrentStatus updates ScooterStatus with given parameters.
func (scdb *ScooterRepoDB) SendCurrentStatus(id, stationID int, lat, lon, battery float64) error {
	var canBeRent bool
	if battery > 10 {
		canBeRent = true
	}

	querySQL := `UPDATE scooter_statuses 
					SET latitude=$1, longitude=$2, battery_remain=$3, can_be_rent=$4, station_id=$5
					WHERE scooter_id=$6`

	row, err := scdb.db.QueryResult(context.Background(), querySQL, lat, lon, battery, canBeRent, stationID, id)
	defer row.Close()
	return err
}
