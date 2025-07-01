package routine

import (
	"context"
	"farma-santi_backend/internal/server/setup"
)

func Init(ctx context.Context) {
	deps := setup.GetDependencies()

	go startActualizarLotesVencidos(ctx, deps.Service.LoteProducto)
}
