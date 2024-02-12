package helpers

import (
	"bytes"
	"fmt"

	"github.com/Projects/Zanjeer/models"
)

func ParseData(data []byte, size int, imei string) (elements []models.Record, err error) {
	reader := bytes.NewBuffer(data)
	// fmt.Println("Reader Size:", reader.Len())

	// Header
	reader.Next(4)                                    // 4 Zero Bytes
	dataLength, err := streamToInt32(reader.Next(4))  // Header
	reader.Next(1)                                    // CodecID
	recordNumber, err := streamToInt8(reader.Next(1)) // Number of Records
	fmt.Println("Length of data:", dataLength)

	elements = make([]models.Record, recordNumber)

	var i int8 = 0
	for i < recordNumber {
		timestamp, err := streamToTime(reader.Next(8))
		if err != nil {
			return elements, fmt.Errorf("")
		} // Timestamp
		reader.Next(1) // Priority

		// GPS Element
		longitudeInt, err := streamToInt32(reader.Next(4)) // Longitude
		fmt.Println("#####################")
		fmt.Println("longitudeInt :", longitudeInt)
		longitude := float64(longitudeInt) // PRECISION
		latitudeInt, err := streamToInt32(reader.Next(4))
		fmt.Println("latitudeInt :", latitudeInt) // Latitude
		fmt.Println("#####################")
		fmt.Println()
		latitude := float64(latitudeInt) // PRECISION

		reader.Next(2)                              // Altitude
		angle, err := streamToInt16(reader.Next(2)) // Angle
		reader.Next(1)                              // Satellites
		speed, err := streamToInt16(reader.Next(2)) // Speed

		if err != nil {
			return elements, fmt.Errorf("")
		}

		elements[i] = models.Record{
			Imei: imei,
			Location: models.Location{Type: "Point",
				Coordinates: []float64{longitude, latitude}},
			Time:  timestamp,
			Angle: angle,
			Speed: speed}

		// IO Events Elements

		reader.Next(1) // ioEventID
		reader.Next(1) // total Elements

		stage := 1
		for stage <= 4 {
			stageElements, err := streamToInt8(reader.Next(1))
			if err != nil {
				break
			}

			var j int8 = 0
			for j < stageElements {
				reader.Next(1) // elementID

				switch stage {
				case 1: // One byte IO Elements
					_, err = streamToInt8(reader.Next(1))
				case 2: // Two byte IO Elements
					_, err = streamToInt16(reader.Next(2))
				case 3: // Four byte IO Elements
					_, err = streamToInt32(reader.Next(4))
				case 4: // Eigth byte IO Elements
					_, err = streamToInt64(reader.Next(8))
				}
				j++
			}
			stage++
		}

		if err != nil {
			return elements, fmt.Errorf("Error while reading IO Elements")
		}

		// fmt.Println("Timestamp:", timestamp)
		// fmt.Println("Longitude:", longitude, "Latitude:", latitude)

		i++
	}

	// Once finished with the records we read the Record Number and the CRC

	_, err = streamToInt8(reader.Next(1))  // Number of Records
	_, err = streamToInt32(reader.Next(4)) // CRC

	return
}
