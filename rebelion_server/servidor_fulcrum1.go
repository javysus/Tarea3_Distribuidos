package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

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

	fmt.Println(vector_r)
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
		return &pb.Rebeldes{Rebeldes: int32(-2), Vector: []int32{0, 0, 0}, Servidor: int32(1)}, nil
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

	return &pb.Rebeldes{Rebeldes: int32(numero_rebeldes), Vector: vector_r, Servidor: int32(1)}, nil

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

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func merge() {

	tdr := time.Tick(90 * time.Second)

	for horaActual := range tdr {
		fmt.Println("La hora es", horaActual)
		conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		//Conexion con Broker
		c := pb.NewInformantesClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		r, err := c.Merge(ctx, &pb.Flag{Flag: "merge"})
		if err != nil {
			log.Fatalf("could not greet s2: %v", err)
		}

		logs_2 := r.GetListaLogs()

		/*
			for _, log := range logs_2 {
				fmt.Println("Reloj: ", log.Reloj)
				fmt.Println("Planeta: ", log.Planeta)
				fmt.Println("Logs: ", log.Logs)
			}
		*/
		conn2, err2 := grpc.Dial("localhost:50054", grpc.WithInsecure(), grpc.WithBlock())
		if err2 != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		//Conexion con Broker
		c2 := pb.NewInformantesClient(conn2)
		ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel2()

		r2, err2 := c2.Merge(ctx2, &pb.Flag{Flag: "merge"})
		if err2 != nil {
			log.Fatalf("could not greet s2: %v", err2)
		}

		logs_3 := r2.GetListaLogs()
		var planetas_s1 []string //Los planetas del servidor 1
		var comandos_finales []string

		for _, planeta := range vectores {
			planetas_s1 = append(planetas_s1, planeta.nombre_planeta)

			var ciudades_agregadas []string
			ciudadEncontrada := 0
			f, err := os.ReadFile("log_" + planeta.nombre_planeta + ".txt")
			if err != nil {
				log.Println(err)
			}
			lines := strings.Split(string(f), "\n")
			lines = lines[:len(lines)-1] //Comandos de planeta de servidor 1

			for _, comando := range lines {
				s := strings.Fields(comando)

				if s[0] == "UpdateName" {
					for i, nombre := range ciudades_agregadas {
						if s[2] == nombre {
							ciudades_agregadas[i] = s[3]
							break
						}
					}
				} else if s[0] == "DeleteCity" {
					for i, nombre := range ciudades_agregadas {
						if s[2] == nombre {
							ciudades_agregadas = remove(ciudades_agregadas, i)
							break
						}
					}
				}
				comandos_finales = append(comandos_finales, comando)
			}

			for _, log := range logs_2 {
				if log.Planeta == planeta.nombre_planeta {
					log_2 := log.Logs

					for _, comando := range log_2 {
						s := strings.Fields(comando)

						if s[0] == "AddCity" {
							for _, ciudad := range ciudades_agregadas {
								if ciudad == s[2] {
									ciudadEncontrada = 1
								}
							}
							if ciudadEncontrada == 0 {
								ciudades_agregadas = append(ciudades_agregadas, s[2])
								comandos_finales = append(comandos_finales, comando)
							}
						} else if s[0] == "UpdateName" {
							for _, ciudad := range ciudades_agregadas {
								if ciudad == s[3] {
									ciudadEncontrada = 1 //Si la ciudad ya se encuentra no se agrega a ciudades y no se agrega el comando a la lista de comandos finales
								}
							}
							if ciudadEncontrada == 0 {
								ciudades_agregadas = append(ciudades_agregadas, s[3])
								comandos_finales = append(comandos_finales, comando)
							}
						} else if s[0] == "DeleteCity" {
							for i, nombre := range ciudades_agregadas {
								if s[2] == nombre { //Se elimina
									ciudades_agregadas = remove(ciudades_agregadas, i)
									break
								}
							}
							comandos_finales = append(comandos_finales, comando)
						} else { //Caso de UpdateNumber
							comandos_finales = append(comandos_finales, comando)
						}
					}

					//Merge de relojes
					reloj_2 := log.Reloj
					planeta.vector[1] = reloj_2[1]
				}
			}

			ciudadEncontrada = 0

			for _, log := range logs_3 {
				if log.Planeta == planeta.nombre_planeta {
					log_3 := log.Logs

					for _, comando := range log_3 {
						s := strings.Fields(comando)

						if s[0] == "AddCity" {
							for _, ciudad := range ciudades_agregadas {
								if ciudad == s[2] {
									ciudadEncontrada = 1
								}
							}
							if ciudadEncontrada == 0 {
								ciudades_agregadas = append(ciudades_agregadas, s[2])
								comandos_finales = append(comandos_finales, comando)
							}
						} else if s[0] == "UpdateName" {
							for _, ciudad := range ciudades_agregadas {
								if ciudad == s[3] {
									ciudadEncontrada = 1 //Si la ciudad ya se encuentra no se agrega a ciudades y no se agrega el comando a la lista de comandos finales
								}
							}
							if ciudadEncontrada == 0 {
								ciudades_agregadas = append(ciudades_agregadas, s[3])
								comandos_finales = append(comandos_finales, comando)
							}
						} else if s[0] == "DeleteCity" {
							for i, nombre := range ciudades_agregadas {
								if s[2] == nombre { //Se elimina
									ciudades_agregadas = remove(ciudades_agregadas, i)
									break
								}
							}
							comandos_finales = append(comandos_finales, comando)
						} else { //Caso de UpdateNumber
							comandos_finales = append(comandos_finales, comando)
						}
					}

					//Merge de relojes
					reloj_3 := log.Reloj
					planeta.vector[2] = reloj_3[2]
				}
			}

			ciudadEncontrada = 0

		}

		//Agregar los planetas de s2 que no esten en el 1
		planetaNuevo := 1
		for _, planeta := range logs_2 {
			for _, planeta_s1 := range planetas_s1 {
				if planeta.Planeta == planeta_s1 {
					planetaNuevo = 0
					break
				}
			}

			if planetaNuevo == 1 {
				planetas_s1 = append(planetas_s1, planeta.Planeta)
				var ciudades_agregadas []string
				reloj_planeta := planeta.Reloj
				ciudadEncontrada := 0
				for _, comando := range planeta.Logs {
					s := strings.Fields(comando)

					if s[0] == "UpdateName" {
						for i, nombre := range ciudades_agregadas {
							if s[2] == nombre {
								ciudades_agregadas[i] = s[3]
								break
							}
						}
					} else if s[0] == "DeleteCity" {
						for i, nombre := range ciudades_agregadas {
							if s[2] == nombre {
								ciudades_agregadas = remove(ciudades_agregadas, i)
								break
							}
						}
					}
					comandos_finales = append(comandos_finales, comando)
				}

				for _, log := range logs_3 {
					if log.Planeta == planeta.Planeta {
						log_3 := log.Logs

						for _, comando := range log_3 {
							s := strings.Fields(comando)

							if s[0] == "AddCity" {
								for _, ciudad := range ciudades_agregadas {
									if ciudad == s[2] {
										ciudadEncontrada = 1
									}
								}
								if ciudadEncontrada == 0 {
									ciudades_agregadas = append(ciudades_agregadas, s[2])
									comandos_finales = append(comandos_finales, comando)
								}
							} else if s[0] == "UpdateName" {
								for _, ciudad := range ciudades_agregadas {
									if ciudad == s[3] {
										ciudadEncontrada = 1 //Si la ciudad ya se encuentra no se agrega a ciudades y no se agrega el comando a la lista de comandos finales
									}
								}
								if ciudadEncontrada == 0 {
									ciudades_agregadas = append(ciudades_agregadas, s[3])
									comandos_finales = append(comandos_finales, comando)
								}
							} else if s[0] == "DeleteCity" {
								for i, nombre := range ciudades_agregadas {
									if s[2] == nombre { //Se elimina
										ciudades_agregadas = remove(ciudades_agregadas, i)
										break
									}
								}
								comandos_finales = append(comandos_finales, comando)
							} else { //Caso de UpdateNumber
								comandos_finales = append(comandos_finales, comando)
							}
						}

						//Merge de relojes
						reloj_3 := log.Reloj
						reloj_planeta[2] = reloj_3[2]
					}
				}
				ciudadEncontrada = 0

				//Agregar al struct del s1

				p := reloj_vector{nombre_planeta: planeta.Planeta, vector: reloj_planeta}
				vectores = append(vectores, p)

			}
		}

		//Revisar planetas nuevos de s3 (que no estaban ni en s1 ni s2)
		planetaNuevo = 1
		for _, planeta := range logs_3 {
			for _, planeta_s1 := range planetas_s1 {
				if planeta.Planeta == planeta_s1 {
					planetaNuevo = 0
					break
				}
			}

			if planetaNuevo == 1 {
				reloj_planeta := planeta.Reloj //[0,0,algo]
				for _, comando := range planeta.Logs {
					comandos_finales = append(comandos_finales, comando)
				}

				//Agregar al struct del s1
				p := reloj_vector{nombre_planeta: planeta.Planeta, vector: reloj_planeta}
				vectores = append(vectores, p)
			}
		}

		fmt.Println("Comandos finales:", comandos_finales)
	}

}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterInformantesServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())
	go merge()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
