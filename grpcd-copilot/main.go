package main

import (
	"context"
	pb "grpcd/canf22g2/grpc"
	cfg "grpcd/config"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func main() {
	grpcAddr := "0.0.0.0:50051"
	httpAddr := "0.0.0.0:8081"
	swaggerDir := "./swagger"

	cfg.Init()
	MqttInit()
	go StartMqttInLoop()
	StartMqttWorker()
	// It is used on emulator
	//LoadConfigDefault()

	// --- Create the main multiplexer ---
	// This will route traffic to either the gRPC gateway or the Swagger UI
	mainMux := http.NewServeMux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// --- Setup gRPC-Gateway ---
	gwmux := grpcRuntime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterLEDServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		Log.Fatalf("LED can not register gateway endpoint: %v", err)
	}
	err = pb.RegisterLuxServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		Log.Fatalf("Lux can not register gateway endpoint: %v", err)
	}
	err = pb.RegisterNetworkInfoServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		Log.Fatalf("Network can not register gateway endpoint: %v", err)
	}
	err = pb.RegisterDeviceInfoServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		Log.Fatalf("System can not register gateway endpoint: %v", err)
	}
	err = pb.RegisterVideoInfoServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		Log.Fatalf("Video can not register gateway endpoint: %v", err)
	}
	err = pb.RegisterWatermarkInfoServiceHandlerFromEndpoint(ctx, gwmux, grpcAddr, opts)
	if err != nil {
		Log.Fatalf("Watermark can not register gateway endpoint: %v", err)
	}

	// Mount the gRPC gateway to the main multiplexer
	mainMux.Handle("/", gwmux)
	// --- Setup Swagger UI ---
	// This handler serves the ioctrl.swagger.json file
	mainMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(swaggerDir, "swagger.json"))
	})
	fs := http.FileServer(http.Dir("./swagger"))
	mainMux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	Log.Printf("Serving Swagger UI at http://%s/swagger/", httpAddr)

	// Start gRPC server
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		Log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()

	// register statusLED server (inject config)
	ledServer := &LEDServer{cfg: &cfg.AppConfig}
	pb.RegisterLEDServiceServer(srv, ledServer)

	// register DeviceInfo server
	deviceInfoServer := &DeviceInfoServer{cfg: &cfg.AppConfig}
	pb.RegisterDeviceInfoServiceServer(srv, deviceInfoServer)

	// register NetworkInfo server
	networkInfoServer := &NetworkInfoServer{cfg: &cfg.AppConfig}
	pb.RegisterNetworkInfoServiceServer(srv, networkInfoServer)

	// register VideoInfo server
	videoInfoServer := &VideoInfoServer{cfg: &cfg.AppConfig}
	pb.RegisterVideoInfoServiceServer(srv, videoInfoServer)

	// register WatermarkInfo server
	watermarkInfoServer := &WatermarkInfoServer{cfg: &cfg.AppConfig}
	pb.RegisterWatermarkInfoServiceServer(srv, watermarkInfoServer)

	// register UnifiedFileTransfer server
	unifiedFileTransferServer := &UnifiedFileTransferServer{}
	pb.RegisterUnifiedFileTransferServer(srv, unifiedFileTransferServer)

	// register Lux server
	luxServer := &LuxServer{cfg: &cfg.AppConfig}
	pb.RegisterLuxServiceServer(srv, luxServer)

	fotaServer := &FotaServer{cfg: &cfg.AppConfig}
	pb.RegisterFotaServiceServer(srv, fotaServer)

	// Register reflection for debugging
	reflection.Register(srv)
	// start gRPC server

	Log.Printf("gRPC server is starting to listen on %s", grpcAddr)
	go func() {
		if err := srv.Serve(lis); err != nil {
			Log.Fatalf("failed to serve: %v", err)
		}
	}()

	Log.Printf("HTTP reverse proxy is starting to listen on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, mainMux); err != nil {
		Log.Fatalf("HTTP Server fails: %v", err)
	}
}

func init() {
	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05.000000",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
	Log.SetOutput(io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename:   "/tmp/logger_storage/APLog/grpcd.log",
			MaxSize:    1, // megabytes
			MaxBackups: 15,
			MaxAge:     1,     //days
			Compress:   false, // disabled by default
		}))
	numProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(numProcs)
	Log.Infoln("GOMAXPROCS set to:", numProcs)
}

type LogDump struct {
	level logrus.Level
}

func (lw LogDump) Write(p []byte) (n int, err error) {
	switch lw.level {
	case logrus.DebugLevel:
		Log.Debug(string(p))
	case logrus.ErrorLevel:
		Log.Error(string(p))
	default:
		Log.Info(string(p))
	}
	return len(p), nil
}
