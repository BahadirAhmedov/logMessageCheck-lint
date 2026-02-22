package zap

type Logger struct{}

func NewProduction() (*Logger, error) {
	return &Logger{}, nil
}

func (l *Logger) Sync() error {
	return nil
}

type Field struct{}

func (l *Logger) Info(msg string, fields ...Field) {}

func (l *Logger) Error(msg string, fields ...Field) {}

func (l *Logger) Debug(msg string, fields ...Field) {}

func (l *Logger) Warn(msg string, fields ...Field) {}