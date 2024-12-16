package main

type DeviceResponse struct {
	Code    int      `json:"code"`
	Count   int      `json:"count"`
	Devices []Device `json:"devices"`
}

type Device struct {
	UDID                string      `json:"UDID"`
	LocationID          int         `json:"locationId"`
	SerialNumber        string      `json:"serialNumber"`
	AssetTag            string      `json:"assetTag"`
	InTrash             bool        `json:"inTrash"`
	Class               string      `json:"class"`
	Model               Model       `json:"model"` // single object, not an array
	OS                  OS          `json:"os"`    // single object, not an array
	Name                string      `json:"name"`
	Owner               Owner       `json:"owner"` // single object, not an array
	IsManaged           bool        `json:"isManaged"`
	IsSupervised        bool        `json:"isSupervised"`
	EnrollType          string      `json:"enrollType"`
	DEPProfile          string      `json:"depProfile"`
	Groups              []string    `json:"groups"`
	BatteryLevel        float64     `json:"batteryLevel"`
	TotalCapacity       float64     `json:"totalCapacity"`
	AvailableCapacity   float64     `json:"availableCapacity"`
	ICloudBackupEnabled bool        `json:"iCloudBackupEnabled"`
	ICloudBackupLatest  int         `json:"iCloudBackupLatest"`
	ITunesStoreLoggedIn bool        `json:"iTunesStoreLoggedIn"`
	Region              Region      `json:"region"` // single object, not an array
	Notes               string      `json:"notes"`
	LastCheckin         string      `json:"lastCheckin"`
	Modified            string      `json:"modified"`
	NetworkInformation  NetworkInfo `json:"networkInformation"`
}

type Model struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

type OS struct {
	Prefix  string `json:"prefix"`
	Version string `json:"version"`
}

type Owner struct {
	ID         int    `json:"id"`
	LocationID int    `json:"locationId"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Name       string `json:"name"`
	VPP        []VPP  `json:"vpp"`
	Notes      string `json:"notes"`
	Modified   string `json:"modified"`
}

type VPP struct {
	Status string `json:"status"`
}

type Region struct {
	String      string `json:"string"`
	Coordinates string `json:"coordinates"`
}

type NetworkInfo struct {
	IPAddress              string `json:"IPAddress"`
	IsNetworkTethered      int    `json:"isNetworkTethered"`
	BluetoothMAC           string `json:"BluetoothMAC"`
	WiFiMAC                string `json:"WiFiMAC"`
	EthernetMACs           string `json:"EthernetMACs"`
	VoiceRoamingEnabled    int    `json:"VoiceRoamingEnabled"`
	DataRoamingEnabled     int    `json:"DataRoamingEnabled"`
	PersonalHotspotEnabled int    `json:"PersonalHotspotEnabled"`
}
