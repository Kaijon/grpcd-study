package config

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

// LoadConfig reads JSON from filePath and unmarshals into the package AppConfig.
func LoadConfig(filePath string) error {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("error reading config file: %w", err)
    }
    if err := json.Unmarshal(data, &AppConfig); err != nil {
        return fmt.Errorf("error parsing config file: %w", err)
    }
    return nil
}

// UpdateConfig runs updateFunc while holding a lock and then persists the config.
func UpdateConfig(updateFunc func(*Config)) error {
    configLock.Lock()
    defer configLock.Unlock()
    updateFunc(&AppConfig)
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
    // no-op for now (device-specific path). Keep for compatibility.
    return nil
}

func LoadConfigDefault() {
    filename := "config.json"
    if runtime.GOARCH == "amd64" {
        filename = "fake.data"
    }

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        AppConfig = Config{ // Use the package AppConfig variable
            LEDs: map[string]LEDConfig{
                "0": {StatusLed: "RED", RecLedOn: true},
                "1": {StatusLed: "GREEN", RecLedOn: false},
            },
            Network: NetworkConfig{
                IPv4: "192.168.5.24",
                IPv6: "3001:0db8:85a3:0000:0000:8a2e:0370:7334",
            },
            System: SystemConfig{
                FWVersion:  "0.0.0.1",
                Time:       "2024-08-22T15:00:02.574UTC+08:00",
                SerialNo:   "WMXNF222G2X0",
                SKUName:    "CA-NF22G2",
                DeviceName: "WMXNF222G2X0",
                MAC:        "02:00:01:00:01:7B",
            },
            Videos: map[string]VideoConfig{
                "0": {Resolution: "2560x1080", StreamFormat: "h265", BitRate: 12, Type: "vbr", Fps: 30,
                    SubResolution: "1280x720", SubStreamFormat: "h264", SubBitRate: 4, SubType: "vbr", SubFps: 30, MirrorAction: "normal"},

                "1": {Resolution: "2560x1080", StreamFormat: "h265", BitRate: 12, Type: "vbr", Fps: 30,
                    SubResolution: "1280x720", SubStreamFormat: "h264", SubBitRate: 4, SubType: "vbr", SubFps: 30, MirrorAction: "normal"},
            },
            Watermarks: map[string]WatermarkConfig{
                "0": {
                    Username:         "DefaultUser",
                    OptionUserName:   true,
                    OptionDeviceName: true,
                    OptionGPS:        true,
                    OptionTime:       true,
                    OptionLogo:       true,
                },
                "1": {
                    Username:         "DefaultUser",
                    OptionUserName:   true,
                    OptionDeviceName: true,
                    OptionGPS:        true,
                    OptionTime:       true,
                    OptionLogo:       true,
                },
            },
            DayNightMode: SenserConfig{
                Mode: "day",
                Lux:  128,
            },
        }

        if err := SaveAppConfigDefault(); err != nil {
            log.Fatalf("Failed to create default config file: %v", err)
        }
    } else {
        if err := LoadConfig(filename); err != nil {
            log.Printf("failed loading config %s: %v", filename, err)
        }
    }
    fmt.Println("Device Name:", AppConfig.System.DeviceName)
}

func Init() {
    // initialize defaults
    AppConfig = Config{ // Use the package AppConfig variable
        LEDs: map[string]LEDConfig{
            "0": {StatusLed: "RED", RecLedOn: true},
            "1": {StatusLed: "GREEN", RecLedOn: false},
        },
        Network: NetworkConfig{
            IPv4: "192.168.5.24",
            IPv6: "3001:0db8:85a3:0000:0000:8a2e:0370:7334",
        },
        System: SystemConfig{
            FWVersion:  "0.0.0.1",
            Time:       "2024-08-22T15:00:02.574UTC+08:00",
            SerialNo:   "WMXNF222G2X0",
            SKUName:    "CA-NF22G2",
            DeviceName: "WMXNF222G2X0",
            MAC:        "02:00:01:00:01:7B",
        },
        Videos: map[string]VideoConfig{
            "0": {Resolution: "2560x1080", StreamFormat: "h265", BitRate: 12, Type: "vbr", Fps: 30,
                SubResolution: "1280x720", SubStreamFormat: "h264", SubBitRate: 4, SubType: "vbr", SubFps: 30, MirrorAction: "normal"},

            "1": {Resolution: "2560x1080", StreamFormat: "h265", BitRate: 12, Type: "vbr", Fps: 30,
                SubResolution: "1280x720", SubStreamFormat: "h264", SubBitRate: 4, SubType: "vbr", SubFps: 30, MirrorAction: "normal"},
        },
        Watermarks: map[string]WatermarkConfig{
            "0": {
                Username:         "DefaultUser",
                OptionUserName:   true,
                OptionDeviceName: true,
                OptionGPS:        true,
                OptionTime:       true,
                OptionLogo:       true,
            },
            "1": {
                Username:         "DefaultUser",
                OptionUserName:   true,
                OptionDeviceName: true,
                OptionGPS:        true,
                OptionTime:       true,
                OptionLogo:       true,
            },
        },
        DayNightMode: SenserConfig{
            Mode: "day",
            Lux:  128,
        },
    }
}
