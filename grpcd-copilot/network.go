package main

import (
	"context"
	"fmt"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"time"
)

type NetworkInfoServer struct {
	pb.UnimplementedNetworkInfoServiceServer
}

func (s *NetworkInfoServer) GetIPv4(ctx context.Context, in *pb.GetIPv4Request) (*pb.GetIPv4Response, error) {
	Log.Info(">>Run")
	var ipv4 string
	cfg.ReadConfig(func(current cfg.Config) {
		ipv4 = current.Network.IPv4
	})
	return &pb.GetIPv4Response{
		IPv4: ipv4,
	}, nil
}

func (s *NetworkInfoServer) GetIPv6(ctx context.Context, in *pb.GetIPv6Request) (*pb.GetIPv6Response, error) {
	Log.Info(">>Run")
	var ipv6 string
	cfg.ReadConfig(func(current cfg.Config) {
		ipv6 = current.Network.IPv6
	})
	return &pb.GetIPv6Response{
		IPv6: ipv6,
	}, nil
}

func (s *NetworkInfoServer) GetAllNetworkInfo(ctx context.Context, in *pb.GetAllNetworkInfoRequest) (*pb.GetAllNetworkInfoResponse, error) {
	Log.Info(">>Run")
	var (
		ipv4 string
		ipv6 string
	)
	cfg.ReadConfig(func(current cfg.Config) {
		ipv4 = current.Network.IPv4
		ipv6 = current.Network.IPv6
	})
	return &pb.GetAllNetworkInfoResponse{
		IPv4: ipv4,
		IPv6: ipv6,
	}, nil
}

func (s *NetworkInfoServer) UpdateIPv4(ctx context.Context, in *pb.UpdateIPv4Request) (*pb.UpdateIPv4Response, error) {
	if err := cfg.UpdateConfig(func(c *cfg.Config) {
		c.Network.IPv4 = in.IPv4
	}); err != nil {
		return nil, err
	}
	strTmp := fmt.Sprintf("{\"IPv4\":\"%v\"}", in.IPv4)
	msg := MqttMessage{
		Topic:   "config/network",
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message: %v", msg)
	}
	return &pb.UpdateIPv4Response{Message: "IPv4 set to " + in.IPv4}, nil
}

// as following proto, please give me a <func>Server for rpc call in golang.
// AppConfig.Network is a global variable that stores the current network configuration
// and UpdateConfig is a function that safely updates the network configurations and provide return message says what the user set.
