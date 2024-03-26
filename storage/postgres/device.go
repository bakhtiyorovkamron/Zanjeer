package postgres

import (
	"fmt"

	"github.com/lib/pq"

	"github.com/Projects/Zanjeer/models"
)

func (p *postgresRepo) SetLocation(req models.Record) error {
	var (
		longitude, latitude = []string{}, []string{}
	)

	query := `call set_location($1,$2,$3)`

	longitude = append(longitude, fmt.Sprintf("%f", req.Longitude))
	latitude = append(latitude, fmt.Sprintf("%f", req.Latitude))

	longitudeArray := pq.Array(longitude)
	latitudeArray := pq.Array(latitude)
	if len(longitude) == 0 || len(latitude) == 0 {
		return fmt.Errorf("empty latitude array")
	}

	_, err := p.Db.Db.Exec(query, req.Imei, longitudeArray, latitudeArray)
	if err != nil {
		fmt.Println("Error executing:", err)
		return err
	}
	return nil
}
