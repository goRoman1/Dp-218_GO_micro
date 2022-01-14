package postgres

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"context"
	"fmt"
	"strings"
)

// ScooterInitRepoDB - struct representing user  ScooterInit
type ScooterInitRepoDB struct {
	db repositories.AnyDatabase
}

// NewScooterInitRepoDB -  ScooterInit repo initialization
func NewScooterInitRepoDB(db repositories.AnyDatabase) *ScooterInitRepoDB {
	return &ScooterInitRepoDB{db}
}

// GetOwnersScooters - get list of all system scooters that related for current supplier from the DB
func (si *ScooterInitRepoDB) GetOwnersScooters() (*models.SuppliersScooterList, error) {
	suppliersScooterList := &models.SuppliersScooterList{}

	idFromStatuses, err := si.getScooterIDFromStatuses()
	if err != nil {
		return suppliersScooterList, err
	}

	querySQL := `SELECT 
		id, serial_number
		FROM scooters WHERE owner_id = $1
		ORDER BY id DESC;`
	rows, err := si.db.QueryResult(context.Background(), querySQL, userId)
	if err != nil {
		return suppliersScooterList, err
	}
	defer rows.Close()
	for rows.Next() {
		var suppliersScooter models.SuppliersScooter
		err := rows.Scan(&suppliersScooter.ID, &suppliersScooter.SerialNumber)
		if err != nil {
			return suppliersScooterList, err
		}

		if si.findStatusInTheList(idFromStatuses, suppliersScooter.ID) {
			continue
		}

		suppliersScooterList.Scooters = append(suppliersScooterList.Scooters, suppliersScooter)
	}
	return suppliersScooterList, nil
}

//  getScooterIDFromStatuses - get scooter IDs from table of scooter statuses
func (si *ScooterInitRepoDB) getScooterIDFromStatuses() (*models.ScooterIDsStatusesList, error) {
	list := &models.ScooterIDsStatusesList{}
	querySQL := `SELECT scooter_id FROM scooter_statuses WHERE can_be_rent=true ;`
	rows, err := si.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var scooterId models.ScooterIDsStatuses
		err := rows.Scan(&scooterId.ID)
		if err != nil {
			return list, err
		}

		list.ScooterIDsStatusesList = append(list.ScooterIDsStatusesList, scooterId)
	}
	return list, nil
}

// findStatusInTheList -  check if there are statuses for current scooter
func (si *ScooterInitRepoDB) findStatusInTheList(scooterIds *models.ScooterIDsStatusesList, scooterId int) bool {
	for _, v := range scooterIds.ScooterIDsStatusesList {
		if v.ID == scooterId {
			return true
		}
	}
	return false
}

// GetActiveStations - get all station data for stations that is active
func (si *ScooterInitRepoDB) GetActiveStations() (*models.StationList, error) {
	list := &models.StationList{}
	querySQL := `SELECT * FROM scooter_stations WHERE is_active=true ORDER BY id;`
	rows, err := si.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var station models.Station
		err := rows.Scan(&station.ID, &station.Name, &station.IsActive, &station.Latitude, &station.Longitude)
		if err != nil {
			return list, err
		}

		list.Station = append(list.Station, station)
	}
	return list, nil
}

//AddStatusesToScooters - add coordinates and statuses to scooter statuses
func (si *ScooterInitRepoDB) AddStatusesToScooters(scooterIds []int, station models.Station) error {
	batteryRemain := 100

	valueStrings := make([]string, 0, len(scooterIds))
	valueArgs := make([]interface{}, 0, len(scooterIds)*6)
	for i, scooter := range scooterIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		valueArgs = append(valueArgs, scooter)
		valueArgs = append(valueArgs, batteryRemain)
		valueArgs = append(valueArgs, station.ID)
		valueArgs = append(valueArgs, station.Latitude)
		valueArgs = append(valueArgs, station.Longitude)
		valueArgs = append(valueArgs, true)
	}

	stmt := fmt.Sprintf("INSERT INTO scooter_statuses(scooter_id, battery_remain, station_id, latitude, longitude, can_be_rent) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := si.db.QueryExec(context.Background(), stmt, valueArgs...); err != nil {
		fmt.Println("Unable to insert due to: ", err)
		return err
	}

	return nil
}
