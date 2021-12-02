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
	address = "localhost:50052"
)

var nombre_planeta string
var nombre_ciudad string
var nuevo_valor int
var nuevo_valor_ciudad string

type registros struct {
	ciudad             string
	reloj_vector       []int
	direccion_servidor string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var comando string
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
		fmt.Print(">>> ")
		scanner.Scan()
		//fmt.Scanln(&comando)
		comando = scanner.Text()
		fmt.Println(comando)

		s := strings.Fields(comando)
		fmt.Println(s)

		if s[0] != "0" && s[0] != "AddCity" && s[0] != "UpdateName" && s[0] != "UpdateNumber" && s[0] != "DeleteCity" {
			fmt.Println("Comando incorrecto imbecil >:c")
			continue
		}

		r, err := c.SolicitarIP(ctx, &pb.Comando{Comando: comando})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
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

			fmt.Println(r.GetVector())
		} else if s[0] == "UpdateName" {
			nombre_planeta = s[1]
			nombre_ciudad = s[2]
			nuevo_valor_ciudad = s[3]

			r, err := c.UpdateName(ctx, &pb.InfoUpdateName{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad, NuevoValor: nuevo_valor_ciudad})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}

			fmt.Println(r.GetVector())

		} else if s[0] == "UpdateNumber" {

		} else if s[0] == "DeleteCity" {

		}
	}
}
