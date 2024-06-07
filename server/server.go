package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/julmust/grpc_file_transfer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PORT        = ":50051"
	FILE_FOLDER = "files/"
)

type TransferServer struct {
	pb.UnimplementedTransferServiceServer
}

func (s *TransferServer) GetBatchTransfer(ctx context.Context, in *pb.TransferInfoMsg) (*pb.GetBatchTransferMsg, error) {
	log.Printf("Received request for: %v", in.GetFilename())
	fp := FILE_FOLDER + in.GetFilename()

	data, err := os.ReadFile(fp)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return &pb.GetBatchTransferMsg{}, status.Errorf(codes.NotFound, "Could not open file: "+in.GetFilename())
	}

	fileSize, err := os.Stat(fp)
	if err != nil {
		log.Print(err)
		return &pb.GetBatchTransferMsg{}, err
	}

	trans := &pb.GetBatchTransferMsg{
		Filename: fp,
		Filesize: fileSize.Size(),
		File:     data,
	}

	return trans, nil
}

func (s *TransferServer) PutBatchTransfer(ctx context.Context, in *pb.PutBatchTransferMsg) (*pb.TransferInfoMsg, error) {
	log.Printf("Recieved upload request for: %v", in.GetFilename())
	fp := FILE_FOLDER + in.GetFilename()

	err := os.WriteFile(fp, in.GetFile(), 0666)
	if err != nil {
		log.Print(err)
		return &pb.TransferInfoMsg{}, status.Errorf(codes.PermissionDenied, "Could not write file: %s", in.GetFilename())
	}

	return &pb.TransferInfoMsg{Filename: fp}, nil
}

func (s *TransferServer) GetFileList(ctx context.Context, in *pb.PutFolderName) (*pb.FileInfo, error) {
	path := in.GetName()
	if string(path[0]) == "/" {
		if len(path) == 1 {
			path = "files/"
		} else {
			path = "files/" + path
		}
	}

	fl, err := os.ReadDir(path)
	if err != nil {
		return &pb.FileInfo{}, err
	}

	ret := &pb.FileInfo{}
	for _, f := range fl {
		ret.Files = append(ret.Files, f.Name())
	}

	return ret, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterTransferServiceServer(s, &TransferServer{})

	log.Printf("server listening at %v", PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
