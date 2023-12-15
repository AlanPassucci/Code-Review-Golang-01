package internal

import (
	"errors"
)

var (
	// ErrServiceVehicleNotFound is returned when no vehicle is found.
	ErrServiceVehiclesNotFound           = errors.New("service: vehicles not found")
	ErrServiceInvalidVehicleBrand        = errors.New("service: invalid vehicle brand")
	ErrServiceInvalidVehicleModel        = errors.New("service: invalid vehicle model")
	ErrServiceInvalidVehicleRegistration = errors.New("service: invalid vehicle registration")
	ErrServiceInvalidVehicleYear         = errors.New("service: invalid vehicle year")
	ErrServiceInvalidVehicleColor        = errors.New("service: invalid vehicle color")
	ErrServiceInvalidVehicleMaxSpeed     = errors.New("service: invalid vehicle max speed")
	ErrServiceInvalidVehicleFuelType     = errors.New("service: invalid vehicle fuel type")
	ErrServiceInvalidVehicleTransmission = errors.New("service: invalid vehicle transmission")
	ErrServiceInvalidVehiclePassengers   = errors.New("service: invalid vehicle passengers")
	ErrServiceInvalidVehicleHeight       = errors.New("service: invalid vehicle height")
	ErrServiceInvalidVehicleWidth        = errors.New("service: invalid vehicle width")
	ErrServiceInvalidVehicleWeight       = errors.New("service: invalid vehicle weight")
	ErrServiceVehicleIdAlreadyExists     = errors.New("service: vehicle id already exists")
	ErrServiceVehicleNotFound            = errors.New("service: vehicle not found")
)

// ServiceVehicle is the interface that wraps the basic methods for a vehicle service.
// - conections with external apis
// - business logic
type ServiceVehicle interface {
	// FindAll returns all vehicles
	FindAll() (v []Vehicle, err error)
	Insert(v Vehicle) (nv Vehicle, err error)
	FindAllByColorAndYear(c string, y int) (v []Vehicle, err error)
	FindAllByBrandAndBetweenYears(b string, sy int, ey int) (v []Vehicle, err error)
	CalculateAverageSpeedByBrand(b string) (avg float64, err error)
	InsertMany(v []Vehicle) (nvs []Vehicle, err error)
	UpdateMaxSpeedById(id int, ms int) (uv Vehicle, err error)
	FindAllByFuelType(ft string) (v []Vehicle, err error)
	Delete(id int) (err error)
	FindAllByTransmission(t string) (v []Vehicle, err error)
	UpdateFuelTypeById(id int, ft string) (uv Vehicle, err error)
	CalculateAverageCapacityByBrand(b string) (avg float64, err error)
	FindAllByDimensions(minH, maxH, minW, maxW float64) (v []Vehicle, err error)
	FindAllByWeight(minW, maxW float64) (v []Vehicle, err error)
}
