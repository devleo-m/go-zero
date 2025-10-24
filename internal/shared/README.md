# 🔧 Utilitários Compartilhados (Shared)

> **🎯 Objetivo:** Centralizar funcionalidades comuns que são utilizadas em múltiplos módulos, evitando duplicação de código e garantindo consistência.

## 📚 O que é a Pasta Shared?

A pasta **Shared** é o **arsenal** da sua aplicação! É aqui que:

- 🔧 **Centralizamos** utilitários comuns
- 🎯 **Padronizamos** comportamentos
- 🚀 **Reutilizamos** código entre módulos
- 🛡️ **Garantimos** consistência

## 🎓 Por que usar Shared?

### ❌ **Problema sem Shared:**
```go
// Código duplicado em cada módulo
// user/handler.go
func CreateUser(c *gin.Context) {
    c.JSON(200, gin.H{
        "success": true,
        "data": user,
    })
}

// product/handler.go  
func CreateProduct(c *gin.Context) {
    c.JSON(200, gin.H{
        "success": true,
        "data": product,
    })
}
```

### ✅ **Solução com Shared:**
```go
// shared/response/response.go
func Success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{
        Success: true,
        Data:    data,
    })
}

// user/handler.go
response.Success(c, user)

// product/handler.go
response.Success(c, product)
```

## 🏗️ Estrutura da Pasta Shared

```
shared/
├── 📄 README.md           # Este arquivo - conceitos de shared
├── 📁 response/           # 📤 Padronização de respostas HTTP
│   └── response.go        # Helpers para respostas
├── 📁 validation/         # ✅ Validações comuns
│   └── validators.go      # Validadores reutilizáveis
├── 📁 pagination/         # 📄 Utilitários de paginação
│   └── pagination.go      # Helpers de paginação
├── 📁 auth/              # 🔐 Autenticação e autorização
├── 📁 cache/             # 🗄️ Cache (Redis)
├── 📁 email/             # 📧 Serviço de email
├── 📁 queue/             # 📬 Filas de processamento
├── 📁 storage/           # 💾 Armazenamento (S3/MinIO)
└── 📁 websocket/         # 🔌 WebSocket hub
```

## 📤 Response - Padronização de Respostas

### **1. Estrutura Padrão**
```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page       int   `json:"page,omitempty"`
    Limit      int   `json:"limit,omitempty"`
    Total      int64 `json:"total,omitempty"`
    TotalPages int   `json:"total_pages,omitempty"`
}
```

### **2. Helpers de Resposta**

#### **✅ Success - Resposta de Sucesso**
```go
func Success(c *gin.Context, data interface{}, message ...string) {
    msg := ""
    if len(message) > 0 {
        msg = message[0]
    }
    
    c.JSON(http.StatusOK, Response{
        Success: true,
        Message: msg,
        Data:    data,
    })
}
```

#### **🆕 Created - Resposta de Criação**
```go
func Created(c *gin.Context, data interface{}, message ...string) {
    msg := "Created successfully"
    if len(message) > 0) {
        msg = message[0]
    }
    
    c.JSON(http.StatusCreated, Response{
        Success: true,
        Message: msg,
        Data:    data,
    })
}
```

#### **❌ Error - Resposta de Erro**
```go
func Error(c *gin.Context, statusCode int, errorCode string, message string) {
    c.JSON(statusCode, Response{
        Success: false,
        Error:   errorCode,
        Message: message,
    })
}

// Helpers específicos
func BadRequest(c *gin.Context, errorCode string, message string) {
    Error(c, http.StatusBadRequest, errorCode, message)
}

func NotFound(c *gin.Context, errorCode string, message string) {
    Error(c, http.StatusNotFound, errorCode, message)
}

func InternalServerError(c *gin.Context, errorCode string, message string) {
    Error(c, http.StatusInternalServerError, errorCode, message)
}
```

#### **📄 Paginated - Resposta Paginada**
```go
func Paginated(c *gin.Context, data interface{}, meta *Meta, message ...string) {
    msg := ""
    if len(message) > 0 {
        msg = message[0]
    }
    
    c.JSON(http.StatusOK, Response{
        Success: true,
        Message: msg,
        Data:    data,
        Meta:    meta,
    })
}
```

## ✅ Validation - Validações Comuns

### **1. Validação de Email**
```go
func ValidateEmail(email string) error {
    if email == "" {
        return ValidationError{Field: "email", Message: "Email is required"}
    }
    
    if !emailRegex.MatchString(email) {
        return ValidationError{Field: "email", Message: "Invalid email format"}
    }
    
    return nil
}
```

### **2. Validação de Senha Forte**
```go
func ValidatePassword(password string) error {
    if password == "" {
        return ValidationError{Field: "password", Message: "Password is required"}
    }
    
    if len(password) < 8 {
        return ValidationError{Field: "password", Message: "Password must be at least 8 characters long"}
    }
    
    // Verificar complexidade
    var hasUpper, hasLower, hasDigit, hasSpecial bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    // Validar cada critério
    if !hasUpper {
        return ValidationError{Field: "password", Message: "Password must contain at least one uppercase letter"}
    }
    // ... outros critérios
    
    return nil
}
```

