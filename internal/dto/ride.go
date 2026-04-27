package dto

type CreateRideRequest struct {
	RiderID   string  `json:"rider_id" binding:"required"`
	PickupLat float64 `json:"pickup_lat" binding:"required"`
	PickupLng float64 `json:"pickup_lng" binding:"required"`
}