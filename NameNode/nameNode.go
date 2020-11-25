package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "lab2"

	"google.golang.org/grpc"
)

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type books struct {
	name  string
	parts int32
	//stored int32
}

var library map[string]*books

func storeInLibrary(book books) {
	library[book.name] = &book
}

// RequestBook envia un bookinfo a un cliente
func (s *server) RequestBook(ctx context.Context, in *pb.BookRequest) (*pb.BookReply1, error) {
	return &pb.BookReply1{Parts: library[in.GetBookNamePart()].parts}, nil
}

// StoreBook guarda un libro en el log
func (s *server) StoreBook(ctx context.Context, in *pb.BookStoreRequest) (*pb.Response, error) {

	tempBook := books{
		name:  in.GetBookName(),
		parts: in.GetTotalParts(),
	}
	// TODO guardar info del libro
	storeInLibrary(tempBook)

	fmt.Println(library)

	return &pb.Response{Response: "Book succesfuly stored"}, nil
}

//ListenToClient listener
func ListenToClient(port string) {
	//--------------------------------------------------------------> Server1
	fmt.Println("Esperando solicitudes")
	lis, err := net.Listen("tcp", port)
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
	library = make(map[string]*books)

	downloads := ":50050"
	DataNodes := ":50051"

	go ListenToClient(downloads) // downloads
	ListenToClient(DataNodes)    // DataNodes

}
