package logger

import "go.uber.org/zap"

// структура логгера
type LoggerZap struct {
	*zap.Logger
}

// создание логгера
func NewLogger() (Logger, error) {
	zapLogger, err := zap.NewDevelopment()

	return &LoggerZap{zapLogger}, err
}

// обертка для функции Error
func (log *LoggerZap) Error(msg string, err error) {
	log.Logger.Error(msg, zap.Error(err))
}

// обертка для функции Fatal
func (log *LoggerZap) Fatal(msg string, err error) {
	log.Logger.Fatal(msg, zap.Error(err))
}

// обертка для функции Info
func (log *LoggerZap) Info(msg string, err error) {
	log.Logger.Info(msg, zap.Error(err))
}

// обертка для функции Sync
func (log *LoggerZap) Sync() {
	log.Logger.Sync()
}
