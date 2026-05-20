package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"sd-broadcast/pkg/protocolo"
)

func main() {
	direccionServidor := os.Getenv("SERVIDOR")
	if direccionServidor == "" {
		direccionServidor = "localhost:4000"
	}

	nombre := os.Getenv("NOMBRE")
	if nombre == "" {
		fmt.Print("Ingrese su nombre: ")
		lector := bufio.NewReader(os.Stdin)
		nombreBytes, _, _ := lector.ReadLine()
		nombre = string(nombreBytes)
	}

	conexion, err := net.Dial("tcp", direccionServidor)
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor: %v", err)
	}
	defer conexion.Close()

	log.Printf("Conectado exitosamente al servidor %s", direccionServidor)

	msgIdentificacion := protocolo.NuevoMensaje(nombre, "", "identificacion")
	err = protocolo.Codificar(conexion, msgIdentificacion)
	if err != nil {
		log.Fatalf("Error al enviar el mensaje de identificación: %v", err)
	}

	go recibirMensajes(conexion)

	lectorTeclado := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	for {
		texto, err := lectorTeclado.ReadString('\n')
		if err != nil {
			log.Printf("Error al leer desde teclado: %v", err)
			break
		}

		texto = strings.TrimSpace(texto)
		if texto == "" {
			fmt.Print("> ")
			continue
		}

		if texto == "/salir" {
			break
		}

		msgBroadcast := protocolo.NuevoMensaje(nombre, texto, "broadcast")
		err = protocolo.Codificar(conexion, msgBroadcast)
		if err != nil {
			log.Printf("Error al enviar mensaje: %v", err)
			break
		}
		fmt.Print("> ")
	}

	log.Println("Cliente finalizado")
}

func recibirMensajes(conexion net.Conn) {
	for {
		mensaje, err := protocolo.Decodificar(conexion)
		if err != nil {
			log.Println("\n[Sistema] Desconectado del servidor.")
			return
		}

		timestampFormateado := mensaje.Timestamp.Format("15:04:05")
		fmt.Printf("\r[%s] %s: %s\n> ", timestampFormateado, mensaje.Emisor, mensaje.Contenido)
	}
}
