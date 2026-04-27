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

	c.JSON(200, ride)
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

	c.JSON(200, gin.H{"message": "ride started"})
}