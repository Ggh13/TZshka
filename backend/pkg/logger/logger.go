package logger

import (
	"go.uber.org/zap"
)

type Logger = zap.Logger

func New() (*zap.Logger, error) {
	return zap.NewProduction()
}
