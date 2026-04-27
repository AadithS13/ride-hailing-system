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

	config.RedisClient.GeoAdd(config.Ctx, "drivers", &redis.GeoLocation{
		Name:      driverID,
		Longitude: input.Lng,
		Latitude:  input.Lat,
	})

	config.RedisClient.Set(config.Ctx, "driver_status:"+driverID, "AVAILABLE", 0)

	c.JSON(200, gin.H{"message": "location updated"})
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
	if !input.Available {
		status = "UNAVAILABLE"
	}

	key := "driver_status:" + driverID
	config.RedisClient.Set(config.Ctx, key, status, 0)

	c.JSON(200, gin.H{"status": status})
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