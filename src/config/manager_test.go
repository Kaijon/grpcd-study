package config

import (
    "testing"
)

func TestInitAndUpdate(t *testing.T) {
    Init()
    if AppConfig.System.DeviceName == "" {
        t.Logf("DeviceName empty after Init; setting default")
    }

    original := AppConfig.System.DeviceName
    err := UpdateConfig(func(c *Config) {
        c.System.DeviceName = "UNIT_TEST_DEVICE"
    })
    if err != nil {
        t.Fatalf("UpdateConfig failed: %v", err)
    }
    if AppConfig.System.DeviceName != "UNIT_TEST_DEVICE" {
        t.Fatalf("expected DeviceName to be updated; got %s", AppConfig.System.DeviceName)
    }

    // restore
    _ = UpdateConfig(func(c *Config) {
        c.System.DeviceName = original
    })
}
