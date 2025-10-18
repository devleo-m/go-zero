Inicialize o módulo Go.

```Bash
go mod init github.com/seu-usuario/go-hexagonal-api # Use seu próprio namespace/usuário
```
Estrutura de Pastas (Base da Arquitetura Hexagonal):

A arquitetura hexagonal se traduz em uma estrutura de pastas que separa as camadas de forma clara:

```
go-hexagonal-api/
├── cmd/                # Ponto de entrada da aplicação (main)
│   └── server/         # Configura e inicia o servidor HTTP
│       └── main.go
├── internal/           # Lógica de domínio e aplicação (Core)
│   ├── core/           # O "Hexágono" - Domain + Use Cases
│   │   ├── domain/     # Entidades, Value Objects, Portas de Saída (Interfaces)
│   │   └── ports/      # Portas de Entrada (Interfaces: HTTP, CLI, etc.)
│   │   └── usecase/    # Implementação da Lógica de Negócio (Use Cases)
│   └── infra/          # Adaptadores e Infraestrutura
│       ├── adapter/    # Implementações das portas (Adapters)
│       │   ├── handler/    # Adaptador de entrada (ex: HTTP handlers/controllers)
│       │   └── repository/ # Adaptador de saída (ex: GORM/DB implementation)
│       ├── config/     # Gerenciamento de Configuração (Viper)
│       ├── database/   # Conexão com DB, Migrações, GORM
│       └── logger/     # Logger (Zap)
├── pkg/                # Código de uso comum que pode ser importado por outros projetos
├── tests/              # Testes E2E e Integração
├── vendor/             # Dependências de módulos (opcional, se go mod vendor for usado)
├── .env                # Variáveis de ambiente
├── go.mod
├── go.sum
└── Dockerfile          # Para containerização
└── docker-compose.yml  # Para ambiente de desenvolvimento
```

Crie a estrutura de pastas principal usando o comando mkdir -p
```
mkdir -p cmd/server internal/core/{domain,ports,usecase} internal/infra/{adapter/{handler,repository},config,database,logger} pkg tests
touch cmd/server/main.go
```

docker-compose.yml

```
version: '3.8'

services:
  app:
    # O build será definido mais tarde com o Dockerfile da aplicação Go
    # Por enquanto, vamos focar no DB
    container_name: go_api_app
    restart: always
    build:
      context: .
      dockerfile: Dockerfile.dev # Usaremos um Dockerfile separado para desenvolvimento com live-reload (Air)
    ports:
      - "8080:8080"
    volumes:
      - .:/app
    depends_on:
      - db
    environment:
      # Variáveis de ambiente da aplicação
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: user
      DB_PASSWORD: password
      DB_NAME: go_hex_db
      APP_PORT: 8080

  db:
    container_name: go_api_db
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: go_hex_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Crie um arquivo .env para centralizar as variáveis (mesmas que estão no docker-compose por enquanto, mas centralizar é bom para o deploy).

```
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_exemple_db
DB_SSLMODE=disable

# Application Configuration
APP_ENV=dev
APP_PORT=8080
```

Instalar Dependências Core (Gin, GORM, PostgreSQL Driver, Viper, Zap)
Ação:

Instale as dependências core.
```
# Gin framework
go get github.com/gin-gonic/gin

# GORM
go get gorm.io/gorm
go get gorm.io/driver/postgres # Driver PostgreSQL

# Configuração (Viper)
go get github.com/spf13/viper

# Logging (Zap)
go get go.uber.org/zap
```

