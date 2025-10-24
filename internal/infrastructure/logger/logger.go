package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

type Config struct {
	Level  string
	Format string
}

func New(config Config) (*Logger, error) {
	zapConfig := zap.NewProductionConfig()

	// Configurar nível de log
	switch config.Level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configurar formato
	if config.Format == "json" {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	}

	// Configurar output
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: logger}, nil
}

func NewFromEnv() (*Logger, error) {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		format = "json"
	}

	return New(Config{
		Level:  level,
		Format: format,
	})
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

// WithFieldsFromMap cria campos zap a partir de um map
func (l *Logger) WithFieldsFromMap(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{Logger: l.Logger.With(zapFields...)}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithFields(zap.String("user_id", userID))
}

func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}

func (l *Logger) WithComponent(component string) *Logger {
	return l.WithFields(zap.String("component", component))
}

func (l *Logger) WithOperation(operation string) *Logger {
	return l.WithFields(zap.String("operation", operation))
}

func (l *Logger) WithDuration(duration time.Duration) *Logger {
	return l.WithFields(zap.Duration("duration", duration))
}

func (l *Logger) WithHTTPRequest(method, path string, statusCode int) *Logger {
	return l.WithFields(
		zap.String("http_method", method),
		zap.String("http_path", path),
		zap.Int("http_status", statusCode),
	)
}

func (l *Logger) WithDatabaseQuery(query string, duration time.Duration) *Logger {
	return l.WithFields(
		zap.String("db_query", query),
		zap.Duration("db_duration", duration),
	)
}

func (l *Logger) WithCacheOperation(operation, key string, hit bool) *Logger {
	return l.WithFields(
		zap.String("cache_operation", operation),
		zap.String("cache_key", key),
		zap.Bool("cache_hit", hit),
	)
}

// Métodos de conveniência para logging comum
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}

// Sync garante que todos os logs sejam escritos
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
