# 🧠 DOMAIN LAYER - GUIA COMPLETO

## **O QUE É O DOMAIN LAYER?**

O **Domain Layer** é o **"cérebro"** da sua aplicação! É onde ficam todas as **regras de negócio** - as regras que fazem sua aplicação funcionar.

### **ANALOGIA SIMPLES:**
Imagine que você está construindo uma **loja online**:

- **Domain Layer** = As **regras da loja** (quem pode comprar, como calcular preço, quando enviar produto)
- **HTTP Layer** = A **vitrine** (onde o cliente vê os produtos)
- **Database Layer** = O **estoque** (onde ficam os produtos guardados)

**O Domain Layer NÃO sabe:**
- ❌ Como os dados são salvos no banco
- ❌ Como as requisições HTTP chegam
- ❌ Como enviar emails
- ❌ Como fazer login

**O Domain Layer SÓ sabe:**
- ✅ Quem é um usuário
- ✅ Como validar uma senha
- ✅ Como calcular um preço
- ✅ Quais são as regras de negócio

---

## **POR QUE EXISTE O DOMAIN LAYER?**

### **PROBLEMA SEM DOMAIN LAYER:**
```go
// ❌ RUIM: Tudo misturado no handler
func CreateUser(c *gin.Context) {
    // Validação de email
    if !isValidEmail(email) { ... }
    
    // Hash da senha
    passwordHash := bcrypt.Hash(password)
    
    // Salvar no banco
    db.Create(&user)
    
    // Enviar email
    sendEmail(user.Email)
    
    // Log
    log.Info("User created")
    
    // Resposta HTTP
    c.JSON(200, user)
}
```

**Problemas:**
- ❌ Regras de negócio espalhadas
- ❌ Difícil de testar
- ❌ Difícil de reutilizar
- ❌ Difícil de manter

### **SOLUÇÃO COM DOMAIN LAYER:**
```go
// ✅ BOM: Separado por responsabilidade
func CreateUser(c *gin.Context) {
    // 1. Converter HTTP para Domain
    user, err := domain.NewUser(email, password, name, role)
    
    // 2. Aplicar regras de negócio
    err = useCase.CreateUser(user)
    
    // 3. Resposta HTTP
    c.JSON(200, user)
}
```

**Benefícios:**
- ✅ Regras centralizadas
- ✅ Fácil de testar
- ✅ Fácil de reutilizar
- ✅ Fácil de manter

---

## **ONDE FICA NA ARQUITETURA?**

```
┌─────────────────────────────────────┐
│           HTTP HANDLERS             │ ← Recebe requests
├─────────────────────────────────────┤
│            USE CASES                │ ← Orquestra lógica
├─────────────────────────────────────┤
│           DOMAIN LAYER              │ ← 🧠 REGRAS DE NEGÓCIO
│  • Entities (User, Product)        │
│  • Value Objects (Email, Money)    │
│  • Domain Services                 │
│  • Repository Interfaces           │
├─────────────────────────────────────┤
│            ADAPTERS                 │ ← Implementa interfaces
└─────────────────────────────────────┘
```

**O Domain Layer fica no MEIO:**
- **Recebe** dados dos Use Cases
- **Aplica** regras de negócio
- **Retorna** resultado para Use Cases

---

## **COMPONENTES DO DOMAIN LAYER**

### **1. 🏗️ ENTITIES (Entidades)**
**O que são:** Objetos que representam coisas importantes do seu negócio.

**Exemplo:** `User` - representa um usuário do sistema

**Características:**
- ✅ Têm **identidade única** (ID)
- ✅ Têm **ciclo de vida** (criado, modificado, deletado)
- ✅ Têm **comportamento** (métodos)

```go
type User struct {
    id    string
    email Email
    // ...
}

// Comportamento da entidade
func (u *User) CanLogin() bool {
    return u.IsActive() && u.IsVerified()
}
```

### **2. 💎 VALUE OBJECTS (Objetos de Valor)**
**O que são:** Objetos que são identificados pelo **valor**, não por ID.

**Exemplo:** `Email` - dois emails são iguais se o valor for igual

**Características:**
- ✅ **Imutáveis** (não mudam)
- ✅ **Identificados pelo valor**
- ✅ **Sem identidade única**

