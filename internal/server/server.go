package server

import (
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

type Server struct {
	authHandler      port.AuthHandler
	rolHandler       port.RolHandler
	usuarioHandler   port.UsuarioHandler
	categoriaHandler port.CategoriaHandler
	proveedorHandler port.ProveedorHandler
}

func NewServer(
	authHandler port.AuthHandler,
	rolHandler port.RolHandler,
	usuarioHandler port.UsuarioHandler,
	categoriaHandler port.CategoriaHandler,
	proveedorHandler port.ProveedorHandler,
) *Server {
	return &Server{
		authHandler:      authHandler,
		rolHandler:       rolHandler,
		usuarioHandler:   usuarioHandler,
		categoriaHandler: categoriaHandler,
		proveedorHandler: proveedorHandler,
	}
}

// Variables globales para la configuración y el estado del servidor
var (
	httpPort string // Almacena el puerto HTTP
	errEnv   error  // Almacena errores relacionados con la carga del archivo .env
)

// InitEnv carga las variables de entorno desde el archivo .env.
func InitEnv(filenames ...string) {
	// Cargar variables de entorno
	errEnv = godotenv.Load(filenames...) // Carga las variables de entorno desde el archivo especificado
	if errEnv != nil {
		log.Println(errEnv)                          // Registra cualquier error al cargar el archivo .env
		log.Fatal("Error al cargar el archivo .env") // Termina la ejecución si hay un error
	}
	httpPort = os.Getenv("HTTP_PORT") // Obtiene el puerto HTTP de las variables de entorno
}

func (s *Server) startServer() {

	// Configura una nueva aplicación Fiber
	app := fiber.New(fiber.Config{
		BodyLimit:             20 << 23,         // Establece el límite del cuerpo de la solicitud a 20 MB
		ReadTimeout:           30 * time.Second, // Tiempo de espera de lectura
		WriteTimeout:          30 * time.Second, // Tiempo de espera de escritura
		IdleTimeout:           30 * time.Second, // Tiempo de espera inactivo
		DisableStartupMessage: true,             // Desactiva el mensaje de inicio
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		Prefork:               false,
		AppName:               "Farma Santi Backend",
	})

	// Configuración de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://127.0.0.1:5173, http://192.168.0.28:5173",
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

	// Definir rutas
	s.initEndPoints(app)

	// Iniciar el servidor en el puerto especificado
	serverAddress := fmt.Sprintf(":%s", httpPort)

	// Inicia el servidor Fiber
	log.Println("🚀 Servidor iniciado en http://localhost" + serverAddress)
	if err := app.Listen(":" + httpPort); err != nil {
		log.Fatalf("Error al iniciar el servidor Fiber: %v", err) // Registra y termina si hay un error al iniciar el servidor
	}
}

// initDB establece la conexión a la base de datos y realiza migraciones.
func (s *Server) initDB() {
	if err := database.Connection(); err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	//if err := db.Migration(); err != nil {
	//	log.Fatalf("Error al realizar la migración de la base de datos: %v", err)
	//}
}

// Initialize Inicializa el servidor
func (s *Server) Initialize() {
	InitEnv("./.env")
	s.initDB()
	s.startServer()
}
