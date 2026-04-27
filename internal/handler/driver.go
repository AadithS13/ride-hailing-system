package handler

import (

	"ride-hailing/internal/config"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func CreateDriver(c *gin.Context) {
	var input struct {
		DriverID string  `json:"driver_id"`
		Lat      float64 `json:"lat"`
		Lng      float64 `json:"lng"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Add to Redis GEO
	err := config.RedisClient.GeoAdd(
		config.Ctx,
		"drivers",
		&redis.GeoLocation{
			Name:      input.DriverID,
			Latitude:  input.Lat,
			Longitude: input.Lng,
		},
	).Err()

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to add driver"})
		return
	}

	// Set default status
	config.RedisClient.Set(
		config.Ctx,
		"driver_status:"+input.DriverID,
		"AVAILABLE",
		0,
	)

	c.JSON(200, gin.H{
		"message":   "driver created",
		"driver_id": input.DriverID,
	})
}

func UpdateDriverLocation(c *gin.Context) {
	driverID := c.Param("id")

	var input struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// OPTIONAL: check if driver exists
	exists, _ := config.RedisClient.Get(
		config.Ctx,
		"driver_status:"+driverID,
	).Result()

	if exists == "" {
		c.JSON(400, gin.H{"error": "driver does not exist"})
		return
	}

	// Update location
	err := config.RedisClient.GeoAdd(
		config.Ctx,
		"drivers",
		&redis.GeoLocation{
			Name:      driverID,
			Latitude:  input.Lat,
			Longitude: input.Lng,
		},
	).Err()

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to update location"})
		return
	}

	c.JSON(200, gin.H{
		"message":   "location updated",
		"driver_id": driverID,
	})
}

func SetAvailability(c *gin.Context) {
	driverID := c.Param("id")

	var input struct {
		Available bool `json:"available"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	status := "AVAILABLE"
	msg := "driver set to AVAILABLE"

	if !input.Available {
		status = "UNAVAILABLE"
		msg = "driver set to UNAVAILABLE"
	}

	key := "driver_status:" + driverID
	config.RedisClient.Set(config.Ctx, key, status, 0)

	c.JSON(200, gin.H{
		"message":   msg,
		"driver_id": driverID,
		"status":    status,
	})
}

// Driver stats
func GetDriverStats(c *gin.Context) {
	keys, _ := config.RedisClient.Keys(config.Ctx, "driver_status:*").Result()

	total := len(keys)
	available := 0
	ongoing := 0
	unavailable := 0

	for _, key := range keys {
		status, _ := config.RedisClient.Get(config.Ctx, key).Result()

		switch status {
		case "AVAILABLE":
			available++
		case "ONGOING":
			ongoing++
		case "UNAVAILABLE":
			unavailable++
		}
	}

	c.JSON(200, gin.H{
		"total":       total,
		"available":   available,
		"ongoing":     ongoing,
		"unavailable": unavailable,
	})
}

func GetDriverList(c *gin.Context) {

	keys, _ := config.RedisClient.Keys(config.Ctx, "driver_status:*").Result()

	var drivers []gin.H

	for _, key := range keys {
		driverID := key[len("driver_status:"):]

		status, _ := config.RedisClient.Get(config.Ctx, key).Result()

		// get location from GEO
		pos, _ := config.RedisClient.GeoPos(config.Ctx, "drivers", driverID).Result()

		var lat, lng float64
		if len(pos) > 0 && pos[0] != nil {
			lng = pos[0].Longitude
			lat = pos[0].Latitude
		}

		drivers = append(drivers, gin.H{
			"driver_id": driverID,
			"status":    status,
			"lat":       lat,
			"lng":       lng,
		})
	}

	c.JSON(200, gin.H{
		"drivers": drivers,
	})
}

func ClearDrivers(c *gin.Context) {

	// delete all driver statuses
	statusKeys, _ := config.RedisClient.Keys(config.Ctx, "driver_status:*").Result()
	if len(statusKeys) > 0 {
		config.RedisClient.Del(config.Ctx, statusKeys...)
	}

	// delete GEO set
	config.RedisClient.Del(config.Ctx, "drivers")

	c.JSON(200, gin.H{
		"message": "all drivers cleared",
	})
}