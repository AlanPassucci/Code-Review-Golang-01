package handler

import (
	"app/internal"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// VehicleJSON is an struct that represents a vehicle in json format.
type VehicleJSON struct {
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
}

type BodyRequestUpdateMaxSpeed struct {
	MaxSpeed int `json:"max_speed"`
}

type BodyRequestUpdateFuelType struct {
	FuelType string `json:"fuel_type"`
}

// NewVehicleDefault returns a new instance of a vehicle handler.
func NewVehicleDefault(sv internal.ServiceVehicle) *VehicleDefault {
	return &VehicleDefault{sv: sv}
}

// VehicleDefault is an struct that contains handlers for vehicle.
type VehicleDefault struct {
	sv internal.ServiceVehicle
}

// GetAll returns all vehicles.
func (hd *VehicleDefault) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// ...

		// process
		// - get all vehicles from the service
		vehicles, err := hd.sv.FindAll()
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"message": "vehicles not found"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			}
			return
		}

		// response
		// - serialize vehicles
		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "success to find vehicles", "data": data})
	}
}

func (hd *VehicleDefault) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var m map[string]any
		if err := ctx.ShouldBindJSON(&m); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		requiredFields := []string{"brand", "model", "registration", "year", "color", "max_speed", "fuel_type", "transmission", "passengers", "height", "width", "weight"}
		for _, field := range requiredFields {
			if _, exists := m[field]; !exists {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("missing required field: %s", field)})
				return
			}
		}

		vehicle := internal.Vehicle{
			Attributes: internal.VehicleAttributes{
				Brand:        m["brand"].(string),
				Model:        m["model"].(string),
				Registration: m["registration"].(string),
				Year:         int(m["year"].(float64)),
				Color:        m["color"].(string),
				MaxSpeed:     int(m["max_speed"].(float64)),
				FuelType:     m["fuel_type"].(string),
				Transmission: m["transmission"].(string),
				Passengers:   int(m["passengers"].(float64)),
				Height:       m["height"].(float64),
				Width:        m["width"].(float64),
				Weight:       m["weight"].(float64),
			},
		}

		newVehicle, err := hd.sv.Insert(vehicle)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleBrand) || errors.Is(err, internal.ErrServiceInvalidVehicleModel) || errors.Is(err, internal.ErrServiceInvalidVehicleRegistration) || errors.Is(err, internal.ErrServiceInvalidVehicleYear) || errors.Is(err, internal.ErrServiceInvalidVehicleColor) || errors.Is(err, internal.ErrServiceInvalidVehicleMaxSpeed) || errors.Is(err, internal.ErrServiceInvalidVehicleFuelType) || errors.Is(err, internal.ErrServiceInvalidVehicleTransmission) || errors.Is(err, internal.ErrServiceInvalidVehiclePassengers) || errors.Is(err, internal.ErrServiceInvalidVehicleHeight) || errors.Is(err, internal.ErrServiceInvalidVehicleWidth) || errors.Is(err, internal.ErrServiceInvalidVehicleWeight):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle"})
			case errors.Is(err, internal.ErrServiceVehicleIdAlreadyExists):
				ctx.JSON(http.StatusConflict, gin.H{"error": "vehicle already exists"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"message": "vehicle created",
			"data": VehicleJSON{
				ID:           newVehicle.ID,
				Brand:        newVehicle.Attributes.Brand,
				Model:        newVehicle.Attributes.Model,
				Registration: newVehicle.Attributes.Registration,
				Year:         newVehicle.Attributes.Year,
				Color:        newVehicle.Attributes.Color,
				MaxSpeed:     newVehicle.Attributes.MaxSpeed,
				FuelType:     newVehicle.Attributes.FuelType,
				Transmission: newVehicle.Attributes.Transmission,
				Passengers:   newVehicle.Attributes.Passengers,
				Height:       newVehicle.Attributes.Height,
				Width:        newVehicle.Attributes.Width,
				Weight:       newVehicle.Attributes.Weight,
			},
		})
	}
}

func (hd *VehicleDefault) GetAllByColorAndYear() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		color := ctx.Param("color")
		year, err := strconv.Atoi(ctx.Param("year"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}

		vehicles, err := hd.sv.FindAllByColorAndYear(color, year)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleColor) || errors.Is(err, internal.ErrServiceInvalidVehicleYear):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that color and year"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "vehicles with that color and year were found",
			"data":    data,
		})
	}
}

func (hd *VehicleDefault) GetAllByBrandAndBetweenYears() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		brand := ctx.Param("brand")
		startYear, err := strconv.Atoi(ctx.Param("start_year"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}
		endYear, err := strconv.Atoi(ctx.Param("end_year"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}

		vehicles, err := hd.sv.FindAllByBrandAndBetweenYears(brand, startYear, endYear)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleBrand) || errors.Is(err, internal.ErrServiceInvalidVehicleYear):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that brand and range of years"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "vehicles with that brand and range of years were found",
			"data":    data,
		})
	}
}

