package service

import (
	"errors"
	"ride-hailing/internal/config"
	"ride-hailing/internal/model"
	"ride-hailing/internal/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RideService struct {
	repo *repository.RideRepository
}

func NewRideService(repo *repository.RideRepository) *RideService {
	return &RideService{repo: repo}
}

// ---------------- CREATE RIDE ----------------

func (s *RideService) CreateRide(riderID string, lat, lng float64) (*model.Ride, error) {

	if config.RedisClient == nil {
		return nil, errors.New("redis not initialized")
	}

	// Find nearby drivers
	drivers, err := config.RedisClient.GeoSearch(
		config.Ctx,
		"drivers",
		&redis.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     5, // 5 km radius
			RadiusUnit: "km",
			Count:      20,
			Sort:       "ASC",
		},
	).Result()

	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		return nil, errors.New("no drivers available")
	}

	var selectedDriver string

	// pick first AVAILABLE driver
	for _, d := range drivers {

		status, err := config.RedisClient.Get(config.Ctx, "driver_status:"+d).Result()
		if err != nil && err != redis.Nil {
			continue
		}

		// default AVAILABLE if not set
		if status == "" || status == "AVAILABLE" {
			selectedDriver = d
			break
		}
	}

	if selectedDriver == "" {
		return nil, errors.New("no available drivers nearby")
	}

	// simple flat fare
	fare := 50.0

	ride := &model.Ride{
		ID:        uuid.New(),
		RiderID:   riderID,
		PickupLat: lat,
		PickupLng: lng,
		Status:    "MATCHED",
		DriverID:  &selectedDriver,
		Fare:      fare,
	}

	err = s.repo.Create(ride)
	if err != nil {
		return nil, err
	}

	// attach driver status
	status, _ := config.RedisClient.Get(config.Ctx, "driver_status:"+selectedDriver).Result()
	ride.DriverStatus = status

	return ride, nil
}

// ---------------- ACCEPT RIDE ----------------

func (s *RideService) AcceptRide(rideID, driverID string) error {

	ride, err := s.repo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.Status != "MATCHED" {
		return errors.New("ride already accepted or invalid state")
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return errors.New("not assigned to this driver")
	}

	if config.RedisClient == nil {
		return errors.New("redis not initialized")
	}

	// check availability
	status, _ := config.RedisClient.Get(config.Ctx, "driver_status:"+driverID).Result()
	if status != "" && status != "AVAILABLE" {
		return errors.New("driver not available")
	}

	// update ride
	ride.Status = "ONGOING"

	err = s.repo.Update(ride)
	if err != nil {
		return err
	}

	// update driver status
	config.RedisClient.Set(config.Ctx, "driver_status:"+driverID, "ONGOING", 0)

	err = config.RedisClient.Set(
		config.Ctx,
		"driver_active_ride:"+driverID,
		rideID,
		0,
	).Err()

	if err != nil {
		return err
	}

	return nil
}

// ---------------- END RIDE ----------------

func (s *RideService) EndRide(rideID string) (*model.Ride, error) {

	ride, err := s.repo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if ride.Status != "ONGOING" {
		return nil, errors.New("ride not ongoing")
	}

	ride.Status = "COMPLETED"

	err = s.repo.Update(ride)
	if err != nil {
		return nil, err
	}

	// free driver
	if ride.DriverID != nil {
		driverID := *ride.DriverID

		config.RedisClient.Set(config.Ctx, "driver_status:"+driverID, "AVAILABLE", 0)

		// remove mapping
		config.RedisClient.Del(
			config.Ctx,
			"driver_active_ride:"+driverID,
		)
	}

	return ride, nil
}