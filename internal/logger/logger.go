package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	*zap.Logger
}

var ZapLog = zapLogger{} // Logger global

// SetupLogger configura el logger global con zap
func SetupLogger() error {
	var err error
	config := zap.NewProductionConfig() // Puedes usar NewDevelopmentConfig si estás en desarrollo
	config.OutputPaths = []string{
		//"stdout",        // Salida en consola
		"./console.log", // Archivo de log
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// Crear el logger con la configuración proporcionada
	ZapLog.Logger, err = config.Build()
	if err != nil {
		return fmt.Errorf("Error al configurar el logger: %v", err)
	}

	// Si la configuración es exitosa, podemos registrar un mensaje
	ZapLog.Logger.Info("Logger configurado correctamente")
	return nil
}

func Shutdown() {
	ZapLog.Logger.Sync() // Hacer flush de los logs
}