```go
type Email struct {
    value string  // Encapsulado!
}

// Dois emails são iguais se o valor for igual
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

### **3. 🔧 DOMAIN SERVICES (Serviços de Domínio)**
**O que são:** Lógica de negócio que **não pertence a uma entidade específica**.

**Exemplo:** `PasswordService` - gera senhas, valida força

**Características:**
- ✅ **Não pertence a uma entidade**
- ✅ **Lógica de negócio pura**
- ✅ **Sem dependências externas**

```go
type PasswordService struct{}

func (s *PasswordService) GenerateRandomPassword(length int) (string, error) {
    // Lógica para gerar senha aleatória
}
```

### **4. 📋 REPOSITORY INTERFACES (Interfaces de Repositórios)**
**O que são:** **Contratos** que definem como acessar dados.

**Exemplo:** `UserRepository` - define como salvar/buscar usuários

**Características:**
- ✅ **Só interfaces** (não implementação)
- ✅ **Definem contratos**
- ✅ **Independentes de tecnologia**

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    // ... mais operações
}
```

### **5. ❌ DOMAIN ERRORS (Erros do Domínio)**
**O que são:** Erros específicos do seu negócio.

**Exemplo:** `ErrUserNotFound` - usuário não encontrado

**Características:**
- ✅ **Códigos padronizados**
- ✅ **Mensagens claras**
- ✅ **Contexto adicional**

```go
var ErrUserNotFound = NewDomainError("USER_NOT_FOUND", "usuário não encontrado", nil)
```

---

## **EXPLICAÇÃO DETALHADA DE CADA ARQUIVO**

### **📁 entities/user.go - ENTIDADE USER**

**O que é:** Representa um usuário do sistema com todas as regras de negócio.

**Por que existe:** Para centralizar todas as regras sobre usuários em um lugar só.

**O que faz:**
- ✅ Valida email e senha
- ✅ Controla status (ativo, inativo, bloqueado)
- ✅ Gerencia verificação de email
- ✅ Controla papéis (admin, client)
- ✅ Implementa soft delete
- ✅ Rastreia login

**Exemplo de uso:**
```go
// Criar usuário
user, err := entities.NewUser("user@example.com", "Password123!", "João Silva", entities.UserRoleClient)

// Verificar se pode fazer login
if user.CanLogin() {
    fmt.Println("Pode fazer login!")
}

// Alterar senha
err = user.ChangePassword("NewPassword123!")
```

**Por que é importante:**
- **Centraliza** todas as regras sobre usuários
- **Facilita** manutenção e evolução
- **Garante** consistência das regras

---

### **📁 valueobjects/email.go - EMAIL VALUE OBJECT**

**O que é:** Email válido com validação rigorosa.

**Por que existe:** Para garantir que emails sejam sempre válidos e normalizados.

**O que faz:**
- ✅ Valida formato do email
- ✅ Normaliza (lowercase, trim)
- ✅ Extrai domínio e parte local
- ✅ Compara emails de forma segura

**Exemplo de uso:**
```go
// Criar email
email, err := valueobjects.NewEmail("USER@EXAMPLE.COM")
// Resultado: "user@example.com"

// Extrair informações
fmt.Println(email.GetDomain())     // "example.com"
fmt.Println(email.GetLocalPart())  // "user"
```

**Por que é importante:**
- **Garante** que emails sejam sempre válidos
- **Evita** problemas de validação espalhados
- **Normaliza** dados automaticamente

---

### **📁 valueobjects/password.go - PASSWORD VALUE OBJECT**

**O que é:** Senha segura com hash Argon2.

**Por que existe:** Para garantir segurança máxima das senhas.

**O que faz:**
- ✅ Valida força da senha
- ✅ Gera hash seguro (Argon2)
- ✅ Verifica senha sem expor hash
- ✅ Resistente a ataques

**Exemplo de uso:**
```go
// Criar senha
password, err := valueobjects.NewPassword("MySecure123!")
if err != nil {
    log.Fatal(err)
}

// Verificar senha
if password.Verify("MySecure123!") {
    fmt.Println("Senha correta!")
}
```

