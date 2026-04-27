package router

import (
	"ride-hailing/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/v1/drivers/:id/location", handler.UpdateLocation)
	r.POST("/v1/rides", handler.CreateRide)
	r.POST("/v1/drivers/:id/accept", handler.AcceptRide)
	r.POST("/v1/trips/:id/end", handler.EndRide)
	r.POST("/v1/payments", handler.CreatePayment)

	return r
}