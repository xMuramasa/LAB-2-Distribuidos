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
	dataNode1 = "dist31:50050"
	dataNode2 = "dist30:50050"
	dataNode3 = "dist32:50050"
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

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// RequestBook envia un bookinfo a un cliente
func (s *server) RequestBook(ctx context.Context, in *pb.BookRequest) (*pb.BookReply1, error) {

	b := fmt.Sprintf("%s**%s**%s**%s**%s",
		library[in.GetBookNamePart()].name,
		library[in.GetBookNamePart()].c1,
		library[in.GetBookNamePart()].c2,
		library[in.GetBookNamePart()].c3,
		library[in.GetBookNamePart()].parts)
	return &pb.BookReply1{Locations: b}, nil
}

func (s *server) WriteRequest(ctx context.Context, in *pb.Message) (*pb.Message, error) {

	log.Println("Received Writing Request.")

	var i int
	var j int
	var k int

	proposal := strings.Split(in.GetM(), "**") //nombre**c1**c2**c3**total

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

	log.Println("Writing log file.")

	//write
	fileName := "log.txt"
	if Exists(fileName) {
		_, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		_, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	buff := library[tempBook.name].name + " " + library[tempBook.name].parts + "\n"

	nChunks, _ := strconv.Atoi(library[tempBook.name].c1)
	for i = 0; i < nChunks; i++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", i) + " " + dataNode1 + "\n"
	}

	nChunks, _ = strconv.Atoi(library[tempBook.name].c2)
	for j = 0; j < nChunks; j++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", j+i) + " " + dataNode2 + "\n"
	}

	nChunks, _ = strconv.Atoi(library[tempBook.name].c3)
	for k = 0; k < nChunks; k++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", k+j+i) + " " + dataNode3 + "\n"
	}

	ioutil.WriteFile(fileName, []byte(buff), os.ModeAppend)
	//defer f.Close()
	log.Println("Accpeting proposal.")
	return &pb.Message{M: "A"}, nil
}

func (s *server) Proposal(ctx context.Context, in *pb.Message) (*pb.Message, error) {

	log.Println("Received proposal.")

	var i int
	var j int
	var k int

	flag := false

	proposal := strings.Split(in.GetM(), "**") //nombre**c1**c2**c3**total

	// proposal = [nombre, c1, c2, c3, total]
	// if c1 down -> proposal[1] = 0,
	//			proposal[2] = proposal[2] + proposal[1],
	//			proposal[3] = proposal[3],

	log.Println("Checking availability of datanodes.")

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

	log.Println("Writing log file.")

	//write
	fileName := "log.txt"
	if Exists(fileName) {
		_, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		_, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	buff := library[tempBook.name].name + " " + library[tempBook.name].parts + "\n"

	nChunks, _ := strconv.Atoi(library[tempBook.name].c1)
	for i = 0; i < nChunks; i++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", i) + " " + dataNode1 + "\n"
	}

	nChunks, _ = strconv.Atoi(library[tempBook.name].c2)
	for j = 0; j < nChunks; j++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", j+i) + " " + dataNode2 + "\n"
	}

	nChunks, _ = strconv.Atoi(library[tempBook.name].c3)
	for k = 0; k < nChunks; k++ {
		buff = buff + library[tempBook.name].name + "_parte_" + fmt.Sprintf("%d", k+j+i) + " " + dataNode3 + "\n"
	}

	ioutil.WriteFile(fileName, []byte(buff), os.ModeAppend)
	//defer f.Close()
	log.Println("Accpeting proposal.")
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
		log.Fatalf("[CONNECT TO DN] could not greet: %s --- %v", dataNode, err)
		return "0"
	}
	return r.GetM()
}

//ListenToClient listener
func ListenToClient(port string) {
	//--------------------------------------------------------------> Server1
	log.Println("Esperando solicitudes en puerto", port)
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

	go ListenToClient(":50051") // DataNode 1
	go ListenToClient(":50052") // DataNode 2
	go ListenToClient(":50053") // DataNode 3
	ListenToClient(":50054")    // Cliente Descargas

}
