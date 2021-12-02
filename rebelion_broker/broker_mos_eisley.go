package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"

	pb "example.com/go-rebelion-grpc/rebelion"
)

const (
	port = ":50052"
)

var direcciones = [3]string{"localhost:50051", "localhost:50051", "localhost:50051"}

type server struct {
	pb.UnimplementedBrokerServer
}

func random_range(min int, max int) int {
	number := (rand.Intn(max-min+1) + min)
	return number
}

func (s *server) SolicitarIP(ctx context.Context, in *pb.Comando) (*pb.IP, error) {
	//Servidor aleatorio (para el broker)
	servidor := random_range(1, 3)
	//Direccion de servidores

	return &pb.IP{Direccion: direcciones[servidor-1]}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBrokerServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
