package internal

// LoadData is an struct that represents the data of file.
type LoadData struct {
	// Data is the slice of vehicles.
	Data []Vehicle
	// lastId is the last id of the slice of vehicles.
	LastId int
}

// Loader is the interface that wraps the basic methods for a vehicle loader.
type Loader interface {
	// Load returns all vehicles
	Load() (d LoadData, err error)
}