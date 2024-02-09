package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

type DeviceData struct {
	TimeDate     string `json:"date"`
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	NumberOfData string
	Altitude     string
	Angle        string
}

func main() {
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

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
			Post(DecToLocation(HexToDec(deviceInfo.Lat)), DecToLocation(HexToDec(deviceInfo.Lng)))

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

func Post(lat, long string) {

	url := fmt.Sprintf("https://bf19-188-113-230-172.ngrok-free.app/longlat?lat=%s&lon=%s", lat, long)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
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
