package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	pb "github.com/julmust/grpc_file_transfer/proto"
	"google.golang.org/grpc"
)

const (
	ADDRESS = "localhost:50051"
)

var c pb.TransferServiceClient

func TestGetBatchTransfer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetBatchTransfer(ctx, &pb.TransferInfoMsg{
		Filename: "test.png",
	})
	if err != nil {
		log.Fatalf("could not create transfer: %v", err)
	}

	// Check that recieved data is the same as saved data
	fs, _ := os.ReadFile("test_files/test.png")
	if !bytes.Equal(fs, res.GetFile()) {
		t.Fatal("Mismatch!")
	}
}

func TestPutBatchTransfer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sd, _ := os.ReadFile("test_files/cnc.png")
	res, err := c.PutBatchTransfer(ctx, &pb.PutBatchTransferMsg{Filename: "test2.png", File: sd})
	if err != nil {
		log.Fatal(err)
	}

	// Check that recieved data is the same as saved data
	newFile, _ := os.ReadFile(res.GetFilename())
	if !bytes.Equal(sd, newFile) {
		t.Fatal("Mismatch!")
	}
}

func TestFileList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.GetFileList(ctx, &pb.PutFolderName{Name: "/"})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)

	// Check that file list is correct
}

func TestMain(m *testing.M) {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed connection: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterTransferServiceServer(s, &TransferServer{})

	// Launching server in goroutine as to not block execution
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server: %v", err)
		}
	}()

	conn, err := grpc.NewClient(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c = pb.NewTransferServiceClient(conn)

	// Setup folders
	os.RemoveAll("files/")
	os.Mkdir("files/", 0755)
	os.Mkdir("files/folder/", 0755)

	// Copy over file
	bf, err := os.ReadFile("test_files/test.png")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("files/test.png", bf, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ev := m.Run()

	os.Exit(ev)
}
