package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/v1/bookmyshow")
	{
		v1.POST("/reservation", createReservation)
		v1.POST("/payment/callback", updateReservation)
		v1.GET("/available", getAvailableSeats)

	}
	router.Run()
}

type SeatReservationDTO struct {
	CustomerId        uint   `json:"customer_id" `
	MovieSeatId       uint   `json:"movie_seat_id" `
	PaymentIdentifier uint   `json:"payment_identifier" `
	PaymentStatus     string `json:"payment_status" `
}

func getAvailableSeats(c *gin.Context) {
	var (
		err        error
		movieSeats []MovieSeat
	)
	err, movieSeats = findAvailableSeats()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusAccepted, gin.H{"available": movieSeats})
	}
}

func createReservation(c *gin.Context) {
	var (
		reservationDTO SeatReservationDTO
		err            error
	)

	if err = c.ShouldBindJSON(&reservationDTO); err == nil {
		if err = persistReservation(&reservationDTO); err == nil {
			c.JSON(http.StatusAccepted, gin.H{"status": "Seat blocked, make payment to reserve."})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"error": "Seat not available."})
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func updateReservation(c *gin.Context) {
	var (
		reservationDTO SeatReservationDTO
		err            error
	)

	if err = c.ShouldBindJSON(&reservationDTO); err == nil {
		if err = updateReservationStatus(&reservationDTO); err == nil {
			c.JSON(http.StatusAccepted, gin.H{"status": "Success. Your seat is confirmed!"})
		} else {
			c.JSON(http.StatusAccepted, gin.H{"error": "Payment unsuccessfull! Please try again!"})
		}
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
