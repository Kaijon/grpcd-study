package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
)

var (
	configLock sync.Mutex
	AppConfig  Config
)

var filename = "/mnt/getac/config/grpcd.cfg"

type Config struct {
	LEDs         map[string]LEDConfig       `json:"LEDs"`
	Network      NetworkConfig              `json:"network"`
	System       SystemConfig               `json:"system"`
	Videos       map[string]VideoConfig     `json:"video"`
	Watermarks   map[string]WatermarkConfig `json:"watermark"`
	DayNightMode SenserConfig               `json:"dayNightMode"`
}

type SenserConfig struct {
	Mode string `json:"Mode"`
	Lux  int    `json:"Lux"`
}

type LEDConfig struct {
	StatusLed string `json:"StatusLed"`
	RecLedOn  bool   `json:"RecLedOn"`
}

type NetworkConfig struct {
	IPv4 string `json:"IPv4"`
	IPv6 string `json:"IPv6"`
}

type SystemConfig struct {
	FWVersion   string `json:"FWVersion"`
	Time        string `json:"Time"`
	SerialNo    string `json:"serialNo"`
	SKUName     string `json:"SKUName"`
	DeviceName  string `json:"deviceName"`
	MAC         string `json:"MAC"`
	AlprEnabled bool   `json:"Enable"`
}

type VideoConfig struct {
	Resolution      string `json:"Resolution"`
	StreamFormat    string `json:"StreamFormat"`
	BitRate         uint32 `json:"BitRate"`
	Type            string `json:"Type"`
	Fps             uint32 `json:"Fps"`
	SubResolution   string `json:"SubResolution"`
	SubStreamFormat string `json:"SubStreamFormat"`
	SubBitRate      uint32 `json:"SubBitRate"`
	SubType         string `json:"SubType"`
	SubFps          uint32 `json:"SubFps"`
	MirrorAction    string `json:"MirrorAction"`
}

type WatermarkConfig struct {
	Username         string `json:"Username"`
	OptionUserName   bool   `json:"OptionUserName"`
	OptionDeviceName bool   `json:"OptionDeviceName"`
	OptionGPS        bool   `json:"OptionGPS"`
	OptionTime       bool   `json:"OptionTime"`
	OptionLogo       bool   `json:"OptionLogo"`
	OptionExposure   bool   `json:"OptionExposure"`
}

func LoadConfig(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := json.Unmarshal(data, &AppConfig); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}

func UpdateConfig(updateFunc func()) error {
	configLock.Lock()
	defer configLock.Unlock()
	updateFunc()
	return SaveAppConfigDefault()
}

func SaveAppConfig(config Config, filename string) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func SaveAppConfigDefault() error {
	package main

	import (
		cfg "grpcd/config"
	)

	// Keep a reference to the shared config package to avoid duplicate
	// declarations in this tree. Other files should import cfg and use
	// cfg.AppConfig directly. This file exists as a small shim only.
	var _ = cfg.AppConfig
	if runtime.GOARCH == "amd64" {
