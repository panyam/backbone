package connectors

import (
	"github.com/panyam/backbone/connectors/gorilla"
)

func CreateServer() Server {
	return gorilla.NewServer()
}
