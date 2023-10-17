package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	BASE_URL      string
	NETWORK_ID    string
	KEY           string
	AUTHORIZATION string
)

func init() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	BASE_URL = os.Getenv("BASE_URL")
	NETWORK_ID = os.Getenv("NETWORKT_ID")
	KEY = os.Getenv("KEY")
}

func getAuthorizationHeader() string {
	// Combine Network ID and Key
	auth := NETWORK_ID + ":" + KEY

	// Base64 encode
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	// Return full header value
	return "Basic " + encodedAuth
}

func sendCommandToDevices(udids []string, client *http.Client, command string) {
	// Send a command to each device
	var i int
	for _, udid := range udids {
		url := fmt.Sprintf("%s/devices/%s/%s", BASE_URL, udid, command)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Add("Authorization", getAuthorizationHeader())
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to %s device with UDID %s: %v", command, udid, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Error %s device with UDID %s. HTTP Status: %d. Response: %s", command, udid, resp.StatusCode, body)
		}
		//log.Printf("%s device with UDID: %s", command, udid)
		body, _ := io.ReadAll(resp.Body)
		log.Println("HTTP Response:", string(body))
		resp.Body.Close()
		i++
	}
	log.Printf("Count: %d", i)

}

func getUDIDs(client *http.Client) []string {

	// Set headers
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/devices?assettag=AppleTV", BASE_URL), nil)
	req.Header.Add("Authorization", getAuthorizationHeader())
	req.Header.Add("Accept", "application/json")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch devices: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error fetching devices. HTTP Status: %d. Response: %s", resp.StatusCode, body)
	}

	// Read and print the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	//fmt.Println("HTTP Response:", string(bodyBytes))

	var deviceResponse DeviceResponse

	err = json.Unmarshal(bodyBytes, &deviceResponse)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	udids := make([]string, len(deviceResponse.Devices))
	for i, device := range deviceResponse.Devices {
		udids[i] = device.UDID
	}

	fmt.Println("UDIDs:", udids)
	fmt.Println("Count:", len(udids))
	return udids

}

func displayHelp() {
	fmt.Println("Usage:")
	fmt.Println("  AppleTVRestarter [--refresh] [--restart]")
	fmt.Println("Options:")
	fmt.Println("  --refresh\tSend refresh command to devices")
	fmt.Println("  --restart\tSend restart command to devices")
	fmt.Println("  --help\tDisplay this help message")
}

func main() {

	client := &http.Client{}
	var udids []string

	refreshFlag := flag.Bool("refresh", false, "Send refresh command to devices")
	restartFlag := flag.Bool("restart", false, "Send restart command to devices")
	helpFlag := flag.Bool("help", false, "Display help message")

	flag.Parse()

	if *helpFlag {
		displayHelp()
		return
	}

	if *refreshFlag || *restartFlag {
		udids = getUDIDs(client)

	}

	if *refreshFlag {
		fmt.Println("Sending refresh command...")
		sendCommandToDevices(udids, client, "refresh")
	}

	if *restartFlag {
		fmt.Println("Sending restart command...")
		sendCommandToDevices(udids, client, "restart")
	}

	// If neither flag is provided, show an error message or help
	if !*refreshFlag && !*restartFlag {
		fmt.Println("Error: No command specified.")
		displayHelp()

	}
}
