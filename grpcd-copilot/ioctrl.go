package main

import (
	"context"
	"fmt"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"time"
)

type LEDServer struct {
	pb.UnimplementedLEDServiceServer
}

func (s *LEDServer) SetLEDs(ctx context.Context, in *pb.SetLEDsRequest) (*pb.SetLEDsResponse, error) {
	channelKey := fmt.Sprintf("%d", in.Channel)
	strStatus := in.StatusLedColor.String()
	if err := cfg.UpdateConfig(func(c *cfg.Config) {
		if c.LEDs == nil {
			c.LEDs = make(map[string]cfg.LEDConfig)
		}
		c.LEDs[channelKey] = cfg.LEDConfig{
			StatusLed: strStatus,
			RecLedOn:  in.RecLedOn,
		}
	}); err != nil {
		return nil, err
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
	var (
		ledConfig cfg.LEDConfig
		ok        bool
	)
	cfg.ReadConfig(func(current cfg.Config) {
		if current.LEDs != nil {
			ledConfig, ok = current.LEDs[channelKey]
		}
	})
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
