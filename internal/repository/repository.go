package repository

import (
	"ride-hailing/internal/model"

	"gorm.io/gorm"
)

type RideRepository struct {
	DB *gorm.DB
}

func NewRideRepository(db *gorm.DB) *RideRepository {
	return &RideRepository{DB: db}
}

func (r *RideRepository) Create(ride *model.Ride) error {
	return r.DB.Create(ride).Error
}

func (r *RideRepository) GetByID(id string) (*model.Ride, error) {
	var ride model.Ride
	err := r.DB.Where("id = ?", id).First(&ride).Error
	return &ride, err
}

func (r *RideRepository) Update(ride *model.Ride) error {
	return r.DB.Save(ride).Error
}