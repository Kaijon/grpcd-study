package main

import (
	"context"
	"fmt"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"time"
)

type VideoInfoServer struct {
	pb.UnimplementedVideoInfoServiceServer
	cfg *cfg.Config
}

func (s *VideoInfoServer) SetVideoSettings(ctx context.Context, in *pb.SetVideoSettingsRequest) (*pb.SetVideoSettingsResponse, error) {
	channelKey := fmt.Sprintf("%d", in.Channel)
	if s.cfg.Videos == nil {
		s.cfg.Videos = make(map[string]cfg.VideoConfig)
	}
	s.cfg.Videos[channelKey] = cfg.VideoConfig{
		Resolution:      in.Resolution,
		StreamFormat:    in.StreamFormat,
		BitRate:         in.BitRate,
		Type:            in.Type,
		Fps:             in.Fps,
		SubResolution:   in.SubResolution,
		SubStreamFormat: in.SubStreamFormat,
		SubBitRate:      in.SubBitRate,
		SubType:         in.SubType,
		SubFps:          in.SubFps,
		MirrorAction:    in.MirrorAction,
	}

	strTmp := fmt.Sprintf("{\"Resolution\":\"%s\",\"StreamFormat\":\"%s\",\"BitRate\":%d,\"Type\":\"%s\",\"Fps\":%d,\"SubResolution\":\"%s\",\"SubStreamFormat\":\"%s\",\"SubBitRate\":%d,\"SubType\":\"%s\",\"SubFps\":%d,\"MirrorAction\":\"%s\"}",
		in.Resolution, in.StreamFormat, in.BitRate, in.Type, in.Fps,
		in.SubResolution, in.SubStreamFormat, in.SubBitRate, in.SubType, in.SubFps, in.MirrorAction)
	topic := fmt.Sprintf("config/video/%d", in.Channel)
	msg := MqttMessage{
		Topic:   topic,
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message to channel: %v", msg)
	}
	return &pb.SetVideoSettingsResponse{Message: "Video settings updated successfully for channel " + channelKey}, nil
}

func (s *VideoInfoServer) GetVideoSettings(ctx context.Context, in *pb.GetVideoSettingsRequest) (*pb.GetVideoSettingsResponse, error) {
	Log.Info(">>Run")
	channelKey := fmt.Sprintf("%d", in.Channel)
	videoConfig, ok := s.cfg.Videos[channelKey]
	if !ok {
		return nil, fmt.Errorf("video settings not found for channel %s", channelKey)
	}
	Log.Debugf("Ch:%d, %v", in.Channel, videoConfig)
	return &pb.GetVideoSettingsResponse{
		Resolution:      videoConfig.Resolution,
		StreamFormat:    videoConfig.StreamFormat,
		BitRate:         videoConfig.BitRate,
		Type:            videoConfig.Type,
		Fps:             videoConfig.Fps,
		SubResolution:   videoConfig.SubResolution,
		SubStreamFormat: videoConfig.SubStreamFormat,
		SubBitRate:      videoConfig.SubBitRate,
		SubType:         videoConfig.SubType,
		SubFps:          videoConfig.SubFps,
		MirrorAction:    videoConfig.MirrorAction,
	}, nil
	//Resolution:      "2560x1080",
	//StreamFormat:    "h265",
	//BitRate:         12,
	//Type:            "vbr",
	//Fps:             30,
	//SubResolution:   "1280x720",
	//SubStreamFormat: "h264",
	//SubBitRate:      4,
	//SubType:         "vbc",
	//SubFps:          14,
	//MirrorAction:    "normal", //normal, flipHori, flipVert, rotate180
	//}, nil

}
