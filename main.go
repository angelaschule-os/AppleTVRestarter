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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	BASE_URL = os.Getenv("BASE_URL")
	NETWORK_ID := os.Getenv("NETWORK_ID")
	KEY := os.Getenv("KEY")

	auth := NETWORK_ID + ":" + KEY
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	AUTHORIZATION = "Basic " + encodedAuth
}

func sendCommandToDevices(udids []string, command string) {
	for _, udid := range udids {
		url := fmt.Sprintf("%s/devices/%s/%s", BASE_URL, udid, command)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		req.Header.Add("Authorization", AUTHORIZATION)
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to %s device with UDID %s: %v", command, udid, err)
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error %s device with UDID %s. HTTP Status: %d. Response: %s", command, udid, resp.StatusCode, body)
		} else {
			log.Println("HTTP Response:", string(body))
		}
	}
}

func getUDIDs() []string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/devices?assettag=AppleTV", BASE_URL), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("Authorization", AUTHORIZATION)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch devices: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var deviceResponse DeviceResponse
	err = json.Unmarshal(bodyBytes, &deviceResponse)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	udids := make([]string, len(deviceResponse.Devices))
	for i, device := range deviceResponse.Devices {
		udids[i] = device.UDID
	}

	return udids
}

func main() {
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
		log.Printf("UDIDs: %v\n", udids)
		log.Println("Count:", len(udids))
	}

	if *refreshFlag {
		log.Println("Sending refresh command...")
		sendCommandToDevices(udids, "refresh")
	}

	if *restartFlag {
		log.Println("Sending restart command...")
		sendCommandToDevices(udids, "restart")
	}

	if *versionFlag {
		fmt.Println("AppleTVRestarter Commit Hash:", GitCommit)
		return
	}
}
