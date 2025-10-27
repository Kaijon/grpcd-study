package main

import (
	"context"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
)

type LuxServer struct {
	pb.UnimplementedLuxServiceServer
	cfg *cfg.Config
}

func (s *LuxServer) GetDayNightMode(ctx context.Context, req *pb.GetDayNightModeRequest) (*pb.GetDayNightModeResponse, error) {
	Log.Info(">>Run")
	dayNightMode := s.cfg.DayNightMode.Mode
	return &pb.GetDayNightModeResponse{DayNightMode: dayNightMode}, nil
}
