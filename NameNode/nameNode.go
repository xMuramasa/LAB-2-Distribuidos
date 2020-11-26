package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "lab2"

	"google.golang.org/grpc"
)

const (
	dataNode1 = ""
	dataNode2 = ""
	dataNode3 = ""
)

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type books struct {
	name  string
	c1    string
	c2    string
	c3    string
	parts string
}

var library map[string]*books

func storeInLibrary(book books) {
	library[book.name] = &book
}

// RequestBook envia un bookinfo a un cliente
func (s *server) RequestBook(ctx context.Context, in *pb.BookRequest) (*pb.BookReply1, error) {
	i, err := strconv.Atoi(library[in.GetBookNamePart()].parts)
	if err != nil {
		log.Fatal("could convert to int")
	}
	return &pb.BookReply1{Parts: int32(i)}, nil
}

/*
// StoreBook guarda un libro en el log
func (s *server) StoreBook(ctx context.Context, in *pb.BookStoreRequest) (*pb.Message, error) {

	tempBook := books{
		name:  in.GetBookName(),
		parts: in.GetTotalParts(),
	}
	// TODO guardar info del libro
	storeInLibrary(tempBook)

	fmt.Println(library)

	return &pb.Message{M: "Book succesfuly stored"}, nil
}
*/

func (s *server) Proposal(ctx context.Context, in *pb.Message) (*pb.Message, error) {

	var i int
	var j int
	var k int

	flag := false

	proposal := strings.Split(in.GetM(), "**") //nombre**c1**c2**c3**total

	// proposal = [nombre, c1, c2, c3, total]
	// if c1 down -> proposal[1] = 0,
	//			proposal[2] = proposal[2] + proposal[1],
	//			proposal[3] = proposal[3],

	// checkear que los 3 DN estan vivos
	if connectToDataNode(dataNode1) == "0" {
		flag = true
		x, _ := strconv.Atoi(proposal[1])
		y, _ := strconv.Atoi(proposal[2])
		proposal[1] = "0"
		proposal[2] = fmt.Sprintf("%d", x+y)
	}
	if connectToDataNode(dataNode2) == "0" {
		flag = true
		x, _ := strconv.Atoi(proposal[2])
		y, _ := strconv.Atoi(proposal[3])
		proposal[2] = "0"
		proposal[3] = fmt.Sprintf("%d", x+y)
	}
	if connectToDataNode(dataNode3) == "0" {
		flag = true
		x, _ := strconv.Atoi(proposal[3])
		y, _ := strconv.Atoi(proposal[1])
		proposal[3] = "0"
		proposal[1] = fmt.Sprintf("%d", x+y)
	}
	if flag {
		b := fmt.Sprintf("%s**%s**%s**%s**%s",
			proposal[0], proposal[1], proposal[2], proposal[3], proposal[4])
		return &pb.Message{M: b}, nil
	}

	//guardamos
	tempBook := books{
		name:  proposal[0],
		c1:    proposal[1],
		c2:    proposal[2],
		c3:    proposal[3],
		parts: proposal[4],
	}
	//guardar info del libro
	storeInLibrary(tempBook)

	//write
	fileName := "log.txt"
	_, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buff := library[tempBook.name].name + " " + library[tempBook.name].parts + "\n"

	nChunks, err := strconv.Atoi(library[tempBook.name].c1)
	for i = 0; i < nChunks; i++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", i) + " " + dataNode1 + "\n"
	}

	nChunks, err = strconv.Atoi(library[tempBook.name].c2)
	for j = 0; j < nChunks; j++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", j+i) + " " + dataNode2 + "\n"
	}

	nChunks, err = strconv.Atoi(library[tempBook.name].c3)
	for k = 0; k < nChunks; k++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", k+j) + " " + dataNode3 + "\n"
	}

	ioutil.WriteFile(fileName, []byte(buff), os.ModeAppend)
	//defer f.Close()

	return &pb.Message{M: "A"}, nil
}

func connectToDataNode(dataNode string) string {
	conn, err := grpc.Dial(dataNode, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Hacer una consulta
	m := &pb.Message{
		M: "0",
	}
	r, err := c.Greeting(ctx, m)
	if err != nil {
		log.Fatalf("could not greet: %s --- %v", dataNode, err)
		return "0"
	}
	return r.GetM()
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

	go ListenToClient(downloads) // downloads
	go ListenToClient(":50051")  // DataNode 1
	go ListenToClient(":50052")  // DataNodes 2
	ListenToClient(":50053")     // DataNodes 3

}
