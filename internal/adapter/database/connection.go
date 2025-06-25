package database

import (
	"context"
	"farma-santi_backend/internal/slog_logger"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"golang.org/x/exp/slog"
)

type DB struct {
	Pool *pgxpool.Pool
}

var DBInstance = DB{}

func Connection() error {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	timezone := os.Getenv("DB_TIMEZONE")
	sslMode := "disable"

	if host == "" || user == "" || password == "" || dbname == "" || port == "" || timezone == "" {
		return fmt.Errorf("una o más variables de entorno están vacías para inicializar la base de datos")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
		user, password, host, port, dbname, sslMode, timezone)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error al parsear la cadena de conexión: %w", err)
	}

	// Configuración de slog
	// Configuración del logger con formato y nivel de log
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	// Asignar tracer con slog adaptado
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   slog_logger.NewLogger(logger),
		LogLevel: tracelog.LogLevelTrace,
	}

	// Configuraciones opcionales del pool
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.HealthCheckPeriod = time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
		return err
	}

	DBInstance.Pool = pool

	logger.Info("Inicializando base de datos",
		slog.String("host", host),
		slog.String("usuario", user),
		slog.String("base de datos", dbname),
	)

	return nil
}
