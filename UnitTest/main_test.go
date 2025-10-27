package main

import (
	pb "grpcd/canf22g2/grpc"
	"io"
	"log"
	"os"
	"runtime"
    //"net"
    //"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/natefinch/lumberjack.v2"
    "google.golang.org/grpc/test/bufconn"
    //"google.golang.org/grpc/credentials/insecure"
)

var Log = logrus.New()

func main() {
//func server(ctx context.Context) (pb.LEDServiceClient, func()) {
	configInit()
	MqttInit()
	go StartMqttInLoop()
	// It is used on emulator
	LoadConfigDefault()

    buffer := 1024 * 1024
    lis := bufconn.Listen(buffer)

	// Start gRPC server
	//lis, err := net.Listen("tcp", "[::]:50051")

	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	srv := grpc.NewServer()

	// register statusLED server
	ledServer := &LEDServer{}
	pb.RegisterLEDServiceServer(srv, ledServer)

    // register NetworkInfo server
	networkInfoServer := &NetworkInfoServer{}
	pb.RegisterNetworkInfoServiceServer(srv, networkInfoServer)

    // register WatermarkInfo server
	watermarkInfoServer := &WatermarkInfoServer{}
	pb.RegisterWatermarkInfoServiceServer(srv, watermarkInfoServer)

    // register VideoInfo server
	videoInfoServer := &VideoInfoServer{}
	pb.RegisterVideoInfoServiceServer(srv, videoInfoServer)

    // register DeviceInfo server
	deviceInfoServer := &DeviceInfoServer{}
	pb.RegisterDeviceInfoServiceServer(srv, deviceInfoServer)

    /*
	// register UnifiedFileTransfer server
	unifiedFileTransferServer := &UnifiedFileTransferServer{}
	pb.RegisterUnifiedFileTransferServer(srv, unifiedFileTransferServer)

	// register Lux server
	luxServer := &LuxServer{}
	pb.RegisterLuxServiceServer(srv, luxServer)

	fotaServer := &FotaServer{}
	pb.RegisterFotaServiceServer(srv, fotaServer)
	*/
	// Register reflection for debugging
	reflection.Register(srv)
	// start gRPC server
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
    /*
    conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		srv.Stop()
	}

	client := pb.NewLEDServiceClient(conn)

	return client, closer
    */
}


func init() {
	Log.SetLevel(logrus.DebugLevel)
	Log.SetReportCaller(true)
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
	})
	Log.SetOutput(io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename:   "/mnt/flash/logger_storage/APLog/grpcd.log",
			MaxSize:    1, // megabytes
			MaxBackups: 5,
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
