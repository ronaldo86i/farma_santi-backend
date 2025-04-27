package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

// DBInstance Instancia global de la conexión
var DBInstance = DB{}

// Connection establece la conexión con la base de datos Postgres usando pgxpool
func Connection() error {
	// Leer variables de entorno
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	timezone := os.Getenv("DB_TIMEZONE")
	sslMode := "disable"

	// Validación
	if host == "" || user == "" || password == "" || dbname == "" || port == "" || timezone == "" {
		return fmt.Errorf("una o más variables de entorno están vacías para inicializar la base de datos")
	}

	// Cadena de conexión
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		user, password, host, port, dbname, sslMode, timezone)

	// Configuración del pool
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error al parsear la cadena de conexión: %w", err)
	}

	// Configuraciones opcionales del pool
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.HealthCheckPeriod = time.Minute

	// Conexión
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
		return err
	}

	// Verificación de conexión
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
		return err
	}

	// Asignar pool a la instancia global
	DBInstance.Pool = pool

	fmt.Println("Conexión exitosa a Postgres.")
	return nil
}
