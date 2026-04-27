package handler

import (
	"net/http"
	"ride-hailing/internal/config"
	"ride-hailing/internal/dto"
	"ride-hailing/internal/repository"
	"ride-hailing/internal/service"

	"github.com/gin-gonic/gin"
)

func CreateRide(c *gin.Context) {
	var input dto.CreateRideRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	ride, err := svc.CreateRide(input.RiderID, input.PickupLat, input.PickupLng)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ride)
}

func EndRideHandler(c *gin.Context) {
	rideID := c.Param("id")

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	ride, err := svc.EndRide(rideID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, ride)
}

func GetRide(c *gin.Context) {
	id := c.Param("id")

	repo := repository.NewRideRepository(config.DB)

	ride, err := repo.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "ride not found"})
		return
	}

	// Attach driver status from Redis
	if ride.DriverID != nil && config.RedisClient != nil {
		status, err := config.RedisClient.Get(config.Ctx, "driver_status:"+*ride.DriverID).Result()
		if err == nil {
			ride.DriverStatus = status
		}
	}

	c.JSON(200, ride)
}

func GetAllRides(c *gin.Context) {
	repo := repository.NewRideRepository(config.DB)

	rides, err := repo.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch rides"})
		return
	}

	// attach driver status
	for _, ride := range rides {
		if ride.DriverID != nil {
			status, _ := config.RedisClient.Get(config.Ctx, "driver_status:"+*ride.DriverID).Result()
			ride.DriverStatus = status
		}
	}

	c.JSON(200, gin.H{"rides": rides})
}

func AcceptRideHandler(c *gin.Context) {
	driverID := c.Param("id")

	var input struct {
		RideID string `json:"ride_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	err := svc.AcceptRide(input.RideID, driverID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated ride
	ride, err := repo.GetByID(input.RideID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch ride"})
		return
	}

	// Attach driver status
	if ride.DriverID != nil {
		status, _ := config.RedisClient.Get(config.Ctx, "driver_status:"+*ride.DriverID).Result()
		ride.DriverStatus = status
	}

	c.JSON(200, ride)
}

func EndRideByDriver(c *gin.Context) {
	driverID := c.Param("id")

	// get active ride
	rideID, err := config.RedisClient.Get(
		config.Ctx,
		"driver_active_ride:"+driverID,
	).Result()

	if err != nil {
		c.JSON(400, gin.H{"error": "no active ride for driver"})
		return
	}

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	ride, err := svc.EndRide(rideID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, ride)
}