package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "example.com/go-rebelion-grpc/rebelion"
)

const (
	port = ":50052"
)

var direcciones = [3]string{"localhost:50051", "localhost:50053", "localhost:50054"}

type server struct {
	pb.UnimplementedBrokerServer
}

func random_range(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	number := (rand.Intn(max-min+1) + min)
	return number
}

func (s *server) SolicitarIP(ctx context.Context, in *pb.Comando) (*pb.IP, error) {
	//Servidor aleatorio (para el broker)
	nombre_planeta := in.GetPlaneta()
	vector_informante := in.GetVector()
	var arr_direcciones []string
	//Direccion de servidores

	//Obtener reloj de servidor 1
	conn, err := grpc.Dial(direcciones[0], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	//Conexion con Servidor 1
	c := pb.NewInformantesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	r, err := c.SolicitarRelojes(ctx, &pb.SolicitudR{Planeta: nombre_planeta})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	vector_s1 := r.GetVector()
	conn.Close()
	cancel()
	//Lo mismo conexion con Servidor 2 y 3

	conn, err = grpc.Dial(direcciones[1], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//Conexion con Servidor 2
	c = pb.NewInformantesClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	r, err = c.SolicitarRelojes(ctx, &pb.SolicitudR{Planeta: nombre_planeta})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	vector_s2 := r.GetVector()
	conn.Close()
	cancel()

	conn, err = grpc.Dial(direcciones[2], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//Conexion con Servidor 3
	c = pb.NewInformantesClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	r, err = c.SolicitarRelojes(ctx, &pb.SolicitudR{Planeta: nombre_planeta})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	vector_s3 := r.GetVector()
	conn.Close()
	cancel()

	//Comparar vectores
	if vector_s1[0] >= vector_informante[0] && vector_s1[1] >= vector_informante[1] && vector_s1[2] >= vector_informante[2] {
		arr_direcciones = append(arr_direcciones, direcciones[0])
	}

	if vector_s2[0] >= vector_informante[0] && vector_s2[1] >= vector_informante[1] && vector_s2[2] >= vector_informante[2] {
		arr_direcciones = append(arr_direcciones, direcciones[1])
	}

	if vector_s3[0] >= vector_informante[0] && vector_s3[1] >= vector_informante[1] && vector_s3[2] >= vector_informante[2] {
		arr_direcciones = append(arr_direcciones, direcciones[2])
	}

	//Repetir eso para los demas servidores

	//Escoger al azar
	servs := len(arr_direcciones)
	servidor := random_range(0, servs-1)

	return &pb.IP{Direccion: arr_direcciones[servidor]}, nil
}

func (s *server) GetNumberRebeldes(ctx context.Context, in *pb.SolicitudLeia) (*pb.Rebeldes, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()
	vector_informante := in.GetVector()
	var arr_direcciones []string

	//Obtener reloj de servidor 1
	conn, err := grpc.Dial(direcciones[0], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	//Conexion con Servidor 1
	c := pb.NewInformantesClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	r, err := c.SolicitarRelojes(ctx, &pb.SolicitudR{Planeta: nombre_planeta})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	vector_s1 := r.GetVector()
	conn.Close()
	cancel()
	//Lo mismo conexion con Servidor 2 y 3

	conn, err = grpc.Dial(direcciones[1], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//Conexion con Servidor 2
	c = pb.NewInformantesClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	r, err = c.SolicitarRelojes(ctx, &pb.SolicitudR{Planeta: nombre_planeta})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	vector_s2 := r.GetVector()
	conn.Close()
	cancel()

	conn, err = grpc.Dial(direcciones[2], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//Conexion con Servidor 3
	c = pb.NewInformantesClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	r, err = c.SolicitarRelojes(ctx, &pb.SolicitudR{Planeta: nombre_planeta})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	vector_s3 := r.GetVector()
	conn.Close()
	cancel()

	//Comparar vectores
	if vector_s1[0] >= vector_informante[0] && vector_s1[1] >= vector_informante[1] && vector_s1[2] >= vector_informante[2] {
		arr_direcciones = append(arr_direcciones, direcciones[0])
	}

	if vector_s2[0] >= vector_informante[0] && vector_s2[1] >= vector_informante[1] && vector_s2[2] >= vector_informante[2] {
		arr_direcciones = append(arr_direcciones, direcciones[1])
	}

	if vector_s3[0] >= vector_informante[0] && vector_s3[1] >= vector_informante[1] && vector_s3[2] >= vector_informante[2] {
		arr_direcciones = append(arr_direcciones, direcciones[2])
	}

	//Repetir eso para los demas servidores

	//Escoger al azar
	servs := len(arr_direcciones)
	servidor := arr_direcciones[random_range(0, servs-1)]

	//Preguntar a servidor

	conn, err = grpc.Dial(servidor, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//Conexion con Servidor
	c = pb.NewInformantesClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	rr, errr := c.SolicitarRebeldes(ctx, &pb.Solicitud{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad})
	if errr != nil {
		log.Fatalf("could not greet: %v", err)
	}
	defer conn.Close()
	defer cancel()

	return &pb.Rebeldes{Rebeldes: rr.GetRebeldes(), Vector: rr.GetVector(), Servidor: rr.GetServidor()}, nil
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
