# Middleware Module / Módulo de Middleware

Intermediate components for HTTP and gRPC request/response processing in Go microservices / Componentes intermediários para processamento de requisições/respostas HTTP e gRPC em microserviços Go

## 🌍 Language / Idioma

• [🇺🇸 English](#english)  
• [🇧🇷 Português](#português)

## English

### 🎯 Purpose

This module provides reusable middleware components that act as intermediaries between HTTP/gRPC requests and responses in Go microservices, implementing cross-cutting concerns with modularity and composability.

### ✨ Key Features

• **Authentication**: JWT validation, API key verification, and OAuth integration  
• **Logging**: Structured request/response logging with correlation IDs  
• **CORS**: Cross-origin resource sharing configuration and handling  
• **Rate Limiting**: Request throttling and quota management  
• **Metrics**: Performance monitoring and observability collection  
• **Recovery**: Panic recovery and graceful error handling  
• **Validation**: Input sanitization and request validation  

### 📁 Typical Structure

```
internal/middleware/
├── auth.go              # Authentication middleware
├── logging.go           # Request/response logging
├── metrics.go           # Performance metrics collection
├── cors.go              # CORS policy enforcement
├── ratelimit.go         # Rate limiting and throttling
├── recovery.go          # Panic recovery handling
├── validation.go        # Input validation middleware
├── tracing.go           # Distributed tracing
└── chain.go             # Middleware composition utilities
```

### 🚀 Basic Usage Example

#### Gin Framework Integration

```go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/galafis/go-data-api-microservices/internal/logger"
)

// Logger middleware with structured logging
func Logger() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return logger.FormatRequest(param)
    })
}

// Auth middleware with JWT validation
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if !validateJWT(token) {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return
        }
        c.Next()
    }
}

// Rate limiter middleware
func RateLimit(requests int, duration time.Duration) gin.HandlerFunc {
    limiter := createRateLimiter(requests, duration)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatusJSON(429, gin.H{"error": "Rate limit exceeded"})
            return
        }
        c.Next()
    }
}
```

#### Middleware Composition

```go
func SetupMiddleware(r *gin.Engine) {
    r.Use(Logger())
    r.Use(Recovery())
    r.Use(CORS())
    r.Use(RateLimit(100, time.Minute))
    
    // Protected routes
    api := r.Group("/api/v1")
    api.Use(Auth())
    {
        api.GET("/users", handlers.GetUsers)
        api.POST("/users", handlers.CreateUser)
    }
}
```

### 📋 Best Practices

#### Modular Design
• **Single Responsibility**: Each middleware handles one specific concern  
• **Configurable**: Allow customization through configuration parameters  
• **Framework Agnostic**: Design for compatibility across different HTTP frameworks  

#### Performance & Composition
• **Order Matters**: Place lightweight middleware first (auth, logging, then heavy processing)  
• **Conditional Application**: Apply middleware only where needed  
• **Context Propagation**: Pass request context through the middleware chain  

#### Extension Guidelines
• **Interface Implementation**: Follow standard middleware patterns  
• **Error Handling**: Implement consistent error responses  
• **Testing**: Create unit tests for each middleware component  

---

## Português

### 🎯 Propósito

Este módulo fornece componentes de middleware reutilizáveis que atuam como intermediários entre requisições/respostas HTTP/gRPC em microserviços Go, implementando preocupações transversais com modularidade e composição.

### ✨ Funcionalidades Principais

• **Autenticação**: Validação JWT, verificação de chave API e integração OAuth  
• **Logging**: Log estruturado de requisições/respostas com IDs de correlação  
• **CORS**: Configuração e manipulação de compartilhamento de recursos entre origens  
• **Rate Limiting**: Throttling de requisições e gerenciamento de cotas  
• **Métricas**: Coleta de monitoramento de performance e observabilidade  
• **Recuperação**: Recuperação de panic e tratamento gracioso de erros  
• **Validação**: Sanitização de entrada e validação de requisições  

### 📁 Estrutura Típica

```
internal/middleware/
├── auth.go              # Middleware de autenticação
├── logging.go           # Log de requisições/respostas
├── metrics.go           # Coleta de métricas de performance
├── cors.go              # Aplicação de políticas CORS
├── ratelimit.go         # Rate limiting e throttling
├── recovery.go          # Manipulação de recuperação de panic
├── validation.go        # Middleware de validação de entrada
├── tracing.go           # Rastreamento distribuído
└── chain.go             # Utilitários de composição de middleware
```

### 🚀 Exemplo de Uso Básico

#### Integração com Framework Gin

```go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/galafis/go-data-api-microservices/internal/logger"
)

// Middleware de logger com logging estruturado
func Logger() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return logger.FormatRequest(param)
    })
}

// Middleware de auth com validação JWT
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if !validateJWT(token) {
            c.AbortWithStatusJSON(401, gin.H{"error": "Não autorizado"})
            return
        }
        c.Next()
    }
}

// Middleware de rate limiter
func RateLimit(requests int, duration time.Duration) gin.HandlerFunc {
    limiter := createRateLimiter(requests, duration)
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.AbortWithStatusJSON(429, gin.H{"error": "Limite de taxa excedido"})
            return
        }
        c.Next()
    }
}
```

#### Composição de Middleware

```go
func SetupMiddleware(r *gin.Engine) {
    r.Use(Logger())
    r.Use(Recovery())
    r.Use(CORS())
    r.Use(RateLimit(100, time.Minute))
    
    // Rotas protegidas
    api := r.Group("/api/v1")
    api.Use(Auth())
    {
        api.GET("/users", handlers.GetUsers)
        api.POST("/users", handlers.CreateUser)
    }
}
```

### 📋 Melhores Práticas

#### Design Modular
• **Responsabilidade Única**: Cada middleware trata uma preocupação específica  
• **Configurável**: Permitir customização através de parâmetros de configuração  
• **Framework Agnóstico**: Design para compatibilidade entre diferentes frameworks HTTP  

#### Performance e Composição
• **Ordem Importa**: Coloque middleware leve primeiro (auth, logging, depois processamento pesado)  
• **Aplicação Condicional**: Aplicar middleware apenas onde necessário  
• **Propagação de Contexto**: Passar contexto de requisição pela cadeia de middleware  

#### Diretrizes de Extensão
• **Implementação de Interface**: Seguir padrões padrão de middleware  
• **Tratamento de Erro**: Implementar respostas de erro consistentes  
• **Testes**: Criar testes unitários para cada componente de middleware  

---

## 👨‍💻 Author / Autor

**Gabriel Demetrios Lafis**  
Building robust and scalable microservices, one module at a time.  
Construindo microserviços robustos e escaláveis, um módulo por vez.

---

## 🚀 Welcome to the Journey!

### English
Welcome to the middleware module! This is where the magic of request processing happens. Every middleware you create is a guardian of your application, ensuring security, performance, and reliability. Think of middleware as the backbone of your microservices - they silently work behind the scenes to make everything flow smoothly. Remember: great middleware is invisible to users but invaluable to developers. Keep coding, keep protecting, and keep building the foundation that matters!

### Português
Bem-vindo ao módulo de middleware! Este é onde a mágica do processamento de requisições acontece. Cada middleware que você cria é um guardião da sua aplicação, garantindo segurança, performance e confiabilidade. Pense no middleware como a espinha dorsal dos seus microserviços - eles trabalham silenciosamente nos bastidores para fazer tudo fluir suavemente. Lembre-se: um grande middleware é invisível para usuários mas inestimável para desenvolvedores. Continue codificando, continue protegendo e continue construindo a fundação que importa!

Happy coding! / Bom código! 🎉
