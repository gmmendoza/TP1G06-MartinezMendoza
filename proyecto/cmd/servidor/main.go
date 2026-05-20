package main

import (
	"log"
	"net"
	"os"

	"sd-broadcast/internal/registro"
	"sd-broadcast/pkg/protocolo"
)

const puertoPorDefecto = "4000"

func main() {
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = puertoPorDefecto
	}

	escuchador, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		log.Fatalf("No se pudo iniciar el escuchador: %v", err)
	}
	defer escuchador.Close()

	log.Printf("Servidor de broadcast escuchando en :%s", puerto)

	registroClientes := registro.NuevoRegistro()

	for {
		conexion, err := escuchador.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}

		go manejarCliente(conexion, registroClientes)
	}
}

func manejarCliente(conexion net.Conn, registroClientes *registro.RegistroClientes) {
	defer conexion.Close()

	msgIdentificacion, err := protocolo.Decodificar(conexion)
	if err != nil {
		log.Printf("Error al leer identificación del cliente: %v", err)
		return
	}

	nombreCliente := msgIdentificacion.Emisor
	log.Printf("Cliente conectado: %s desde %s", nombreCliente, conexion.RemoteAddr())

	registroClientes.Agregar(nombreCliente, conexion)

	msgUnion := protocolo.NuevoMensaje("Sistema", nombreCliente+" se unió", "sistema")
	difundirMensaje(registroClientes, msgUnion, nombreCliente)

	defer func() {
		registroClientes.Eliminar(nombreCliente)
		msgDesconexion := protocolo.NuevoMensaje("Sistema", nombreCliente+" se desconectó", "sistema")
		difundirMensaje(registroClientes, msgDesconexion, nombreCliente)
		log.Printf("Cliente desconectado: %s", nombreCliente)
	}()

	for {
		mensaje, err := protocolo.Decodificar(conexion)
		if err != nil {
			break
		}

		log.Printf("[Servidor] Mensaje recibido de %s (Tipo: %s)", mensaje.Emisor, mensaje.Tipo)

		if mensaje.Tipo == "broadcast" {
			difundirMensaje(registroClientes, mensaje, nombreCliente)
		}
	}
}

func difundirMensaje(registroClientes *registro.RegistroClientes, mensaje protocolo.Mensaje, exceptoEmisor string) {
	conexiones := registroClientes.ObtenerConexiones()

	for _, conexion := range conexiones {
		err := protocolo.Codificar(conexion, mensaje)
		if err != nil {
			log.Printf("Error enviando broadcast por socket: %v", err)
		}
	}
}
