package main

import (
	"context"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type SetAllWatermarkInfoRequest struct {
	Username         string                `json:"username"`
	OptionUserName   bool                  `json:"OptionUserName"`
	OptionDeviceName bool                  `json:"OptionDeviceName"`
	OptionGPS        bool                  `json:"OptionGPS"`
	OptionTime       bool                  `json:"OptionTime"`
	OptionLogo       *wrapperspb.BoolValue `json:"OptionLogo"`
	OptionExposure   *wrapperspb.BoolValue `json:"OptionExposure"`
}

func TestWatermarkInfoServer_SetAllWatermarkInfo(t *testing.T) {
	cfg.LoadConfigDefault()
	server := &WatermarkInfoServer{}

	tests := []struct {
		name     string
		request  *pb.SetAllWatermarkInfoRequest
		wantResp *pb.SetAllWatermarkInfoResponse
		wantErr  bool
	}{
		{
			name: "Set watermark on channel 0",
			request: &pb.SetAllWatermarkInfoRequest{
				Username:         "User0",
				OptionUserName:   false,
				OptionDeviceName: true,
				OptionGPS:        false,
				OptionTime:       true,
				OptionLogo:       wrapperspb.Bool(false),
				Channel:          0,
			},
			wantResp: &pb.SetAllWatermarkInfoResponse{Message: "Watermark information updated successfully for channel 0"},
			wantErr:  false,
		},
		{
			name: "Set watermark on channel 1",
			request: &pb.SetAllWatermarkInfoRequest{
				Username:         "User1",
				OptionUserName:   true,
				OptionDeviceName: false,
				OptionGPS:        true,
				OptionTime:       false,
				OptionLogo:       wrapperspb.Bool(true),
				Channel:          1,
			},
			wantResp: &pb.SetAllWatermarkInfoResponse{Message: "Watermark information updated successfully for channel 1"},
			wantErr:  false,
		},
	}
	/*
	       for _, tt := range tests {
	   		resp, err := server.SetAllWatermarkInfo(context.Background(), tt.request)
	   		if (err != nil) != tt.wantErr {
	   			t.Errorf("SetAllWatermarkInfo() error = %v, wantErr %v", err, tt.wantErr)
	   			continue
	   		}
	   		if !tt.wantErr && resp.Message != tt.wantResp.Message {
	   			t.Errorf("SetAllWatermarkInfo() gotResp = %v, want %v", resp.Message, tt.wantResp.Message)
	   		}
	   	}
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.SetAllWatermarkInfo(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetAllWatermarkInfo() error = %v, wantErr %v", err, tt.wantErr)
				//continue
			}
			if !tt.wantErr && resp.Message != tt.wantResp.Message {
				t.Errorf("SetAllWatermarkInfo() gotResp = %v, want %v", resp.Message, tt.wantResp.Message)
			}
		})
	}
}

func TestWatermarkInfoServer_GetAllWatermarkInfo(t *testing.T) {
	//LoadConfigDefault()
	server := &WatermarkInfoServer{}

	tests := []struct {
		name     string
		request  *pb.GetAllWatermarkInfoRequest
		wantResp *pb.GetAllWatermarkInfoResponse
		wantErr  bool
	}{
		{
			name: "Get watermark on channel 0",
			request: &pb.GetAllWatermarkInfoRequest{
				Channel: 0,
			},
			wantResp: &pb.GetAllWatermarkInfoResponse{
				Username:         "User0",
				OptionUserName:   false,
				OptionDeviceName: true,
				OptionGPS:        false,
				OptionTime:       true,
				OptionLogo:       false,
			},
			wantErr: false,
		},
		{
			name: "Get watermark on channel 1",
			request: &pb.GetAllWatermarkInfoRequest{
				Channel: 1,
			},
			wantResp: &pb.GetAllWatermarkInfoResponse{
				Username:         "User1",
				OptionUserName:   true,
				OptionDeviceName: false,
				OptionGPS:        true,
				OptionTime:       false,
				OptionLogo:       true,
			},
			wantErr: false,
		},
	}
	/*
	       for _, tt := range tests {
	   		resp, err := server.GetAllWatermarkInfo(context.Background(), tt.request)
	   		if (err != nil) != tt.wantErr {
	   			t.Errorf("GetAllWatermarkInfo() error = %v, wantErr %v", err, tt.wantErr)
	   			continue
	   		}

	           if !reflect.DeepEqual(resp, tt.wantResp) {
	   			t.Errorf("GetAllWatermarkInfo() gotResp = %v, want %v", resp, tt.wantResp)
	   		}
	   	}
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetAllWatermarkInfo(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllWatermarkInfo() error = %v, wantErr %v", err, tt.wantErr)
				//continue
			}

			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("GetAllWatermarkInfo() gotResp = %v, want %v", resp, tt.wantResp)
			}
		})
	}

}