**Por que é importante:**
- **Argon2** é mais seguro que bcrypt
- **Nunca** expõe senha em texto plano
- **Valida** força automaticamente

---

### **📁 valueobjects/money.go - MONEY VALUE OBJECT**

**O que é:** Valor monetário com precisão decimal.

**Por que existe:** Para evitar problemas de centavos com float.

**O que faz:**
- ✅ Usa `big.Int` para precisão
- ✅ Suporte a múltiplas moedas
- ✅ Operações matemáticas seguras
- ✅ Formatação para exibição

**Exemplo de uso:**
```go
// Criar valor
price, err := valueobjects.NewMoney(99.99, "BRL")
if err != nil {
    log.Fatal(err)
}

// Operações (simples!)
total := price.Multiply(2)      // R$ 199,98
sum := price.Add(price)         // R$ 199,98

fmt.Println(price.Format())     // "99,99 BRL"
```

**Por que é importante:**
- **Evita** bugs de centavos
- **Garante** precisão decimal
- **Suporta** múltiplas moedas

---

### **📁 services/password_service.go - PASSWORD SERVICE**

**O que é:** Serviços relacionados a senhas.

**Por que existe:** Para centralizar lógica de senhas que não pertence à entidade User.

**O que faz:**
- ✅ Gera senhas aleatórias
- ✅ Valida força da senha
- ✅ Gera hash de senha
- ✅ Verifica senha

**Exemplo de uso:**
```go
service := services.NewPasswordService()

// Gerar senha aleatória
password, err := service.GenerateRandomPassword(12)
if err != nil {
    log.Fatal(err)
}

// Validar força
isValid, errors := service.ValidatePasswordStrength("MyPass123!")
if !isValid {
    fmt.Println("Erros:", errors)
}
```

**Por que é importante:**
- **Centraliza** lógica de senhas
- **Reutilizável** em diferentes contextos
- **Especializado** em senhas

---

### **📁 repositories/user_repository.go - USER REPOSITORY INTERFACE**

**O que é:** Interface que define como acessar dados de usuários.

**Por que existe:** Para desacoplar regras de negócio de persistência de dados.

**O que faz:**
- ✅ Define contrato para CRUD
- ✅ Define operações específicas
- ✅ Suporta paginação
- ✅ Suporta soft delete

**Exemplo de uso:**
```go
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id string) (*entities.User, error)
    GetByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id string) error
    // ... mais operações
}
```

**Por que é importante:**
- **Desacopla** domain de persistência
- **Facilita** testes (mock)
- **Permite** trocar banco de dados

---

### **📁 errors/domain_errors.go - DOMAIN ERRORS**

**O que é:** Erros específicos do domínio com códigos e contexto.

**Por que existe:** Para padronizar e facilitar tratamento de erros.

**O que faz:**
- ✅ Define erros padronizados
- ✅ Inclui códigos e contexto
- ✅ Facilita identificação
- ✅ Melhora debugging

**Exemplo de uso:**
```go
// Erro específico
err := errors.ErrUserNotFound

// Erro com contexto
err := errors.NewDomainError("INVALID_EMAIL", "email inválido", map[string]interface{}{
    "email": "invalid-email",
    "field": "email",
})
```

**Por que é importante:**
- **Padroniza** tratamento de erros
- **Facilita** debugging
- **Melhora** experiência do desenvolvedor

---

### **📁 domain.go - PONTO DE ENTRADA**

**O que é:** Arquivo que facilita imports do domain layer.

**Por que existe:** Para simplificar o uso do domain layer.

**O que faz:**
- ✅ Re-exporta tipos principais
- ✅ Re-exporta constantes
- ✅ Re-exporta erros comuns
- ✅ Facilita imports

**Exemplo de uso:**
```go
import "github.com/devleo-m/go-zero/internal/modules/ecommerce/domain"

// Usar tipos re-exportados
var user domain.User
var email domain.Email
```

**Por que é importante:**
- **Simplifica** imports
- **Centraliza** exports
- **Facilita** uso

---

## **PRINCÍPIOS DO DOMAIN LAYER**

