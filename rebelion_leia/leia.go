package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "example.com/go-rebelion-grpc/rebelion"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

type registros struct {
	planeta      string
	reloj_vector []int32
	servidor     int32
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
	fmt.Println("Leia Organa")
	fmt.Println("Introduzca su comando o 0 para salir: ")

	for opcion == 1 {
		nuevo_planeta := 1 //Flag para distinguir si Leia ha leido del planeta antes
		fmt.Print(">>> ")
		scanner.Scan()
		comando = scanner.Text()
		fmt.Println(comando)

		s := strings.Fields(comando)
		fmt.Println(s)

		if s[0] != "0" && s[0] != "GetNumberRebelds" {
			fmt.Println("Comando incorrecto imbecil >:c")
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

		nombre_planeta := s[1]
		nombre_ciudad := s[2]
		r, err := c.GetNumberRebeldes(ctx, &pb.SolicitudLeia{NombrePlaneta: nombre_planeta, NombreCiudad: nombre_ciudad, Vector: vector})
		if err != nil {
			log.Fatalf("could not greet broker: %v", err)
		}

		//Actualizar datos

		//Verificar si planeta esta en el registros
		reloj_planeta := r.GetVector()
		if nuevo_planeta == 1 {
			//Append
			p := registros{planeta: nombre_planeta, reloj_vector: reloj_planeta, servidor: r.GetServidor()}
			vectores = append(vectores, p)
		} else {
			//Actualizar cambios
			vectores[pos_planeta].reloj_vector = reloj_planeta
			vectores[pos_planeta].servidor = r.GetServidor()
		}

		cant_rebeldes := int(r.GetRebeldes())
		//Mostrar la cantidad de rebeldes?

		if cant_rebeldes == -1 {
			fmt.Println("La ciudad no existe")
		} else if cant_rebeldes == -2 {
			fmt.Println("El planeta no existe")
		} else {
			fmt.Println("Numero de rebeldes: ", cant_rebeldes)
		}
	}
}
