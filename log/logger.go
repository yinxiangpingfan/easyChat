package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

var logPath string

func InitLogger() {
	logger, _ = zap.NewProduction()
}