### **3. Validação de UUID**
```go
func ValidateUUID(uuid string) error {
    if uuid == "" {
        return ValidationError{Field: "id", Message: "ID is required"}
    }
    
    if !uuidRegex.MatchString(uuid) {
        return ValidationError{Field: "id", Message: "Invalid ID format"}
    }
    
    return nil
}
```

### **4. Sanitização de Strings**
```go
func SanitizeString(input string) string {
    return strings.TrimSpace(input)
}
```

## 📄 Pagination - Utilitários de Paginação

### **1. Estrutura de Parâmetros**
```go
type Params struct {
    Page  int
    Limit int
    Sort  string
    Order string
}

type Result struct {
    Data       interface{} `json:"data"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    Total      int64       `json:"total"`
    TotalPages int         `json:"total_pages"`
    HasNext    bool        `json:"has_next"`
    HasPrev    bool        `json:"has_prev"`
}
```

### **2. Parse de Query String**
```go
func ParseFromQuery(c *gin.Context) *Params {
    page := parseInt(c.Query("page"), 1)
    limit := parseInt(c.Query("limit"), 10)
    sort := c.Query("sort")
    order := c.Query("order")
    
    // Validações básicas
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }
    if order != "asc" && order != "desc" {
        order = "asc"
    }
    
    return &Params{
        Page:  page,
        Limit: limit,
        Sort:  sort,
        Order: order,
    }
}
```

### **3. Cálculo de Offset**
```go
func (p *Params) Offset() int {
    return (p.Page - 1) * p.Limit
}
```

### **4. Criação de Resultado**
```go
func NewResult(data interface{}, total int64, params *Params) *Result {
    totalPages := int(total) / params.Limit
    if int(total)%params.Limit > 0 {
        totalPages++
    }
    
    return &Result{
        Data:       data,
        Page:       params.Page,
        Limit:      params.Limit,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    params.Page < totalPages,
        HasPrev:    params.Page > 1,
    }
}
```

## 🎯 Padrões Aplicados

### **1. DRY (Don't Repeat Yourself)**
- ✅ **Código centralizado** - Uma implementação, múltiplos usos
- ✅ **Manutenção única** - Mudança em um lugar afeta todos
- ✅ **Consistência** - Comportamento padronizado

### **2. Single Responsibility**
- ✅ **Response** - Apenas formatação de respostas
- ✅ **Validation** - Apenas validações
- ✅ **Pagination** - Apenas paginação

### **3. Reusability**
- ✅ **Módulos** - Podem ser usados em qualquer lugar
- ✅ **Flexibilidade** - Parâmetros configuráveis
- ✅ **Extensibilidade** - Fácil adicionar novas funcionalidades

## 🧪 Como Testar

### **Teste de Response:**
```go
func TestSuccessResponse(t *testing.T) {
    // Arrange
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    
    data := map[string]string{"name": "João"}
    
    // Act
    response.Success(c, data, "User created")
    
    // Assert
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "success")
    assert.Contains(t, w.Body.String(), "João")
}
```

### **Teste de Validation:**
```go
func TestValidateEmail(t *testing.T) {
    // Test cases
    testCases := []struct {
        email   string
        wantErr bool
    }{
        {"test@example.com", false},
        {"invalid-email", true},
        {"", true},
    }
    
    for _, tc := range testCases {
        err := validation.ValidateEmail(tc.email)
        if tc.wantErr {
            assert.Error(t, err)
        } else {
            assert.NoError(t, err)
        }
    }
}
```

## 🎓 Benefícios do Shared

### **1. Consistência**
- ✅ **Respostas** - Mesmo formato em toda aplicação
- ✅ **Validações** - Mesmas regras em todos os módulos
- ✅ **Paginação** - Comportamento uniforme

### **2. Manutenibilidade**
- ✅ **Mudanças centralizadas** - Uma alteração afeta tudo
- ✅ **Debugging** - Mais fácil encontrar problemas
- ✅ **Evolução** - Fácil adicionar novas funcionalidades

### **3. Produtividade**
- ✅ **Reutilização** - Não precisa reescrever código
- ✅ **Padrões** - Desenvolvedores seguem as mesmas práticas
- ✅ **Qualidade** - Código testado e validado

## 🚀 Próximos Passos

1. **Explore cada utilitário** individualmente
2. **Entenda** como são utilizados nos módulos
3. **Veja** como facilitam o desenvolvimento
4. **Pratique** criando novos utilitários

---

> **💡 Dica:** A pasta Shared deve conter apenas código **genérico** e **reutilizável**. Evite colocar lógica específica de negócio aqui!
