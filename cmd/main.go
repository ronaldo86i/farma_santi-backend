package main

import (
	"context"
	"farma-santi_backend/internal/postgresql/routine"
	"farma-santi_backend/internal/server"
	"farma-santi_backend/internal/server/setup"
)

func main() {
	// Inicializar contenedor de dependencias, variables de entorno y conexión a base de datos
	setup.Init()

	// Inicializar rutinas del servidor
	routine.Init(context.Background())

	deps := setup.GetDependencies()

	// Inicializar el servidor HTTP
	httpServer := server.NewServer(deps.Handler)

	// Iniciar el servidor
	httpServer.Initialize()
}
