package logger

import (
	"log"

	"go.uber.org/zap"
)

// Logger é a instância global do logger
var Logger *zap.Logger

// InitLogger inicializa o logger baseado no ambiente
func InitLogger(env string) {
	var err error
	if env == "production" {
		Logger, err = zap.NewProduction()
	} else {
		Logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}

	// Garante que o buffer do log seja descarregado ao encerrar
	defer Logger.Sync() 
}

// Métodos de atalho para evitar a dependência direta do zap em todo o código
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Error(msg string, err error, fields ...zap.Field) {
	if Logger != nil {
		// Adiciona o erro ao campo do log
		fields = append(fields, zap.Error(err)) 
		Logger.Error(msg, fields...)
	}
}

func Fatal(msg string, err error, fields ...zap.Field) {
	if Logger != nil {
		fields = append(fields, zap.Error(err))
		Logger.Fatal(msg, fields...)
	}
}

// Adicione outros níveis (Debug, Warn) conforme a necessidade