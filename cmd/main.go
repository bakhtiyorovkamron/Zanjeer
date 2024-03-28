package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fileAccessTest()
}

func inputTrigger() {
	fmt.Println("Paste full 'Codec 8' packet to parse it or:")
	fmt.Println("Type SERVER to start the server or:")
	fmt.Println("Type EXIT to stop the program")
	deviceIMEI := "default_IMEI"
	var userInput string
	fmt.Print("waiting for input: ")
	fmt.Scanln(&userInput)
	userInput = strings.ToUpper(userInput)
	if userInput == "EXIT" {
		fmt.Println("exiting program............")
		os.Exit(0)
	} else if userInput == "SERVER" {
		startServerTrigger()
	} else {
		if codec8eChecker(strings.ReplaceAll(userInput, " ", "")) == false {
			fmt.Println("Wrong input or invalid Codec8 packet")
			fmt.Println()
			inputTrigger()
		} else {
			codecParserTrigger(userInput, deviceIMEI, "USER")
		}
	}
}

func crc16Arc(data string) bool {
	dataPartLengthCRC, _ := strconv.ParseInt(data[8:16], 16, 64)
	dataPartForCRC := data[16 : 16+2*dataPartLengthCRC]
	crc16ArcFromRecord := data[16+len(dataPartForCRC)*2 : 24+len(dataPartForCRC)*2]
	crc := 0
	for _, byte := range dataPartForCRC {
		crc ^= int(byte)
		for i := 0; i < 8; i++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	if strings.ToUpper(crc16ArcFromRecord) == strings.ToUpper(fmt.Sprintf("%08X", crc)) {
		fmt.Println("CRC check passed!")
		fmt.Printf("Record length: %d characters // %d bytes\n", len(data), len(data)/2)
		return true
	} else {
		fmt.Println("CRC check Failed!")
		return false
	}
}

func codec8eChecker(codec8Packet string) bool {
	if strings.ToUpper(codec8Packet[16:18]) != "8E" && strings.ToUpper(codec8Packet[16:18]) != "08" {
		fmt.Println()
		fmt.Println("Invalid packet!!!!!!!!!!!!!!!!!!!")
		return false
	} else {
		return crc16Arc(codec8Packet)
	}
}

func codecParserTrigger(codec8Packet string, deviceIMEI string, props string) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error occured: %v enter proper Codec8 packet or EXIT!!!\n", r)
			inputTrigger()
		}
	}()
	return codec8EParser(strings.ReplaceAll(codec8Packet, " ", ""), deviceIMEI, props)
}

func imeiChecker(hexIMEI string) bool {
	imeiLength, _ := strconv.ParseInt(hexIMEI[:4], 16, 64)
	if imeiLength != int64(len(hexIMEI[4:])/2) {
		return false
	} else {
		asciiIMEI := asciiIMEIConverter(hexIMEI)
		fmt.Printf("IMEI received = %s\n", asciiIMEI)
		if _, err := strconv.Atoi(asciiIMEI); err != nil || len(asciiIMEI) != 15 {
			fmt.Println("Not an IMEI - is not numeric or wrong length!")
			return false
		} else {
			return true
		}
	}
}

func asciiIMEIConverter(hexIMEI string) string {
	decoded, _ := hex.DecodeString(hexIMEI[4:])
	return string(decoded)
}

func startServerTrigger() {
	fmt.Println("Starting server!")
	ln, _ := net.Listen("tcp", ":1234")
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		conn.SetDeadline(time.Now().Add(20 * time.Second))
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	deviceIMEI := "default_IMEI"
	for {
		data := make([]byte, 1280)
		_, err := conn.Read(data)
		if err != nil {
			break
		} else if imeiChecker(fmt.Sprintf("%X", data)) != false {
			deviceIMEI = asciiIMEIConverter(fmt.Sprintf("%X", data))
			imeiReply := []byte{1}
			conn.Write(imeiReply)
			fmt.Printf("-- %s sending reply = %X\n", timeStamper(), imeiReply)
		} else if codec8eChecker(strings.ReplaceAll(fmt.Sprintf("%X", data), " ", "")) != false {
			recordNumber := codecParserTrigger(fmt.Sprintf("%X", data), deviceIMEI, "SERVER")
			fmt.Printf("received records %d\n", recordNumber)
			fmt.Printf("from device IMEI = %s\n\n", deviceIMEI)
			recordResponse := make([]byte, 4)
			binary.BigEndian.PutUint32(recordResponse, uint32(recordNumber))
			conn.Write(recordResponse)
			fmt.Printf("// %s // response sent = %X\n", timeStamper(), recordResponse)
		} else {
			fmt.Printf("// %s // no expected DATA received - dropping connection\n", timeStamper())
			break
		}
	}
}

