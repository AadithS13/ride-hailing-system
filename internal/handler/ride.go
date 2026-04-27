package handler

import (
	"encoding/json"
	"net/http"
	"time"

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

// ---------------- END RIDE BY ID ----------------

func EndRideHandler(c *gin.Context) {
	rideID := c.Param("id")

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	ride, err := svc.EndRide(rideID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Invalidate cache
	if config.RedisClient != nil {
		config.RedisClient.Del(config.Ctx, "ride:"+rideID)
	}

	c.JSON(200, ride)
}

// ---------------- GET RIDE (WITH CACHE) ----------------

func GetRide(c *gin.Context) {
	id := c.Param("id")

	cacheKey := "ride:" + id

	// CACHE HIT
	if config.RedisClient != nil {
		cached, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
		if err == nil {
			c.Data(http.StatusOK, "application/json", []byte(cached))
			return
		}
	}

	repo := repository.NewRideRepository(config.DB)

	ride, err := repo.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "ride not found"})
		return
	}

	// Attach driver status
	if ride.DriverID != nil && config.RedisClient != nil {
		status, err := config.RedisClient.Get(config.Ctx, "driver_status:"+*ride.DriverID).Result()
		if err == nil {
			ride.DriverStatus = status
		}
	}

	// CACHE SET
	if config.RedisClient != nil {
		bytes, _ := json.Marshal(ride)
		config.RedisClient.Set(config.Ctx, cacheKey, bytes, 10*time.Second)
	}

	c.JSON(200, ride)
}

// ---------------- GET ALL RIDES ----------------

func GetAllRides(c *gin.Context) {
	repo := repository.NewRideRepository(config.DB)

	rides, err := repo.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch rides"})
		return
	}

	// attach driver status
	for _, ride := range rides {
		if ride.DriverID != nil && config.RedisClient != nil {
			status, _ := config.RedisClient.Get(config.Ctx, "driver_status:"+*ride.DriverID).Result()
			ride.DriverStatus = status
		}
	}

	c.JSON(200, gin.H{"rides": rides})
}

// ---------------- ACCEPT RIDE ----------------

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

	// Invalidate cache
	if config.RedisClient != nil {
		config.RedisClient.Del(config.Ctx, "ride:"+input.RideID)
	}

	// Fetch updated ride
	ride, err := repo.GetByID(input.RideID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch ride"})
		return
	}

	// Attach driver status
	if ride.DriverID != nil && config.RedisClient != nil {
		status, _ := config.RedisClient.Get(config.Ctx, "driver_status:"+*ride.DriverID).Result()
		ride.DriverStatus = status
	}

	c.JSON(200, ride)
}

// ---------------- END RIDE BY DRIVER ----------------

func EndRideByDriver(c *gin.Context) {
	driverID := c.Param("id")

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

	// Invalidate cache
	if config.RedisClient != nil {
		config.RedisClient.Del(config.Ctx, "ride:"+rideID)
	}

	c.JSON(200, ride)
}