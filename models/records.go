package models

import "time"

type Record struct {
	Imei      string
	Location  Location
	Time      time.Time
	Angle     int16
	Speed     int16
	Longitude string
	Latitude  string
}
type Location struct {
	Type        string
	Coordinates []int32
}
