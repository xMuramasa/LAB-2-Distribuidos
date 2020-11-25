// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	pb "lab2"

	"google.golang.org/grpc"
)

type books struct {
	name  string
	parts int32
	//stored int32
}

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// Recibe un chunk y lo guarda en un archivo en el disco
func (s *server) ReceiveChunk(ctx context.Context, in *pb.StoreRequest) (*pb.StoreReply, error) {
	log.Printf("Received: chunk 250kb. From: %v", in.GetClientName())

	// write to disk
	fileName := "./stored/" + in.GetFileName() + "_part_" + in.GetChunkPart()
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// write/save buffer to disk
	ioutil.WriteFile(fileName, in.GetChunk(), os.ModeAppend)

	// Set up a connection to the NameNode.
	// Contact the server and print out its response.
	conn, err := grpc.Dial("dist29:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Hacer una consulta
	tempBook := &pb.BookStoreRequest{
		BookName:   in.GetFileName(),
		TotalParts: in.GetPart(),
	}
	r, err := c.StoreBook(ctx, tempBook)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r)

	return &pb.StoreReply{Message: "Received chunk & stored in disk"}, nil
}

// RequestChunk envia un chunk a un cliente
func (s *server) RequestChunk(ctx context.Context, in *pb.BookRequest) (*pb.BookReply2, error) {

	fileName := "./stored/" + in.GetBookNamePart()

	// open and read file
	toSend, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &pb.BookReply2{Chunk: toSend}, nil
}

//ListenToClient listener
func ListenToClient(puerto string) {
	//--------------------------------------------------------------> Server1
	fmt.Println("Esperando solicitudes")
	lis, err := net.Listen("tcp", puerto)
	if err != nil {
		log.Fatalf("failed to listen1: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve s1: %v", err)
	}

}

func main() {
	go ListenToClient(":50051") // clientes descargas
	ListenToClient(":50052")    // clientes cargas

}
