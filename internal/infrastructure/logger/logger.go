package logger

import (
	"time"

	"go.uber.org/zap"
)

// Config configurações do logger
type Config struct {
	Level  string
	Format string
}

var (
	// Logger é a instância global do logger
	Logger *zap.Logger
)

// InitLogger inicializa o logger
func InitLogger(cfg Config) error {
	var config zap.Config

	// Configurar baseado no ambiente
	if cfg.Format == "json" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Configurar nível de log
	switch cfg.Level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Criar logger
	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	// Substituir logger global
	zap.ReplaceGlobals(Logger)

	return nil
}

// Sync sincroniza o logger
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// Info log de informação
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error log de erro
func Error(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	Logger.Error(msg, fields...)
}

// Fatal log fatal
func Fatal(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	Logger.Fatal(msg, fields...)
}

// Debug log de debug
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn log de warning
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// LogHTTPRequest log de requisição HTTP
func LogHTTPRequest(method, path string, statusCode int, latency time.Duration, userID string) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.Duration("latency", latency),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	// Determinar nível baseado no status code
	switch {
	case statusCode >= 500:
		Logger.Error("HTTP Request", fields...)
	case statusCode >= 400:
		Logger.Warn("HTTP Request", fields...)
	default:
		Logger.Info("HTTP Request", fields...)
	}
}

// LogDatabaseQuery log de query do banco
func LogDatabaseQuery(query string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("query", query),
		zap.Duration("duration", duration),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		Logger.Error("Database Query", fields...)
	} else {
		Logger.Debug("Database Query", fields...)
	}
}

// LogBusinessEvent log de evento de negócio
func LogBusinessEvent(event string, userID string, data map[string]interface{}) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.String("user_id", userID),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	// Adicionar dados do evento
	for key, value := range data {
		fields = append(fields, zap.Any(key, value))
	}

	Logger.Info("Business Event", fields...)
}

// LogSecurityEvent log de evento de segurança
func LogSecurityEvent(event string, userID string, ip string, userAgent string, err error) {
	fields := []zap.Field{
		zap.String("event", event),
		zap.String("user_id", userID),
		zap.String("ip", ip),
		zap.String("user_agent", userAgent),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		Logger.Warn("Security Event", fields...)
	} else {
		Logger.Info("Security Event", fields...)
	}
}

// GetLogger retorna a instância do logger
func GetLogger() *zap.Logger {
	return Logger
}

// SetLogger define a instância do logger
func SetLogger(logger *zap.Logger) {
	Logger = logger
	zap.ReplaceGlobals(logger)
}