func (hd *VehicleDefault) CalculateAverageSpeedByBrand() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		brand := ctx.Param("brand")

		avg, err := hd.sv.CalculateAverageSpeedByBrand(brand)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleBrand):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that brand"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("the average max speed of %s vehicles is %.2f", brand, avg),
		})
	}
}

func (hd *VehicleDefault) CreateMany() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ms []map[string]any
		if err := ctx.ShouldBindJSON(&ms); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		vhToInsert := make([]internal.Vehicle, 0)
		requiredFields := []string{"brand", "model", "registration", "year", "color", "max_speed", "fuel_type", "transmission", "passengers", "height", "width", "weight"}
		for i := range ms {
			for _, field := range requiredFields {
				if _, exists := ms[i][field]; !exists {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("missing required field: %s", field)})
					return
				}
			}

			vehicle := internal.Vehicle{
				Attributes: internal.VehicleAttributes{
					Brand:        ms[i]["brand"].(string),
					Model:        ms[i]["model"].(string),
					Registration: ms[i]["registration"].(string),
					Year:         int(ms[i]["year"].(float64)),
					Color:        ms[i]["color"].(string),
					MaxSpeed:     int(ms[i]["max_speed"].(float64)),
					FuelType:     ms[i]["fuel_type"].(string),
					Transmission: ms[i]["transmission"].(string),
					Passengers:   int(ms[i]["passengers"].(float64)),
					Height:       ms[i]["height"].(float64),
					Width:        ms[i]["width"].(float64),
					Weight:       ms[i]["weight"].(float64),
				},
			}
			vhToInsert = append(vhToInsert, vehicle)
		}

		newVehicles, err := hd.sv.InsertMany(vhToInsert)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleBrand) || errors.Is(err, internal.ErrServiceInvalidVehicleModel) || errors.Is(err, internal.ErrServiceInvalidVehicleRegistration) || errors.Is(err, internal.ErrServiceInvalidVehicleYear) || errors.Is(err, internal.ErrServiceInvalidVehicleColor) || errors.Is(err, internal.ErrServiceInvalidVehicleMaxSpeed) || errors.Is(err, internal.ErrServiceInvalidVehicleFuelType) || errors.Is(err, internal.ErrServiceInvalidVehicleTransmission) || errors.Is(err, internal.ErrServiceInvalidVehiclePassengers) || errors.Is(err, internal.ErrServiceInvalidVehicleHeight) || errors.Is(err, internal.ErrServiceInvalidVehicleWidth) || errors.Is(err, internal.ErrServiceInvalidVehicleWeight):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "some vehicles are invalid"})
			case errors.Is(err, internal.ErrServiceVehicleIdAlreadyExists):
				ctx.JSON(http.StatusConflict, gin.H{"error": "some vehicles already exists"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		newVehiclesJSON := make([]VehicleJSON, 0)
		for _, v := range newVehicles {
			newVehiclesJSON = append(newVehiclesJSON, VehicleJSON{
				ID:           v.ID,
				Brand:        v.Attributes.Brand,
				Model:        v.Attributes.Model,
				Registration: v.Attributes.Registration,
				Year:         v.Attributes.Year,
				Color:        v.Attributes.Color,
				MaxSpeed:     v.Attributes.MaxSpeed,
				FuelType:     v.Attributes.FuelType,
				Transmission: v.Attributes.Transmission,
				Passengers:   v.Attributes.Passengers,
				Height:       v.Attributes.Height,
				Width:        v.Attributes.Width,
				Weight:       v.Attributes.Weight,
			})
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "vehicles created",
			"data":    newVehiclesJSON,
		})
	}
}

func (hd *VehicleDefault) UpdateMaxSpeedById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier"})
			return
		}

		var body BodyRequestUpdateMaxSpeed
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		uv, err := hd.sv.UpdateMaxSpeedById(id, body.MaxSpeed)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleMaxSpeed):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle max speed"})
			case errors.Is(err, internal.ErrServiceVehicleNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		uvJSON := VehicleJSON{
			ID:           uv.ID,
			Brand:        uv.Attributes.Brand,
			Model:        uv.Attributes.Model,
			Registration: uv.Attributes.Registration,
			Year:         uv.Attributes.Year,
			Color:        uv.Attributes.Color,
			MaxSpeed:     uv.Attributes.MaxSpeed,
			FuelType:     uv.Attributes.FuelType,
			Transmission: uv.Attributes.Transmission,
			Passengers:   uv.Attributes.Passengers,
			Height:       uv.Attributes.Height,
			Width:        uv.Attributes.Width,
			Weight:       uv.Attributes.Weight,
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "updated max speed of vehicle",
			"data":    uvJSON,
		})
	}
}

