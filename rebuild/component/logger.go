package component

type Logger struct {
}

func (l *Logger) Log(level int, message string)  {}
func (l *Logger) Logf(format string, a ...interface{}) {}

func (l *Logger) Info(message string) {}
func (l *Logger) Debug(message string) {}
func (l *Logger) Warn(message string) {}
func (l *Logger) Error(message string) {}
func (l *Logger) Fatal(message string) {}
func (l *Logger) Panic(message string) {}

func (l *Logger) Infof(format string, args ...interface{}) {}
func (l *Logger) Debugf(format string, args ...interface{}) {}
func (l *Logger) Warnf(format string, args ...interface{}) {}
func (l *Logger) Errorf(format string, args ...interface) {}

func NewLogger() *Logger {
	return &Logger{}
}
