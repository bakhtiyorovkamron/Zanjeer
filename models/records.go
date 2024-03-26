package models

import "time"

type Record struct {
	Imei      string
	Location  Location
	Time      time.Time
	Angle     int64
	Speed     int64
	Longitude float64
	Latitude  float64
	Altitude  int64
	Size      uint8
}
type Location struct {
	Type        string
	Coordinates []int32
}
