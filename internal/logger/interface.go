package logger

// интерфейс логгера
type Logger interface {
	Error(msg string, err error)
	Fatal(msg string, err error)
	Info(msg string, err error)
	Sync()
}
