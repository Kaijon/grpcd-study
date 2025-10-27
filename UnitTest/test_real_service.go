package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	localgrpc "canf22g2/grpc" // prevent conflict

	"google.golang.org/grpc"
)

type uploadType int

const (
	firmwareUpload uploadType = iota
	aiFileUpload
)

func uploadFile(client localgrpc.UnifiedFileTransferClient, path string, uType uploadType) {
	var stream localgrpc.UnifiedFileTransfer_UploadFirmwareClient
	var err error

	switch uType {
	case firmwareUpload:
		stream, err = client.UploadFirmware(context.Background())
	case aiFileUpload:
		stream, err = client.UploadAIFile(context.Background())
	default:
		log.Fatalf("unknown upload type: %v", uType)
	}

	if err != nil {
		log.Fatalf("failed to initiate upload: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			if n > 0 {
				if err := stream.Send(&localgrpc.UnifiedChunk{
					Filename: filepath.Base(path),
					Content:  buffer[:n],
				}); err != nil {
					log.Fatalf("failed to send final chunk to server: %v", err)
				}
			}
			break
		}
		if err != nil {
			log.Fatalf("failed to read chunk from file: %v", err)
		}

		if err := stream.Send(&localgrpc.UnifiedChunk{
			Filename: filepath.Base(path),
			Content:  buffer[:n],
		}); err != nil {
			log.Fatalf("failed to send chunk to server: %v", err)
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		if err == io.EOF {
			log.Println("Stream closed by server unexpectedly")
		} else {
			log.Fatalf("failed to receive upload status: %v", err)
		}
	} else {
		logType := "Upload"
		if uType == aiFileUpload {
			logType = "AI file upload"
		}
		log.Printf("%s finished with status: %v", logType, status)
	}
}

func main() {
	const testFilename = "test_real_service.go"
	const AIFilename = "alpr_update.tar.gz"

	conn, err := grpc.Dial("192.168.5.48:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := localgrpc.NewUnifiedFileTransferClient(conn)
	uploadFile(client, testFilename, firmwareUpload)
	uploadFile(client, AIFilename, aiFileUpload)
}
