package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "example.com/go-rebelion-grpc/rebelion"
	"google.golang.org/grpc"
)

const (
	address = "dist40:50052" //Conexion con broker
)

var nombre_planeta string
var nombre_ciudad string
var nuevo_valor int
var nuevo_valor_ciudad string

type registros struct {
	planeta            string
	reloj_vector       []int32
	direccion_servidor string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var comando string
	var vectores []registros //Reloj de vectores que guarda el informante
	var vector []int32

	opcion := 1

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//Conexion con Broker
	c := pb.NewBrokerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	fmt.Println("Informante Ahsoka Tano")
	fmt.Println("Introduzca su comando o 0 para salir: ")
	for opcion == 1 {
		nuevo_planeta := 1 //Flag para distinguir cuando se agrega un nuevo planeta
		fmt.Print(">>> ")
		scanner.Scan()
		comando = scanner.Text()

		s := strings.Fields(comando)

		if s[0] != "0" && s[0] != "AddCity" && s[0] != "UpdateName" && s[0] != "UpdateNumber" && s[0] != "DeleteCity" {
			fmt.Println("Comando incorrecto")
			continue
		}

		var pos_planeta int
		vector = []int32{0, 0, 0}
		for i, reloj := range vectores { //Buscar el reloj del planeta
			if reloj.planeta == s[1] {
				//Vector
				pos_planeta = i
				vector = reloj.reloj_vector
				nuevo_planeta = 0
			}
		}
		r, err := c.SolicitarIP(ctx, &pb.Comando{Comando: comando, Vector: vector, Planeta: s[1]})
		if err != nil {
			log.Fatalf("could not greet broker: %v", err)
		}

		direccion := r.GetDireccion()
		log.Printf("%s", direccion)

		//Conexion con servidores fulcrum
		conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewInformantesClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if s[0] == "0" {
			opcion = 0
		} else if s[0] == "AddCity" {
			nombre_planeta = s[1]
			nombre_ciudad = s[2]

			if len(s) == 3 {
				nuevo_valor = 0
			} else {
				nuevo_valor, _ = strconv.Atoi(s[3])
			}
			r, err := c.AddCity(ctx, &pb.Info{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad, NuevoValor: int32(nuevo_valor)})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			//Verificar si es un nuevo planeta o no

			if nuevo_planeta == 1 {
				//Append
				p := registros{planeta: nombre_planeta, reloj_vector: r.GetVector(), direccion_servidor: direccion}
				vectores = append(vectores, p)
			} else {
				//Actualizar cambios
				vectores[pos_planeta].reloj_vector = r.GetVector()
				vectores[pos_planeta].direccion_servidor = direccion
			}

			//Printear vector
			fmt.Println("Planeta ", nombre_planeta, " con reloj ", r.GetVector())

		} else if s[0] == "UpdateName" {
			nombre_planeta = s[1]
			nombre_ciudad = s[2]
			nuevo_valor_ciudad = s[3]

			r, err := c.UpdateName(ctx, &pb.InfoUpdateName{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad, NuevoValor: nuevo_valor_ciudad})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			//Actualizar cambios
			vectores[pos_planeta].reloj_vector = r.GetVector()
			vectores[pos_planeta].direccion_servidor = direccion

			//Printear vector
			fmt.Println("Planeta ", nombre_planeta, " con reloj ", r.GetVector())

		} else if s[0] == "UpdateNumber" {
			nombre_planeta = s[1]
			nombre_ciudad = s[2]
			nuevo_valor, _ = strconv.Atoi(s[3])

			r, err := c.UpdateNumber(ctx, &pb.Info{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad, NuevoValor: int32(nuevo_valor)})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			//Actualizar cambios
			vectores[pos_planeta].reloj_vector = r.GetVector()
			vectores[pos_planeta].direccion_servidor = direccion

			//Printear vector
			fmt.Println("Planeta ", nombre_planeta, " con reloj ", r.GetVector())

		} else if s[0] == "DeleteCity" {
			nombre_planeta = s[1]
			nombre_ciudad = s[2]

			r, err := c.DeleteCity(ctx, &pb.InfoDelete{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			//Actualizar cambios
			vectores[pos_planeta].reloj_vector = r.GetVector()
			vectores[pos_planeta].direccion_servidor = direccion

			//Printear vector
			fmt.Println("Planeta ", nombre_planeta, " con reloj ", r.GetVector())
		}
	}
}