func codec8EParser(codec8EPacket string, deviceIMEI string, props string) int {
	ioDictRaw := make(map[string]interface{})
	ioDictRaw["device_IMEI"] = deviceIMEI
	ioDictRaw["server_time"] = timeStamperForJSON()
	ioDictRaw["data_length"] = fmt.Sprintf("Record length: %d characters // %d bytes", len(codec8EPacket), len(codec8EPacket)/2)
	ioDictRaw["_raw_data__"] = codec8EPacket
	jsonPrinterRawData(ioDictRaw, deviceIMEI)
	zeroBytes := codec8EPacket[:8]
	fmt.Printf("\nzero bytes = %s\n", zeroBytes)
	dataFieldLength, _ := strconv.ParseInt(codec8EPacket[8:16], 16, 64)
	fmt.Printf("data field length = %d bytes\n", dataFieldLength)
	codecType := codec8EPacket[16:18]
	fmt.Printf("codec type = %s\n", codecType)
	dataStep := 4
	if strings.ToUpper(codecType) == "08" {
		dataStep = 2
	}
	numberofRecords, _ := strconv.ParseInt(codec8EPacket[18:20], 16, 64)
	fmt.Printf("number of records = %d\n", numberofRecords)
	recordNumber := 1
	avlDataStart := codec8EPacket[20:]
	dataFieldPosition := 0
	for int64(dataFieldPosition) < (2*dataFieldLength - 6) {
		ioDict := make(map[string]interface{})
		ioDict["device_IMEI"] = deviceIMEI
		ioDict["server_time"] = timeStamperForJSON()
		fmt.Printf("\ndata from record %d\n", recordNumber)
		timestamp := avlDataStart[dataFieldPosition : dataFieldPosition+16]
		ioDict["_timestamp_"] = deviceTimeStamper(timestamp)
		fmt.Printf("timestamp = %s\n", deviceTimeStamper(timestamp))
		ioDict["_rec_delay_"] = recordDelayCounter(timestamp)
		dataFieldPosition += len(timestamp)
		priority := avlDataStart[dataFieldPosition : dataFieldPosition+2]
		priorityInt, _ := strconv.ParseInt(priority, 16, 64)
		ioDict["priority"] = priorityInt
		fmt.Printf("record priority = %d\n", priorityInt)
		dataFieldPosition += len(priority)
		longitude := avlDataStart[dataFieldPosition : dataFieldPosition+8]
		longitudeInt := int32(binary.BigEndian.Uint32([]byte(longitude)))
		ioDict["longitude"] = longitudeInt
		fmt.Printf("longitude = %d\n", longitudeInt)
		dataFieldPosition += len(longitude)
		latitude := avlDataStart[dataFieldPosition : dataFieldPosition+8]
		latitudeInt := int32(binary.BigEndian.Uint32([]byte(latitude)))
		ioDict["latitude"] = latitudeInt
		fmt.Printf("latitude = %d\n", latitudeInt)
		dataFieldPosition += len(latitude)
		altitude := avlDataStart[dataFieldPosition : dataFieldPosition+4]
		altitudeInt, _ := strconv.ParseInt(altitude, 16, 64)
		ioDict["altitude"] = altitudeInt
		fmt.Printf("altitude = %d\n", altitudeInt)
		dataFieldPosition += len(altitude)
		angle := avlDataStart[dataFieldPosition : dataFieldPosition+4]
		angleInt, _ := strconv.ParseInt(angle, 16, 64)
		ioDict["angle"] = angleInt
		fmt.Printf("angle = %d\n", angleInt)
		dataFieldPosition += len(angle)
		satellites := avlDataStart[dataFieldPosition : dataFieldPosition+2]
		satellitesInt, _ := strconv.ParseInt(satellites, 16, 64)
		ioDict["satellites"] = satellitesInt
		fmt.Printf("satellites = %d\n", satellitesInt)
		dataFieldPosition += len(satellites)
		speed := avlDataStart[dataFieldPosition : dataFieldPosition+4]
		speedInt, _ := strconv.ParseInt(speed, 16, 64)
		ioDict["speed"] = speedInt
		fmt.Printf("speed = %d\n", speedInt)
		dataFieldPosition += len(speed)
		eventIOID := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
		eventIOIDInt, _ := strconv.ParseInt(eventIOID, 16, 64)
		ioDict["eventID"] = eventIOIDInt
		fmt.Printf("event ID = %d\n", eventIOIDInt)
		dataFieldPosition += len(eventIOID)
		totalIOElements := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
		totalIOElementsParsed, _ := strconv.ParseInt(totalIOElements, 16, 64)
		fmt.Printf("total I/O elements in record %d = %d\n", recordNumber, totalIOElementsParsed)
		dataFieldPosition += len(totalIOElements)
		byte1IONumber := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
		byte1IONumberParsed, _ := strconv.ParseInt(byte1IONumber, 16, 64)
		fmt.Printf("1 byte io count = %d\n", byte1IONumberParsed)
		dataFieldPosition += len(byte1IONumber)
		if byte1IONumberParsed > 0 {
			i := int64(1)
			for i <= byte1IONumberParsed {
				key := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
				dataFieldPosition += len(key)
				value := avlDataStart[dataFieldPosition : dataFieldPosition+2]
				ioDict[fmt.Sprintf("%d", key)] = sortingHat(key, value)
				dataFieldPosition += len(value)
				fmt.Printf("avl_ID: %d : %v\n", key, ioDict[fmt.Sprintf("%d", key)])
				i++
			}
		}
		byte2IONumber := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
		byte2IONumberParsed, _ := strconv.ParseInt(byte2IONumber, 16, 64)
		fmt.Printf("2 byte io count = %d\n", byte2IONumberParsed)
		dataFieldPosition += len(byte2IONumber)
		if byte2IONumberParsed > 0 {
			i := int64(1)
			for i <= byte2IONumberParsed {
				key := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
				dataFieldPosition += len(key)
				value := avlDataStart[dataFieldPosition : dataFieldPosition+4]
				ioDict[fmt.Sprintf("%d", key)] = sortingHat(key, value)
				dataFieldPosition += len(value)
				fmt.Printf("avl_ID: %d : %v\n", key, ioDict[fmt.Sprintf("%d", key)])
				i++
			}
		}
		byte4IONumber := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
		byte4IONumberParsed, _ := strconv.ParseInt(byte4IONumber, 16, 64)
		fmt.Printf("4 byte io count = %d\n", byte4IONumberParsed)
		dataFieldPosition += len(byte4IONumber)
		if byte4IONumberParsed > 0 {
			i := int64(1)
			for i <= byte4IONumberParsed {
				key := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
				dataFieldPosition += len(key)
				value := avlDataStart[dataFieldPosition : dataFieldPosition+8]
				ioDict[fmt.Sprintf("%d", key)] = sortingHat(key, value)
				dataFieldPosition += len(value)
				fmt.Printf("avl_ID: %d : %v\n", key, ioDict[fmt.Sprintf("%d", key)])
				i++
			}
		}
		byte8IONumber := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
		byte8IONumberParsed, _ := strconv.ParseInt(byte8IONumber, 16, 64)
		fmt.Printf("8 byte io count = %d\n", byte8IONumberParsed)
		dataFieldPosition += len(byte8IONumber)
		if byte8IONumberParsed > 0 {
			i := int64(1)
			for i <= byte8IONumberParsed {
				key := avlDataStart[dataFieldPosition : dataFieldPosition+dataStep]
				dataFieldPosition += len(key)
				value := avlDataStart[dataFieldPosition : dataFieldPosition+16]
				ioDict[fmt.Sprintf("%d", key)] = sortingHat(key, value)
				dataFieldPosition += len(value)
				fmt.Printf("avl_ID: %d : %v\n", key, ioDict[fmt.Sprintf("%d", key)])
				i++
			}
		}
		if strings.ToUpper(codecType) == "8E" {
			byteXIONumber := avlDataStart[dataFieldPosition : dataFieldPosition+4]
			byteXIONumberParsed, _ := strconv.ParseInt(byteXIONumber, 16, 64)
			fmt.Printf("X byte io count = %d\n", byteXIONumberParsed)
			dataFieldPosition += len(byteXIONumber)
			if byteXIONumberParsed > 0 {
				i := int64(1)
				for i <= byteXIONumberParsed {
					key := avlDataStart[dataFieldPosition : dataFieldPosition+4]
					dataFieldPosition += len(key)
					valueLength := avlDataStart[dataFieldPosition : dataFieldPosition+4]
					dataFieldPosition += 4
					valueLengthInt, err := strconv.Atoi(valueLength)
					if err != nil {
						fmt.Println("Error converting FATTAL")
					}

					value := avlDataStart[dataFieldPosition : dataFieldPosition+2*valueLengthInt]
					ioDict[fmt.Sprintf("%d", key)] = sortingHat(key, value)
					dataFieldPosition += len(value)
					fmt.Printf("avl_ID: %d : %v\n", key, ioDict[fmt.Sprintf("%d", key)])
					i++
				}
			}
		}
		recordNumber++
		jsonPrinter(ioDict, deviceIMEI)
	}
	if props == "SERVER" {
		totalRecordsParsed, _ := strconv.ParseInt(avlDataStart[dataFieldPosition:dataFieldPosition+2], 16, 64)
		fmt.Printf("\ntotal parsed records = %d\n\n", totalRecordsParsed)
	} else {
		totalRecordsParsed, _ := strconv.ParseInt(avlDataStart[dataFieldPosition:dataFieldPosition+2], 16, 64)
		fmt.Printf("\ntotal parsed records = %d\n\n", totalRecordsParsed)
		inputTrigger()
	}
	return recordNumber
}

