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

func EndRide(c *gin.Context) {
	rideID := c.Param("id")

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	ride, err := svc.EndRide(rideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ride)
}