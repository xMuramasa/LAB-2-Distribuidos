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
	name   string
	parts  int32
	titles []string
	chunks []string
}

var storage map[string]*books

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type bookPart struct {
	name string
	part []byte
}

var library map[string]*bookPart

func store(b bookPart) {
	if library[b.name] == nil {
		library[b.name] = &b
	}
}

func storeInStorage(b books, title string, chunk string) {
	if storage[b.name] == nil {
		storage[b.name] = &b
	} else {
		storage[b.name].titles = append(storage[b.name].titles, title)
		storage[b.name].chunks = append(storage[b.name].chunks, chunk)
	}
}

// Recibe un chunk y lo guarda en un archivo en el disco
func (s *server) ReceiveChunk(ctx context.Context, in *pb.StoreRequest) (*pb.StoreReply, error) {
	log.Printf("Received: chunk 250kb. From: %v", in.GetClientName())

	//save chunk to memory
	tempBook := books{
		name:   in.GetFileName(),
		parts:  in.GetPart(),
		titles: []string{in.GetFileName() + "_part_" + in.GetChunkPart()},
		chunks: []string{string(in.GetChunk())},
	}
	storeInStorage(tempBook, tempBook.titles[0], tempBook.chunks[0])

	//if got all chunks -> debo avisar a nnode y repartir entre dnodes
	lenTitles := len(storage[tempBook.name].titles)
	lenChunks := len(storage[tempBook.name].chunks)
	if int(storage[tempBook.name].parts) == lenTitles && lenTitles == lenChunks {
		//prepare Proposal
		partsPerNode := int(int(storage[tempBook.name].parts) / 3)
		rest := int(storage[tempBook.name].parts) % 3
		message := fmt.Sprintf("%s**%d**%d**%d**%d", tempBook.name,
			partsPerNode+rest,
			partsPerNode,
			partsPerNode,
			tempBook.parts)

		//send proposal
		conn, err := grpc.Dial("dist29:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.Proposal(ctx, &pb.Message{
			M: message})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		//deal with response in r
		if r.GetM() == "A" {
			//propuesta aceptada -> repartir con message
		} else {
			//propuesta rechazada -> repartir con r
		}

	}

	// write to disk
	//fileName := "./stored/" + in.GetFileName() + "_part_" + in.GetChunkPart()
	//_, err := os.Create(fileName)
	/*
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/
	// write/save buffer to disk
	//ioutil.WriteFile(fileName, in.GetChunk(), os.ModeAppend)

	//if libro completo, realizar solicitud

	return &pb.StoreReply{Message: "Received & stored chunk"}, nil
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

func (s *server) Greeting(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	return &pb.Message{M: "1"}, nil
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
	library = make(map[string]*bookPart)
	storage = make(map[string]*books)

	go ListenToClient(":50051") // clientes descargas
	go ListenToClient(":50052") // clientes cargas
	go ListenToClient(":50053") // datanode 2
	ListenToClient(":50054")    // datanode 3

}
