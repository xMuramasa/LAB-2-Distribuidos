// Package chunks implements a client for Greeter service.
package chunks

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	pb "lab2"

	"google.golang.org/grpc"
	//pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// constantes de puertos y nombres de instancias
const (
	address    = "dist29:50051"
	clientName = "Sender"
)

func createChunksForFile(fileName string, c pb.GreeterClient) {

	fileToBeChunked := "./" + fileName

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 250000 // 1 * (1 << 20) // 1 MB, change this to your requirement

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// write to disk
		fileName := "./out/" + fileName + "_part_" + strconv.FormatUint(i, 10)
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		//ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
		//fmt.Println("Split to : ", fileName)

		// generacion de orden
		storeRequest := &pb.StoreRequest{
			ChunkPart:  strconv.FormatUint(i, 10),
			FileName:   fileName,
			ClientName: clientName,
			Chunk:      string(partBuffer),
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

/*
	hacerOrden()
	hace pedidos automaticamente para un cliente desde una lista esperando un intervalo de tiempo
	Input: lista p, conexion c, tiempo de espera waitingTime
	returns: nada

func hacerOrden(p *list.List, c pb.GreeterClient, waitingTime int) {

	for {
		if p.Front() == nil {
			fmt.Println("No hay mas ordenes que enviar.")
			break

		} else {
			// waiting time
			time.Sleep(time.Duration(waitingTime) * time.Second)

			front := p.Front()
			itemI := Items(front.Value.(Items))

			tipo := "none"

			if itemI.prioridad == "0" {
				tipo = "Normal"
			} else if itemI.prioridad == "1" {
				tipo = "prioritario"
			} else {
				tipo = "retail"
			}

			// generacion de orden
			orden := &pb.OrderRequest{
				OrderID:      itemI.id,
				ProductName:  itemI.name,
				ProductValue: itemI.data[0],
				Src:          itemI.data[1],
				Dest:         itemI.data[2],
				Priority:     itemI.prioridad,
				ProductType:  tipo,
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			// Hacer una consulta
			r, err := c.MakeOrder(ctx, orden)
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			// mostrar codigo seguimiento por pantalla
			fmt.Print("\n<--------------- INFORMATION --------------->\n")
			log.Printf("Order Tracking Code: %s\n", r.GetMessage())
			fmt.Println("\n<--------------- INFORMATION --------------->")

			p.Remove(front)

		}
	}
}
*/
func main() {

	// Set up a connection to the server.
	// Contact the server and print out its response.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	createChunksForFile("Don_Quijote_de_la_Mancha-Cervantes_Miguel.pdf", c)
}
