package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePayment(c *gin.Context) {
	var input struct {
		RideID string `json:"ride_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// mock success
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"ride_id": input.RideID,
		"amount":  input.Amount,
	})
}