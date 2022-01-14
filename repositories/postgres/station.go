package postgres

import (
	"Dp-218_GO_micro/models"
	"Dp-218_GO_micro/repositories"
	"context"
)

type StationRepoDB struct {
	db repositories.AnyDatabase
}

func NewStationRepoDB(db repositories.AnyDatabase) *StationRepoDB {
	return &StationRepoDB{db}
}

func (pg *StationRepoDB) GetAllStations() (*models.StationList, error) {
	list := &models.StationList{}

	querySQL := `SELECT * FROM scooter_stations ORDER BY id;`
	rows, err := pg.db.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

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

func (pg *StationRepoDB) AddStation(station *models.Station) error {
	var id int
	querySQL := `INSERT INTO scooter_stations(id, name, is_active, latitude, longitude) 
		VALUES($1, $2, $3, $4, $5)
		RETURNING id;`
	err := pg.db.QueryResultRow(context.Background(), querySQL, station.ID, station.Name, station.IsActive, &station.Latitude, &station.Longitude).Scan(&id)
	if err != nil {
		return err
	}
	station.ID = id
	return nil
}

func (pg *StationRepoDB) GetStationById(stationId int) (models.Station, error) {
	station := models.Station{}

	querySQL := `SELECT * FROM scooter_stations WHERE id = $1;`
	row := pg.db.QueryResultRow(context.Background(), querySQL, stationId)
	err := row.Scan(&station.ID, &station.Name, &station.IsActive, &station.Latitude, &station.Longitude)

	return station, err
}

func (pg *StationRepoDB) DeleteStation(stationId int) error {
	querySQL := `DELETE FROM scooter_stations WHERE id = $1;`
	_, err := pg.db.QueryExec(context.Background(), querySQL, stationId)
	return err
}

func (pg *StationRepoDB) UpdateStation(stationId int, stationData models.Station) (models.Station, error) {
	station := models.Station{}
	querySQL := `UPDATE scooter_stations 
		SET is_active=$1, name=$2, latitude=$3, longitude=$4
		WHERE id=$5
		RETURNING id, is_active, name, latitude, longitude;`
	err := pg.db.QueryResultRow(context.Background(), querySQL, stationData.IsActive, stationData.Name, stationData.Latitude, stationData.Longitude, stationId).Scan(&station.ID, &station.IsActive, &station.Name, &station.Latitude, &station.Longitude)
	if err != nil {
		return station, err
	}
	return station, nil
}
