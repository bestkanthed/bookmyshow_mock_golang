package main

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const unblockTime = 300 // Seconds
type Status int

const (
	Available Status = 1
	Blocked   Status = 2
	Confirmed Status = 3
)

type MovieSeat struct {
	Id            uint      `json:"_id" `
	MovieId       uint      `json:"movie_id" `
	SeatId        uint      `json:"seat_id" `
	ScreenId      uint      `json:"screen_id" `
	TheaterId     uint      `json:"theater_id" `
	ShowStartTime time.Time `json:"show_start_time" gorm:"type:datetime"`
	ShowEndTime   time.Time `json:"show_end_time" gorm:"type:datetime"`

	CustomerId        uint `json:"customer_id" `
	PaymentIdentifier uint `json:"payment_identifier" `
	Status            Status
}

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/bookmyshow?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println("Logging err:", err)
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&MovieSeat{})

	// If database is empty create dummy movie seats
	var movieSeat MovieSeat
	db.First(&movieSeat, 1)
	if movieSeat.Id == 0 {
		startTime := time.Now().Local().Add(time.Hour * time.Duration(24))
		endTime := time.Now().Local().Add(time.Hour * time.Duration(26))
		for i := 1; i <= 10; i++ {
			db.Create(&MovieSeat{Id: uint(i), MovieId: 2, SeatId: 23 + uint(i), ScreenId: 4, TheaterId: 104, ShowStartTime: startTime, ShowEndTime: endTime, Status: Available})
		}
	}

}

func findAvailableSeats() (error, []MovieSeat) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error, nil
	}
	var movieSeats []MovieSeat
	tx.Where(&MovieSeat{Status: Available}).Find(&movieSeats)
	fmt.Println(movieSeats)
	return tx.Commit().Error, movieSeats
}

func persistReservation(res *SeatReservationDTO) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	//TODO : Check that there is no overlapping reservation
	var movieSeat MovieSeat
	tx.Where(&MovieSeat{Id: res.MovieSeatId, Status: Available}).First(&movieSeat)
	if movieSeat.Id == 0 {
		return errors.New("Invalid Seat")
	}

	tx.Model(&movieSeat).Update(MovieSeat{CustomerId: res.CustomerId, PaymentIdentifier: res.PaymentIdentifier, Status: Blocked})

	go func() {
		time.Sleep(unblockTime * time.Second)
		rb := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				rb.Rollback()
			}
		}()
		if rb.Error != nil {
			fmt.Println(rb.Error)
		}
		rb.First(&movieSeat, movieSeat.Id)
		if movieSeat.Status != Confirmed {
			rb.Model(&movieSeat).Update(map[string]interface{}{"customer_id": 0, "payment_identifier": 0, "status": Available})
		}
		rb.Commit()
	}()

	return tx.Commit().Error
}

func updateReservationStatus(reservationDTO *SeatReservationDTO) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	var movieSeat MovieSeat
	if reservationDTO.PaymentStatus == "success" {
		tx.Where(&MovieSeat{
			Id:                reservationDTO.MovieSeatId,
			CustomerId:        reservationDTO.CustomerId,
			PaymentIdentifier: reservationDTO.PaymentIdentifier,
			Status:            Blocked,
		}).First(&movieSeat)
		if movieSeat.Status == Blocked {
			db.Model(&movieSeat).Update("status", Confirmed)
			return tx.Commit().Error
		} else {
			return errors.New("Invalid payment callback")
		}
	} else {
		return errors.New("Payment not successfull")
	}
}
