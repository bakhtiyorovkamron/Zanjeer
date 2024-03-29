package main

import (
	"log"
	"net/http"

	"github.com/Projects/Zanjeer/config"
	"github.com/Projects/Zanjeer/models"
	"github.com/Projects/Zanjeer/pkg/db"
	"github.com/Projects/Zanjeer/pkg/logger"
	"github.com/Projects/Zanjeer/storage/postgres"
	"github.com/gin-gonic/gin"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	IMEI      string  `json:"imei"`
}
type H struct {
	storage postgres.PostgresI
}

func (h *H) handleLocation(c *gin.Context) {
	var loc Location
	if err := c.BindJSON(&loc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Do something with the location data, for example, print it
	log.Printf("Received location - Longitude: %f, Latitude: %f, IMEI: %s\n", loc.Longitude, loc.Latitude, loc.IMEI)

	h.storage.SetLocation(models.Record{
		Imei:      loc.IMEI,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
	})

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Location received successfully"})
}

func main() {

	cfg := config.Load()
	logger := logger.New(cfg.LogLevel)

	db, err := db.New(cfg)
	if err != nil {
		logger.Error("Error while connecting to database", err)
	} else {
		logger.Info("Successfully connected to database")
	}

	pg := postgres.New(db, logger, cfg)

	h := H{
		storage: pg,
	}

	r := gin.Default()

	r.POST("/location", h.handleLocation)

	log.Println("Server is listening on port 1234...")
	log.Fatal(r.Run(":1234"))
}
