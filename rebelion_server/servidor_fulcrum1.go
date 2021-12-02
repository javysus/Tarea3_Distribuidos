package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	pb "example.com/go-rebelion-grpc/rebelion"
)

const (
	port = ":50051"
)

type reloj_vector struct {
	nombre_planeta string
	vector         []int32
}

var posicion int
var vectores []reloj_vector //Arreglo de struct reloj_vector, para guardar el planeta y su reloj de vector asociado

type server struct {
	pb.UnimplementedInformantesServer
}

func modificarReloj(planeta string) []int32 {
	nuevoPlaneta := 1    //Para saber si el reloj se debe agregar o modificar
	var vector_r []int32 //Vector modificado o creado para retornar
	for i, reloj := range vectores {
		if reloj.nombre_planeta == planeta { //El planeta ya tiene un reloj
			nuevoPlaneta = 0
			vectores[i].vector = []int32{reloj.vector[0] + 1, reloj.vector[1], reloj.vector[2]} //Se modifica
			vector_r = vectores[i].vector
		}
	}

	if nuevoPlaneta == 1 { //Si es un nuevo planeta se agrega
		p := reloj_vector{nombre_planeta: planeta, vector: []int32{1, 0, 0}}
		vectores = append(vectores, p)
		vector_r = p.vector
	}
	return vector_r
}

func (s *server) AddCity(ctx context.Context, in *pb.Info) (*pb.Respuesta, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()
	nuevo_valor := strconv.Itoa(int(in.GetNuevoValor()))

	f, err := os.OpenFile(nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor + "\n"); err != nil {
		log.Println(err)
	}

	f, err = os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("AddCity " + nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor + "\n"); err != nil {
		log.Println(err)
	}

	vector_r := modificarReloj(nombre_planeta)
	return &pb.Respuesta{Vector: vector_r}, nil
}

func (s *server) UpdateName(ctx context.Context, in *pb.InfoUpdateName) (*pb.Respuesta, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()
	nuevo_valor := in.GetNuevoValor()

	f, err := os.ReadFile(nombre_planeta + ".txt")
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(string(f), "\n")
	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			line_data := strings.Fields(line)
			lines[i] = nombre_planeta + " " + nuevo_valor + " " + line_data[2]
		}
	}

	a, er := os.OpenFile(nombre_planeta+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		log.Println(err)
	}

	for _, line := range lines {
		if _, err := a.WriteString(line + "\n"); err != nil {
			log.Println(err)
		}
	}
	defer a.Close()

	a, er = os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer a.Close()
	if _, err := a.WriteString("UpdateName " + nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor + "\n"); err != nil {
		log.Println(err)
	}

	vector_r := modificarReloj(nombre_planeta)
	return &pb.Respuesta{Vector: vector_r}, nil

}

func (s *server) UpdateNumber(ctx context.Context, in *pb.Info) (*pb.Respuesta, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()
	nuevo_valor := in.GetNuevoValor()

	f, err := os.ReadFile(nombre_planeta + ".txt")
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(string(f), "\n")
	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			lines[i] = nombre_planeta + " " + nombre_ciudad + " " + strconv.Itoa(int(nuevo_valor))
		}
	}

	a, er := os.OpenFile(nombre_planeta+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		log.Println(err)
	}

	for _, line := range lines {
		if _, err := a.WriteString(line + "\n"); err != nil {
			log.Println(err)
		}
	}
	defer a.Close()

	a, er = os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer a.Close()
	if _, err := a.WriteString("UpdateNumber " + nombre_planeta + " " + nombre_ciudad + " " + strconv.Itoa(int(nuevo_valor)) + "\n"); err != nil {
		log.Println(err)
	}

	vector_r := modificarReloj(nombre_planeta)
	return &pb.Respuesta{Vector: vector_r}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterInformantesServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
