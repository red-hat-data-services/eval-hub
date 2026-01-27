package storage

import (
	"fmt"
	"log/slog"

	"github.com/eval-hub/eval-hub/internal/abstractions"
	"github.com/eval-hub/eval-hub/internal/config"
	"github.com/eval-hub/eval-hub/internal/storage/storage_sql"
)

func NewStorage(serviceConfig *config.Config, logger *slog.Logger) (abstractions.Storage, error) {
	// search for the first enabled database configuration
	for name, sqlConfig := range serviceConfig.Database.SQL {
		if sqlConfig.Enabled {
			logger.Info("Using SQL database configuration", "name", name)
			return storage_sql.NewSQLStorage(&sqlConfig, logger)
		}
	}
	for name, jsonConfig := range serviceConfig.Database.JSON {
		if jsonConfig.Enabled {
			// return storage_json.NewJSONStorage(jsonConfig)
			return nil, fmt.Errorf("JSON database configuration %s not supported yet", name)
		}
	}
	for name, dbConfig := range serviceConfig.Database.Other {
		if dbConfig.Enabled {
			// return storage_json.NewJSONStorage(jsonConfig)
			return nil, fmt.Errorf("Other database configuration %s not supported yet", name)
		}
	}
	// if no other database configuration is enabled, use the fallback one
	for name, sqlConfig := range serviceConfig.Database.SQL {
		if sqlConfig.Fallback {
			logger.Info("Using fallback SQL database configuration", "name", name)
			return storage_sql.NewSQLStorage(&sqlConfig, logger)
		}
	}
	return nil, fmt.Errorf("failed to find a supported and enabled database configuration")
}
