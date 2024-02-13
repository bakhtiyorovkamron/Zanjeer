package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Projects/Zanjeer/models"
)

func Write(data []byte) {

	_ = ioutil.WriteFile("../data/fmb_location.json", data, 0644)
}
func Read() []models.Record {
	jsonFile, err := os.Open("../data/fmb_location.json")
	if err != nil {
		fmt.Println("error while reading from json")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var location []models.Record

	if json.Unmarshal(byteValue, &location) == nil {
		return location
	}
	return location
}
