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

func (s *RideService) CreateRide(riderID string, lat, lng float64) (*model.Ride, error) {

	if config.RedisClient == nil {
		return nil, errors.New("redis not initialized")
	}

	// Find nearest drivers using Redis GEO
	drivers, err := config.RedisClient.GeoSearch(
		config.Ctx,
		"drivers",
		&redis.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     5,
			RadiusUnit: "km",
			Count:      1,
			Sort:       "ASC",
		},
	).Result()

	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		return nil, errors.New("no drivers available")
	}

	driverID := drivers[0]

	// Check availability
	availabilityKey := "driver_available:" + driverID

	available, err := config.RedisClient.Get(config.Ctx, availabilityKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if available == "false" {
		return nil, errors.New("driver not available")
	}

	// Simple fare
	fare := 50.0

	ride := &model.Ride{
		ID:        uuid.New(),
		RiderID:   riderID,
		PickupLat: lat,
		PickupLng: lng,
		Status:    "MATCHED",
		DriverID:  &driverID,
		Fare:      fare,
	}

	err = s.repo.Create(ride)
	if err != nil {
		return nil, err
	}

	return ride, nil
}

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

	ride.Status = "ONGOING"

	return s.repo.Update(ride)
}

func (s *RideService) EndRide(rideID string) (*model.Ride, error) {
	ride, err := s.repo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	ride.Status = "COMPLETED"

	return ride, s.repo.Update(ride)
}