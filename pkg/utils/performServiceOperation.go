package utils

import (
	"context"
	"time"
	"github.com/aragornz325/piloto-api/pkg/logger"
	"go.uber.org/zap"
)

func PerformServiceOperation(opts PerformServiceOperationFunc ) error {
	start := time.Now()

	logger.Log.Info("🧪 Iniciando operación", 
		zap.String("op", opts.Name),
	)

	err := opts.Operation()

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error("💥 Error en operación", 
			zap.String("op", opts.Name),
			zap.Error(err),
			zap.Duration("duración", duration),
		)
	} else {
		logger.Log.Info("✅ Operación exitosa", 
			zap.String("op", opts.Name),
			zap.Duration("duración", duration),
		)
	}

	return err
}

///-------------structs------------------///
type PerformServiceOperationFunc struct {
	Ctx	  context.Context
	Name  string
	Operation func() error
}