### **1. PURE BUSINESS LOGIC**
```go
// ✅ BOM: Só regras de negócio
func (u *User) CanLogin() bool {
    return u.IsActive() && u.IsVerified()
}

// ❌ RUIM: Dependência de framework
func (u *User) CanLogin(c *gin.Context) bool {
    // Gin não deveria estar aqui!
}
```

### **2. NO EXTERNAL DEPENDENCIES**
```go
// ✅ BOM: Só tipos Go nativos
type User struct {
    id    string
    email Email
}

// ❌ RUIM: Dependência de GORM
type User struct {
    gorm.Model  // Não deveria estar aqui!
    id    string
    email Email
}
```

### **3. RICH DOMAIN MODEL**
```go
// ✅ BOM: Entidade com comportamento
func (u *User) IsActive() bool {
    return u.status == "active"
}

func (u *User) CanLogin() bool {
    return u.IsActive() && u.IsVerified()
}

// ❌ RUIM: Entidade anêmica (só dados)
type User struct {
    id     string
    email  string
    status string
    // Sem métodos!
}
```

---

## **BENEFÍCIOS DO DOMAIN LAYER**

### **1. TESTABILIDADE**
```go
func TestUser_CanLogin(t *testing.T) {
    user := NewUser("test@example.com", "Password123!", "John Doe", UserRoleClient)
    
    assert.False(t, user.CanLogin()) // Não verificado
    
    user.VerifyEmail()
    assert.True(t, user.CanLogin()) // Agora pode!
}
```

### **2. REUTILIZAÇÃO**
```go
// Mesma lógica funciona em HTTP, CLI, gRPC
user := NewUser(email, password, name, role)
err := user.Validate()
// Usado em qualquer lugar!
```

### **3. MANUTENIBILIDADE**
```go
// Mudança na regra de senha? Só muda no Domain!
func (u *User) ChangePassword(newPassword string) error {
    if len(newPassword) < 12 { // Era 8, agora é 12
        return errors.New("senha muito curta")
    }
    // ...
}
```

### **4. INDEPENDÊNCIA DE FRAMEWORKS**
```go
// Domain não sabe se é Gin, Echo, ou HTTP puro
type User struct {
    // Só tipos Go nativos!
}
```

---

## **COMMON MISTAKES (Evitar!)**

### **❌ MISTAKE 1: Dependências Externas**
```go
// ❌ RUIM: GORM no domain
type User struct {
    gorm.Model
    id string
}

// ✅ BOM: Só tipos Go
type User struct {
    id string
}
```

### **❌ MISTAKE 2: Lógica de Infraestrutura**
```go
// ❌ RUIM: HTTP no domain
func (u *User) ToJSON() string {
    return gin.H{"id": u.ID}.String()
}

// ✅ BOM: Só regras de negócio
func (u *User) IsActive() bool {
    return u.status == "active"
}
```

### **❌ MISTAKE 3: Entidades Anêmicas**
```go
// ❌ RUIM: Só dados
type User struct {
    id     string
    email  string
    status string
}

// ✅ BOM: Com comportamento
type User struct {
    id     string
    email  string
    status string
}

func (u *User) IsActive() bool { return u.status == "active" }
func (u *User) CanLogin() bool { return u.IsActive() }
```

---

## **RESUMO**

### **O QUE É:**
- **Cérebro** da aplicação
- **Regras de negócio** puras
- **Independente** de frameworks

### **O QUE FAZ:**
- **Valida** dados
- **Aplica** regras de negócio
- **Garante** consistência

### **POR QUE É IMPORTANTE:**
- **Centraliza** regras de negócio
- **Facilita** manutenção
- **Melhora** testabilidade
- **Aumenta** reutilização

### **COMO USAR:**
1. **Criar** entidades com comportamento
2. **Usar** value objects para validação
3. **Definir** interfaces para persistência
4. **Aplicar** regras de negócio no domain
5. **Testar** tudo isoladamente

---

## **PRÓXIMOS PASSOS**

1. **Use Cases Layer** - Orquestrar lógica de negócio
2. **Adapters Layer** - Implementar repositórios
3. **HTTP Layer** - Expor via API REST

---

**O Domain Layer é o coração da sua aplicação! Mantenha-o puro e focado nas regras de negócio.** 🧠
