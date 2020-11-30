package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "lab2"

	"google.golang.org/grpc"
)

const (
	clientName = "clientDownloader"
)

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

//downloadBookInfo descarga informacion de un libro desde el namenode
func downloadBookInfo(fileName string, c pb.GreeterClient) string {
	// generacion de bookRequest para nameNode
	storeRequest := &pb.BookRequest{
		BookNamePart: fileName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Hacer una consulta
	r, err := c.RequestBook(ctx, storeRequest)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return r.GetLocations()
}

//downloadChunks descarga chunks desde el datanode correspondiente en base a la coneccion c
func downloadChunks(fileName string, c pb.GreeterClient) {
	// generacion de bookRequest para dataNode
	storeRequest := &pb.BookRequest{
		BookNamePart: fileName,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Hacer una consulta
	r, err := c.RequestChunk(ctx, storeRequest)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// write to disk
	file := "./downloads/" + fileName
	_, err = os.Create(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ioutil.WriteFile(file, r.GetChunk(), os.ModeAppend)
}

//JoinFile ensambla un archivo separado en chunks, asume que todas las partes estan descargadas
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
		currentChunkFileName := "./downloads/" + fileName + "_part_" + strconv.FormatUint(j, 10)

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

		//fmt.Println("Appending at position : [", writePosition, "] bytes")
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

//connectAndDownload se conecta a un datanode y descarga los chunks
func connectAndDownload(fileName string, host string, nChunks int, start int) int {
	//conexion con dataNode
	var i int
	conn, err := grpc.Dial(host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	for i = start; i < nChunks; i++ {
		downloadChunks(fileName+"_part_"+fmt.Sprint(i), c)
	}
	return i
}

func main() {

	var book string
	fmt.Println("Ingresa nombre del libro que quieres descargar: ")
	fmt.Scan(&book)

	//conexion con nameNodes := strings.Split("a,b,c", ",")
	// Set up a connection to the server. Contact the server and print out its response.
	conn, err := grpc.Dial("dist29:50054", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	bookData := strings.Split(downloadBookInfo(book, c), "**")

	c1, _ := strconv.Atoi(bookData[1])
	c2, _ := strconv.Atoi(bookData[2])
	c3, _ := strconv.Atoi(bookData[3])
	totalParts, _ := strconv.Atoi(bookData[4])

	i := connectAndDownload(book, "dist31:50051", c1, 0)
	j := connectAndDownload(book, "dist30:50051", c2+i, i)
	_ = connectAndDownload(book, "dist32:50051", c3+j, j)

	// union archivo descargado
	joinFile(book, int32(totalParts))
	fmt.Println("Libro Descargado")

}
