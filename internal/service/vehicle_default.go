package service

import (
	"app/internal"
	"errors"
	"fmt"
)

// NewDefault returns a new instance of a vehicle service.
func NewDefault(rp internal.RepositoryVehicle) *Default {
	return &Default{rp: rp}
}

// Default is an struct that represents a vehicle service.
type Default struct {
	rp internal.RepositoryVehicle
}

// FindAll returns all vehicles.
func (sv *Default) FindAll() (v []internal.Vehicle, err error) {
	// get all vehicles from the repository
	v, err = sv.rp.FindAll()
	if err != nil {
		if errors.Is(err, internal.ErrRepositoryVehiclesNotFound) {
			err = fmt.Errorf("%w. %v", internal.ErrServiceVehiclesNotFound, err)
			return
		}
		return
	}

	return
}

func (sv *Default) Insert(v internal.Vehicle) (nv internal.Vehicle, err error) {
	if v.Attributes.Brand == "" {
		err = internal.ErrServiceInvalidVehicleBrand
		return
	}
	if v.Attributes.Model == "" {
		err = internal.ErrServiceInvalidVehicleModel
		return
	}
	if v.Attributes.Registration == "" {
		err = internal.ErrServiceInvalidVehicleRegistration
		return
	}
	if v.Attributes.Year < 1887 {
		err = internal.ErrServiceInvalidVehicleYear
		return
	}
	if v.Attributes.Color == "" {
		err = internal.ErrServiceInvalidVehicleColor
		return
	}
	if v.Attributes.MaxSpeed < 1 || v.Attributes.MaxSpeed > 999 {
		err = internal.ErrServiceInvalidVehicleMaxSpeed
		return
	}
	if v.Attributes.FuelType == "" || (v.Attributes.FuelType != "gasoline" && v.Attributes.FuelType != "gas" && v.Attributes.FuelType != "diesel" && v.Attributes.FuelType != "biodiesel") {
		err = internal.ErrServiceInvalidVehicleFuelType
		return
	}
	if v.Attributes.Transmission == "" || (v.Attributes.Transmission != "automatic" && v.Attributes.Transmission != "semi-automatic" && v.Attributes.Transmission != "manual") {
		err = internal.ErrServiceInvalidVehicleTransmission
		return
	}
	if v.Attributes.Passengers < 1 || v.Attributes.Passengers > 6 {
		err = internal.ErrServiceInvalidVehiclePassengers
		return
	}
	if v.Attributes.Height < 1 {
		err = internal.ErrServiceInvalidVehicleHeight
		return
	}
	if v.Attributes.Width < 1 {
		err = internal.ErrServiceInvalidVehicleWidth
		return
	}
	if v.Attributes.Weight < 1 {
		err = internal.ErrServiceInvalidVehicleWeight
		return
	}

	nv, err = sv.rp.Insert(v)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehicleIdAlreadyExists):
			return internal.Vehicle{}, internal.ErrServiceVehicleIdAlreadyExists
		default:
			return internal.Vehicle{}, err
		}
	}

	return nv, nil
}

func (sv *Default) FindAllByColorAndYear(c string, y int) (v []internal.Vehicle, err error) {
	if c == "" {
		err = internal.ErrServiceInvalidVehicleColor
		return
	}
	if y < 1887 {
		err = internal.ErrServiceInvalidVehicleYear
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return nil, internal.ErrServiceVehiclesNotFound
		default:
			return nil, err
		}
	}

	vehiclesWithColorAndYear := make([]internal.Vehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Color == c && vehicle.Attributes.Year == y {
			vehiclesWithColorAndYear = append(vehiclesWithColorAndYear, vehicle)
		}
	}

	if len(vehiclesWithColorAndYear) == 0 {
		return nil, internal.ErrServiceVehiclesNotFound
	}

	return vehiclesWithColorAndYear, nil
}

func (sv *Default) FindAllByBrandAndBetweenYears(b string, sy int, ey int) (v []internal.Vehicle, err error) {
	if b == "" {
		err = internal.ErrServiceInvalidVehicleBrand
		return
	}
	if sy < 1887 || ey < sy {
		err = internal.ErrServiceInvalidVehicleYear
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return nil, internal.ErrServiceVehiclesNotFound
		default:
			return nil, err
		}
	}

	vehiclesWithColorAndBetweenYears := make([]internal.Vehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Brand == b && vehicle.Attributes.Year > sy && vehicle.Attributes.Year < ey {
			vehiclesWithColorAndBetweenYears = append(vehiclesWithColorAndBetweenYears, vehicle)
		}
	}

	if len(vehiclesWithColorAndBetweenYears) == 0 {
		return nil, internal.ErrServiceVehiclesNotFound
	}

	return vehiclesWithColorAndBetweenYears, nil
}

