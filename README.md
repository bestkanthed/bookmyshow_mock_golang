# Requirements
1) Go Lang 
2) MYSQL v8+

# MYSQL
$ mysql
> CREATE DATABASE bookmyshow;
> CREATE TABLE movie_seats;

# GO Setup
// Place the folder containing this file in $GOPATH/src/
$ go build
$ ./bookmyshow_mock

# Get available seats api
curl -X GET 127.0.0.1:8080/v1/bookmyshow/available

# Create booking api
curl -d '{"movie_seat_id": 1, "payment_identifier": 61, "customer_id":122 }' -H "Content-Type: application/json"  -X POST 127.0.0.1:8080/v1/bookmyshow/reservation

# Confirm payment callback api
curl -d '{"movie_seat_id": 7, "payment_identifier": 61, "customer_id":122, "payment_status": "success" }' -H "Content-Type: application/json"  -X POST 127.0.0.1:8080/v1/bookmyshow/payment/callback