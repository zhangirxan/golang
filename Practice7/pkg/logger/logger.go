package logger

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct{}

func New(_ string) *Logger { return &Logger{} }

func (l *Logger) Debug(msg interface{}, args ...interface{}) {}
func (l *Logger) Info(msg string, args ...interface{})       {}
func (l *Logger) Warn(msg string, args ...interface{})       {}
func (l *Logger) Error(msg interface{}, args ...interface{}) {}
func (l *Logger) Fatal(msg interface{}, args ...interface{}) {}
