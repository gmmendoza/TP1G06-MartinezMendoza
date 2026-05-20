package registro

import (
	"net"
	"sync"
)

type RegistroClientes struct {
	mu       sync.RWMutex
	clientes map[string]net.Conn
}

func NuevoRegistro() *RegistroClientes {
	return &RegistroClientes{
		clientes: make(map[string]net.Conn),
	}
}

func (r *RegistroClientes) Agregar(nombre string, conexion net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clientes[nombre] = conexion
}

func (r *RegistroClientes) Eliminar(nombre string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clientes, nombre)
}

func (r *RegistroClientes) ObtenerConexiones() []net.Conn {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conexiones := make([]net.Conn, 0, len(r.clientes))
	for _, conn := range r.clientes {
		conexiones = append(conexiones, conn)
	}
	return conexiones
}

func (r *RegistroClientes) Cantidad() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clientes)
}

func (r *RegistroClientes) Nombres() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nombres := make([]string, 0, len(r.clientes))
	for nombre := range r.clientes {
		nombres = append(nombres, nombre)
	}
	return nombres
}
