package utils

import (
	"context"
	"time"

	"github.com/aragornz325/piloto-api/pkg/logger"
	"go.uber.org/zap"
)

func PerformServiceOperation(opts PerformServiceOperationFunc) error {
	start := time.Now().UTC()

	logger.Log.Info("🧪 start "+opts.ServiceName+" service operation",
		zap.String("op", opts.Name),
	)

	err := opts.Operation()

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error("💥 An error occurred in "+opts.ServiceName+" service operation",
			zap.String("op", opts.Name),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
	} else {
		logger.Log.Info("✅ Operation "+opts.ServiceName+" service operation success",
			zap.String("op", opts.Name),
			zap.Duration("duration", duration),
		)
	}

	return err
}

// /-------------structs------------------///
type PerformServiceOperationFunc struct {
	Ctx         context.Context
	Name        string
	ServiceName string
	Operation   func() error
}