func jsonPrinter(ioDict map[string]interface{}, deviceIMEI string) {
	jsonData, _ := json.MarshalIndent(ioDict, "", "    ")
	dataPath := "./data/" + deviceIMEI
	jsonFile := deviceIMEI + "_data.json"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.MkdirAll(dataPath, os.ModePerm)
	}
	if _, err := os.Stat(dataPath + "/" + jsonFile); os.IsNotExist(err) {
		file, _ := os.Create(dataPath + "/" + jsonFile)
		file.Write(jsonData)
		file.Close()
	} else {
		file, _ := os.OpenFile(dataPath+"/"+jsonFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		file.Write(jsonData)
		file.Close()
	}
}

func jsonPrinterRawData(ioDictRaw map[string]interface{}, deviceIMEI string) {
	jsonData, _ := json.MarshalIndent(ioDictRaw, "", "    ")
	dataPath := "./data/" + deviceIMEI
	jsonFile := deviceIMEI + "_RAWdata.json"
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.MkdirAll(dataPath, os.ModePerm)
	}
	if _, err := os.Stat(dataPath + "/" + jsonFile); os.IsNotExist(err) {
		file, _ := os.Create(dataPath + "/" + jsonFile)
		file.Write(jsonData)
		file.Close()
	} else {
		file, _ := os.OpenFile(dataPath+"/"+jsonFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		file.Write(jsonData)
		file.Close()
	}
}

