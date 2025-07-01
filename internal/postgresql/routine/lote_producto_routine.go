package routine

import (
	"context"
	"farma-santi_backend/internal/core/port"
	"log"
	"time"
)

func startActualizarLotesVencidos(ctx context.Context, service port.LoteProductoService) {
	go func() {
		// Ejecutar inmediatamente al iniciar
		if err := service.ActualizarLotesVencidos(ctx); err != nil {
			log.Printf("Error al actualizar lotes vencidos inicialmente: %v", err)
		}

		// Calcular duración hasta la próxima 00:00
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		timeUntilNext := time.Until(next)

		// Esperar hasta las 00:00
		select {
		case <-time.After(timeUntilNext):
		case <-ctx.Done():
			return
		}

		// Ejecutar cada 24h a partir de las 00:00
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			if err := service.ActualizarLotesVencidos(ctx); err != nil {
				log.Printf("Error al actualizar lotes vencidos: %v", err)
			}

			select {
			case <-ticker.C:
				// Esperar al siguiente día
			case <-ctx.Done():
				return
			}
		}
	}()
}
