# 📋 Revisão: Arquitetura Hexagonal e Estrutura do Projeto

## 1. O Conceito Central: Arquitetura Hexagonal (Ports & Adapters)

Essa é a fundação da nossa API profissional.

### A Metáfora do Hexágono

Imagine o **Core** (o coração, o cérebro) da sua aplicação. O Core contém apenas as regras de negócio puras. Ele não se importa se o banco de dados é PostgreSQL, se a rota é feita com Gin, ou se os logs são salvos em um arquivo.

O Hexágono é um **isolamento de segurança**. Ele só se comunica com o mundo exterior através de **Portas** (Interfaces Go) e usa **Adaptadores** (Implementações concretas) para realizar o trabalho sujo de infraestrutura.

| Componente              | Conceito                                  | Função em Go                                |
|-------------------------|-------------------------------------------|---------------------------------------------|
| **Domain Layer**        | O Core, o Cérebro. Regras de negócio.     | Estruturas (structs), métodos de validação. |
| **Porta (Port)**        | Contrato. O que o Core precisa.           | Uma interface Go em `/core/repository`.     |
| **Adaptador (Adapter)** | A Implementação. Como o trabalho é feito. | Uma struct que implementa a interface em `/infra/repository`. |

> **Por que isso é PROFISSIONAL?** Porque protege o ativo mais valioso: a lógica de negócio. Se decidirmos mudar de PostgreSQL para MySQL, só trocamos o Adaptador, o Core nem fica sabendo e continua funcionando perfeitamente. Zero "gambiarra"!

## 2. Estrutura de Pastas (Onde o Hexágono e os Adaptadores Vivem)

A estrutura de pastas reflete diretamente a Arquitetura Hexagonal. Cada diretório tem uma responsabilidade estrita.

### A. `/cmd` (Comando)

**Onde está:** `main.go` (o nosso ponto de partida).

**Por que existe:** No Go, `cmd` é o local padrão para os executáveis. Ele é o Entrypoint da sua aplicação.

**O que faz:**
- Lê as configurações (`.env`)
- Inicializa o Logger e o Banco de Dados
- Cria as instâncias dos Adaptadores (`PostgresUserRepository`)
- Injeta esses Adaptadores nos Casos de Uso
- Inicia o servidor (chama o `router.Run()`)

> **Regra:** O código aqui deve ser mínimo. É apenas o ponto de montagem da aplicação.

### B. `/internal` (Encapsulamento)

**Por que existe:** É uma convenção da comunidade Go. O Go garante que nenhum outro projeto externo possa importar código que esteja dentro de um diretório chamado `internal`.

**O que faz:** Garante que toda a nossa lógica interna de negócio e infraestrutura seja mantida privada e não possa ser usada acidentalmente por terceiros. É um alto nível de encapsulamento.

### C. `/internal/core` (O Hexágono - Lógica Pura)

**Função:** Contém o domínio. Não pode ter dependências externas.

**Subdiretórios:**
- **`/entity`:** Nossas structs de negócio (`User`, `Product`). Pense nelas como a representação dos objetos do mundo real.
- **`/repository`:** As Portas de Saída. Interfaces (contratos) que dizem o que precisamos fazer (ex: `Save(User)`), mas não como fazer.
- **`/usecase`:** A camada de aplicação. Orquestra as Entidades e chama as Portas.

### D. `/internal/infra` (A Sua Seleção - Adaptadores)

Este é o diretório que você selecionou. Ele é a camada de Adaptadores, o elo entre o Core (puro) e o mundo exterior (tecnologias sujas).

| Subdiretório | Função | Tipo de Adaptador | Exemplo Tecnológico |
|--------------|--------|-------------------|-------------------|
| `/config`    | Lida com variáveis de ambiente (`.env`) e carregamento de configurações. | Configuração Externa | Viper (biblioteca Go) |
| `/logger`    | Lida com a escrita de logs (formatos JSON, níveis de severidade). | Saída Externa | Zap (biblioteca Go) |
| `/database`  | Lida com a conexão, configuração do pool e, no nosso caso, as Migrations (ferramenta externa). | Conexão Externa | gorm, postgres/migrate |
| `/repository` | O Adaptador de Saída principal. Implementa as interfaces de `/core/repository`, usando GORM para interagir com o PostgreSQL. | Implementação de Porta | PostgresUserRepository |
| `/web`        | O Adaptador de Entrada principal. Recebe as requisições HTTP, chama os Casos de Uso e formata as respostas JSON. | Porta de Entrada | Gin (framework web) |

> **Em resumo:** o `/infra` é o local onde dizemos "como" faremos o que o `/core` disse que precisa ser feito.

## 3. Configuração do Go Modules (`go.mod`)

**O que é:** O `go.mod` é o coração do gerenciamento de dependências do Go. Ele lista todas as bibliotecas externas que nosso projeto usa (Gin, GORM, Zap, etc.).

- `module github.com/devleo-m/go-zero`: Define o nome do nosso módulo. Isso é crucial, pois é o caminho que usamos para importar nossos próprios pacotes, como `github.com/devleo-m/go-zero/internal/core/entity`.