func timeStamper() string {
	currentServerTime := time.Now()
	serverTimeStamp := currentServerTime.Format("15:04:05 02-01-2006")
	return serverTimeStamp
}

func timeStamperForJSON() string {
	currentServerTime := time.Now()
	timestampUTC := time.Now().UTC()
	serverTimeStamp := fmt.Sprintf("%s (local) / %s (utc)", currentServerTime.Format("15:04:05 02-01-2006"), timestampUTC.Format("15:04:05 02-01-2006"))
	return serverTimeStamp
}

func deviceTimeStamper(timestamp string) string {
	timestampMS, _ := strconv.ParseInt(timestamp, 16, 64)
	timestampUTC := time.Unix(timestampMS/1000, 0).UTC()
	utcOffset := time.Unix(timestampMS/1000, 0).Sub(time.Unix(timestampMS/1000, 0).UTC())
	timestampLocal := timestampUTC.Add(utcOffset)
	formattedTimestampLocal := timestampLocal.Format("15:04:05 02-01-2006")
	formattedTimestampUTC := timestampUTC.Format("15:04:05 02-01-2006")
	formattedTimestamp := fmt.Sprintf("%s (local) / %s (utc)", formattedTimestampLocal, formattedTimestampUTC)
	return formattedTimestamp
}

