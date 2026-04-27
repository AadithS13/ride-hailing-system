package service

import (
	"errors"

	"ride-hailing/internal/config"
	"ride-hailing/internal/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RideRepositoryInterface interface {
	Create(*model.Ride) error
	GetByID(string) (*model.Ride, error)
	Update(*model.Ride) error
}

type RideService struct {
	repo RideRepositoryInterface
}

func NewRideService(repo RideRepositoryInterface) *RideService {
	return &RideService{repo: repo}
}

// Create Ride
func (s *RideService) CreateRide(riderID string, lat, lng float64) (*model.Ride, error) {

	// Redis required for matching
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
			Radius:     5,
			RadiusUnit: "km",
			Count:      5,
			Sort:       "ASC",
		},
	).Result()

	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		return nil, errors.New("no drivers available")
	}

	// Pick first AVAILABLE driver
	var selectedDriver string

	for _, d := range drivers {
		status, err := config.RedisClient.Get(config.Ctx, "driver_status:"+d).Result()

		// default = AVAILABLE
		if err == redis.Nil {
			status = "AVAILABLE"
		}

		if status == "AVAILABLE" {
			selectedDriver = d
			break
		}
	}

	if selectedDriver == "" {
		return nil, errors.New("no available drivers")
	}

	// Simple fare
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

	return ride, nil
}

// Accept Ride
func (s *RideService) AcceptRide(rideID, driverID string) error {
	ride, err := s.repo.GetByID(rideID)
	if err != nil {
		return err
	}

	if ride.Status != "MATCHED" {
		return errors.New("invalid state")
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return errors.New("not assigned driver")
	}

	ride.Status = "ONGOING"

	// Mark driver busy
	if config.RedisClient != nil {
		config.RedisClient.Set(config.Ctx, "driver_status:"+driverID, "ONGOING", 0)
	}

	return s.repo.Update(ride)
}

// End Ride
func (s *RideService) EndRide(rideID string) (*model.Ride, error) {
	ride, err := s.repo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	ride.Status = "COMPLETED"

	// Mark driver available again
	if ride.DriverID != nil && config.RedisClient != nil {
		config.RedisClient.Set(config.Ctx, "driver_status:"+*ride.DriverID, "AVAILABLE", 0)
	}

	return ride, s.repo.Update(ride)
}