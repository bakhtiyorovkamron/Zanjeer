package helpers

import (
	"fmt"
	"strconv"

	"github.com/Projects/Zanjeer/models"
)

func ReadMessage(message, imei string) (models.Record, error) {

	if len(message) < 68 {
		return models.Record{}, fmt.Errorf("invalid data")
	}

	longitude, err := strconv.ParseInt(message[38:46], 16, 64)
	if err != nil {
		return models.Record{}, err
	}
	latitude, err := strconv.ParseInt(message[46:54], 16, 64)
	if err != nil {
		return models.Record{}, err
	}
	// altitude, err := strconv.ParseInt(message[54:58], 16, 64)
	// if err != nil {
	// 	return models.Record{}, err
	// }

	// angle, err := strconv.ParseInt(message[58:62], 16, 64)
	// if err != nil {
	// 	return models.Record{}, err
	// }
	// speed, err := strconv.ParseInt(message[64:68], 16, 64)
	// if err != nil {
	// 	return models.Record{}, err
	// }
	size, err := StringToUint8(message)
	if err != nil {
		return models.Record{}, err
	}

	return models.Record{
		Imei:      imei,
		Longitude: float64(longitude) * 0.0000001,
		Latitude:  float64(latitude) * 0.0000001,
		// Altitude:  altitude,
		// Angle:     angle,
		// Speed:     speed,
		Size: uint8(size),
	}, nil
}