func (sv *Default) CalculateAverageSpeedByBrand(b string) (avg float64, err error) {
	if b == "" {
		err = internal.ErrServiceInvalidVehicleBrand
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return 0.0, internal.ErrServiceVehiclesNotFound
		default:
			return 0.0, err
		}
	}

	vehiclesBrand := make([]internal.Vehicle, 0)
	maxSpeedSum := 0
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Brand == b {
			vehiclesBrand = append(vehiclesBrand, vehicle)
			maxSpeedSum += vehicle.Attributes.MaxSpeed
		}
	}

	if len(vehiclesBrand) == 0 {
		return 0.0, internal.ErrServiceVehiclesNotFound
	}

	avg = float64(maxSpeedSum) / float64(len(vehiclesBrand))

	return
}

func (sv *Default) InsertMany(v []internal.Vehicle) (nvs []internal.Vehicle, err error) {
	for _, vh := range v {
		if vh.Attributes.Brand == "" {
			err = internal.ErrServiceInvalidVehicleBrand
			return
		}
		if vh.Attributes.Model == "" {
			err = internal.ErrServiceInvalidVehicleModel
			return
		}
		if vh.Attributes.Registration == "" {
			err = internal.ErrServiceInvalidVehicleRegistration
			return
		}
		if vh.Attributes.Year < 1887 {
			err = internal.ErrServiceInvalidVehicleYear
			return
		}
		if vh.Attributes.Color == "" {
			err = internal.ErrServiceInvalidVehicleColor
			return
		}
		if vh.Attributes.MaxSpeed < 1 || vh.Attributes.MaxSpeed > 999 {
			err = internal.ErrServiceInvalidVehicleMaxSpeed
			return
		}
		if vh.Attributes.FuelType == "" || (vh.Attributes.FuelType != "gasoline" && vh.Attributes.FuelType != "gas" && vh.Attributes.FuelType != "diesel" && vh.Attributes.FuelType != "biodiesel") {
			err = internal.ErrServiceInvalidVehicleFuelType
			return
		}
		if vh.Attributes.Transmission == "" || (vh.Attributes.Transmission != "automatic" && vh.Attributes.Transmission != "semi-automatic" && vh.Attributes.Transmission != "manual") {
			err = internal.ErrServiceInvalidVehicleTransmission
			return
		}
		if vh.Attributes.Passengers < 1 || vh.Attributes.Passengers > 6 {
			err = internal.ErrServiceInvalidVehiclePassengers
			return
		}
		if vh.Attributes.Height < 1 {
			err = internal.ErrServiceInvalidVehicleHeight
			return
		}
		if vh.Attributes.Width < 1 {
			err = internal.ErrServiceInvalidVehicleWidth
			return
		}
		if vh.Attributes.Weight < 1 {
			err = internal.ErrServiceInvalidVehicleWeight
			return
		}
	}

	nvs, err = sv.rp.InsertMany(v)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehicleIdAlreadyExists):
			return nil, internal.ErrServiceVehicleIdAlreadyExists
		default:
			return nil, err
		}
	}

	return nvs, nil
}

func (sv *Default) UpdateMaxSpeedById(id int, ms int) (uv internal.Vehicle, err error) {
	if ms < 1 || ms > 999 {
		err = internal.ErrServiceInvalidVehicleMaxSpeed
		return
	}

	uv, err = sv.rp.UpdateMaxSpeedById(id, ms)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehicleNotFound):
			return internal.Vehicle{}, internal.ErrServiceVehicleNotFound
		default:
			return internal.Vehicle{}, err
		}
	}
	return uv, nil
}

func (sv *Default) FindAllByFuelType(ft string) (v []internal.Vehicle, err error) {
	if ft == "" || (ft != "gasoline" && ft != "gas" && ft != "diesel" && ft != "biodiesel") {
		err = internal.ErrServiceInvalidVehicleFuelType
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return nil, internal.ErrServiceVehiclesNotFound
		default:
			return nil, err
		}
	}

	vehiclesByFuelType := make([]internal.Vehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.Attributes.FuelType == ft {
			vehiclesByFuelType = append(vehiclesByFuelType, vehicle)
		}
	}

	if len(vehiclesByFuelType) == 0 {
		return nil, internal.ErrServiceVehiclesNotFound
	}

	return vehiclesByFuelType, nil
}

