package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Projects/Zanjeer/helpers"
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
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to read data into
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		_ = n
		fmt.Println(string(buffer[:n]))
		encodedString := hex.EncodeToString(buffer)
		deviceInfo, err := Decoder(encodedString)
		if err == nil {
			Post(DeviceData{
				IMEI: "359633103869421",
				Lat:  DecToLocation(HexToDec(deviceInfo.Lat)),
				Lng:  DecToLocation(HexToDec(deviceInfo.Lng)),
			})
			// Post(DecToLocation(HexToDec(deviceInfo.Lat)), DecToLocation(HexToDec(deviceInfo.Lng)))

			fmt.Println("##############################################")
			hexToTime, _ := hexToTime(deviceInfo.TimeDate)
			// fmt.Println("encodedString :", encodedString)
			fmt.Println("Time:", hexToTime)
			fmt.Println("Device Info :", time.Now())
			fmt.Println("Time:", deviceInfo.TimeDate)
			fmt.Println("Longitude:", deviceInfo.Lng)
			fmt.Println("Latitude:", deviceInfo.Lat)
			fmt.Println("NUmber of data:", deviceInfo.NumberOfData)
			fmt.Println("##############################################")
			fmt.Println()
			// fmt.Println
			decodedByteArray, err := hex.DecodeString("000000" + deviceInfo.NumberOfData)

			if err != nil {
				fmt.Println("Unable to convert hex to byte. ", err)
			}
			fmt.Println("Number of data:", decodedByteArray)
			conn.Write(decodedByteArray)
			fmt.Println("sent ", decodedByteArray)
			// conn.Write([]byte{0x13})

			return
		}

		fmt.Println(string(buffer[:n]))
		nullOne := []byte{0x01}
		conn.Write(nullOne)
		fmt.Println("sent 01 :", nullOne)

	}
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

const (
	url    = "http://5.39.92.50:8081/longlat"
	method = "POST"
)

func Post(d DeviceData) {

	payload := strings.NewReader(fmt.Sprintf(`{
		"imei": "%s",
		"lon": "%s",
		"lat": "%s"
	}`, d.IMEI, d.Lng, d.Lat))

	_, err := helpers.SendHTTPRequest(url, method, payload)
	if err != nil {
		log.Fatal("Error while updating device location", err)
	}

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
