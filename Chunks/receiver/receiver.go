// Package main implements a server for Greeter service.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	pb "lab2"

	"google.golang.org/grpc"
)

type books struct {
	name   string
	parts  int32
	stored int32
}

var library map[string]*books

func storeInLibrary(book books) {
	if library[book.name] != nil {
		fmt.Printf("book++ %s\n", book.name)
		library[book.name].stored++
	} else {
		fmt.Printf("stored new book %s\n", book.name)
		library[book.name] = &book
	}
}

// constantes de los puertos
const (
	portSender = ":50051"
)

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// Recibe un chunk y lo guarda en un archivo en el disco
func (s *server) ReceiveChunk(ctx context.Context, in *pb.StoreRequest) (*pb.StoreReply, error) {
	log.Printf("Received: chunk 250kb. From: %v", in.GetClientName())

	// write to disk
	fileName := "./out/" + in.GetFileName() + "_part_" + in.GetChunkPart()
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tempBook := books{
		name:   in.GetFileName(),
		parts:  in.GetPart(),
		stored: 0,
	}

	storeInLibrary(tempBook)
	// write/save buffer to disk
	ioutil.WriteFile(fileName, in.GetChunk(), os.ModeAppend)
	//fmt.Println("Split to : ", fileName)

	if library[in.GetFileName()].parts == library[in.GetFileName()].stored {
		fmt.Println("ready to join book")
		joinFile(in.GetFileName(), library[in.GetFileName()].parts)
	}

	return &pb.StoreReply{Message: "Received chunk & stored in disk"}, nil
}

//JoinFile ensambla un archivo separado en chunks
func joinFile(fileName string, totalParts int32) {
	totalPartsNum := uint64(totalParts)

	// just for fun, let's recombine back the chunked files in a new file

	newFileName := "./restored/" + fileName
	_, err := os.Create(newFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set the newFileName file to APPEND MODE!!
	// open files r and w

	file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
	// defer file.Close()

	// just information on which part of the new file we are appending
	var writePosition int64 = 0

	for j := uint64(0); j < totalPartsNum; j++ {

		//read a chunk
		currentChunkFileName := "./out/" + fileName + "_part_" + strconv.FormatUint(j, 10)

		newFileChunk, err := os.Open(currentChunkFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer newFileChunk.Close()

		chunkInfo, err := newFileChunk.Stat()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// calculate the bytes size of each chunk
		// we are not going to rely on previous data and constant

		var chunkSize int64 = chunkInfo.Size()
		chunkBufferBytes := make([]byte, chunkSize)

		fmt.Println("Appending at position : [", writePosition, "] bytes")
		writePosition = writePosition + chunkSize

		// read into chunkBufferBytes
		reader := bufio.NewReader(newFileChunk)
		_, err = reader.Read(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
		// write/save buffer to disk
		//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

		_, err = file.Write(chunkBufferBytes)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file.Sync() //flush to disk

		// free up the buffer for next cycle
		// should not be a problem if the chunk size is small, but
		// can be resource hogging if the chunk size is huge.
		// also a good practice to clean up your own plate after eating

		chunkBufferBytes = nil // reset or empty our buffer

		//fmt.Println("Written ", n, " bytes")

		//fmt.Println("Recombining part [", j, "] into : ", newFileName)
	}

	// now, we close the newFileName
	file.Close()
}

//ListenToSender listener
func ListenToSender() {
	//--------------------------------------------------------------> Server1
	fmt.Print("Waitin for my Clientes...")
	lis, err := net.Listen("tcp", portSender)
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
	ListenToSender()

}
