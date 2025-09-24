# Logger Package / Pacote de Logging

## English

### Purpose
The `logger` package provides a centralized, configurable logging system for the Go Data API Microservices project. It's built on top of the [Logrus](https://github.com/sirupsen/logrus) library and offers structured logging with multiple output formats and levels.

### Description
This package serves as the unified logging solution across all microservices, providing consistent log formatting, configurable output destinations, and specialized middleware for HTTP request logging with the Gin framework.

### Features
- **Structured Logging**: JSON and Text format support
- **Configurable Output**: stdout, stderr, or file output
- **Multiple Log Levels**: Debug, Info, Warn, Error, Fatal, Panic
- **Gin Middleware**: Built-in HTTP request logging
- **Context Support**: Request-scoped logging with fields
- **Thread-Safe**: Safe for concurrent use across goroutines

### Expected Structure
```
logger/
├── logger.go    # Main logger implementation
└── README.md    # This documentation
```

### Configuration
The logger accepts a `Config` struct with the following fields:
- `Level`: Log level (debug, info, warn, error, fatal, panic)
- `Format`: Output format (json, text)
- `Output`: Output destination (stdout, stderr, or file path)
- `TimeFormat`: Timestamp format (RFC3339, etc.)

### Usage

#### Basic Usage
```go
import "github.com/galafis/go-data-api-microservices/pkg/logger"

// Basic logging
logger.Info("Application started")
logger.Errorf("Failed to connect: %v", err)

// With fields
logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action":  "login",
}).Info("User authenticated")
```

#### Configuration
```go
config := &logger.Config{
    Level:      "info",
    Format:     "json",
    Output:     "stdout",
    TimeFormat: time.RFC3339,
}
logger.Configure(config)
```

#### Gin Middleware
```go
router := gin.New()
router.Use(logger.GinLogger())
```

### Best Practices
- Use appropriate log levels (Debug for development, Info for production)
- Include relevant context with `WithFields()` for better traceability
- Avoid logging sensitive information (passwords, tokens)
- Use structured logging in production environments

---

## Português

### Propósito
O pacote `logger` fornece um sistema de logging centralizado e configurável para o projeto Go Data API Microservices. É construído sobre a biblioteca [Logrus](https://github.com/sirupsen/logrus) e oferece logging estruturado com múltiplos formatos de saída e níveis.

### Descrição
Este pacote serve como a solução unificada de logging em todos os microsserviços, fornecendo formatação consistente de logs, destinos de saída configuráveis e middleware especializado para logging de requisições HTTP com o framework Gin.

### Funcionalidades
- **Logging Estruturado**: Suporte aos formatos JSON e Text
- **Saída Configurável**: stdout, stderr ou arquivo
- **Múltiplos Níveis de Log**: Debug, Info, Warn, Error, Fatal, Panic
- **Middleware Gin**: Logging integrado de requisições HTTP
- **Suporte a Contexto**: Logging com escopo de requisição e campos
- **Thread-Safe**: Seguro para uso concorrente entre goroutines

### Estrutura Esperada
```
logger/
├── logger.go    # Implementação principal do logger
└── README.md    # Esta documentação
```

### Como Usar no Projeto

#### Uso Básico
```go
import "github.com/galafis/go-data-api-microservices/pkg/logger"

// Logging básico
logger.Info("Aplicação iniciada")
logger.Errorf("Falha na conexão: %v", err)

// Com campos
logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action":  "login",
}).Info("Usuário autenticado")
```

#### Configuração
```go
config := &logger.Config{
    Level:      "info",
    Format:     "json",
    Output:     "stdout",
    TimeFormat: time.RFC3339,
}
logger.Configure(config)
```

### Boas Práticas
- Use níveis apropriados de log (Debug para desenvolvimento, Info para produção)
- Inclua contexto relevante com `WithFields()` para melhor rastreabilidade
- Evite logar informações sensíveis (senhas, tokens)
- Use logging estruturado em ambientes de produção

---

## Author / Autor
**Gabriel Demetrios Lafis**

## Onboarding Message / Mensagem de Integração
🔍 **Logging is the window into your application's soul** / **O logging é a janela para a alma da sua aplicação**

Effective logging is crucial for debugging, monitoring, and maintaining microservices. This package provides you with powerful, flexible logging capabilities that scale with your application's needs.

O logging efetivo é crucial para debug, monitoramento e manutenção de microsserviços. Este pacote fornece capacidades de logging poderosas e flexíveis que escalam com as necessidades da sua aplicação.
