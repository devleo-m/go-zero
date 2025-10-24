# 🌐 Camada de Apresentação (Presentation Layer)

> **🎯 Objetivo:** Gerenciar a interface com o mundo externo, convertendo requisições HTTP em chamadas para use cases e formatando respostas adequadas.

## 📚 O que é a Camada de Apresentação?

A **Presentation Layer** é a **fachada** da sua aplicação! É aqui que:

- 🌐 **Recebemos** requisições HTTP
- 🔄 **Convertemos** dados de entrada/saída
- 🎯 **Orquestramos** chamadas para use cases
- 📤 **Formatamos** respostas para o cliente

## 🎓 Por que separar Presentation?

### ❌ **Problema sem separação:**
```go
// Lógica de negócio misturada com HTTP
func CreateUser(c *gin.Context) {
    // Validação HTTP
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Lógica de negócio aqui
    if req.Email == "" {
        c.JSON(400, gin.H{"error": "email required"})
        return
    }
    
    // Persistência aqui
    user := &User{Name: req.Name, Email: req.Email}
    db.Save(user)
    
    // Resposta HTTP
    c.JSON(201, user)
}
```

### ✅ **Solução com Presentation Layer:**
```go
// Handler foca apenas em HTTP
func (h *Handler) CreateUser(c *gin.Context) {
    // 1️⃣ Validar entrada
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "INVALID_REQUEST", err.Error())
        return
    }
    
    // 2️⃣ Converter para use case
    input := application.CreateUserInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
    
    // 3️⃣ Executar use case
    result, err := h.createUserUseCase.Execute(c.Request.Context(), input)
    if err != nil {
        response.BadRequest(c, "CREATE_USER_FAILED", err.Error())
        return
    }
    
    // 4️⃣ Converter resposta
    response.Created(c, toUserResponse(result.User), result.Message)
}
```

## 🏗️ Estrutura da Camada

```
presentation/
├── 📄 README.md           # Este arquivo - conceitos de apresentação
└── http/                  # 🌐 Handlers HTTP
    ├── handler.go         # Controllers HTTP
    ├── routes.go          # Definição de rotas
    └── dto.go            # Data Transfer Objects
```

## 🌐 Handlers HTTP - Análise Detalhada

### **1. Estrutura do Handler**
```go
type Handler struct {
    createUserUseCase *application.CreateUserUseCase
    getUserUseCase    *application.GetUserUseCase
    listUsersUseCase  *application.ListUsersUseCase
    updateUserUseCase *application.UpdateUserUseCase
    deleteUserUseCase *application.DeleteUserUseCase
}
```

### **2. CreateUser - Criar Usuário**

#### **📝 Fluxo Completo:**
```go
func (h *Handler) CreateUser(c *gin.Context) {
    // 1️⃣ Validar entrada HTTP
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "INVALID_REQUEST", err.Error())
        return
    }
    
    // 2️⃣ Sanitização de dados
    var phone *string
    if req.Phone != "" {
        phone = &req.Phone
    }
    
    // 3️⃣ Converter para use case
    input := application.CreateUserInput{
        Name:     validation.SanitizeString(req.Name),
        Email:    validation.SanitizeString(req.Email),
        Password: req.Password,
        Phone:    phone,
    }
    
    // 4️⃣ Executar use case
    result, err := h.createUserUseCase.Execute(c.Request.Context(), input)
    if err != nil {
        response.BadRequest(c, "CREATE_USER_FAILED", err.Error())
        return
    }
    
    // 5️⃣ Converter resposta
    response.Created(c, toUserResponse(result.User), result.Message)
}
```

### **3. GetUser - Buscar Usuário**

#### **🔍 Validação e Busca:**
```go
func (h *Handler) GetUser(c *gin.Context) {
    // 1️⃣ Validar UUID
    idStr := c.Param("id")
    if err := validation.ValidateUUID(idStr); err != nil {
        response.BadRequest(c, "INVALID_ID", err.Error())
        return
    }
    
    // 2️⃣ Converter para use case
    id, _ := uuid.Parse(idStr)
    input := application.GetUserInput{ID: id}
    
    // 3️⃣ Executar use case
    result, err := h.getUserUseCase.Execute(c.Request.Context(), input)
    if err != nil {
        if err == domain.ErrUserNotFound {
            response.NotFound(c, "USER_NOT_FOUND", "User not found")
            return
        }
        response.InternalServerError(c, "GET_USER_FAILED", err.Error())
        return
    }
    
    // 4️⃣ Retornar resposta
    response.Success(c, toUserResponse(result.User))
}
```

