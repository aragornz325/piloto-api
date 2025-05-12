package utils

import (
	"context"
	"time"

	"github.com/aragornz325/piloto-api/pkg/logger"
	"go.uber.org/zap"
)

func PerformHandlerOperation(opts PerformHandlerOperationFunc) error {
	start := time.Now().UTC()

	logger.Log.Info("🧪 start "+opts.HandlerName+" handler operation",
		zap.String("op", opts.Name),
	)

	err := opts.Operation()

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error("💥 An error occurred in "+opts.Name+" handler operation",
			zap.String("op", opts.Name),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
	} else {
		logger.Log.Info("✅ Operation "+opts.Name+" Success",
			zap.String("op", opts.Name),
			zap.Duration("duration", duration),
		)
	}

	return err
}

// /-------------structs------------------///
type PerformHandlerOperationFunc struct {
	Ctx         context.Context
	Name        string
	HandlerName string
	Operation   func() error
}
