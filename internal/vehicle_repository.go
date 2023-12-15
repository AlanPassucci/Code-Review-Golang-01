package internal

import (
	"errors"
)

var (
	// ErrRepositoryVehicleNotFound is returned when a vehicle is not found.
	ErrRepositoryVehiclesNotFound       = errors.New("repository: vehicles not found")
	ErrRepositoryVehicleIdAlreadyExists = errors.New("repository: vehicle id already exists")
	ErrRepositoryVehicleNotFound        = errors.New("repository: vehicle not found")
)

// RepositoryVehicle is the interface that wraps the basic methods for a vehicle repository.
type RepositoryVehicle interface {
	// FindAll returns all vehicles
	FindAll() (v []Vehicle, err error)
	Insert(v Vehicle) (nv Vehicle, err error)
	InsertMany(v []Vehicle) (nvs []Vehicle, err error)
	UpdateMaxSpeedById(id int, ms int) (uv Vehicle, err error)
	Delete(id int) (err error)
	UpdateFuelTypeById(id int, ft string) (uv Vehicle, err error)
}
