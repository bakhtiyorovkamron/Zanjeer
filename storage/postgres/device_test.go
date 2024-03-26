package postgres

import (
	"testing"

	"github.com/Projects/Zanjeer/config"
	"github.com/Projects/Zanjeer/models"
	"github.com/Projects/Zanjeer/pkg/db"
	"github.com/Projects/Zanjeer/pkg/logger"
)

func TestSetLocation(t *testing.T) {

	cfg := config.Load()
	logger := logger.New(cfg.LogLevel)

	db, err := db.New(cfg)
	if err != nil {
		logger.Error("Error while connecting to database", err)
	} else {
		logger.Info("Successfully connected to database")
	}

	pg := New(db, logger, cfg)

	data := models.Record{
		// {
		// 	Imei:      "359633103869421",
		// 	Longitude: 94898293,
		// 	Latitude:  123444589,
		// },
		// {
		// 	Imei:      "359633103869421",
		// 	Longitude: 94898293,
		// 	Latitude:  123444589,
		// },
		// {
		// 	Imei:      "359633103869421",
		// 	Longitude: 94898293,
		// 	Latitude:  123444589,
		// },
		// {
		// 	Imei:      "359633103869421",
		// 	Longitude: 94898293,
		// 	Latitude:  123444589,
		// },
	}

	if pg.SetLocation(data) != nil {
		panic(err)
	}

	// err = pg.Postgres().SetLocation()

}
