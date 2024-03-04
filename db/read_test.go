package db

import (
	"encoding/json"
	"testing"

	"github.com/Projects/Zanjeer/models"
)

func TestWriteandRead(t *testing.T) {
	data := Read()
	data = append(data, models.Record{
		Imei: "fuckyou",
	})
	file, _ := json.MarshalIndent(data, "", " ")
	Write(file)

}
