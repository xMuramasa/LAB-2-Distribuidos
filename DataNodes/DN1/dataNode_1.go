// Package main implements a server for Greeter service.
package main

import (
	"container/list"
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
	clientName = "DATANODE1"
)

var algoritmo bool //true = centralizado; false: distribuido

var writing = list.New() // escrituras pendientes

// server is used to implement lab2.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type books struct {
	name   string
	parts  int32
	titles []string
	chunks []string
}

var storage map[string]*books

func storeInStorage(b books, title string, chunk string) {
	if storage[b.name] == nil {
		storage[b.name] = &b
	} else {
		storage[b.name].titles = append(storage[b.name].titles, title)
		storage[b.name].chunks = append(storage[b.name].chunks, chunk)
	}
}

func dataNodeProposal(ip string, mensaje string) bool {
	status := true
	log.Println("[DATANODE P] Proposal", mensaje)
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("[DATANODE P] did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.Greeting(ctx, &pb.Message{
		M: mensaje})
	if err != nil {
		log.Printf("[DATANODE P] could not greet: %v\n", err)
		status = false
	}
	return status
}

// RichardAgrawala resuleve colisiones posibles
func (s *server) RichardAgrawala(ctx context.Context, in *pb.Conflict) (*pb.Conflict, error) {
	log.Println("------------------------------->[MESAGE RECEIVED]")

	if writing.Front() == nil {
		return &pb.Conflict{
			ClientName: "dist31",
			Time:       "inf",
		}, nil
	}

	now := time.Now()
	t, _ := time.Parse(time.ANSIC, in.GetTime())
	e := writing.Front() // First element
	message := string(e.Value.(string))

	// tiempo del servidor es menor para escribir
	if now.Before(t) {
		//wait para escribir
		log.Println("[WRITE REQUEST] Proposal", message)
		conn, err := grpc.Dial("dist29:50051", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("[WRITE REQUEST] did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = c.WriteRequest(ctx, &pb.Message{
			M: message})
		if err != nil {
			log.Printf("[WRITE REQUEST] could not greet: %v\n", err)
		}
		writing.Remove(e)
	}

	return &pb.Conflict{
		ClientName: "dist31",
		Time:       "inf",
	}, nil

}

// CallRichardAgrawalla para escribir en NameNode
func CallRichardAgrawalla(ip string) string {
	log.Println("[WRITE REQUEST RECEIVE CHUNK]")

	now := time.Now()
	t := now.Format("Mon Jan _2 15:04:05 2006")

	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("[WRITE REQUEST RECEIVE CHUNK] did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r, err := c.RichardAgrawala(ctx, &pb.Conflict{
		ClientName: clientName,
		Time:       t,
	})
	if err != nil {
		log.Printf("[WRITE REQUEST RECEIVE CHUNK] could not greet: %v\n", err)
	}
	if r.GetTime() == "inf" {
		return "ok"
	}
	return ""
}

