package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	pb "lab2"

	"google.golang.org/grpc"
)

func uploadBook(fileName string, c pb.GreeterClient) {

	fileToBeChunked := "./" + fileName

	file, err := os.Open(fileToBeChunked)

	fmt.Println("File opened")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()
	const fileChunk = 250000
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// generacion de orden
		storeRequest := &pb.StoreRequest{
			ChunkPart:  strconv.FormatUint(i, 10),
			FileName:   fileName,
			ClientName: clientName,
			Chunk:      partBuffer,
			Part:       int32(totalPartsNum),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Hacer una consulta
		r, err := c.ReceiveChunk(ctx, storeRequest)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("%s", r)
	}
}

const (
	clientName = "clientUploader"
)

// Select selecciona aleatoreamente un dataNode a conectar
func Select() string {
	in := [3]string{"dist30:50051", "dist31:50052", "dist32:50053"}
	randomIndex := rand.Intn(len(in))
	pick := in[randomIndex]

	return pick
}

func main() {

	var book string
	fmt.Println("Ingresa nombre del libro que quieres subir: ")
	fmt.Scan(&book)

	// Set up a connection to the server.
	// Contact the server and print out its response.
	node := Select()
	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect to %s: %v\n", node, err)
		for {
			node = Select()
			conn, err = grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock())
			if err == nil {
				break
			} else {
				log.Printf("did not connect to %s: %v\n", node, err)
			}
		}
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	uploadBook(book, c)

}
