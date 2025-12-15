package port

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"io"

	"github.com/gofiber/fiber/v2"
)

type BackupService interface {
	GenerateBackup(ctx context.Context) (io.ReadCloser, error)
	ListarBackups(ctx context.Context) (*[]domain.BackupInfo, error)
	GetBackupFile(ctx context.Context, filename string) (io.ReadCloser, error)
}

type BackupHandler interface {
	DownloadBackup(c *fiber.Ctx) error
	ListarBackups(c *fiber.Ctx) error
	DownloadBackupFile(c *fiber.Ctx) error
}
