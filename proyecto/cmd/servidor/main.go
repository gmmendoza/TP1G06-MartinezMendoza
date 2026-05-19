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

	// El servidor SIEMPRE inicia con Listen (Socket pasivo)
	escuchador, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		log.Fatalf("No se pudo iniciar el escuchador: %v", err)
	}
	defer escuchador.Close()

	log.Printf("Servidor de broadcast escuchando en :%s", puerto)

	// Inicializa el mapa o estructura de registro de clientes
	registroClientes := registro.NuevoRegistro()

	for {
		// Accept se bloquea esperando que un cliente remoto haga Dial
		conexion, err := escuchador.Accept()
		if err != nil {
			log.Printf("Error al aceptar conexión: %v", err)
			continue
		}

		// Una goroutine por cada conexión entrante para que sea concurrente
		go manejarCliente(conexion, registroClientes)
	}
}

func manejarCliente(conexion net.Conn, registroClientes *registro.RegistroClientes) {
	defer conexion.Close()

	// 1. Leer el mensaje JSON de identificación
	msgIdentificacion, err := protocolo.Decodificar(conexion)
	if err != nil {
		log.Printf("Error al leer identificación del cliente: %v", err)
		return
	}

	nombreCliente := msgIdentificacion.Emisor
	log.Printf("Cliente conectado: %s desde %s", nombreCliente, conexion.RemoteAddr())

	// 2. Registrar la conexión
	registroClientes.Agregar(nombreCliente, conexion)

	// 3. Avisar a los demás que se unió
	msgUnion := protocolo.NuevoMensaje("Sistema", nombreCliente+" se unió", "sistema")
	difundirMensaje(registroClientes, msgUnion, nombreCliente)

	// 4. Defer para cuando el socket se cierre (Desconexión limpia)
	defer func() {
		registroClientes.Eliminar(nombreCliente)
		msgDesconexion := protocolo.NuevoMensaje("Sistema", nombreCliente+" se desconectó", "sistema")
		difundirMensaje(registroClientes, msgDesconexion, nombreCliente)
		log.Printf("Cliente desconectado: %s", nombreCliente)
	}()

	// 5. Escuchar continuamente los mensajes de este cliente
	for {
		mensaje, err := protocolo.Decodificar(conexion)
		if err != nil {
			break // Salir si el cliente cerró el socket de forma abrupta o normal
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
		// Le enviamos el mensaje a todas las conexiones activas de la lista
		err := protocolo.Codificar(conexion, mensaje)
		if err != nil {
			log.Printf("Error enviando broadcast por socket: %v", err)
		}
	}
}