func (hd *VehicleDefault) GetAllByFuelType() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ft := ctx.Param("type")

		vehicles, err := hd.sv.FindAllByFuelType(ft)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleFuelType):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle fuel type"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that fuel type"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "vehicles with that fuel type were found",
			"data":    data,
		})
	}
}

func (hd *VehicleDefault) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier"})
			return
		}

		if err := hd.sv.Delete(id); err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceVehicleNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func (hd *VehicleDefault) GetAllByTransmission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := ctx.Param("type")

		vehicles, err := hd.sv.FindAllByTransmission(t)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleTransmission):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle transmission"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that transmission"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "vehicles with that transmission were found",
			"data":    data,
		})
	}
}

func (hd *VehicleDefault) UpdateFuelTypeById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid identifier"})
			return
		}

		var body BodyRequestUpdateFuelType
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		uv, err := hd.sv.UpdateFuelTypeById(id, body.FuelType)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleFuelType):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle fuel type"})
			case errors.Is(err, internal.ErrServiceVehicleNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		uvJSON := VehicleJSON{
			ID:           uv.ID,
			Brand:        uv.Attributes.Brand,
			Model:        uv.Attributes.Model,
			Registration: uv.Attributes.Registration,
			Year:         uv.Attributes.Year,
			Color:        uv.Attributes.Color,
			MaxSpeed:     uv.Attributes.MaxSpeed,
			FuelType:     uv.Attributes.FuelType,
			Transmission: uv.Attributes.Transmission,
			Passengers:   uv.Attributes.Passengers,
			Height:       uv.Attributes.Height,
			Width:        uv.Attributes.Width,
			Weight:       uv.Attributes.Weight,
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "updated fuel type of vehicle",
			"data":    uvJSON,
		})
	}
}

func (hd *VehicleDefault) CalculateAverageCapacityByBrand() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		brand := ctx.Param("brand")

		avg, err := hd.sv.CalculateAverageCapacityByBrand(brand)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleBrand):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that brand"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("the average capacity of %s vehicles is %.2f", brand, avg),
		})
	}
}

func (hd *VehicleDefault) GetAllByDimensions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		height := ctx.Query("height")
		splitH := strings.Split(height, "-")
		if len(splitH) != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid height"})
			return
		}
		width := ctx.Query("width")
		splitW := strings.Split(width, "-")
		if len(splitW) != 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid width"})
			return
		}
		minH, err := strconv.ParseFloat(splitH[0], 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid min height"})
			return
		}
		maxH, err := strconv.ParseFloat(splitH[1], 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid max height"})
			return
		}
		minW, err := strconv.ParseFloat(splitW[0], 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid min width"})
			return
		}
		maxW, err := strconv.ParseFloat(splitW[1], 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid max width"})
			return
		}

		vehicles, err := hd.sv.FindAllByDimensions(minH, maxH, minW, maxW)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleWidth) || errors.Is(err, internal.ErrServiceInvalidVehicleHeight):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid dimensions"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that dimensions"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "vehicles with that dimensions were found",
			"data":    data,
		})
	}
}

func (hd *VehicleDefault) GetAllByWeights() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		minW, err := strconv.ParseFloat(ctx.Query("min"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid min weight"})
			return
		}
		maxW, err := strconv.ParseFloat(ctx.Query("max"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid max weight"})
			return
		}

		vehicles, err := hd.sv.FindAllByWeight(minW, maxW)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrServiceInvalidVehicleWeight):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid weight"})
			case errors.Is(err, internal.ErrServiceVehiclesNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": "there are not any vehicles with that weight"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
			}
			return
		}

		data := make([]VehicleJSON, len(vehicles))
		for i, vehicle := range vehicles {
			data[i] = VehicleJSON{
				ID:           vehicle.ID,
				Brand:        vehicle.Attributes.Brand,
				Model:        vehicle.Attributes.Model,
				Registration: vehicle.Attributes.Registration,
				Year:         vehicle.Attributes.Year,
				Color:        vehicle.Attributes.Color,
				MaxSpeed:     vehicle.Attributes.MaxSpeed,
				FuelType:     vehicle.Attributes.FuelType,
				Transmission: vehicle.Attributes.Transmission,
				Passengers:   vehicle.Attributes.Passengers,
				Height:       vehicle.Attributes.Height,
				Width:        vehicle.Attributes.Width,
				Weight:       vehicle.Attributes.Weight,
			}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "vehicles with that weight were found",
			"data":    data,
		})
	}
}
