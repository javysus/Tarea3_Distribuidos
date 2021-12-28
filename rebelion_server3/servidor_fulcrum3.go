package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	pb "example.com/go-rebelion-grpc/rebelion"
)

const (
	port = ":50054"
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

func remove(s []reloj_vector, i int) []reloj_vector {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func modificarReloj(planeta string) []int32 {
	nuevoPlaneta := 1    //Para saber si el reloj se debe agregar o modificar
	var vector_r []int32 //Vector modificado o creado para retornar
	for i, reloj := range vectores {
		if reloj.nombre_planeta == planeta { //El planeta ya tiene un reloj
			nuevoPlaneta = 0
			vectores[i].vector = []int32{reloj.vector[0], reloj.vector[1], reloj.vector[2] + 1} //Se modifica
			vector_r = vectores[i].vector
		}
	}

	if nuevoPlaneta == 1 { //Si es un nuevo planeta se agrega
		p := reloj_vector{nombre_planeta: planeta, vector: []int32{0, 0, 1}}
		vectores = append(vectores, p)
		vector_r = p.vector
	}
	return vector_r
}

func (s *server) PropagarCambios(ctx context.Context, in *pb.InfoActualizada) (*pb.Flag, error) {
	comandos := in.GetComandosFinales()
	data_actualizada := in.GetDataActualizada()

	vectores_copia := vectores
	for _, reloj := range vectores_copia {
		//Eliminar del struct
		vectores = remove(vectores, 0)
		//Eliminar archivo de planeta y su log, y volver a crearlo
		e := os.Remove(reloj.nombre_planeta + ".txt")
		if e != nil {
			log.Fatal(e)
		}

		e = os.Remove("log_" + reloj.nombre_planeta + ".txt")
		if e != nil {
			log.Fatal(e)
		}
	}

	//Actualizar relojes
	for _, data := range data_actualizada {
		v := reloj_vector{nombre_planeta: data.Planeta, vector: data.Reloj}
		vectores = append(vectores, v)

		f, err := os.OpenFile("log_"+data.Planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(data.Planeta, ": ", data.Reloj)
		f.Close()
	}
	for _, comando := range comandos {
		s := strings.Fields(comando)

		if s[0] == "AddCity" {
			f, err := os.OpenFile(s[1]+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			if _, err := f.WriteString(s[1] + " " + s[2] + " " + s[3] + "\n"); err != nil {
				log.Println(err)
			}
			f.Close()

		} else if s[0] == "UpdateNumber" {
			f, err := os.ReadFile(s[1] + ".txt")
			if err != nil {
				log.Println(err)
			}

			lines := strings.Split(string(f), "\n")
			lines = lines[:len(lines)-1]

			for i, line := range lines {
				if strings.Contains(line, s[2]) {
					lines[i] = s[1] + " " + s[2] + " " + s[3]
				}
			}

			e := os.Remove(s[1] + ".txt")
			if e != nil {
				log.Fatal(e)
			}

			a, er := os.OpenFile(s[1]+".txt", os.O_CREATE|os.O_WRONLY, 0644)
			if er != nil {
				log.Println(err)
			}

			for _, line := range lines {
				if _, err := a.WriteString(line + "\n"); err != nil {
					log.Println(err)
				}
			}
			a.Close()

		} else if s[0] == "UpdateName" {
			nombre_planeta := s[1]
			nombre_ciudad := s[2]
			nuevo_valor := s[3]

			f, err := os.ReadFile(nombre_planeta + ".txt")
			if err != nil {
				log.Println(err)
			}

			lines := strings.Split(string(f), "\n")
			lines = lines[:len(lines)-1]
			for i, line := range lines {
				if strings.Contains(line, nombre_ciudad) {
					line_data := strings.Fields(line)
					lines[i] = nombre_planeta + " " + nuevo_valor + " " + line_data[2]
				}
			}

			e := os.Remove(nombre_planeta + ".txt")
			if e != nil {
				log.Fatal(e)
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
			a.Close()

		} else if s[0] == "DeleteCity" {
			nombre_planeta := s[1]
			nombre_ciudad := s[2]

			f, err := os.ReadFile(nombre_planeta + ".txt")
			if err != nil {
				log.Println(err)
			}

			lines := strings.Split(string(f), "\n")
			lines = lines[:len(lines)-1]
			var newlines []string
			for i, line := range lines {
				if strings.Contains(line, nombre_ciudad) {
					continue
				}
				newlines = append(newlines, lines[i])
			}

			e := os.Remove(nombre_planeta + ".txt")
			if e != nil {
				log.Fatal(e)
			}

			a, er := os.OpenFile(nombre_planeta+".txt", os.O_CREATE|os.O_WRONLY, 0644)
			if er != nil {
				log.Println(err)
			}

			for _, line := range newlines {
				if _, err := a.WriteString(line + "\n"); err != nil {
					log.Println(err)
				}
			}
			a.Close()
		}
	}

	return &pb.Flag{Flag: "Cambios propagados"}, nil

}

func (s *server) Merge(ctx context.Context, in *pb.Flag) (*pb.InfoMerge, error) {

	var infos []*pb.DataMerge

	for _, reloj := range vectores {

		//Leer archivo
		f, err := os.ReadFile("log_" + reloj.nombre_planeta + ".txt")
		if err != nil {
			log.Println(err)
		}
		lines := strings.Split(string(f), "\n")
		lines = lines[:len(lines)-1]

		//Agregar a nuestros arreglos
		/*planetas = append(planetas, reloj.nombre_planeta)
		relojes = append(relojes, reloj.vector)
		logs = append(logs, lines)*/

		//inf := DataMerge{nombre_planeta: reloj.nombre_planeta, vector: reloj.vector, logs: lines}
		inf := &pb.DataMerge{Reloj: reloj.vector, Planeta: reloj.nombre_planeta, Logs: lines}
		infos = append(infos, inf)
	}

	return &pb.InfoMerge{ListaLogs: infos}, nil
}

//Funcion para retornar el reloj asociado al planeta solicitado por el Broker
func (s *server) SolicitarRelojes(ctx context.Context, in *pb.SolicitudR) (*pb.Respuesta, error) {
	nombre_planeta := in.GetPlaneta()
	nuevo_planeta := 1
	var vector_r []int32

	for _, reloj := range vectores { //Recorrer el arreglo para encontrar el vector del planeta
		if reloj.nombre_planeta == nombre_planeta {
			vector_r = reloj.vector
			nuevo_planeta = 0
		}
	}

	if nuevo_planeta == 1 { //Nuevo planeta asi que no tiene cambios
		vector_r = []int32{0, 0, 0}
	}

	return &pb.Respuesta{Vector: vector_r}, nil
}

func (s *server) SolicitarRebeldes(ctx context.Context, in *pb.Solicitud) (*pb.Rebeldes, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()
	numero_rebeldes := -1
	var vector_r []int32
	//Leer archivo
	f, err := os.ReadFile(nombre_planeta + ".txt")
	if err != nil {
		return &pb.Rebeldes{Rebeldes: int32(-2), Vector: []int32{0, 0, 0}, Servidor: int32(3)}, nil
	}

	lines := strings.Split(string(f), "\n")
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			line_data := strings.Fields(line)
			numero_rebeldes, _ = strconv.Atoi(line_data[2])
		}
	}

	for _, reloj := range vectores { //Recorrer el arreglo para encontrar el vector del planeta, asumiendo que existe
		if reloj.nombre_planeta == nombre_planeta {
			fmt.Println("reloj.vector:", reloj.vector)
			vector_r = reloj.vector
		}
	}

	return &pb.Rebeldes{Rebeldes: int32(numero_rebeldes), Vector: vector_r, Servidor: int32(3)}, nil

}

func (s *server) AddCity(ctx context.Context, in *pb.Info) (*pb.Respuesta, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()
	nuevo_valor := strconv.Itoa(int(in.GetNuevoValor()))

	f, err := os.OpenFile(nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if _, err := f.WriteString(nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor + "\n"); err != nil {
		log.Println(err)
	}
	f.Close()

	f, err = os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if _, err := f.WriteString("AddCity " + nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor + "\n"); err != nil {
		log.Println(err)
	}
	f.Close()

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
	lines = lines[:len(lines)-1]
	fmt.Println(lines)
	fmt.Println(len(lines))
	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			line_data := strings.Fields(line)
			fmt.Println("Esto es el line_data:", line_data)
			lines[i] = nombre_planeta + " " + nuevo_valor + " " + line_data[2]
		}
	}
	fmt.Println(lines)

	e := os.Remove(nombre_planeta + ".txt")
	if e != nil {
		log.Fatal(e)
	}

	a, er := os.OpenFile(nombre_planeta+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		log.Println(err)
	}

	for _, line := range lines {
		fmt.Println("Linea a escribir: ", line)
		if _, err := a.WriteString(line + "\n"); err != nil {
			log.Println(err)
		}
	}
	a.Close()

	b, erro := os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if erro != nil {
		log.Println(err)
	}
	if _, err := b.WriteString("UpdateName " + nombre_planeta + " " + nombre_ciudad + " " + nuevo_valor + "\n"); err != nil {
		log.Println(err)
	}

	b.Close()

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
	lines = lines[:len(lines)-1]
	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			lines[i] = nombre_planeta + " " + nombre_ciudad + " " + strconv.Itoa(int(nuevo_valor))
		}
	}

	e := os.Remove(nombre_planeta + ".txt")
	if e != nil {
		log.Fatal(e)
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
	a.Close()

	b, erro := os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if erro != nil {
		log.Println(err)
	}
	if _, err := b.WriteString("UpdateNumber " + nombre_planeta + " " + nombre_ciudad + " " + strconv.Itoa(int(nuevo_valor)) + "\n"); err != nil {
		log.Println(err)
	}

	b.Close()

	vector_r := modificarReloj(nombre_planeta)
	return &pb.Respuesta{Vector: vector_r}, nil
}

func (s *server) DeleteCity(ctx context.Context, in *pb.InfoDelete) (*pb.Respuesta, error) {
	nombre_planeta := in.GetNombrePlaneta()
	nombre_ciudad := in.GetNombreCiudad()

	f, err := os.ReadFile(nombre_planeta + ".txt")
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(string(f), "\n")
	lines = lines[:len(lines)-1]
	var newlines []string
	for i, line := range lines {
		if strings.Contains(line, nombre_ciudad) {
			continue
		}
		newlines = append(newlines, lines[i])
	}

	e := os.Remove(nombre_planeta + ".txt")
	if e != nil {
		log.Fatal(e)
	}

	a, er := os.OpenFile(nombre_planeta+".txt", os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		log.Println(err)
	}

	for _, line := range newlines {
		if _, err := a.WriteString(line + "\n"); err != nil {
			log.Println(err)
		}
	}
	a.Close()

	b, erro := os.OpenFile("log_"+nombre_planeta+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if erro != nil {
		log.Println(err)
	}
	if _, err := b.WriteString("DeleteCity " + nombre_planeta + " " + nombre_ciudad + "\n"); err != nil {
		log.Println(err)
	}

	b.Close()

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
