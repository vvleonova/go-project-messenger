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
func (log *LoggerZap) Error(msg string, err error, fields ...any) {
	allFields := append(fields, "err", err)
	log.Logger.Sugar().Errorw(msg, allFields...)
}

// обертка для функции Fatal
func (log *LoggerZap) Fatal(msg string, err error, fields ...any) {
	allFields := append(fields, "err", err)
	log.Logger.Sugar().Fatalw(msg, allFields...)
}

// обертка для функции Info
func (log *LoggerZap) Info(msg string, err error, fields ...any) {
	allFields := append(fields, "err", err)
	log.Logger.Sugar().Infow(msg, allFields...)
}

// обертка для функции Sync
func (log *LoggerZap) Sync() {
	log.Logger.Sync()
}
