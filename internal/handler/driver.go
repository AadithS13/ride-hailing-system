package handler

import (
	"net/http"
	"ride-hailing/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9" 

	"ride-hailing/internal/repository" 
	"ride-hailing/internal/service"
)

func UpdateLocation(c *gin.Context) {
	driverID := c.Param("id")

	var input struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Redis GEO
	config.RedisClient.GeoAdd(config.Ctx, "drivers", &redis.GeoLocation{
		Name:      driverID,
		Longitude: input.Lng,
		Latitude:  input.Lat,
	})

	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
}

func AcceptRide(c *gin.Context) {
	driverID := c.Param("id")

	var input struct {
		RideID string `json:"ride_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := repository.NewRideRepository(config.DB)
	svc := service.NewRideService(repo)

	err := svc.AcceptRide(input.RideID, driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ride started"})
}