# 🎯 O QUE VAMOS FAZER NA ETAPA 3?
- Vamos criar um sistema profissional de configurações e logs para nossa aplicação. É como preparar a "casa" antes de começar a construir os móveis!
## 📚 CONCEITOS QUE VAMOS APRENDER
1. Variáveis de Ambiente 🌍
O QUE É: Configurações que mudam entre ambientes (desenvolvimento, produção, teste)
ANALOGIA: Pense como as "configurações do seu celular":
- Em casa: WiFi da sua casa
- No trabalho: WiFi da empresa
- No café: WiFi público
POR QUE USAR:
* ✅ Segurança (senhas não ficam no código)
* ✅ Flexibilidade (mesmo código, ambientes diferentes)
* ✅ Boas práticas (12-factor app)

2. Sistema de Logs Estruturados 📝
O QUE É: Logs organizados que facilitam debugging e monitoramento
ANALOGIA: Pense como um "diário detalhado":
* ❌ Log ruim: "Erro aconteceu"
* ✅ Log bom: "ERRO: Falha ao criar usuário - ID: 123 - Tempo: 2.3s - IP: 192.168.1.1"
POR QUE USAR:
* ✅ Debug mais fácil
* ✅ Monitoramento em produção
* ✅ Rastreamento de problemas
3. Validação de Configurações ✅
O QUE É: Garantir que todas as configurações necessárias estão presentes
ANALOGIA: Como verificar se você tem tudo antes de sair de casa:
* Chaves? ✅
* Carteira? ✅
* Celular? ✅

🏗️ ONDE ISSO SE ENCAIXA NA ARQUITETURA?
```
┌─────────────────────────────────────┐
│        INFRASTRUCTURE LAYER         │
│  ┌─────────────┐  ┌─────────────┐   │
│  │   CONFIG    │  │   LOGGER    │   │
│  │             │  │             │   │
│  │ • .env      │  │ • Zap       │   │
│  │ • Viper     │  │ • Levels    │   │
│  │ • Validate  │  │ • JSON      │   │
│  └─────────────┘  └─────────────┘   │
└─────────────────────────────────────┘
```

# 🎯 O QUE VAMOS CRIAR?
1. Sistema de Configuração (.env + Viper)
```
// Estrutura que vamos criar
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
}
```

2. Sistema de Logs (Zap Logger)
```
// Logs estruturados que vamos implementar
logger.Info("servidor iniciado",
    zap.String("port", "8080"),
    zap.String("env", "development"),
)

logger.Error("falha ao conectar banco",
    zap.Error(err),
    zap.String("host", "localhost"),
)
```

3. Validação de Configs
```
// Verificar se tudo está configurado
func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return errors.New("DATABASE_HOST é obrigatório")
    }
    // ... outras validações
}
```

# 🔄 FLUXO QUE VAMOS IMPLEMENTAR
1. Aplicação inicia
   │
   ▼
2. Carrega .env
   │
   ▼  
3. Valida configurações
   │
   ▼
4. Inicializa logger
   │
   ▼
5. Conecta banco/redis
   │
   ▼
6. Inicia servidor

# 🎓 POR QUE ISSO É IMPORTANTE?
SEM um sistema de configs:
* ❌ Senhas no código (inseguro!)
* ❌ Diferentes ambientes = código diferente
* ❌ Debugging difícil
* ❌ Deploy manual e propenso a erros
COM um sistema de configs:
* ✅ Segurança (senhas em variáveis)
* ✅ Mesmo código, ambientes diferentes
* ✅ Logs detalhados para debug
* ✅ Deploy automatizado e confiável