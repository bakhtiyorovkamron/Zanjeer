package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Projects/Zanjeer/config"
	"github.com/Projects/Zanjeer/helpers"
	"github.com/Projects/Zanjeer/pkg/db"
	"github.com/Projects/Zanjeer/pkg/logger"
	"github.com/Projects/Zanjeer/storage"
)

type DeviceData struct {
	IMEI         string `json:"imei"`
	TimeDate     string `json:"date"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	NumberOfData string
	Altitude     string
	Angle        string
}

const port = "1234"

func main() {

	cfg := config.Load()

	logger := logger.New(cfg.LogLevel)

	db, err := db.New(cfg)
	if err != nil {
		logger.Error("Error while connecting to database", err)
	} else {
		logger.Info("Successfully connected to database")
	}

	// Listen for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port " + port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn, db, logger, cfg)
	}
}

func handleClient(conn net.Conn, db *db.Postgres, log *logger.Logger, cfg config.Config) {

	pg := storage.New(db, log, cfg)

	defer conn.Close()

	var messageTrans = map[int]func(step *int, imei *string, msg string, conn net.Conn){}

	messageTrans[1] = takeImei

	// Create a buffer to read data into
	buffer := make([]byte, 1024)
	var (
		imeiTaken bool    = true
		step      *int    = new(int)
		imei      *string = new(string)
	)
	*step = 1

	for {
		// Read data from the client
		size, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		if imeiTaken {

			message := hex.EncodeToString(buffer[:size])

			if helpers.Imei(buffer) {
				*imei = string(buffer[2:17])
			}

			switch *step {
			case 1:
				messageTrans[*step](step, imei, message, conn)
			case 2:

				data, err := helpers.ReadMessage(message, *imei)
				if err == nil {
					if pg.Postgres().SetLocation(data) != nil {
						fmt.Println("ERROR while Setting location!", err)
					}
				} else {
					fmt.Println("ERROR while Reading Message :", err)
				}

				d, _ := json.MarshalIndent(data, "", " ")
				fmt.Println(string(d))
				if err == nil {
					fmt.Println("Size sent : ", data.Size)
					conn.Write([]byte{0, 0, 0, uint8(size)})
				}
			}
		} else {
			b := []byte{0} // 0x00 if we decline the message
			conn.Write(b)
			break
		}

	}
}

func takeImei(step *int, imei *string, msg string, conn net.Conn) {
	firstReply := []byte{1}
	*step = 2
	conn.Write(firstReply)
}

func Decoder(enCode string) (DeviceData, error) {

	if len(enCode) < 58 {
		return DeviceData{}, fmt.Errorf("Minimum packet size is 45 Bytes, got %v", len(enCode))
	}

	// zeroBytes := (enCode[0:8])
	// dataFieldLength := (enCode[8:16])
	codecId := (enCode[16:18])
	numberOfData := (enCode[18:20])
	timestamp := (enCode[20:36])
	// priority := (enCode[36:38])
	longitude := (enCode[38:46])
	latitude := (enCode[46:54])

	altitude := (enCode[54:56])
	angle := (enCode[56:58])

	if codecId != "08" && codecId != "8E" {
		return DeviceData{}, fmt.Errorf("Invalid Codec ID, want 0x08 or 0x8E, get %v", codecId)
	}

	return DeviceData{
		TimeDate:     timestamp,
		Lat:          latitude,
		Lng:          longitude,
		NumberOfData: numberOfData,
		Altitude:     altitude,
		Angle:        angle,
	}, nil

}

func DecToLocation(dec int64) string {
	if dec <= 0 {
		return ""
	}
	location := fmt.Sprintf("%d", dec)

	return location[0:2] + "." + location[2:]
}
func HexToDec(hex string) int64 {
	num, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0
	}
	return num
}

func hexToTime(hexTimestamp string) (time.Time, error) {
	// Decode the hex string into a byte slice
	hexBytes, err := hex.DecodeString(hexTimestamp)
	if err != nil {
		return time.Time{}, err
	}

	// Parse the byte slice as a Unix timestamp
	// Assuming that the timestamp is a 64-bit integer (8 bytes)
	unixTimestamp := int64(0)
	for _, b := range hexBytes {
		unixTimestamp = (unixTimestamp << 8) | int64(b)
	}

	// Convert Unix timestamp to time.Time
	timestamp := time.Unix(unixTimestamp, 0)

	return timestamp, nil
}
