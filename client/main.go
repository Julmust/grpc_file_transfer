package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/julmust/grpc_file_transfer/proto"
	"google.golang.org/grpc"
)

const (
	ADDRESS = "localhost:50051"
)

type TransferInfo struct {
	Filename string
}

func main() {
	conn, err := grpc.NewClient(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTransferServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.CreateTransfer(ctx, &pb.TransferInfo{
		Filename: "test2.png",
	})
	if err != nil {
		log.Fatalf("could not create transfer: %v", err)
	}

	log.Printf(`
		Filename: %s
		Filesize: %f
	`, res.GetFilename(), res.GetFilesize())

	err = os.WriteFile("out/test.png", res.GetFile(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
