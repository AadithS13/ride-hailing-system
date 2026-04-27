package handler

import (

	"ride-hailing/internal/config"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func UpdateLocation(c *gin.Context) {
	driverID := c.Param("id")

	var input struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Add/update location
	config.RedisClient.GeoAdd(config.Ctx, "drivers", &redis.GeoLocation{
		Name:      driverID,
		Longitude: input.Lng,
		Latitude:  input.Lat,
	})

	// Initialize status if not present
	key := "driver_status:" + driverID
	exists, _ := config.RedisClient.Exists(config.Ctx, key).Result()
	if exists == 0 {
		config.RedisClient.Set(config.Ctx, key, "AVAILABLE", 0)
	}

	c.JSON(200, gin.H{
		"message":   "driver added/updated successfully",
		"driver_id": driverID,
		"location": gin.H{
			"lat": input.Lat,
			"lng": input.Lng,
		},
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