func (sv *Default) Delete(id int) (err error) {
	if err = sv.rp.Delete(id); err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehicleNotFound):
			return internal.ErrServiceVehicleNotFound
		default:
			return err
		}
	}
	return nil
}

func (sv *Default) FindAllByTransmission(t string) (v []internal.Vehicle, err error) {
	if t == "" || (t != "automatic" && t != "semi-automatic" && t != "manual") {
		err = internal.ErrServiceInvalidVehicleTransmission
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return nil, internal.ErrServiceVehiclesNotFound
		default:
			return nil, err
		}
	}

	vehiclesByTransmission := make([]internal.Vehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Transmission == t {
			vehiclesByTransmission = append(vehiclesByTransmission, vehicle)
		}
	}

	if len(vehiclesByTransmission) == 0 {
		return nil, internal.ErrServiceVehiclesNotFound
	}

	return vehiclesByTransmission, nil
}

func (sv *Default) UpdateFuelTypeById(id int, ft string) (uv internal.Vehicle, err error) {
	if ft == "" || (ft != "gasoline" && ft != "gas" && ft != "diesel" && ft != "biodiesel") {
		err = internal.ErrServiceInvalidVehicleFuelType
		return
	}

	uv, err = sv.rp.UpdateFuelTypeById(id, ft)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehicleNotFound):
			return internal.Vehicle{}, internal.ErrServiceVehicleNotFound
		default:
			return internal.Vehicle{}, err
		}
	}
	return uv, nil
}

func (sv *Default) CalculateAverageCapacityByBrand(b string) (avg float64, err error) {
	if b == "" {
		err = internal.ErrServiceInvalidVehicleBrand
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return 0.0, internal.ErrServiceVehiclesNotFound
		default:
			return 0.0, err
		}
	}

	vehiclesBrand := make([]internal.Vehicle, 0)
	passengersSum := 0
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Brand == b {
			vehiclesBrand = append(vehiclesBrand, vehicle)
			passengersSum += vehicle.Attributes.Passengers
		}
	}

	if len(vehiclesBrand) == 0 {
		return 0.0, internal.ErrServiceVehiclesNotFound
	}

	avg = float64(passengersSum) / float64(len(vehiclesBrand))

	return
}

func (sv *Default) FindAllByDimensions(minH, maxH, minW, maxW float64) (v []internal.Vehicle, err error) {
	if minH < 1 || maxH < minH {
		err = internal.ErrServiceInvalidVehicleHeight
		return
	}
	if minW < 1 || maxW < minW {
		err = internal.ErrServiceInvalidVehicleWidth
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return nil, internal.ErrServiceVehiclesNotFound
		default:
			return nil, err
		}
	}

	vehiclesWithDimensions := make([]internal.Vehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Height > minH && vehicle.Attributes.Height < maxH && vehicle.Attributes.Width > minW && vehicle.Attributes.Width < maxW {
			vehiclesWithDimensions = append(vehiclesWithDimensions, vehicle)
		}
	}

	if len(vehiclesWithDimensions) == 0 {
		return nil, internal.ErrServiceVehiclesNotFound
	}

	return vehiclesWithDimensions, nil
}

func (sv *Default) FindAllByWeight(minW, maxW float64) (v []internal.Vehicle, err error) {
	if minW < 1 || maxW < minW {
		err = internal.ErrServiceInvalidVehicleWeight
		return
	}

	vehicles, err := sv.FindAll()
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRepositoryVehiclesNotFound):
			return nil, internal.ErrServiceVehiclesNotFound
		default:
			return nil, err
		}
	}

	vehiclesWithWeight := make([]internal.Vehicle, 0)
	for _, vehicle := range vehicles {
		if vehicle.Attributes.Weight > minW && vehicle.Attributes.Weight < maxW {
			vehiclesWithWeight = append(vehiclesWithWeight, vehicle)
		}
	}

	if len(vehiclesWithWeight) == 0 {
		return nil, internal.ErrServiceVehiclesNotFound
	}

	return vehiclesWithWeight, nil
}
