package logger

// интерфейс логгера
type Logger interface {
	Error(msg string, err error, fields ...any)
	Fatal(msg string, err error, fields ...any)
	Info(msg string, err error, fields ...any)
	Sync()
}
