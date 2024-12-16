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
	AUTHORIZATION string
	client        = &http.Client{}
	GitCommit     string // This will hold the git commit hash
)

func init() {
	log.Println("Initializing environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	BASE_URL = os.Getenv("BASE_URL")
	NETWORK_ID := os.Getenv("NETWORK_ID")
	KEY := os.Getenv("KEY")

	if BASE_URL == "" || NETWORK_ID == "" || KEY == "" {
		log.Fatalf("Missing required environment variables (BASE_URL, NETWORK_ID, or KEY)")
	}

	log.Printf("BASE_URL: %s", BASE_URL)

	auth := NETWORK_ID + ":" + KEY
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	AUTHORIZATION = "Basic " + encodedAuth
	log.Println("Authorization header initialized.")
}

func sendCommandToDevices(udids []string, command string) {
	log.Printf("Sending '%s' command to %d devices.\n", command, len(udids))
	for _, udid := range udids {
		url := fmt.Sprintf("%s/devices/%s/%s", BASE_URL, udid, command)
		log.Printf("Preparing request for device UDID: %s with URL: %s", udid, url)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Printf("Error creating request for device UDID %s: %v", udid, err)
			continue
		}

		req.Header.Add("Authorization", AUTHORIZATION)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to send '%s' command to device with UDID %s: %v", command, udid, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error sending '%s' command to device with UDID %s. HTTP Status: %d. Response: %s", command, udid, resp.StatusCode, body)
		} else {
			log.Printf("Successfully sent '%s' command to device with UDID %s. Response: %s", command, udid, body)
		}
	}
}

func getUDIDs() []string {
	log.Println("Fetching UDIDs from API...")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices?assettag=AppleTV", BASE_URL), nil)
	if err != nil {
		log.Fatalf("Error creating request to fetch devices: %v", err)
	}

	req.Header.Add("Authorization", AUTHORIZATION)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch devices: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Response Status: %s", resp.Status)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	log.Printf("Raw Response Body: %s", bodyBytes)

	var deviceResponse DeviceResponse
	err = json.Unmarshal(bodyBytes, &deviceResponse)
	if err != nil {
		log.Fatalf("Error decoding JSON response: %v", err)
	}

	log.Printf("Found %d devices. Extracting UDIDs...", len(deviceResponse.Devices))

	udids := make([]string, len(deviceResponse.Devices))
	for i, device := range deviceResponse.Devices {
		udids[i] = device.UDID
		log.Printf("Extracted UDID: %s", device.UDID)
	}

	return udids
}

func main() {
	log.Println("Starting application...")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s (version %s):\n", os.Args[0], GitCommit)
		flag.PrintDefaults()
	}
	refreshFlag := flag.Bool("refresh", false, "Send refresh command to devices")
	restartFlag := flag.Bool("restart", false, "Send restart command to devices")
	versionFlag := flag.Bool("version", false, "Print version information and exit")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	var udids []string
	if *refreshFlag || *restartFlag {
		udids = getUDIDs()
		log.Printf("UDIDs fetched: %v", udids)
		log.Printf("Number of UDIDs: %d", len(udids))
	}

	if *refreshFlag {
		log.Println("Sending refresh commands...")
		sendCommandToDevices(udids, "refresh")
	}

	if *restartFlag {
		log.Println("Sending restart commands...")
		sendCommandToDevices(udids, "restart")
	}

	if *versionFlag {
		fmt.Println("AppleTVRestarter Commit Hash:", GitCommit)
		return
	}
}
