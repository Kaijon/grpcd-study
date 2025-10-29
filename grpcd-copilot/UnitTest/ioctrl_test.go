package main

import (
	"context"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"testing"
)

func TestSetLEDs(t *testing.T) {
	cfg.LoadConfigDefault()

	server := &LEDServer{}
	tests := []struct {
		name     string
		request  *pb.SetLEDsRequest
		wantResp *pb.SetLEDsResponse
		wantErr  bool
	}{
		{
			name: "Set RED LED on channel 0",
			request: &pb.SetLEDsRequest{
				StatusLedColor: pb.Color_RED,
				RecLedOn:       true,
				Channel:        0,
			},
			wantResp: &pb.SetLEDsResponse{Message: "LED and channel settings updated successfully"},
			wantErr:  false,
		},
		{
			name: "Set GREEN LED on channel 1",
			request: &pb.SetLEDsRequest{
				StatusLedColor: pb.Color_GREEN,
				RecLedOn:       false,
				Channel:        1,
			},
			wantResp: &pb.SetLEDsResponse{Message: "LED and channel settings updated successfully"},
			wantErr:  false,
		},
	}
	/*
		for _, tt := range tests {
			resp, err := server.SetLEDs(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLEDs() error = %v, wantErr %v", err, tt.wantErr)
				continue
			}
			if !tt.wantErr && resp.Message != tt.wantResp.Message {
				t.Errorf("SetLEDs() gotResp = %v, want %v", resp.Message, tt.wantResp.Message)
			}
		}
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.SetLEDs(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLEDs() error = %v, wantErr %v", err, tt.wantErr)
				//continue
			}
			if !tt.wantErr && resp.Message != tt.wantResp.Message {
				t.Errorf("SetLEDs() gotResp = %v, want %v", resp.Message, tt.wantResp.Message)
			}
		})
	}
}

func TestGetLEDs(t *testing.T) {
	//cfg.LoadConfigDefault()
	server := &LEDServer{}
	tests := []struct {
		name     string
		request  *pb.GetLEDsRequest
		wantResp *pb.GetLEDsResponse
		wantErr  bool
	}{
		{
			name: "Get RED LED on channel 0",
			request: &pb.GetLEDsRequest{
				Channel: 0,
			},
			wantResp: &pb.GetLEDsResponse{
				StatusLedColor: pb.Color_RED,
				RecLedOn:       true,
				Channel:        0,
			},
			wantErr: false,
		},
		{
			name: "Get GREEN LED on channel 1",
			request: &pb.GetLEDsRequest{
				Channel: 1,
			},
			wantResp: &pb.GetLEDsResponse{
				StatusLedColor: pb.Color_GREEN,
				RecLedOn:       false,
				Channel:        1,
			},
			wantErr: false,
		},
	}
	/*
		for _, tt := range tests {
			resp, err := server.GetLEDs(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLEDs() error = %v, wantErr %v", err, tt.wantErr)
				continue
			}
			if !tt.wantErr && (resp.StatusLedColor != tt.wantResp.StatusLedColor || resp.RecLedOn != tt.wantResp.RecLedOn || resp.Channel != tt.wantResp.Channel) {
				t.Errorf("GetLEDs() gotResp = %v, want %v", resp, tt.wantResp)
			}
		}
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetLEDs(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLEDs() error = %v, wantErr %v", err, tt.wantErr)
				//continue
			}
			if !tt.wantErr && (resp.StatusLedColor != tt.wantResp.StatusLedColor || resp.RecLedOn != tt.wantResp.RecLedOn || resp.Channel != tt.wantResp.Channel) {
				t.Errorf("GetLEDs() gotResp = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}
