package main

import (
	"context"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"testing"
)

func TestUpdateIPv4(t *testing.T) {
	cfg.LoadConfigDefault()
	newIPv4 := "192.168.5.99"
	//AppConfig.Network.IPv4 = newIPv4
	server := &NetworkInfoServer{}
	_, err := server.UpdateIPv4(context.Background(), &pb.UpdateIPv4Request{IPv4: newIPv4})
	if err != nil {
		t.Errorf("UpdateIPv4 fail, return: %v", err)
	}
	cfg.ReadConfig(func(current cfg.Config) {
		if current.Network.IPv4 != newIPv4 {
			t.Errorf("Expcted IPv4 update to: %s, but get %s", newIPv4, current.Network.IPv4)
		}
	})
}

func TestGetIPv4(t *testing.T) {
	//LoadConfigDefault()
	server := &NetworkInfoServer{}
	response, err := server.GetIPv4(context.Background(), &pb.GetIPv4Request{})
	if err != nil {
		t.Errorf("GetIPv4 fail, return: %v", err)
	}
	cfg.ReadConfig(func(current cfg.Config) {
		if response.IPv4 != current.Network.IPv4 {
			t.Errorf("Expected IPv4: %s, but get %s", current.Network.IPv4, response.IPv4)
		}
	})
}

func TestGetIPv6(t *testing.T) {
	cfg.LoadConfigDefault()
	server := &NetworkInfoServer{}
	response, err := server.GetIPv6(context.Background(), &pb.GetIPv6Request{})
	if err != nil {
		t.Errorf("GetIPv6 fail, return: %v", err)
	}
	cfg.ReadConfig(func(current cfg.Config) {
		if response.IPv6 != current.Network.IPv6 {
			t.Errorf("Expcted IPv6: %s, but get %s", current.Network.IPv6, response.IPv6)
		}
	})
}

/*
func TestUpdateIPv6(t *testing.T) {
	LoadConfigDefault()
	newIPv6 := "fd00::1"
	AppConfig.Network.IPv6 = newIPv6
	server := &NetworkInfoServer{}
	_, err := server.UpdateIPv6(context.Background(), &pb.UpdateIPv6Request{IPv6: newIPv6})
	if err != nil {
		t.Errorf("UpdateIPv6 fail, return: %v", err)
	}
	if AppConfig.Network.IPv6 != newIPv6 {
        t.Errorf("Expected IPv6 update to: %s, but get %s", newIPv6, AppConfig.Network.IPv6)
	}
}
*/

func TestGetAllNetworkInfo(t *testing.T) {
	//LoadConfigDefault()
	server := &NetworkInfoServer{}
	response, err := server.GetAllNetworkInfo(context.Background(), &pb.GetAllNetworkInfoRequest{})
	if err != nil {
		t.Errorf("GetAllNetworkInfo fail, return: %v", err)
	}
	cfg.ReadConfig(func(current cfg.Config) {
		if response.IPv4 != current.Network.IPv4 || response.IPv6 != current.Network.IPv6 {
			t.Errorf("Expected results are not equal to response, IPv4 Expcted/Response: %s/%s, IPv6 Expcted/Exponse: %s/%s", current.Network.IPv4, response.IPv4, current.Network.IPv6, response.IPv6)
		}
	})
}
