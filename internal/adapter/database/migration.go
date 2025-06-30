package database

import (
	"context"
	"fmt"
	"os"
)

func (db DB) Migration() error {
	// Ejecutar migration.sql
	sqlMigrationBytes, err := os.ReadFile("./sql/migration.sql")
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo de migración: %w", err)
	}

	sqlMigrationContent := string(sqlMigrationBytes)
	_, err = db.Pool.Exec(context.Background(), sqlMigrationContent)
	if err != nil {
		return fmt.Errorf("error al ejecutar la migración SQL: %w", err)
	}

	// Ejecutar functions.sql
	sqlFunctionsBytes, err := os.ReadFile("./sql/functions.sql")
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo de funciones: %w", err)
	}

	sqlFunctionsContent := string(sqlFunctionsBytes)
	_, err = db.Pool.Exec(context.Background(), sqlFunctionsContent)
	if err != nil {
		return fmt.Errorf("error al ejecutar las funciones SQL: %w", err)
	}

	logger.Info("✅ Migración y funciones ejecutadas correctamente.")
	return nil
}
