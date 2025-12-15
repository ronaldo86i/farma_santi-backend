package handler

import (
	"errors"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BackupHandler struct {
	backupService port.BackupService
}

func NewBackupHandler(backupService port.BackupService) *BackupHandler {
	return &BackupHandler{backupService: backupService}
}

// DownloadBackup genera el backup en caliente y lo envía como stream al cliente
func (h *BackupHandler) DownloadBackup(c *fiber.Ctx) error {
	// 1. Llamar al servicio
	backupStream, err := h.backupService.GenerateBackup(c.UserContext())
	if err != nil {
		// En producción, es bueno loguear el error real: log.Println(err)
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(util.NewMessage("Error al generar el respaldo de la base de datos"))
	}

	// Nota: No cerramos el stream manualmente aquí; Fiber lo cerrará al terminar de enviar la respuesta.

	// 2. Configurar cabeceras de descarga para el navegador
	filename := fmt.Sprintf("farma_santi_backup_%s.dump", time.Now().Format("2006-01-02_15-04-05"))

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	// 3. Enviar el stream directamente (eficiente para archivos grandes)
	return c.SendStream(backupStream)
}

// ListarBackups retorna el listado de archivos existentes en el directorio de backups
func (h *BackupHandler) ListarBackups(c *fiber.Ctx) error {
	backups, err := h.backupService.ListarBackups(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(util.NewMessage("Error al listar los backups disponibles"))
	}

	return c.JSON(backups)
}

func (h *BackupHandler) DownloadBackupFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(util.NewMessage("El nombre del archivo es requerido"))
	}

	fileStream, err := h.backupService.GetBackupFile(c.Context(), filename)
	if err != nil {
		log.Print(err.Error())
		var errorResponse *datatype.ErrorResponse
		if errors.As(err, &errorResponse) {
			return c.Status(errorResponse.Code).JSON(util.NewMessage(errorResponse.Message))
		}
		return c.Status(http.StatusInternalServerError).JSON(util.NewMessage(err.Error()))
	}

	// Configurar headers para descarga
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	return c.SendStream(fileStream)
}

var _ port.BackupHandler = (*BackupHandler)(nil)
