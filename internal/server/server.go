package server

import (
	"farma-santi_backend/internal/core/util"
	"farma-santi_backend/internal/server/setup"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/exp/slog"
	"log"
	"os"
	"strings"
	"time"
)

var httpPort = "8080"

type Server struct {
	handlers setup.Handler
}

func NewServer(
	handlers setup.Handler,
) *Server {
	return &Server{
		handlers,
	}
}

func (s *Server) startServer() {
	app := fiber.New(fiber.Config{
		BodyLimit:             20 << 23,
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           30 * time.Second,
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		Prefork:               false,
		AppName:               "Farma Santi Backend",
		ErrorHandler:          util.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://localhost:4173,http://127.0.0.1:5173",
		AllowHeaders: strings.Join([]string{
			fiber.HeaderOrigin,
			fiber.HeaderContentType,
			fiber.HeaderAuthorization,
			fiber.HeaderXDownloadOptions,
			fiber.HeaderReferrerPolicy,
		}, ","),
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Use(logger.New(logger.Config{
		Format: "${ip} - ${method} ${path} - ${status} - ${latency}\n",
	}))

	s.initEndPoints(app)

	serverAddr := fmt.Sprintf(":%s", httpPort)
	slog.Info("🚀 Servidor iniciado", "url", "http://localhost"+serverAddr)

	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Error al iniciar el servidor Fiber: %v", err)
	}
}

func (s *Server) Initialize() {
	portFromEnv := os.Getenv("HTTP_PORT")
	if portFromEnv != "" {
		httpPort = portFromEnv
	} else {
		slog.Info("Puerto HTTP no definido en .env, usando puerto por defecto", "port", httpPort)
	}
	s.startServer()
}
