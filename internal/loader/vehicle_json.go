package loader

import (
	"app/internal"
	"encoding/json"
	"os"
)

// LoadDataJSON is an struct that represents the data of file.
type LoadDataJSON struct {
	Data []struct {
		ID           int     `json:"id"`
		Brand        string  `json:"brand"`
		Model        string  `json:"model"`
		Registration string  `json:"registration"`
		Year         int     `json:"year"`
		Color        string  `json:"color"`
		MaxSpeed     int     `json:"max_speed"`
		FuelType     string  `json:"fuel_type"`
		Transmission string  `json:"transmission"`
		Passengers   int     `json:"passengers"`
		Height       float64 `json:"height"`
		Width        float64 `json:"width"`
		Weight       float64 `json:"weight"`
	} `json:"data"`
	LastId int `json:"last_id"`
}

// NewVehicleJSON returns a new instance of a vehicle loader.
func NewVehicleJSON(path string) *VehicleJSON {
	return &VehicleJSON{Path: path}
}

// VehicleJSON is an struct that implements the LoaderVehicle interface.
type VehicleJSON struct {
	Path string
}

// Load returns all vehicles.
func (l *VehicleJSON) Load() (d internal.LoadData, err error) {
	// open file
	f, err := os.Open(l.Path)
	if err != nil {
		return
	}
	defer f.Close()

	// read file
	var loadDataJSON LoadDataJSON
	err = json.NewDecoder(f).Decode(&loadDataJSON)
	if err != nil {
		return
	}

	// serialize load data
	// - data
	d.Data = make([]internal.Vehicle, len(loadDataJSON.Data))
	for i, vehicle := range loadDataJSON.Data {
		d.Data[i] = internal.Vehicle{
			ID:           vehicle.ID,
			Attributes:   internal.VehicleAttributes{
				Brand:        vehicle.Brand,
				Model:        vehicle.Model,
				Registration: vehicle.Registration,
				Year:         vehicle.Year,
				Color:        vehicle.Color,
				MaxSpeed:     vehicle.MaxSpeed,
				FuelType:     vehicle.FuelType,
				Transmission: vehicle.Transmission,
				Passengers:   vehicle.Passengers,
				Height:       vehicle.Height,
				Width:        vehicle.Width,
				Weight:       vehicle.Weight,
			},
		}
	}
	// - last id
	d.LastId = loadDataJSON.LastId

	return
}
