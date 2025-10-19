# Usa uma imagem base Go
FROM golang:1.25.1-alpine

# Define o diretório de trabalho
WORKDIR /app

# Instala Git (para go get) e Air
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

# Copia os arquivos de dependência Go e baixa os módulos
# Isso otimiza o cache do Docker. Se go.mod e go.sum não mudarem, o cache é usado.
COPY go.mod go.sum ./
RUN go mod download

# Copia os arquivos
COPY . .

# Expõe a porta (Gin padrão é 8080)
EXPOSE 8080

# Comando para rodar a aplicação com Air (será executado pelo docker-compose)
# CMD ["air", "-c", ".air.toml"]