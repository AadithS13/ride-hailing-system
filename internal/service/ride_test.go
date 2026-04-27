package service

import (
	"errors"
	"testing"

	"ride-hailing/internal/config"
	"ride-hailing/internal/model"

	"github.com/google/uuid"
)

type MockRideRepo struct {
	Ride *model.Ride
}

func (m *MockRideRepo) Create(ride *model.Ride) error {
	m.Ride = ride
	return nil
}

func (m *MockRideRepo) GetByID(id string) (*model.Ride, error) {
	if m.Ride == nil {
		return nil, errors.New("not found")
	}
	return m.Ride, nil
}

func (m *MockRideRepo) Update(ride *model.Ride) error {
	m.Ride = ride
	return nil
}

// ---------------- TESTS ----------------

func TestCreateRide_NoRedis(t *testing.T) {
	repo := &MockRideRepo{}
	service := NewRideService(repo)

	// simulate missing Redis
	config.RedisClient = nil

	_, err := service.CreateRide("r1", 12.9, 77.5)

	if err == nil {
		t.Errorf("expected error when redis not initialized")
	}
}

func TestAcceptRide_Success(t *testing.T) {
	driverID := "1"

	repo := &MockRideRepo{
		Ride: &model.Ride{
			ID:       uuid.New(),
			Status:   "MATCHED",
			DriverID: &driverID,
		},
	}

	service := NewRideService(repo)

	err := service.AcceptRide(repo.Ride.ID.String(), driverID)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if repo.Ride.Status != "ONGOING" {
		t.Errorf("expected ONGOING, got %s", repo.Ride.Status)
	}
}

func TestAcceptRide_AlreadyAccepted(t *testing.T) {
	driverID := "1"

	repo := &MockRideRepo{
		Ride: &model.Ride{
			ID:       uuid.New(),
			Status:   "ONGOING",
			DriverID: &driverID,
		},
	}

	service := NewRideService(repo)

	err := service.AcceptRide(repo.Ride.ID.String(), driverID)

	if err == nil {
		t.Errorf("expected error for already accepted ride")
	}
}

func TestEndRide(t *testing.T) {
	driverID := "1"

	repo := &MockRideRepo{
		Ride: &model.Ride{
			ID:       uuid.New(),
			Status:   "ONGOING",
			DriverID: &driverID,
		},
	}

	service := NewRideService(repo)

	ride, err := service.EndRide(repo.Ride.ID.String())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ride.Status != "COMPLETED" {
		t.Errorf("expected COMPLETED, got %s", ride.Status)
	}
}