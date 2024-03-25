package postgres

import (
	"fmt"

	"github.com/lib/pq"

	"github.com/Projects/Zanjeer/models"
)

func (p *postgresRepo) SetLocation(req []models.Record) error {
	var (
		longitude, latitude = []string{}, []string{}
		imei                string
	)

	query := `call set_location($1,$2,$3)`

	for _, v := range req {
		if v.Imei != "" {
			imei = v.Imei
		}
		longitude = append(longitude, fmt.Sprintf("%f", v.Longitude))
		latitude = append(latitude, fmt.Sprintf("%f", v.Latitude))
	}

	longitudeArray := pq.Array(longitude)
	latitudeArray := pq.Array(latitude)

	result, err := p.Db.Db.Exec(query, imei, longitudeArray, latitudeArray)
	if err != nil {
		fmt.Println("Error executing:", err)
		return err
	}
	fmt.Println("Result:", result)
	return nil
}
