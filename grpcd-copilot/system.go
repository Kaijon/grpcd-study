package main

import (
	"context"
	"fmt"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeviceInfoServer struct {
	pb.UnimplementedDeviceInfoServiceServer
	cfg *cfg.Config
}

func ConvertTZ(tz string) string {
    if !strings.HasPrefix(tz, "UTC") {
        return tz
    }

    // Remove "UTC"
    offsetStr := tz[3:]

    // Remove "DST" or "DST,j0,j366"
    if idx := strings.Index(offsetStr, "DST"); idx != -1 {
        offsetStr = offsetStr[:idx]
    }

    sign := offsetStr[0]
    offsetParts := strings.Split(offsetStr[1:], ":")
    if len(offsetParts) != 2 {
        return tz
    }

    hours := offsetParts[0]
    minutes := offsetParts[1]

    newSign := '+'
    if sign == '+' {
        newSign = '-'
    } else if sign == '-' {
        newSign = '+'
    }

    result := fmt.Sprintf("%c%s:%s", newSign, hours, minutes)

    return result
}

func parseUTCOffset(tz string) (int, error) {
    offsetStr := tz
    sign := offsetStr[0]
    offsetParts := strings.Split(offsetStr[1:], ":")
    if len(offsetParts) != 2 {
        return 0, fmt.Errorf("invalid TZ offset format")
    }

    hours, err := strconv.Atoi(offsetParts[0])
    if err != nil {
        return 0, err
    }

    minutes, err := strconv.Atoi(offsetParts[1])
    if err != nil {
        return 0, err
    }

    totalMinutes := hours*60 + minutes
    if sign == '-' {
        totalMinutes = -totalMinutes
    }

    return totalMinutes * 60, nil
}

func getCurrentTimeStr() (string, error) {
	tz := os.Getenv("TZ")
	Log.Infof("tz via Getenv: %s", tz)
	if tz == "" {
		// Using a more specific error message is better than a generic one.
		return "", fmt.Errorf("environment variable TZ is empty; cannot determine timezone")
	}	
	
	hasDST := strings.Contains(tz, "DST")

	// change sign and remove DST.
	// UTC-08:00 to UTC+08:00, UTC-08:00DST to UTC+08:00
	convertedTZ := ConvertTZ(tz)
	Log.Infof("Converted timezone: %s", convertedTZ)

	offsetSeconds, err := parseUTCOffset(convertedTZ)
	if err != nil {
		Log.Debugf("Error parsing timezone offset: %s", err)
		return "", err
	}

	loc := time.FixedZone(convertedTZ, offsetSeconds)
	currentTime := time.Now().In(loc)
	if hasDST {
		currentTime = currentTime.Add(time.Hour)
	}
	currentTimeStr := currentTime.Format("2006-01-02T15:04:05.000") + convertedTZ + func() string {
		if hasDST {
			return "DST"
		}
		return ""
	}()
	return currentTimeStr, nil
}

func (s *DeviceInfoServer) GetAllSystemInfo(ctx context.Context, in *pb.GetAllSystemInfoRequest) (*pb.GetAllSystemInfoResponse, error) {
	Log.Info(">>Run")

	currentTimeStr, err := getCurrentTimeStr()
	if err != nil {
		Log.Debugf("Error formatting current time: %s", err)
		return nil, err
	}
	Log.Infof("Get current time : %s", currentTimeStr)

	return &pb.GetAllSystemInfoResponse{
		FWVersion:  s.cfg.System.FWVersion,
		Time:       currentTimeStr ,
		SerialNo:   s.cfg.System.SerialNo,
		SKUName:    s.cfg.System.SKUName,
		DeviceName: s.cfg.System.DeviceName,
		MAC:        s.cfg.System.MAC,
	}, nil
}

func (s *DeviceInfoServer) SetTime(ctx context.Context, in *pb.SetTimeRequest) (*pb.SetTimeResponse, error) {
	s.cfg.System.Time = in.Time
	strTmp := fmt.Sprintf("{\"time\":\"%s\"}", in.Time)
	msg := MqttMessage{
		Topic:   "config/system/time",
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message to channel: %v", msg)
	}
	return &pb.SetTimeResponse{Message: "System time configuration updated successfully"}, nil
}

func (s *DeviceInfoServer) RunCmd(ctx context.Context, in *pb.RunCmdRequest) (*pb.RunCmdResponse, error) {
	strTmp := fmt.Sprintf("{\"value\":\"%s\"}", in.Cmd)
	msg := MqttMessage{
		Topic:   "config/system/command",
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message: %v", msg)
	}
	return &pb.RunCmdResponse{Message: "System time configuration updated successfully"}, nil
}

func (s *DeviceInfoServer) SetAlprStatus(ctx context.Context, in *pb.SetAlprRequest) (*pb.SetAlprResponse, error) {
	s.cfg.System.AlprEnabled = in.IsEnabled
	strTmp := fmt.Sprintf("{\"Enable\":\"%v\"}", in.IsEnabled)
	msg := MqttMessage{
		Topic:   "config/alpr",
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message: %v", msg)
	}
	return &pb.SetAlprResponse{IsEnabled: s.cfg.System.AlprEnabled}, nil
}

func (s *DeviceInfoServer) GetAlprStatus(ctx context.Context, in *pb.GetAlprRequest) (*pb.GetAlprResponse, error) {
	Log.Infof("ALPR enable=%v", s.cfg.System.AlprEnabled)
	return &pb.GetAlprResponse{IsEnabled: s.cfg.System.AlprEnabled}, nil
}
