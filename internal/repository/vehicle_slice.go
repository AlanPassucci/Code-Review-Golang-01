package repository

import "app/internal"

// NewVehicleSlice returns a new instance of a vehicle repository in an slice.
func NewVehicleSlice(db []internal.Vehicle, lastId int) *VehicleSlice {
	return &VehicleSlice{
		db:     db,
		lastId: lastId,
	}
}

// VehicleSlice is an struct that represents a vehicle repository in an slice.
type VehicleSlice struct {
	// db is the database of vehicles.
	db []internal.Vehicle
	// lastId is the last id of the database.
	lastId int
}

// FindAll returns all vehicles
func (r *VehicleSlice) FindAll() (v []internal.Vehicle, err error) {
	// check if the database is empty
	if len(r.db) == 0 {
		err = internal.ErrRepositoryVehiclesNotFound
		return
	}

	// make a copy of the database
	v = make([]internal.Vehicle, len(r.db))
	copy(v, r.db)
	return
}

func (r *VehicleSlice) Insert(v internal.Vehicle) (nv internal.Vehicle, err error) {
	r.lastId++
	if v.ID == r.lastId {
		r.lastId--
		err = internal.ErrRepositoryVehicleIdAlreadyExists
		return
	}
	v.ID = r.lastId
	r.db = append(r.db, v)
	nv = v
	return nv, nil
}

func (r *VehicleSlice) InsertMany(v []internal.Vehicle) (nvs []internal.Vehicle, err error) {
	nvs = make([]internal.Vehicle, 0)
	for _, vehicle := range v {
		r.lastId++
		if vehicle.ID == r.lastId {
			r.lastId--
			err = internal.ErrRepositoryVehicleIdAlreadyExists
			return
		}
		vehicle.ID = r.lastId
		r.db = append(r.db, vehicle)
		nvs = append(nvs, vehicle)
	}
	return nvs, nil
}

func (r *VehicleSlice) UpdateMaxSpeedById(id int, ms int) (uv internal.Vehicle, err error) {
	for i := range r.db {
		if r.db[i].ID == id {
			r.db[i].Attributes.MaxSpeed = ms
			uv = r.db[i]
			return
		}
	}
	return internal.Vehicle{}, internal.ErrRepositoryVehicleNotFound
}

func (r *VehicleSlice) Delete(id int) (err error) {
	index := -1
	for i := range r.db {
		if r.db[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return internal.ErrRepositoryVehicleNotFound
	}

	r.db = append(r.db[:index], r.db[index+1:]...)

	return nil
}

func (r *VehicleSlice) UpdateFuelTypeById(id int, ft string) (uv internal.Vehicle, err error) {
	for i := range r.db {
		if r.db[i].ID == id {
			r.db[i].Attributes.FuelType = ft
			uv = r.db[i]
			return
		}
	}
	return internal.Vehicle{}, internal.ErrRepositoryVehicleNotFound
}
