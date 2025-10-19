package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger é a instância global do logger
var Logger *zap.Logger

// Config representa as configurações do logger
type Config struct {
	Level  string
	Format string
}

// InitLogger inicializa o logger com configurações personalizadas
func InitLogger(cfg Config) error {
	// 1. Configurar o nível de log
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("nível de log inválido: %w", err)
	}

	// 2. Configurar o encoder
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(getEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(getEncoderConfig())
	}

	// 3. Configurar o core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// 4. Criar o logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// parseLevel converte string para zapcore.Level
func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn", "warning":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("nível '%s' não reconhecido", level)
	}
}

// getEncoderConfig retorna configuração do encoder
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Sync garante que todos os logs sejam escritos
func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}

// ===========================================
// MÉTODOS DE LOGGING ESTRUTURADO
// ===========================================

// Debug registra mensagem de debug
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

// Info registra mensagem informativa
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

// Warn registra mensagem de aviso
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

// Error registra mensagem de erro
func Error(msg string, err error, fields ...zap.Field) {
	if Logger != nil {
		allFields := append(fields, zap.Error(err))
		Logger.Error(msg, allFields...)
	}
}

// Fatal registra mensagem fatal e encerra a aplicação
func Fatal(msg string, err error, fields ...zap.Field) {
	if Logger != nil {
		allFields := append(fields, zap.Error(err))
		Logger.Fatal(msg, allFields...)
	}
}

// ===========================================
// MÉTODOS DE LOGGING COM CONTEXTO
// ===========================================

// WithRequest cria logger com contexto de request
func WithRequest(requestID, method, path string) *zap.Logger {
	if Logger == nil {
		return zap.NewNop()
	}
	return Logger.With(
		zap.String("request_id", requestID),
		zap.String("method", method),
		zap.String("path", path),
	)
}

// WithUser cria logger com contexto de usuário
func WithUser(userID, email string) *zap.Logger {
	if Logger == nil {
		return zap.NewNop()
	}
	return Logger.With(
		zap.String("user_id", userID),
		zap.String("user_email", email),
	)
}

// WithDuration cria logger com duração
func WithDuration(duration time.Duration) *zap.Logger {
	if Logger == nil {
		return zap.NewNop()
	}
	return Logger.With(
		zap.Duration("duration", duration),
	)
}

// ===========================================
// MÉTODOS DE LOGGING DE NEGÓCIO
// ===========================================

// LogUserAction registra ação do usuário
func LogUserAction(action, userID, details string) {
	Info("user_action",
		zap.String("action", action),
		zap.String("user_id", userID),
		zap.String("details", details),
		zap.Time("timestamp", time.Now()),
	)
}

// LogDatabaseOperation registra operação no banco
func LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	if err != nil {
		Error("database_operation_failed",
			err,
			zap.String("operation", operation),
			zap.String("table", table),
			zap.Duration("duration", duration),
		)
	} else {
		Info("database_operation_success",
			zap.String("operation", operation),
			zap.String("table", table),
			zap.Duration("duration", duration),
		)
	}
}

// LogHTTPRequest registra requisição HTTP
func LogHTTPRequest(method, path string, statusCode int, duration time.Duration, userID string) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	}

	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	Info("http_request", fields...)
}

// LogExternalService registra chamada para serviço externo
func LogExternalService(service, operation string, duration time.Duration, err error) {
	if err != nil {
		Error("external_service_failed",
			err,
			zap.String("service", service),
			zap.String("operation", operation),
			zap.Duration("duration", duration),
		)
	} else {
		Info("external_service_success",
			zap.String("service", service),
			zap.String("operation", operation),
			zap.Duration("duration", duration),
		)
	}
}