### **4. ListUsers - Listar com Paginação**

#### **📄 Paginação e Filtros:**
```go
func (h *Handler) ListUsers(c *gin.Context) {
    // 1️⃣ Extrair parâmetros de query
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")
    
    // 2️⃣ Converter para inteiros
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        limit = 10
    }
    
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        offset = 0
    }
    
    // 3️⃣ Validar paginação
    if err := validation.ValidatePagination(offset/limit+1, limit); err != nil {
        response.BadRequest(c, "INVALID_PAGINATION", err.Error())
        return
    }
    
    // 4️⃣ Executar use case
    input := application.ListUsersInput{
        Limit:  limit,
        Offset: offset,
    }
    
    result, err := h.listUsersUseCase.Execute(c.Request.Context(), input)
    if err != nil {
        response.InternalServerError(c, "LIST_USERS_FAILED", err.Error())
        return
    }
    
    // 5️⃣ Converter lista
    users := make([]UserResponse, len(result.Users))
    for i, user := range result.Users {
        users[i] = toUserResponse(user)
    }
    
    // 6️⃣ Calcular metadados de paginação
    page := (offset / limit) + 1
    meta := response.NewMeta(page, limit, int64(result.Total))
    
    // 7️⃣ Retornar resposta paginada
    response.Paginated(c, map[string]interface{}{
        "users": users,
    }, meta)
}
```

## 📋 DTOs (Data Transfer Objects)

### **1. Request DTOs - Entrada**
```go
// CreateUserRequest - Dados de entrada para criação
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Phone    string `json:"phone,omitempty"`
}

// UpdateUserRequest - Dados de entrada para atualização
type UpdateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Phone string `json:"phone,omitempty"`
}
```

### **2. Response DTOs - Saída**
```go
// UserResponse - Dados de saída (sem senha)
type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     *string   `json:"phone,omitempty"`
    Role      string    `json:"role"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### **3. Conversão Domain → Response**
```go
func toUserResponse(user *domain.User) UserResponse {
    return UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Phone:     user.Phone,
        Role:      user.Role,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
        // Password nunca é serializado!
    }
}
```

## 🎯 Padrões Aplicados

### **1. Handler Pattern**
- ✅ **Responsabilidade única** - Cada handler uma operação
- ✅ **Injeção de dependência** - Use cases injetados
- ✅ **Tratamento de erro** - Consistente em todos os handlers

### **2. DTO Pattern**
- ✅ **Separação** - Request/Response separados do domain
- ✅ **Validação** - Tags de validação no DTO
- ✅ **Segurança** - Senha nunca exposta

### **3. Response Pattern**
- ✅ **Consistência** - Todas as respostas seguem o mesmo padrão
- ✅ **Metadados** - Paginação e informações extras
- ✅ **Códigos HTTP** - Semântica correta

## 🧪 Como Testar

### **Teste de Handler:**
```go
func TestCreateUserHandler(t *testing.T) {
    // Arrange
    mockUseCase := &MockCreateUserUseCase{}
    handler := NewHandler(mockUseCase, nil, nil, nil, nil)
    
    req := CreateUserRequest{
        Name:     "João",
        Email:    "joao@test.com",
        Password: "12345678",
    }
    
    // Act
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/users", jsonBody(req))
    
    handler.CreateUser(c)
    
    // Assert
    assert.Equal(t, 201, w.Code)
    assert.Contains(t, w.Body.String(), "João")
}
```

## 🎓 Benefícios da Separação

### **1. Testabilidade**
- ✅ **Handlers** podem ser testados isoladamente
- ✅ **Use cases** podem ser mockados
- ✅ **DTOs** podem ser validados separadamente

### **2. Flexibilidade**
- ✅ **Múltiplas interfaces** - HTTP, gRPC, CLI
- ✅ **Formato de resposta** - JSON, XML, etc.
- ✅ **Validação** - Diferente por interface

### **3. Segurança**
- ✅ **Sanitização** - Dados limpos antes do processamento
- ✅ **Validação** - Entrada validada antes do use case
- ✅ **Exposição** - Apenas dados necessários na resposta

## 🚀 Próximos Passos

1. **Explore os handlers** individualmente
2. **Entenda** como os DTOs funcionam
3. **Veja** como as validações são aplicadas
4. **Pratique** criando novos endpoints

---

> **💡 Dica:** A camada de apresentação deve ser **fininha** - apenas converter dados e orquestrar use cases. Toda lógica de negócio fica nas outras camadas!