func recordDelayCounter(timestamp string) string {
	timestampMS, _ := strconv.ParseInt(timestamp, 16, 64)
	currentServerTime := time.Now().Unix()
	return fmt.Sprintf("%d seconds", currentServerTime-timestampMS/1000)
}

func parseDataInteger(data string) interface{} {
	parsedData, _ := strconv.ParseInt(data, 16, 64)
	return parsedData
}

func intMultiply01(data string) interface{} {
	parsedData, _ := strconv.ParseInt(data, 16, 64)
	return float64(parsedData) * 0.1
}

func intMultiply001(data string) interface{} {
	parsedData, _ := strconv.ParseInt(data, 16, 64)
	return float64(parsedData) * 0.01
}

func intMultiply0001(data string) interface{} {
	parsedData, _ := strconv.ParseInt(data, 16, 64)
	return float64(parsedData) * 0.001
}

func signedNoMultiply(data string) interface{} {
	parsedData, _ := strconv.ParseInt(data, 16, 64)
	return int32(parsedData)
}

func sortingHat(key string, value string) interface{} {
	parseFunctionsDictionary := map[string]func(string) interface{}{
		"F0":  parseDataInteger,
		"EF":  parseDataInteger,
		"50":  parseDataInteger,
		"15":  parseDataInteger,
		"C8":  parseDataInteger,
		"45":  parseDataInteger,
		"B5":  intMultiply01,
		"B6":  intMultiply01,
		"42":  intMultiply0001,
		"18":  parseDataInteger,
		"CD":  parseDataInteger,
		"CE":  parseDataInteger,
		"43":  intMultiply0001,
		"44":  intMultiply0001,
		"F1":  parseDataInteger,
		"12B": parseDataInteger,
		"10":  parseDataInteger,
		"1":   parseDataInteger,
		"9":   parseDataInteger,
		"B3":  parseDataInteger,
		"C":   intMultiply0001,
		"D":   intMultiply001,
		"11":  signedNoMultiply,
		"12":  signedNoMultiply,
		"13":  signedNoMultiply,
		"B":   parseDataInteger,
		"A":   parseDataInteger,
		"2":   parseDataInteger,
		"3":   parseDataInteger,
		"6":   intMultiply0001,
		"B4":  parseDataInteger,
	}
	if parseFunction, ok := parseFunctionsDictionary[key]; ok {
		return parseFunction(value)
	} else {
		return fmt.Sprintf("0x%s", value)
	}
}

func fileAccessTest() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v\n", r)
			inputTrigger()
		}
	}()
	testDict := make(map[string]interface{})
	testDict["_Writing_Test_"] = "Writing_Test"
	testDict["Script_Started"] = timeStamperForJSON()
	jsonPrinter(testDict, "file_Write_Test")
	fmt.Println("---")
	inputTrigger()
}

// package main

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"net"
// 	"strconv"
// 	"time"

// 	"github.com/Projects/Zanjeer/config"
// 	"github.com/Projects/Zanjeer/helpers"
// 	"github.com/Projects/Zanjeer/pkg/db"
// 	"github.com/Projects/Zanjeer/pkg/logger"
// 	"github.com/Projects/Zanjeer/storage"
// )

// type DeviceData struct {
// 	IMEI         string `json:"imei"`
// 	TimeDate     string `json:"date"`
// 	Lat          string `json:"lat"`
// 	Lng          string `json:"lng"`
// 	NumberOfData string
// 	Altitude     string
// 	Angle        string
// }

// const port = "1234"

// func main() {

// 	cfg := config.Load()

// 	logger := logger.New(cfg.LogLevel)

// 	db, err := db.New(cfg)
// 	if err != nil {
// 		logger.Error("Error while connecting to database", err)
// 	} else {
// 		logger.Info("Successfully connected to database")
// 	}

// 	// Listen for incoming connections
// 	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	defer listener.Close()
// 	fmt.Println("Server is listening on port " + port)

// 	for {
// 		// Accept incoming connections
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			continue
// 		}

// 		// Handle client connection in a goroutine
// 		go handleClient(conn, db, logger, cfg)
// 	}
// }

// func handleClient(conn net.Conn, db *db.Postgres, log *logger.Logger, cfg config.Config) {

// 	pg := storage.New(db, log, cfg)

// 	defer conn.Close()

