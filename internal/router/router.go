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

	r.POST("/v1/drivers/:id/location", handler.UpdateLocation)
	r.POST("/v1/rides", handler.CreateRide)
	r.POST("/v1/drivers/:id/accept", handler.AcceptRide)
	r.POST("/v1/trips/:id/end", handler.EndRide)
	r.POST("/v1/payments", handler.CreatePayment)
	r.GET("/v1/rides/:id", handler.GetRide)
	r.POST("/v1/drivers/:id/availability", handler.SetAvailability)
	r.GET("/v1/drivers/count", handler.GetDriverCount)

	return r
}