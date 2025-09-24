# Handlers Module / Módulo de Handlers

HTTP and gRPC request/response handlers for Go microservices / Manipuladores de requisições/respostas HTTP e gRPC para microserviços Go

## 🌍 Language / Idioma

• [🇺🇸 English](#english)  
• [🇧🇷 Português](#português)

---

## English

### 🎯 Purpose

This module provides HTTP and gRPC handler functions to manage requests and responses in Go microservices, implementing clean architecture patterns with proper separation of concerns.

### ✨ Features

• **HTTP Handlers**: RESTful API endpoints with proper HTTP status codes and responses  
• **gRPC Handlers**: High-performance RPC service implementations  
• **Middleware Integration**: Authentication, logging, rate limiting, and CORS support  
• **Input Validation**: Request payload validation and sanitization  
• **Error Handling**: Standardized error responses with proper HTTP status codes  
• **Context Propagation**: Full context support for timeouts, cancellation, and tracing  
• **Framework Integration**: Compatible with Gin, Echo, Fiber, and standard net/http

### 📁 Expected Structure

```
internal/handlers/
├── http/                  # HTTP REST handlers
│   ├── middleware/        # HTTP middleware components
│   ├── routes.go         # Route definitions and setup
│   └── server.go         # HTTP server configuration
├── grpc/                 # gRPC service handlers
│   ├── interceptors/     # gRPC interceptors
│   ├── services/         # Service implementations
│   └── server.go         # gRPC server configuration
├── auth.go              # Authentication handlers
├── user.go              # User management handlers
├── analytics.go         # Analytics and metrics handlers
├── dataset.go           # Data processing handlers
├── query.go             # Query execution handlers
├── validation.go        # Input validation helpers
└── response.go          # Response formatting utilities
```

### 🚀 Basic Usage Example

#### HTTP Handler with Gin Framework

```go
package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/galafis/go-data-api-microservices/internal/services"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID",
        })
        return
    }
    
    user, err := h.userService.GetByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch user",
        })
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request payload",
        })
        return
    }
    
    user, err := h.userService.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create user",
        })
        return
    }
    
    c.JSON(http.StatusCreated, user)
}
```

#### Route Setup

```go
func SetupRoutes(r *gin.Engine, userHandler *UserHandler) {
    api := r.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.GET("/:id", userHandler.GetUser)
            users.POST("/", userHandler.CreateUser)
        }
    }
}
```

### 📋 Best Practices

#### Request Handling
• **Validate Input**: Always validate and sanitize request payloads  
• **Use Context**: Propagate context for timeouts and cancellation  
• **Handle Errors Gracefully**: Return appropriate HTTP status codes with meaningful messages

#### Response Management
• **Consistent Format**: Use standardized response structures  
• **Proper Status Codes**: Return appropriate HTTP status codes  
• **Security Headers**: Include necessary security headers in responses

#### Performance & Scalability
• **Middleware Optimization**: Use efficient middleware for common operations  
• **Connection Pooling**: Properly manage database and external service connections  
• **Async Processing**: Use goroutines for non-blocking operations when appropriate

---

## Português

### 🎯 Propósito

Este módulo fornece funções de manipulação HTTP e gRPC para gerenciar requisições e respostas em microserviços Go, implementando padrões de arquitetura limpa com separação adequada de responsabilidades.

### ✨ Funcionalidades

• **Handlers HTTP**: Endpoints de API RESTful com códigos de status HTTP e respostas adequadas  
• **Handlers gRPC**: Implementações de serviços RPC de alta performance  
• **Integração de Middleware**: Suporte a autenticação, logging, rate limiting e CORS  
• **Validação de Entrada**: Validação e sanitização de payloads de requisição  
• **Tratamento de Erros**: Respostas padronizadas de erro com códigos de status HTTP apropriados  
• **Propagação de Contexto**: Suporte completo a contexto para timeouts, cancelamento e tracing  
• **Integração com Frameworks**: Compatível com Gin, Echo, Fiber e net/http padrão

### 📁 Estrutura Esperada

```
internal/handlers/
├── http/                  # Handlers HTTP REST
│   ├── middleware/        # Componentes de middleware HTTP
│   ├── routes.go         # Definições e configuração de rotas
│   └── server.go         # Configuração do servidor HTTP
├── grpc/                 # Handlers de serviços gRPC
│   ├── interceptors/     # Interceptadores gRPC
│   ├── services/         # Implementações de serviços
│   └── server.go         # Configuração do servidor gRPC
├── auth.go              # Handlers de autenticação
├── user.go              # Handlers de gerenciamento de usuários
├── analytics.go         # Handlers de analytics e métricas
├── dataset.go           # Handlers de processamento de dados
├── query.go             # Handlers de execução de consultas
├── validation.go        # Auxiliares de validação de entrada
└── response.go          # Utilitários de formatação de resposta
```

### 🚀 Exemplo de Uso Básico

#### Handler HTTP com Framework Gin

```go
package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/galafis/go-data-api-microservices/internal/services"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

// GetUser trata GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID de usuário inválido",
        })
        return
    }
    
    user, err := h.userService.GetByID(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Falha ao buscar usuário",
        })
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// CreateUser trata POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Payload de requisição inválido",
        })
        return
    }
    
    user, err := h.userService.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Falha ao criar usuário",
        })
        return
    }
    
    c.JSON(http.StatusCreated, user)
}
```

#### Configuração de Rotas

```go
func SetupRoutes(r *gin.Engine, userHandler *UserHandler) {
    api := r.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.GET("/:id", userHandler.GetUser)
            users.POST("/", userHandler.CreateUser)
        }
    }
}
```

### 📋 Melhores Práticas

#### Manipulação de Requisições
• **Validar Entrada**: Sempre validar e sanitizar payloads de requisição  
• **Usar Contexto**: Propagar contexto para timeouts e cancelamento  
• **Tratar Erros Graciosamente**: Retornar códigos de status HTTP apropriados com mensagens significativas

#### Gerenciamento de Respostas
• **Formato Consistente**: Usar estruturas de resposta padronizadas  
• **Códigos de Status Adequados**: Retornar códigos de status HTTP apropriados  
• **Headers de Segurança**: Incluir headers de segurança necessários nas respostas

#### Performance e Escalabilidade
• **Otimização de Middleware**: Usar middleware eficiente para operações comuns  
• **Pool de Conexões**: Gerenciar adequadamente conexões de banco de dados e serviços externos  
• **Processamento Assíncrono**: Usar goroutines para operações não-bloqueantes quando apropriado

---

## 👨‍💻 Author / Autor

**Gabriel Demetrios Lafis**  
Building robust and scalable microservices, one module at a time.  
Construindo microserviços robustos e escaláveis, um módulo por vez.

---

## 🚀 Welcome to the Journey!

### English

Welcome to the handlers module! This is where your APIs come to life. Every endpoint you create, every request you handle, and every response you craft is a bridge connecting users to your application's power. Remember: great handlers don't just process data – they create experiences. Keep coding, keep learning, and keep building bridges that matter!

### Português

Bem-vindo ao módulo de handlers! Este é o lugar onde suas APIs ganham vida. Cada endpoint que você cria, cada requisição que manipula e cada resposta que constrói é uma ponte conectando usuários ao poder da sua aplicação. Lembre-se: grandes handlers não apenas processam dados – eles criam experiências. Continue codificando, continue aprendendo e continue construindo pontes que importam!

**Happy coding! / Bom código!** 🎉
