package main

import (
	"context"
	"fmt"
	cfg "grpcd/config"
	pb "grpcd/canf22g2/grpc"
	"time"
)

type LEDServer struct {
	pb.UnimplementedLEDServiceServer
	cfg *cfg.Config
}

func (s *LEDServer) SetLEDs(ctx context.Context, in *pb.SetLEDsRequest) (*pb.SetLEDsResponse, error) {
	channelKey := fmt.Sprintf("%d", in.Channel)
	strStatus := in.StatusLedColor.String()
	if s.cfg.LEDs == nil {
		s.cfg.LEDs = make(map[string]cfg.LEDConfig)
	}
	s.cfg.LEDs[channelKey] = cfg.LEDConfig{
		StatusLed: strStatus,
		RecLedOn:  in.RecLedOn,
	}
	strTmp := fmt.Sprintf("{\"StatusLed\":\"%s\", \"RecLedOn\":%v}", strStatus, in.RecLedOn)
	topic := fmt.Sprintf("config/io/led/%d", in.Channel)
	msg := MqttMessage{
		Topic:   topic,
		Payload: strTmp,
	}
	select {
	case MqttPublishChannel <- msg:
	case <-time.After(50 * time.Millisecond):
		Log.Warnf("Timed out sending message to channel: %v", msg)
	}
	return &pb.SetLEDsResponse{Message: "LED and channel settings updated successfully"}, nil
}

func (s *LEDServer) GetLEDs(ctx context.Context, in *pb.GetLEDsRequest) (*pb.GetLEDsResponse, error) {
	Log.Info(">>Run")
	channelKey := fmt.Sprintf("%d", in.Channel)
	ledConfig, ok := s.cfg.LEDs[channelKey]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", channelKey)
	}
	colorValue, ok := pb.Color_value[ledConfig.StatusLed]
	if !ok {
		return nil, fmt.Errorf("invalid color value: %s", ledConfig.StatusLed)
	}

	Log.Debugf("Ch:%d, %v", in.Channel, ledConfig)
	return &pb.GetLEDsResponse{
		StatusLedColor: pb.Color(colorValue),
		RecLedOn:       ledConfig.RecLedOn,
		Channel:        in.Channel,
	}, nil
}
