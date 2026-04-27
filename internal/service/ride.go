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

func (s *RideService) CreateRide(riderID string, lat, lng float64) (*model.Ride, error) {

	// 🔥 Find nearest drivers using Redis GEO
	drivers, err := config.RedisClient.GeoSearch(
		config.Ctx,
		"drivers",
		&redis.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     5, // km
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

	// Create ride
	ride := &model.Ride{
		ID:        uuid.New(),
		RiderID:   riderID,
		PickupLat: lat,
		PickupLng: lng,
		Status:    "MATCHED",
		DriverID:  &driverID,
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