# 🔧 Diretório: `infra`

Este diretório representa a **camada de infraestrutura** da nossa aplicação. Ele contém os **"Adaptadores"** que se conectam ao **"Core"** através das **"Portas"** (interfaces).

## 🎯 Propósito

O código aqui é responsável por lidar com todas as tecnologias e detalhes externos, como:

- **Banco de Dados** (configuração, conexão, implementações de repositório)
- **Servidor Web** (handlers HTTP, middlewares, roteamento)
- **Logging**
- **Gerenciamento de configuração**
- **Clientes de APIs de terceiros**

## 📂 Subdiretórios

### `/config`
Carregamento e gerenciamento de configurações (ex: Viper).

### `/database`
Conexão com o banco de dados e migrations.

### `/logger`
Implementação do logger (ex: Zap).

### `/repository`
Implementações concretas das interfaces de repositório definidas em `/core/repository` (ex: `PostgresUserRepository`).

### `/web`
Tudo relacionado ao servidor web (ex: Gin), como handlers, DTOs e middlewares.