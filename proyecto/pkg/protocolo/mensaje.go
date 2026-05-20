package protocolo

import (
	"encoding/json"
	"io"
	"time"
)

type Mensaje struct {
	Emisor    string    `json:"emisor"`
	Contenido string    `json:"contenido"`
	Tipo      string    `json:"tipo"`
	Timestamp time.Time `json:"timestamp"`
}

func NuevoMensaje(emisor, contenido, tipo string) Mensaje {
	return Mensaje{
		Emisor:    emisor,
		Contenido: contenido,
		Tipo:      tipo,
		Timestamp: time.Now(),
	}
}

func Codificar(escritor io.Writer, mensaje Mensaje) error {
	codificador := json.NewEncoder(escritor)
	return codificador.Encode(mensaje)
}

func Decodificar(lector io.Reader) (Mensaje, error) {
	var mensaje Mensaje
	decodificador := json.NewDecoder(lector)
	err := decodificador.Decode(&mensaje)
	return mensaje, err
}
