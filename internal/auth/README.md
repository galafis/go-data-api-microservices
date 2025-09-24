# Auth Package | Pacote de Autenticação

🔐 Go Data API Microservices - Authentication Package 🔐

## 🇺🇸 English

### 📖 Overview

The Auth Package is a comprehensive authentication and authorization library designed specifically for Go microservices. It provides secure JWT-based authentication, password hashing utilities, and token management functionality, making it perfect for implementing robust security in distributed systems.

### ✨ Key Features

• **JWT Management**: Complete JWT token generation, validation, and parsing
• **Password Security**: Secure password hashing using bcrypt
• **Token Lifecycle**: Access and refresh token management
• **Claims Validation**: Custom JWT claims validation and extraction
• **Security Best Practices**: Industry-standard security implementations
• **Microservice Ready**: Optimized for distributed architecture patterns

### 🏗️ File Structure

```
internal/auth/
├── jwt.go       # JWT token generation, validation, and parsing
├── password.go  # Password hashing and verification utilities
└── README.md    # This documentation file
```

### 🚀 Quick Start

#### Basic Usage

```go
import "github.com/galafis/go-data-api-microservices/internal/auth"

// Password operations
hashed, err := auth.HashPassword("mySecurePassword123!")
if err != nil {
    log.Fatal(err)
}

valid := auth.CheckPassword("mySecurePassword123!", hashed)
if !valid {
    log.Fatal("Invalid password")
}

// JWT operations
token, err := auth.GenerateToken(userID, "user", time.Hour*24)
if err != nil {
    log.Fatal(err)
}

claims, err := auth.ValidateToken(token)
if err != nil {
    log.Fatal("Invalid token")
}
```

#### Middleware Integration

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := auth.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### 🎯 Best Practices

1. **Token Expiry**: Use short-lived access tokens with refresh token rotation
2. **Password Strength**: Enforce strong password policies in your application
3. **Secure Storage**: Never store plain-text passwords or tokens
4. **Token Validation**: Always validate tokens on protected routes
5. **Error Handling**: Implement proper error handling for authentication failures

## 🇧🇷 Português

### 📖 Visão Geral

O Pacote Auth é uma biblioteca abrangente de autenticação e autorização projetada especificamente para microserviços Go. Fornece autenticação segura baseada em JWT, utilitários de hash de senha e funcionalidade de gerenciamento de token, sendo perfeita para implementar segurança robusta em sistemas distribuídos.

### ✨ Funcionalidades Principais

• **Gerenciamento JWT**: Geração, validação e análise completa de tokens JWT
• **Segurança de Senha**: Hash seguro de senhas usando bcrypt
• **Ciclo de Vida do Token**: Gerenciamento de tokens de acesso e refresh
• **Validação de Claims**: Validação e extração de claims JWT personalizados
• **Melhores Práticas de Segurança**: Implementações de segurança padrão da indústria
• **Pronto para Microserviços**: Otimizado para padrões de arquitetura distribuída

### 🏗️ Estrutura de Arquivos

```
internal/auth/
├── jwt.go       # Geração, validação e análise de tokens JWT
├── password.go  # Utilitários de hash e verificação de senhas
└── README.md    # Este arquivo de documentação
```

### 🚀 Início Rápido

#### Uso Básico

```go
import "github.com/galafis/go-data-api-microservices/internal/auth"

// Operações de senha
hashed, err := auth.HashPassword("minhaSenhaSegura123!")
if err != nil {
    log.Fatal(err)
}

valid := auth.CheckPassword("minhaSenhaSegura123!", hashed)
if !valid {
    log.Fatal("Senha inválida")
}

// Operações JWT
token, err := auth.GenerateToken(userID, "user", time.Hour*24)
if err != nil {
    log.Fatal(err)
}

claims, err := auth.ValidateToken(token)
if err != nil {
    log.Fatal("Token inválido")
}
```

#### Integração com Middleware

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Cabeçalho de autorização obrigatório"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := auth.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Token inválido"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### 🎯 Melhores Práticas

1. **Expiração de Token**: Use tokens de acesso de curta duração com rotação de refresh token
2. **Força da Senha**: Aplique políticas de senhas fortes em sua aplicação
3. **Armazenamento Seguro**: Nunca armazene senhas ou tokens em texto simples
4. **Validação de Token**: Sempre valide tokens em rotas protegidas
5. **Tratamento de Erros**: Implemente tratamento adequado de erros para falhas de autenticação

## 🔧 Advanced Usage | Uso Avançado

### Custom Claims | Claims Personalizados

```go
// Custom claims structure | Estrutura de claims personalizados
type CustomClaims struct {
    UserID      string   `json:"user_id"`
    Role        string   `json:"role"`
    Permissions []string `json:"permissions"`
    jwt.StandardClaims
}
```

### Token Refresh Flow | Fluxo de Refresh Token

```go
// Refresh token validation | Validação de refresh token
func RefreshToken(refreshToken string) (string, error) {
    claims, err := auth.ValidateRefreshToken(refreshToken)
    if err != nil {
        return "", err
    }
    
    // Generate new access token | Gerar novo token de acesso
    return auth.GenerateAccessToken(claims.UserID, claims.Role, time.Hour)
}
```

## 👨‍💻 Author | Autor

**Gabriel Demetrios Lafis**
🚀 Passionate Go Developer & Microservices Architect  
🌟 Desenvolvedor Go Apaixonado & Arquiteto de Microserviços

## 💡 Welcome Message | Mensagem de Boas-Vindas

### 🇺🇸 To Fellow Developers

Welcome to the Auth Package! 🎉

Security is the foundation of trust in any system. This package was crafted with security-first principles and battle-tested patterns to protect your microservices and users.

Why you'll love this package:
• 🔒 **Security First**: Industry-standard encryption and security practices
• 🚀 **Production Ready**: Battle-tested in high-traffic environments
• 🔧 **Developer Friendly**: Clean APIs and comprehensive documentation
• 📈 **Scalable**: Built for distributed microservice architectures

Remember: Security is not a feature, it's a requirement. Let's build secure systems together! 💪

### 🇧🇷 Para Desenvolvedores Companheiros

Bem-vindos ao Pacote Auth! 🎉

Segurança é a base da confiança em qualquer sistema. Este pacote foi criado com princípios de segurança em primeiro lugar e padrões testados em batalha para proteger seus microserviços e usuários.

Por que você vai amar este pacote:
• 🔒 **Segurança em Primeiro Lugar**: Práticas de criptografia e segurança padrão da indústria
• 🚀 **Pronto para Produção**: Testado em batalha em ambientes de alto tráfego
• 🔧 **Amigável ao Desenvolvedor**: APIs limpas e documentação abrangente
• 📈 **Escalável**: Construído para arquiteturas de microserviços distribuídos

Lembre-se: Segurança não é uma funcionalidade, é um requisito. Vamos construir sistemas seguros juntos! 💪

Happy Coding! | Codificação Feliz! 🚀✨