// 	var messageTrans = map[int]func(step *int, imei *string, msg string, conn net.Conn){}

// 	messageTrans[1] = takeImei

// 	// Create a buffer to read data into
// 	buffer := make([]byte, 1024)
// 	var (
// 		imeiTaken bool    = true
// 		step      *int    = new(int)
// 		imei      *string = new(string)
// 	)
// 	*step = 1

// 	for {
// 		// Read data from the client
// 		size, err := conn.Read(buffer)
// 		if err != nil {
// 			fmt.Println("Error conn.Read :", err)
// 			b := []byte{0} // 0x00 if we decline the message
// 			conn.Write(b)
// 			break
// 		}

// 		if imeiTaken {

// 			message := hex.EncodeToString(buffer[:size])

// 			if helpers.Imei(buffer) {
// 				*imei = string(buffer[2:17])
// 			}

// 			fmt.Println("Message :", message)

// 			switch *step {
// 			case 1:
// 				messageTrans[*step](step, imei, message, conn)
// 				b := []byte{0,0,0,1} // 0x00 if we decline the message
// 				conn.Write(b)
// 			case 2:

// 				data, err := helpers.ReadMessage(message, *imei)
// 				if err == nil {
// 					if pg.Postgres().SetLocation(data) != nil {
// 						fmt.Println("ERROR while Setting location!", err)
// 					}
// 				} else {
// 					fmt.Println("ERROR while Reading Message :", err)
// 				}

// 				d, _ := json.MarshalIndent(data, "", " ")
// 				fmt.Println(string(d))
// 				// if err == nil {
// 				fmt.Println("Size sent : ", data.Size)
// 				conn.Write([]byte{0, 0, 0, (data.Size)})
// 				// }
// 			}
// 		} else {
// 			b := []byte{0} // 0x00 if we decline the message
// 			conn.Write(b)
// 			break
// 		}

// 	}
// }

// func takeImei(step *int, imei *string, msg string, conn net.Conn) {
// 	firstReply := []byte{1}
// 	*step = 2
// 	conn.Write(firstReply)
// }

// func Decoder(enCode string) (DeviceData, error) {

// 	if len(enCode) < 58 {
// 		return DeviceData{}, fmt.Errorf("Minimum packet size is 45 Bytes, got %v", len(enCode))
// 	}

// 	// zeroBytes := (enCode[0:8])
// 	// dataFieldLength := (enCode[8:16])
// 	codecId := (enCode[16:18])
// 	numberOfData := (enCode[18:20])
// 	timestamp := (enCode[20:36])
// 	// priority := (enCode[36:38])
// 	longitude := (enCode[38:46])
// 	latitude := (enCode[46:54])

// 	altitude := (enCode[54:56])
// 	angle := (enCode[56:58])

// 	if codecId != "08" && codecId != "8E" {
// 		return DeviceData{}, fmt.Errorf("Invalid Codec ID, want 0x08 or 0x8E, get %v", codecId)
// 	}

// 	return DeviceData{
// 		TimeDate:     timestamp,
// 		Lat:          latitude,
// 		Lng:          longitude,
// 		NumberOfData: numberOfData,
// 		Altitude:     altitude,
// 		Angle:        angle,
// 	}, nil

// }

// func DecToLocation(dec int64) string {
// 	if dec <= 0 {
// 		return ""
// 	}
// 	location := fmt.Sprintf("%d", dec)

// 	return location[0:2] + "." + location[2:]
// }
// func HexToDec(hex string) int64 {
// 	num, err := strconv.ParseInt(hex, 16, 64)
// 	if err != nil {
// 		return 0
// 	}
// 	return num
// }

// func hexToTime(hexTimestamp string) (time.Time, error) {
// 	// Decode the hex string into a byte slice
// 	hexBytes, err := hex.DecodeString(hexTimestamp)
// 	if err != nil {
// 		return time.Time{}, err
// 	}

// 	// Parse the byte slice as a Unix timestamp
// 	// Assuming that the timestamp is a 64-bit integer (8 bytes)
// 	unixTimestamp := int64(0)
// 	for _, b := range hexBytes {
// 		unixTimestamp = (unixTimestamp << 8) | int64(b)
// 	}

// 	// Convert Unix timestamp to time.Time
// 	timestamp := time.Unix(unixTimestamp, 0)

// 	return timestamp, nil
// }
