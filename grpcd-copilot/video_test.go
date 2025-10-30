package main

import (
	"context"
	"reflect"
	"testing"

	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
)

func TestVideoInfoServer_SetVideoSettings(t *testing.T) {
	cfg.LoadConfigDefault()
	server := &VideoInfoServer{}
	tests := []struct {
		name     string
		request  *pb.SetVideoSettingsRequest
		wantResp *pb.SetVideoSettingsResponse
		wantErr  bool
	}{
		{
			name: "Set video settings on channel 0",
			request: &pb.SetVideoSettingsRequest{
				Resolution:      "2560x1080",
				StreamFormat:    "h264",
				BitRate:         12,
				Type:            "vbr",
				Fps:             30,
				Channel:         0,
				SubResolution:   "1280x544",
				SubStreamFormat: "h264",
				SubBitRate:      4,
				SubType:         "vbr",
				SubFps:          30,
			},
			wantResp: &pb.SetVideoSettingsResponse{Message: "Video settings updated successfully for channel 0"},
			wantErr:  false,
		},
		{
			name: "Set video settings on channel 1",
			request: &pb.SetVideoSettingsRequest{
				Resolution:      "2560x1440",
				StreamFormat:    "h264",
				BitRate:         12,
				Type:            "vbr",
				Fps:             30,
				Channel:         1,
				SubResolution:   "1280x720",
				SubStreamFormat: "h264",
				SubBitRate:      4,
				SubType:         "vbr",
				SubFps:          30,
			},
			wantResp: &pb.SetVideoSettingsResponse{Message: "Video settings updated successfully for channel 1"},
			wantErr:  false,
		},
	}
	/*
	       for _, tt := range tests {
	   		resp, err := server.SetVideoSettings(context.Background(), tt.request)
	   		if (err != nil) != tt.wantErr {
	   			t.Errorf("SetVideoSettings() error = %v, wantErr %v", err, tt.wantErr)
	   			continue
	   		}
	   		if !tt.wantErr && resp.Message != tt.wantResp.Message {
	   			t.Errorf("SetVideoSettings() gotResp = %v, want %v", resp.Message, tt.wantResp.Message)
	   		}
	   	}
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.SetVideoSettings(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetVideoSettings() error = %v, wantErr %v", err, tt.wantErr)
				//continue
			}
			if !tt.wantErr && resp.Message != tt.wantResp.Message {
				t.Errorf("SetVideoSettings() gotResp = %v, want %v", resp.Message, tt.wantResp.Message)
			}
		})
	}
}

func TestVideoInfoServer_GetVideoSettings(t *testing.T) {
	//LoadConfigDefault()
	server := &VideoInfoServer{}

	tests := []struct {
		name     string
		request  *pb.GetVideoSettingsRequest
		wantResp *pb.GetVideoSettingsResponse
		wantErr  bool
	}{
		{
			name: "Get video settings on channel 0",
			request: &pb.GetVideoSettingsRequest{
				Channel: 0,
			},
			wantResp: &pb.GetVideoSettingsResponse{
				Resolution:      "2560x1080",
				StreamFormat:    "h264",
				BitRate:         12,
				Type:            "vbr",
				Fps:             30,
				SubResolution:   "1280x544",
				SubStreamFormat: "h264",
				SubBitRate:      4,
				SubType:         "vbr",
				SubFps:          30,
			},
			wantErr: false,
		},
		{
			name: "Get video settings on channel 1",
			request: &pb.GetVideoSettingsRequest{
				Channel: 1,
			},
			wantResp: &pb.GetVideoSettingsResponse{
				Resolution:      "2560x1440",
				StreamFormat:    "h264",
				BitRate:         12,
				Type:            "vbr",
				Fps:             30,
				SubResolution:   "1280x720",
				SubStreamFormat: "h264",
				SubBitRate:      4,
				SubType:         "vbr",
				SubFps:          30,
			},
			wantErr: false,
		},
	}
	/*
	       for _, tt := range tests {
	   		resp, err := server.GetVideoSettings(context.Background(), tt.request)
	   		if (err != nil) != tt.wantErr {
	   			t.Errorf("GetVideoSettings() error = %v, wantErr %v", err, tt.wantErr)
	   			continue
	   		}
	           if !reflect.DeepEqual(resp, tt.wantResp) {
	   			t.Errorf("GetVideoSettings() gotResp = %v, want %v", resp, tt.wantResp)
	   		}
	   	}
	*/
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetVideoSettings(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideoSettings() error = %v, wantErr %v", err, tt.wantErr)
				//continue
			}
			if !reflect.DeepEqual(resp, tt.wantResp) {
				t.Errorf("GetVideoSettings() gotResp = %v, want %v", resp, tt.wantResp)
			}
		})
	}
}
