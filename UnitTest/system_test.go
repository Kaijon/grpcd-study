package main

import (
	pb "grpcd/canf22g2/grpc"
	"context"
	"reflect"
	"testing"
    "os"
    "log"
    "time"
    "strings"
)

func TestDeviceInfoServer_SetTime(t *testing.T){
	LoadConfigDefault()
    server := &DeviceInfoServer{}
    //Set mock TZ env
    os.Setenv("TZ", "UTC-08:00")
    currentTimeStr, err := getCurrentTimeStr()
    if err != nil {
		Log.Debugf("Error formatting current time: %s", err)
		//return nil, err
	}
	Log.Infof("Get current time : %s", currentTimeStr)
    resp, err := server.SetTime(context.Background(), &pb.SetTimeRequest{Time: currentTimeStr})
    if err != nil && resp.Message != "System time configuration updated successfully"{
		t.Errorf("SetTime fail, return: %v", err)
	}

	if AppConfig.System.Time != currentTimeStr {
        t.Errorf("Expcted Time update to: %s, but get %s", currentTimeStr, AppConfig.System.Time)
	}
}

func TestDeviceInfoServer_GetAllSystemInfo(t *testing.T) {
	server := &DeviceInfoServer{}

    //Set mock TZ env
    os.Setenv("TZ", "UTC-08:00")
    currentTimeStr, err := getCurrentTimeStr()
    if err != nil {
		Log.Debugf("Error formatting current time: %s", err)
		//return nil, err
	}
	Log.Infof("Get current time : %s", currentTimeStr)

	resp, err := server.GetAllSystemInfo(context.Background(), &pb.GetAllSystemInfoRequest{})
	if err != nil {
		t.Fatalf("GetAllSystemInfo failed: %v", err)
	}

	expected := &pb.GetAllSystemInfoResponse{
		FWVersion:  AppConfig.System.FWVersion,
		Time:       currentTimeStr,
		SerialNo:   AppConfig.System.SerialNo,
		SKUName:    AppConfig.System.SKUName,
		DeviceName: AppConfig.System.DeviceName,
        MAC:        AppConfig.System.MAC,
	}

	if !reflect.DeepEqual(resp, expected) {
        // add check time diff is within 5ms
        // Replace `UTC+00:00` with `Z` to make it a valid ISO 8601 format
        respTimeStr := strings.Replace(resp.Time, "UTC+08:00", "Z", 1)
        log.Printf("respTimeStr: %s", respTimeStr)

        // Parse resp.Time and currentTimeStr into time.Time
        layout := "2006-01-02T15:04:05.000Z07:00"

        respTime, err := time.Parse(layout, respTimeStr)
        if err != nil {
            log.Fatalf("Failed to parse resp.Time: %v", err)
        }
        log.Printf("respTime: %v", respTime)

        currentTimeStrReplace := strings.Replace(currentTimeStr, "UTC+08:00", "Z", 1)
        currentTime, err := time.Parse(layout, currentTimeStrReplace)
        if err != nil {
            log.Fatalf("Failed to parse currentTimeStr: %v", err)
        }
        log.Printf("currentTime: %v", currentTime)

        // Calculate the time difference
        timeDiff := respTime.Sub(currentTime)
        log.Printf("time diff: %v", timeDiff)

        // if time diff is >= 5ms, fail test
        if timeDiff >= 5 * time.Millisecond {
		    t.Errorf("Expected response %v, got %v", expected, resp)
        }
	}
}


