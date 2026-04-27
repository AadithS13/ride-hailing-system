package model

import (
	"time"

	"github.com/google/uuid"
)

type Ride struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RiderID   string
	PickupLat float64
	PickupLng float64
	Status    string
	DriverID  *string
	CreatedAt time.Time
}