package logger

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

// Any log
func Any(k string, v interface{}) Field {
	return Field{
		Key:   k,
		Value: v,
	}
}

type Logger interface {
	Infof(format string, fields ...Field)
	Warnf(format string, fields ...Field)
	Errorf(format string, fields ...Field)
}

func NewZapLogger(appName string) (Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	defer log.Sync()

	log.With(zap.Any("app_name", appName))

	return &zapLogger{logger: log}, nil
}

func (l *zapLogger) Infof(format string, fields ...Field) {
	l.logger.Info(format, convert(fields)...)
}

func (l *zapLogger) Warnf(format string, fields ...Field) {
	l.logger.Warn(format, convert(fields)...)

}

func (l *zapLogger) Errorf(format string, fields ...Field) {
	l.logger.Error(format, convert(fields)...)
}

func convert(fields []Field) []zap.Field {
	var result []zap.Field

	for _, fl := range fields {
		result = append(result, zap.Any(fl.Key, fl.Value))
	}

	return result
}
