package myzone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Aircon struct {
	Aircons map[string]struct {
		Info struct {
			AaAutoFanModeEnabled        bool    `json:"aaAutoFanModeEnabled"`
			ActivationCodeStatus        string  `json:"activationCodeStatus"`
			AirconErrorCode             string  `json:"airconErrorCode"`
			CbFWRevMajor                int     `json:"cbFWRevMajor"`
			CbFWRevMinor                int     `json:"cbFWRevMinor"`
			CbType                      int     `json:"cbType"`
			ClimateControlModeEnabled   bool    `json:"climateControlModeEnabled"`
			ClimateControlModeIsRunning bool    `json:"climateControlModeIsRunning"`
			Constant1                   int     `json:"constant1"`
			Constant2                   int     `json:"constant2"`
			Constant3                   int     `json:"constant3"`
			CountDownToOff              int     `json:"countDownToOff"`
			CountDownToOn               int     `json:"countDownToOn"`
			DbFWRevMajor                int     `json:"dbFWRevMajor"`
			DbFWRevMinor                int     `json:"dbFWRevMinor"`
			Fan                         string  `json:"fan"`
			FilterCleanStatus           int     `json:"filterCleanStatus"`
			FreshAirStatus              string  `json:"freshAirStatus"`
			Mode                        string  `json:"mode"`
			MyAutoModeCurrentSetMode    string  `json:"myAutoModeCurrentSetMode"`
			MyAutoModeEnabled           bool    `json:"myAutoModeEnabled"`
			MyAutoModeIsRunning         bool    `json:"myAutoModeIsRunning"`
			MyZone                      int     `json:"myZone"`
			Name                        string  `json:"name"`
			NoOfConstants               int     `json:"noOfConstants"`
			NoOfZones                   int     `json:"noOfZones"`
			QuietNightModeIsRunning     bool    `json:"quietNightModeIsRunning"`
			RfFWRevMajor                int     `json:"rfFWRevMajor"`
			RfSysID                     int     `json:"rfSysID"`
			SetTemp                     float64 `json:"setTemp"`
			State                       string  `json:"state"`
			Uid                         string  `json:"uid"`
			UnitType                    int     `json:"unitType"`
		} `json:"info"`
		Zones map[string]struct {
			Error           int      `json:"error"`
			Followers       []string `json:"followers"`
			Following       int      `json:"following"`
			MaxDamper       int      `json:"maxDamper"`
			MeasuredTemp    float64  `json:"measuredTemp"`
			MinDamper       int      `json:"minDamper"`
			Motion          int      `json:"motion"`
			MotionConfig    int      `json:"motionConfig"`
			Name            string   `json:"name"`
			Number          int      `json:"number"`
			Rssi            int      `json:"rssi"`
			SensorUid       string   `json:"SensorUid"`
			SetTemp         float64  `json:"setTemp"`
			State           string   `json:"state"`
			TempSensorClash bool     `json:"tempSensorClash"`
			Type            int      `json:"type"`
			Value           int      `json:"value"`
		} `json:"zones"`
	} `json:"aircons"`
}

func Fetch(url string) (*Aircon, error) {
	// Create a new HTTP client with a timeout of 10 seconds
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send a GET request to the provided URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful response status code (200 OK)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %v", resp.Status)
	}

	// Initialize a variable to hold the JSON data
	var airconData Aircon

	// Decode the JSON response into the struct
	err = json.NewDecoder(resp.Body).Decode(&airconData)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	// Return the populated struct
	return &airconData, nil
}
