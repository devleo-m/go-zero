# 🏗️ Entendendo a Arquitetura Hexagonal (Ports & Adapters)

## O que é e por que usar?

A **Arquitetura Hexagonal**, também conhecida como "Ports and Adapters" (Portas e Adaptadores), é um padrão de design de software que tem um objetivo principal: **isolar a lógica de negócio principal da sua aplicação (o "Core") de todas as dependências externas**, como banco de dados, APIs de terceiros, interface de usuário, etc.

Pense no **Core** da sua aplicação como o "cérebro" que contém as regras mais importantes e puras do seu negócio (ex: "um usuário não pode ter dois CPFs iguais", "o valor de um agendamento não pode ser negativo"). Esse cérebro não deve saber como os dados são salvos (é PostgreSQL? É um arquivo de texto?) ou como um usuário interage com ele (é uma API REST? É uma linha de comando?).

## Por que usar?

### ✅ **Testabilidade**
O Core pode ser testado de forma isolada, sem a necessidade de um banco de dados ou um servidor web rodando. Isso torna os testes mais rápidos e confiáveis.

### ✅ **Flexibilidade Tecnológica**
Hoje usamos PostgreSQL. Amanhã, podemos querer mudar para o MySQL ou até mesmo para um banco NoSQL. Com a Arquitetura Hexagonal, essa troca se torna muito mais simples, pois só precisamos trocar o "adaptador" sem tocar na lógica de negócio.

### ✅ **Manutenibilidade**
O código fica mais organizado e com responsabilidades bem definidas. Uma mudança na API não quebra a regra de negócio e vice-versa.

### ✅ **Longevidade**
A lógica de negócio tende a mudar com menos frequência do que as tecnologias externas. Isolar o Core protege o ativo mais valioso do seu software.

## Camadas e Fluxo de Dados

A arquitetura se organiza em torno de **Portas** (interfaces) e **Adaptadores** (implementações).

### Camadas Principais

#### 1. **Domain Layer** (O Hexágono / Core)

**O que contém:** As Entidades (ex: `User`, `Patient`), Value Objects (ex: `Email`, `Password`) e as Regras de Negócio Puras.

> **Regra de Ouro:** Este código não depende de NADA externo. É Go puro.

**Portas de Saída (Output Ports):** Define as interfaces que o Core precisa para se comunicar com o mundo exterior. Exemplo: `UserRepository`, que define métodos como `Save(user *User)` e `FindByID(id string)`, mas não diz como salvar.

#### 2. **Application/Use Case Layer**

**O que contém:** Orquestra a lógica de negócio. É o maestro. Um caso de uso como "Criar Novo Usuário" recebe os dados, usa as entidades do domínio para validar as regras e chama a porta do repositório para persistir os dados.

**Fluxo:** Recebe dados brutos, valida, interage com o Domínio e usa as Portas para executar ações.

#### 3. **Adapters Layer** (Infraestrutura)

**O que contém:** As implementações concretas das portas e tudo que lida com o mundo exterior.

- **Adaptadores de Entrada (Driving Adapters):** Iniciam a interação. Exemplo: Um `UserHandler` (controller HTTP) que recebe um request, chama o `CreateUserUseCase` e devolve uma resposta JSON.

- **Adaptadores de Saída (Driven Adapters):** São "chamados" pelo Core através das portas. Exemplo: `PostgresUserRepository`, que implementa a interface `UserRepository` usando GORM para salvar os dados no PostgreSQL.

## Fluxo de Dados Típico (Ex: Criar um Usuário via API)

1. Um request `POST /users` chega ao **Adaptador de Entrada** (`UserHandler` do Gin)
2. O Handler extrai os dados do request (JSON) e chama o **Caso de Uso** `CreateUserUseCase`
3. O UseCase usa as **Entidades do Domínio** para validar as regras de negócio (ex: a senha é forte?)
4. Se tudo estiver OK, o UseCase chama o método `Save()` da **Porta de Saída** (`UserRepository`)
5. O Go injeta a implementação concreta, que é o **Adaptador de Saída** (`PostgresUserRepository`)
6. O `PostgresUserRepository` usa o GORM para salvar o usuário no banco de dados
7. O resultado volta pelo mesmo caminho até o Handler, que envia uma resposta `201 Created` para o cliente