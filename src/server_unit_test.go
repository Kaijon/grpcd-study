package main

import (
	"context"
	"testing"

	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func drainMqtt() {
	for {
		select {
		case <-MqttPublishChannel:
		default:
			return
		}
	}
}

func TestLEDServer_SetGet(t *testing.T) {
	cfg.Init()
	led := &LEDServer{cfg: &cfg.AppConfig}
	// set
	_, err := led.SetLEDs(context.Background(), &pb.SetLEDsRequest{Channel: 0, StatusLedColor: pb.Color_RED, RecLedOn: true})
	if err != nil {
		t.Fatalf("SetLEDs failed: %v", err)
	}
	// get
	resp, err := led.GetLEDs(context.Background(), &pb.GetLEDsRequest{Channel: 0})
	if err != nil {
		t.Fatalf("GetLEDs failed: %v", err)
	}
	if resp.RecLedOn != true {
		t.Fatalf("expected RecLedOn true, got %v", resp.RecLedOn)
	}
	drainMqtt()
}

func TestNetwork_UpdateGet(t *testing.T) {
	cfg.Init()
	netSrv := &NetworkInfoServer{cfg: &cfg.AppConfig}
	// update
	_, err := netSrv.UpdateIPv4(context.Background(), &pb.UpdateIPv4Request{IPv4: "10.0.0.5"})
	if err != nil {
		t.Fatalf("UpdateIPv4 failed: %v", err)
	}
	// get
	resp, err := netSrv.GetIPv4(context.Background(), &pb.GetIPv4Request{})
	if err != nil {
		t.Fatalf("GetIPv4 failed: %v", err)
	}
	if resp.IPv4 != "10.0.0.5" {
		t.Fatalf("expected IPv4 to be updated, got %s", resp.IPv4)
	}
	drainMqtt()
}

func TestVideoServer_SetGet(t *testing.T) {
	cfg.Init()
	videoSrv := &VideoInfoServer{cfg: &cfg.AppConfig}
	req := &pb.SetVideoSettingsRequest{
		Channel:         1,
		Resolution:      "1920x1080",
		StreamFormat:    "h265",
		BitRate:         4096,
		Type:            "vbr",
		Fps:             30,
		SubResolution:   "1280x720",
		SubStreamFormat: "h264",
		SubBitRate:      1024,
		SubType:         "vbr",
		SubFps:          15,
		MirrorAction:    "normal",
	}
	if _, err := videoSrv.SetVideoSettings(context.Background(), req); err != nil {
		t.Fatalf("SetVideoSettings failed: %v", err)
	}
	resp, err := videoSrv.GetVideoSettings(context.Background(), &pb.GetVideoSettingsRequest{Channel: 1})
	if err != nil {
		t.Fatalf("GetVideoSettings failed: %v", err)
	}
	if resp.Resolution != req.Resolution || resp.StreamFormat != req.StreamFormat || resp.BitRate != req.BitRate {
		t.Fatalf("unexpected video response: %+v", resp)
	}
	drainMqtt()
}

func TestWatermarkServer_SetGet(t *testing.T) {
	cfg.Init()
	wmSrv := &WatermarkInfoServer{cfg: &cfg.AppConfig}
	req := &pb.SetAllWatermarkInfoRequest{
		Channel:          0,
		Username:         "tester",
		OptionUserName:   true,
		OptionDeviceName: true,
		OptionGPS:        false,
		OptionTime:       true,
		OptionLogo:       wrapperspb.Bool(true),
		OptionExposure:   wrapperspb.Bool(false),
	}
	if _, err := wmSrv.SetAllWatermarkInfo(context.Background(), req); err != nil {
		t.Fatalf("SetAllWatermarkInfo failed: %v", err)
	}
	resp, err := wmSrv.GetAllWatermarkInfo(context.Background(), &pb.GetAllWatermarkInfoRequest{Channel: 0})
	if err != nil {
		t.Fatalf("GetAllWatermarkInfo failed: %v", err)
	}
	if resp.Username != req.Username || resp.OptionDeviceName != req.OptionDeviceName {
		t.Fatalf("watermark response mismatch: %+v", resp)
	}
	drainMqtt()
}

func TestDeviceInfoServer_SetGetAlpr(t *testing.T) {
	cfg.Init()
	devSrv := &DeviceInfoServer{cfg: &cfg.AppConfig}
	if _, err := devSrv.SetAlprStatus(context.Background(), &pb.SetAlprRequest{IsEnabled: true}); err != nil {
		t.Fatalf("SetAlprStatus failed: %v", err)
	}
	resp, err := devSrv.GetAlprStatus(context.Background(), &pb.GetAlprRequest{})
	if err != nil {
		t.Fatalf("GetAlprStatus failed: %v", err)
	}
	if !resp.IsEnabled {
		t.Fatalf("expected ALPR enabled")
	}
	drainMqtt()
}
