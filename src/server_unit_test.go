package main

import (
	"context"
	cfg "grpcd/config"
	pb "grpcd/canf22g2/grpc"
	"testing"
)

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
}
