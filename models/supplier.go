package models

type ScooterModel struct {
	ID          int         `json:"id"`
	PaymentType PaymentType `json:"payment_type"`
	ModelName   string      `json:"model_name"`
	MaxWeight   int         `json:"max_weight"`
	Speed       int         `json:"speed"`
}

type ScooterModelList struct {
	ScooterModels []ScooterModel `json:"scooter_models"`
}

type DbScooters struct {
	ID           int    `json:"id"`
	ModelId      int    `json:"model_id"`
	OwnerId      int    `json:"owner_id"`
	SerialNumber string `json:"serial_number"`
}

type DbScootersList struct {
	DbScooters []DbScooters `json:"db_scooters"`
}

type DbScootersStatuses struct {
	ID            int     `json:"id"`
	BatteryRemain int     `json:"battery_remain"`
	StationId     int     `json:"station_id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	CanBeRent     bool    `json:"can_be_rent"`
}

type DbScootersStatusesList struct {
	DbScootersStatuses []DbScootersStatuses `json:"db_scooters"`
}

type SuppliersScooter struct {
	ID           int    `json:"id"`
	ModelId      int    `json:"model_id"`
	SerialNumber string `json:"serial_number"`
}

type UploadedScooters struct {
	SerialNumber string `json,scv:"serial_number"`
}

type SuppliersScooterList struct {
	Scooters []SuppliersScooter `json:"scooters"`
}

type SupplierPrices struct {
	ID          int         `json:"id"`
	Price       int         `json:"price"`
	PaymentType PaymentType `json:"payment_type"`
	User        User        `json:"user"`
}

type SupplierPricesList struct {
	SupplierPrices []SupplierPrices `json:"supplier_prices_list"`
}

type ScooterModelDTO struct {
	ID                int                  `json:"id"`
	Price             int                  `json:"price"`
	ModelName         string               `json:"model_name"`
	MaxWeight         int                  `json:"max_weight"`
	Speed             int                  `json:"speed"`
	SuppliersScooters SuppliersScooterList `json:"scooters"`
}

type ScooterModelDTOList struct {
	ScooterModelsDTO []ScooterModelDTO `json:"models_dto"`
}

type SupplierPricesDTO struct {
	ID            int `json:"id"`
	Price         int `json:"price"`
	PaymentTypeID int `json:"payment_type_id"`
	UserId        int `json:"user_id"`
}

type SupplierPricesDTOList struct {
	SupplierPricesDTO []SupplierPricesDTO `json:"supplier_prices_odt_list"`
}

type SuppliersScooterStatusesDTO struct {
	ID            int    `json:"id"`
	SerialNumber  string `json:"serial_number"`
	BatteryRemain string `json:"battery_remain"`
	StationId     string `json:"station_id"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	CanBeRent     string `json:"can_be_rent"`
}

type SuppliersScooterStatusesDTOList struct {
	SuppliersScooterStatusesDTOList []SuppliersScooterStatusesDTO `json:"suppliers_scooters_statuses_dto"`
}

type ScootersStationsAllocation struct {
	StationList          StationList          `json:"station_list"`
	SuppliersScooterList SuppliersScooterList `json:"suppliers_scooter"`
}

type ScooterIDsStatuses struct {
	ID int `json:"id"`
}

type ScooterIDsStatusesList struct {
	ScooterIDsStatusesList []ScooterIDsStatuses `json:"scooter_ids_statuses_list"`
}