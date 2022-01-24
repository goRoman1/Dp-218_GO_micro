package postgres

import (
	"Dp218GO/models"
	"Dp218GO/repositories"
	"context"
)

// SupMicroRepoDB - struct representing user  ScooterInit
type SupMicroRepoDB struct {
	db repositories.AnyDatabase
}

// NewSupMicroRepoDB -  ScooterInit repo initialization
func NewSupMicroRepoDB(db repositories.AnyDatabase) *SupMicroRepoDB {
	return &SupMicroRepoDB{db}
}

func (repo *SupMicroRepoDB) GetLocations() (*models.LocationList, error) {
	list := &models.LocationList{}

	querySQL := `SELECT * FROM locations ORDER BY id;`
	rows, err := repo.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var location models.Location
		err := rows.Scan(&location.ID, &location.Latitude, &location.Longitude, &location.Label)
		if err != nil {
			return list, err
		}

		list.Location = append(list.Location, location)
	}
	return list, nil
}

func (repo *SupMicroRepoDB) GetStations() (*models.StationList, error) {
	query := `SELECT * FROM locations ORDER BY id;`
	row, err := repo.db.QueryResult(context.Background(), query)

	var result *models.StationList
	if err != nil {
		return result, err
	}
	defer row.Close()
	for row.Next() {
		var station models.Station
		err := row.Scan(&station.ID, &station.IsActive, &station.Latitude, &station.Longitude)
		if err != nil {
			return result, err
		}
		result.Station = append(result.Station, station)
	}
	return result, nil
}

func (repo *SupMicroRepoDB) CreateStationInLocation(station *models.Station, location *models.Location) error {
	query := `INSERT INTO scooter_stations (name, is_active, latitude, longitude)
	VALUES($1, $2, $3, $4)
	RETURNING id`
	row := repo.db.QueryResultRow(context.Background(), query, station.Name, station.IsActive, location.Latitude, location.Longitude)
	err := row.Scan(&station.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SupMicroRepoDB) CreateStation(station *models.Station) error {
	query := `INSERT INTO scooter_stations(name, is_active, latitude, longitude)
	VALUES($1, $2, $3, $4) WHERE 
	RETURNING id`
	row := repo.db.QueryResultRow(context.Background(), query, station.Name, station.IsActive, station.Latitude, station.Longitude)
	err := row.Scan(&station.ID)
	if err != nil {
		return err
	}

	return nil
}
