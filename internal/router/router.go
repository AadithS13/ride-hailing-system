package router

import (
	"ride-hailing/internal/handler"

	"github.com/gin-gonic/gin"
	"ride-hailing/internal/config"

	 "github.com/gin-contrib/cors"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	r.Use(func(c *gin.Context) {
		if config.NRApp == nil {
			c.Next()
			return
		}

		txn := config.NRApp.StartTransaction(c.Request.Method + " " + c.Request.URL.Path)
		defer txn.End()

		txn.SetWebRequestHTTP(c.Request)

		c.Set("txn", txn)

		c.Next()

		txn.SetWebResponse(nil)
	})



	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/v1/drivers", handler.CreateDriver)
	r.POST("/v1/drivers/:id/location", handler.UpdateDriverLocation)
	r.POST("/v1/drivers/:id/availability", handler.SetAvailability)
	r.GET("/v1/drivers/stats", handler.GetDriverStats)
	r.GET("/v1/drivers", handler.GetDriverList)
	r.DELETE("/v1/drivers", handler.ClearDrivers)

	// ride routes
	r.POST("/v1/rides", handler.CreateRide)
	r.GET("/v1/rides/:id", handler.GetRide)
	r.GET("/v1/rides", handler.GetAllRides)
	r.POST("/v1/drivers/:id/end", handler.EndRideByDriver)

	// driver accepts ride
	r.POST("/v1/drivers/:id/accept", handler.AcceptRideHandler)

	// end trip
	r.POST("/v1/trips/:id/end", handler.EndRideHandler)	

	return r
}