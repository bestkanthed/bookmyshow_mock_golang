# BookMyShow Mock Seat Blocking

This is a very simplified implementation of seat booking to enhance user experience by blocking seats for users so that they can complete their payment.
If, user is unable to complete the payment, the seats become available to other users after the specified timeout.

## Requirements
1) Go Lang  
2) MYSQL v8+ 

## MYSQL Setup
$ mysql  
mysql> CREATE DATABASE bookmyshow;  
mysql> CREATE TABLE movie_seats;  

## GO Setup
> Place the folder containing this file in $GOPATH/src/  
 
$ go build  
$ ./bookmyshow_mock  

## Get Available Seats API
curl -X GET 127.0.0.1:8080/v1/bookmyshow/available

## Create Booking API
curl -d '{"movie_seat_id": 1, "payment_identifier": 61, "customer_id":122 }' -H "Content-Type: application/json"  -X POST 127.0.0.1:8080/v1/bookmyshow/reservation

## Confirm Payment Callback API
curl -d '{"movie_seat_id": 7, "payment_identifier": 61, "customer_id":122, "payment_status": "success" }' -H "Content-Type: application/json"  -X POST 127.0.0.1:8080/v1/bookmyshow/payment/callback


## Reference
Hands-On Software Architecture with Golang: Design and Architect Highly Scalable and Robust Applications Using Go  
By Jyotiswarup Raiturkar