// Recibe un chunk y lo guarda en un archivo en el disco
func (s *server) ReceiveChunk(ctx context.Context, in *pb.StoreRequest) (*pb.StoreReply, error) {
	log.Println("------------------------------->[MESAGE RECEIVED]")

	log.Printf("Received: chunk 250kb. From: %v", in.GetClientName())
	//save chunk to memory
	tempBook := books{
		name:   in.GetFileName(),
		parts:  in.GetPart(),
		titles: []string{in.GetFileName() + "_part_" + in.GetChunkPart()},
		chunks: []string{string(in.GetChunk())},
	}
	storeInStorage(tempBook, tempBook.titles[0], tempBook.chunks[0])

	//if libro completo, realizar solicitud
	//if got all chunks -> debo avisar a nnode y repartir entre dnodes
	lenTitles := len(storage[tempBook.name].titles)
	lenChunks := len(storage[tempBook.name].chunks)
	if int(storage[tempBook.name].parts) == lenTitles && lenTitles == lenChunks {
		//prepare Proposal
		partsPerNode := int(int(storage[tempBook.name].parts) / 3)
		rest := int(storage[tempBook.name].parts) % 3
		message := fmt.Sprintf("%s**%d**%d**%d**%d",
			tempBook.name,
			partsPerNode+rest,
			partsPerNode,
			partsPerNode,
			tempBook.parts)

		c1 := 0
		c2 := 0
		c3 := 0
		var t []string
		var i int
		var j int
		log.Println(c1, c2, c3)

		// CENTRALIZADO
		if algoritmo == true {
			//send proposal
			log.Println("[RECEIVE CHUNK] Proposal", message)
			conn, err := grpc.Dial("dist29:50051", grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("[RECEIVE CHUNK] did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.Proposal(ctx, &pb.Message{
				M: message})
			if err != nil {
				log.Fatalf("[RECEIVE CHUNK] could not greet: %v", err)
			}

			//deal with response in r
			if r.GetM() == "A" {
				//propuesta aceptada -> repartir con message
				t = strings.Split(message, "**")
				c1, _ = strconv.Atoi(t[1])
				c2, _ = strconv.Atoi(t[2])
				c3, _ = strconv.Atoi(t[3])

			} else {
				//propesta rechazada -> repartir con r
				t = strings.Split(r.GetM(), "**")
				c1, _ = strconv.Atoi(t[1])
				c2, _ = strconv.Atoi(t[2])
				c3, _ = strconv.Atoi(t[3])
			}

		} else { // DISTRIBUIDO
			// Algo distribuido
			t = strings.Split(message, "**")
			c1, _ = strconv.Atoi(t[1])
			c2, _ = strconv.Atoi(t[2])
			c3, _ = strconv.Atoi(t[3])
			//porposal a dn2
			stat2 := dataNodeProposal("dist30:50053", message)

			//porposal a dn3
			stat3 := dataNodeProposal("dist32:50053", message)

			if stat2 == false {
				c1 = c1 + c2
				c2 = 0
			}
			if stat3 == false {
				c1 = c1 + c3
				c3 = 0
			}

			//enviar tiempo, msg a dn
			writing.PushBack(fmt.Sprintf("%s**%d**%d**%d**%d", tempBook.name, c1, c2, c3, tempBook.parts))

			calls := []string{"", ""}

			for {
				if calls[0] != "" && calls[1] != "" {
					break
				}
				calls[0] = CallRichardAgrawalla("dist30:50053") //dn1
				calls[1] = CallRichardAgrawalla("dist32:50053") //dn2
			}

			e := writing.Front()

			log.Println("[RECEIVE CHUNK] Proposal", message)
			conn, err := grpc.Dial("dist29:50051", grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("[RECEIVE CHUNK] did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewGreeterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err = c.WriteRequest(ctx, &pb.Message{
				M: string(e.Value.(string))})
			if err != nil {
				log.Fatalf("[RECEIVE CHUNK] could not greet: %v", err)
			}
			writing.Remove(e)
		}
		log.Println(c1, c2, c3)

		//ENVIO DE CHUNKS A LOS DN
		//soy dn1 guardo c1 chunks
		if c1 != 0 {
			for i = 0; i < c1; i++ {
				fileName := "./stored/" + in.GetFileName() + "_part_" + fmt.Sprintf("%d", i)
				_, err := os.Create(fileName)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				ioutil.WriteFile(fileName, []byte(storage[tempBook.name].chunks[i]), os.ModeAppend)
			}
		}

		// enviar c2 chunks a dn2
		if c2 != 0 {
			j = SendToDataNode(i, c2+i, "dist30:50053", tempBook.name)
		}

		// enviar c3 chunks a dn3
		if c3 != 0 {
			j = SendToDataNode(j, c3+j, "dist32:50053", tempBook.name)
		}
	}

	delete(storage, tempBook.name)
	return &pb.StoreReply{Message: "Received & stored chunk"}, nil
}

func (s *server) StoreChunk(ctx context.Context, in *pb.StoreRequest) (*pb.StoreReply, error) {
	log.Println("------------------------------->[MESAGE RECEIVED]")
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

	return &pb.StoreReply{Message: "Received chunk & stored in disk"}, nil
}

//SendToDataNode sends chunks to specified datanode at ip
func SendToDataNode(initalIt int, endIt int, ip string, bookName string) int {
	var j int
	for j = initalIt; j < endIt; j++ {
		conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("[SEND CHUNK] did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.StoreChunk(ctx, &pb.StoreRequest{
			ChunkPart:  fmt.Sprintf("%d", j),
			FileName:   bookName,
			ClientName: clientName,
			Chunk:      []byte(storage[bookName].chunks[j]),
			Part:       storage[bookName].parts})
		if err != nil {
			log.Fatalf("[SEND CHUNK] could not greet: %v", err)
		}
		log.Printf("[SEND CHUNK] Stored Chunk: %s", r)
	}
	return j
}

// RequestChunk envia un chunk a un cliente
func (s *server) RequestChunk(ctx context.Context, in *pb.BookRequest) (*pb.BookReply2, error) {
	log.Println("------------------------------->[MESAGE RECEIVED]")

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
	log.Println("Esperando solicitudes en puerto", puerto)
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
	storage = make(map[string]*books)

	var a string
	fmt.Println("Selecciona el tipo de algoritmo que deseas utilizar (numero): ")
	fmt.Println("[1] Centralizado \n[2] Distribuido")
	fmt.Print("Seleccion: ")
	fmt.Scan(&a)
	algoritmo = true
	if a == "2" {
		algoritmo = false
	}

	go ListenToClient(":50050") // NameNode
	go ListenToClient(":50051") // clientes descargas
	go ListenToClient(":50052") // clientes cargas
	go ListenToClient(":50053") // datanode 2
	ListenToClient(":50054")    // datanode 3

}
