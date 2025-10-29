package main

import (
	"context"
	"fmt"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"time"
)

type WatermarkInfoServer struct {
	pb.UnimplementedWatermarkInfoServiceServer
}

func (s *WatermarkInfoServer) GetAllWatermarkInfo(ctx context.Context, req *pb.GetAllWatermarkInfoRequest) (*pb.GetAllWatermarkInfoResponse, error) {
	Log.Info(">>Run")
	channelKey := fmt.Sprintf("%d", req.Channel)
	var (
		watermarkConfig cfg.WatermarkConfig
		ok              bool
	)
	cfg.ReadConfig(func(current cfg.Config) {
		if current.Watermarks != nil {
			watermarkConfig, ok = current.Watermarks[channelKey]
		}
	})
	if !ok {
		return nil, fmt.Errorf("watermark settings not found for channel %s", channelKey)
	}

	Log.Debugf("Ch:%d, %v", req.Channel, watermarkConfig)

	return &pb.GetAllWatermarkInfoResponse{
		Username:         watermarkConfig.Username,
		OptionUserName:   watermarkConfig.OptionUserName,
		OptionDeviceName: watermarkConfig.OptionDeviceName,
		OptionGPS:        watermarkConfig.OptionGPS,
		OptionTime:       watermarkConfig.OptionTime,
		OptionLogo:       watermarkConfig.OptionLogo,
		OptionExposure:   watermarkConfig.OptionExposure,
	}, nil
}

func (s *WatermarkInfoServer) SetAllWatermarkInfo(ctx context.Context, req *pb.SetAllWatermarkInfoRequest) (*pb.SetAllWatermarkInfoResponse, error) {
	channelKey := fmt.Sprintf("%d", req.Channel)

	var logo, expo bool
	var previous cfg.WatermarkConfig
	cfg.ReadConfig(func(current cfg.Config) {
		if current.Watermarks != nil {
			previous = current.Watermarks[channelKey]
		}
	})

	if req.OptionLogo != nil {
		logo = req.OptionLogo.GetValue()
	} else {
		logo = previous.OptionLogo
	}

	if req.OptionExposure != nil {
		expo = req.OptionExposure.GetValue()
	} else {
		expo = previous.OptionExposure
	}

	strTmp := fmt.Sprintf("{\"Username\":\"%s\", \"OptionUserName\":%v, \"OptionDeviceName\":%v,\"OptionGPS\":%v,\"OptionTime\":%v,\"OptionLogo\":%v,\"OptionExposure\":%v}",
		req.Username, req.OptionUserName, req.OptionDeviceName, req.OptionGPS, req.OptionTime, logo, expo)

	if err := cfg.UpdateConfig(func(c *cfg.Config) {
		if c.Watermarks == nil {
			c.Watermarks = make(map[string]cfg.WatermarkConfig)
		}
		c.Watermarks[channelKey] = cfg.WatermarkConfig{
			Username:         req.Username,
			OptionUserName:   req.OptionUserName,
			OptionDeviceName: req.OptionDeviceName,
			OptionGPS:        req.OptionGPS,
			OptionTime:       req.OptionTime,
			OptionLogo:       logo,
			OptionExposure:   expo,
		}
	}); err != nil {
		return nil, err
	}

	topic := fmt.Sprintf("config/watermark/%d", req.Channel)
	msg := MqttMessage{
		Topic:   topic,
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message to channel: %v", msg)
	}
	return &pb.SetAllWatermarkInfoResponse{Message: "Watermark information updated successfully for channel " + channelKey}, nil
}
