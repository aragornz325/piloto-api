package logger

import (
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	env := os.Getenv("APP_ENV")

	var err error
	if env == "dev" {
		Log, err = zap.NewDevelopment()
	} else {
		Log, err = zap.NewProduction()
	}

	if err != nil {
		panic("💥 No se pudo iniciar el logger: " + err.Error())
	}
}
