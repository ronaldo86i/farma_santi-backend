package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

type BackupService struct{}

// GenerateBackup crea un archivo .sql en la carpeta backups en formato de texto plano (SQL)
func (s *BackupService) GenerateBackup(ctx context.Context) (io.ReadCloser, error) {
	// 1. Obtener credenciales y configuración
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbContainer := os.Getenv("DB_CONTAINER_NAME")

	// Directorio donde se guardarán los backups
	backupDir := "./backups"

	// Validación básica
	if dbHost == "" {
		return nil, datatype.NewInternalServerError("Configuración de base de datos no encontrada (DB_HOST)")
	}

	// 2. Asegurar que el directorio existe
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio de backups: %v", err)
	}

	// 3. Generar nombre de archivo único (.sql)
	filename := fmt.Sprintf("backup_%s.sql", time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(backupDir, filename)

	// 4. Crear el archivo vacío donde escribiremos
	outFile, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error creando archivo de backup en disco: %v", err)
	}

	// 5. Preparar el comando (Estrategia Dual)
	var cmd *exec.Cmd
	// El formato de salida es PLAIN (-F p)
	const format = "p"

	if dbContainer != "" {
		dockerArgs := []string{"exec", "-i"}
		dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("PGPASSWORD=%s", dbPass))
		// Usamos -F p (Plain)
		dockerArgs = append(dockerArgs, dbContainer, "pg_dump", "-U", dbUser, "-d", dbName, "-F", format, "-b", "-v")

		cmd = exec.CommandContext(ctx, "docker", dockerArgs...)
	} else {
		// --- ESTRATEGIA NATIVA (Local) ---
		pgDumpPath := os.Getenv("PG_DUMP_PATH")
		if pgDumpPath == "" {
			pgDumpPath = "pg_dump"
		}

		cmd = exec.CommandContext(ctx, pgDumpPath,
			"-h", dbHost,
			"-p", dbPort,
			"-U", dbUser,
			"-d", dbName,
			"-F", format,
			"-b",
			"-v",
		)
		// Inyectar contraseña en el entorno local
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPass))
	}

	// Redirigir stdout al archivo
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	// 6. Ejecutar el backup
	if err := cmd.Run(); err != nil {
		err := outFile.Close()
		if err != nil {
			log.Printf("error closing backup file: %v", err)
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		err = os.Remove(fullPath)
		if err != nil {
			log.Printf("error removing backup file: %v", err)
			return nil, datatype.NewInternalServerErrorGeneric()
		}

		errMsg := fmt.Sprintf("error ejecutando backup: %v", err)
		if dbContainer != "" {
			errMsg += " (Verifica que Docker esté corriendo y el nombre del contenedor sea correcto)"
		} else {
			errMsg += " (Verifica que pg_dump esté instalado y en el PATH)"
		}
		return nil, fmt.Errorf(errMsg)
	}

	// Cerrar archivo de escritura
	err = outFile.Close()
	if err != nil {
		log.Printf("error closing backup file: %v", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// 7. Abrir el archivo nuevamente en modo lectura para enviarlo al cliente
	readFile, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("backup creado en %s, pero falló al abrir para descarga: %v", fullPath, err)
	}

	return readFile, nil
}

// ListarBackups lee la carpeta ./backups y retorna la información de los archivos
func (s *BackupService) ListarBackups(_ context.Context) (*[]domain.BackupInfo, error) {
	backupDir := "./backups"

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return &[]domain.BackupInfo{}, nil
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("error leyendo directorio de backups: %v", err)
	}

	var backups []domain.BackupInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			backups = append(backups, domain.BackupInfo{
				Name: entry.Name(),
				Size: info.Size(),
				Date: info.ModTime(),
			})
		}
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Date.After(backups[j].Date)
	})

	return &backups, nil
}

// GetBackupFile abre un archivo existente en la carpeta de backups
func (s *BackupService) GetBackupFile(_ context.Context, filename string) (io.ReadCloser, error) {
	backupDir := "./backups"

	cleanName := filepath.Base(filename)
	fullPath := filepath.Join(backupDir, cleanName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, datatype.NewNotFoundError("El archivo de respaldo no existe")
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de respaldo: %v", err)
	}

	return file, nil
}

func NewBackupService() *BackupService {
	return &BackupService{}
}

var _ port.BackupService = (*BackupService)(nil)
