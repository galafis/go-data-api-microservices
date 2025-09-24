# Validator Package | Pacote Validator

🚀 **Go Data API Microservices - Validation Package** 🚀

---

## 🇺🇸 English

### 📖 Overview

The **Validator Package** is a robust and extensible validation library designed specifically for Go microservices. It provides comprehensive data validation capabilities with custom validators and user-friendly error messages, making it perfect for API data validation in distributed systems.

### ✨ Key Features

- **Struct Validation**: Complete validation of Go structs with detailed error reporting
- **Field Validation**: Individual field validation with flexible tag support
- **Custom Validators**: Pre-built validators for common use cases:
  - 🆔 **UUID**: Validates universally unique identifiers
  - 📝 **Alpha Space**: Validates text with letters and spaces only
  - 📞 **Phone**: Validates international phone number formats
  - 🔐 **Password**: Validates strong password requirements
- **Error Handling**: Rich error messages with field-specific details
- **Microservice Ready**: Optimized for distributed architecture patterns

### 🚀 Quick Start

#### Installation

```go
import "github.com/galafis/go-data-api-microservices/pkg/validator"
```

#### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/galafis/go-data-api-microservices/pkg/validator"
)

type User struct {
    ID       string `json:"id" validate:"required,uuid"`
    Name     string `json:"name" validate:"required,alpha_space,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Phone    string `json:"phone" validate:"phone"`
    Password string `json:"password" validate:"required,password"`
    Role     string `json:"role" validate:"required,oneof=admin user guest"`
}

func main() {
    user := User{
        ID:       "123e4567-e89b-12d3-a456-426614174000",
        Name:     "John Doe",
        Email:    "john@example.com",
        Phone:    "+1234567890",
        Password: "SecurePass123!",
        Role:     "user",
    }

    if err := validator.Validate(user); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }

    fmt.Println("Validation successful!")
}
```

#### Variable Validation

```go
// Validate individual variables
email := "invalid-email"
if err := validator.ValidateVar(email, "email"); err != nil {
    fmt.Printf("Email validation failed: %v\n", err)
}

password := "weak"
if err := validator.ValidateVar(password, "password"); err != nil {
    fmt.Printf("Password validation failed: %v\n", err)
}
```

### 🏗️ File Structure

```
pkg/validator/
├── validator.go     # Main validation logic and custom validators
└── README.md        # This documentation file
```

### 🎯 Best Practices

1. **Consistent Tags**: Use consistent validation tags across your microservices
2. **Error Handling**: Always handle validation errors gracefully in your APIs
3. **Custom Messages**: Leverage the built-in error messages for better UX
4. **Performance**: The validator is optimized for high-throughput scenarios
5. **Security**: Use the password validator for secure authentication flows

---

## 🇧🇷 Português

### 📖 Visão Geral

O **Pacote Validator** é uma biblioteca de validação robusta e extensível, projetada especificamente para microserviços Go. Fornece capacidades abrangentes de validação de dados com validadores personalizados e mensagens de erro amigáveis, sendo perfeita para validação de dados de API em sistemas distribuídos.

### ✨ Funcionalidades Principais

- **Validação de Structs**: Validação completa de structs Go com relatórios detalhados de erros
- **Validação de Campos**: Validação individual de campos com suporte flexível a tags
- **Validadores Personalizados**: Validadores pré-construídos para casos de uso comuns:
  - 🆔 **UUID**: Valida identificadores únicos universais
  - 📝 **Alpha Space**: Valida texto apenas com letras e espaços
  - 📞 **Phone**: Valida formatos de números de telefone internacionais
  - 🔐 **Password**: Valida requisitos de senhas seguras
- **Tratamento de Erros**: Mensagens de erro ricas com detalhes específicos por campo
- **Pronto para Microserviços**: Otimizado para padrões de arquitetura distribuída

### 🚀 Início Rápido

#### Instalação

```go
import "github.com/galafis/go-data-api-microservices/pkg/validator"
```

#### Uso Básico

```go
package main

import (
    "fmt"
    "github.com/galafis/go-data-api-microservices/pkg/validator"
)

type Usuario struct {
    ID       string `json:"id" validate:"required,uuid"`
    Nome     string `json:"nome" validate:"required,alpha_space,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Telefone string `json:"telefone" validate:"phone"`
    Senha    string `json:"senha" validate:"required,password"`
    Papel    string `json:"papel" validate:"required,oneof=admin usuario convidado"`
}

func main() {
    usuario := Usuario{
        ID:       "123e4567-e89b-12d3-a456-426614174000",
        Nome:     "João Silva",
        Email:    "joao@exemplo.com",
        Telefone: "+5511987654321",
        Senha:    "SenhaSegura123!",
        Papel:    "usuario",
    }

    if err := validator.Validate(usuario); err != nil {
        fmt.Printf("Validação falhou: %v\n", err)
        return
    }

    fmt.Println("Validação bem-sucedida!")
}
```

#### Validação de Variáveis

```go
// Validar variáveis individuais
email := "email-inválido"
if err := validator.ValidateVar(email, "email"); err != nil {
    fmt.Printf("Validação de email falhou: %v\n", err)
}

senha := "fraca"
if err := validator.ValidateVar(senha, "password"); err != nil {
    fmt.Printf("Validação de senha falhou: %v\n", err)
}
```

### 🏗️ Estrutura de Arquivos

```
pkg/validator/
├── validator.go     # Lógica principal de validação e validadores personalizados
└── README.md        # Este arquivo de documentação
```

### 🎯 Melhores Práticas

1. **Tags Consistentes**: Use tags de validação consistentes em seus microserviços
2. **Tratamento de Erros**: Sempre trate erros de validação graciosamente em suas APIs
3. **Mensagens Personalizadas**: Aproveite as mensagens de erro integradas para melhor UX
4. **Performance**: O validador é otimizado para cenários de alto throughput
5. **Segurança**: Use o validador de senha para fluxos de autenticação seguros

---

## 🔧 Advanced Usage | Uso Avançado

### Custom Validation Tags | Tags de Validação Personalizadas

```go
// Available custom validators | Validadores personalizados disponíveis:
// - uuid: Validates UUID format | Valida formato UUID
// - alpha_space: Letters and spaces only | Apenas letras e espaços
// - phone: International phone format | Formato de telefone internacional
// - password: Strong password requirements | Requisitos de senha forte
```

### Error Response Structure | Estrutura de Resposta de Erro

```go
type ValidationError struct {
    Field   string `json:"field"`   // Field name | Nome do campo
    Tag     string `json:"tag"`     // Validation tag | Tag de validação
    Value   string `json:"value"`   // Invalid value | Valor inválido
    Message string `json:"message"` // Human-readable message | Mensagem legível
}
```

---

## 👨‍💻 Author | Autor

**Gabriel Demetrios Lafis**  
🚀 *Passionate Go Developer & Microservices Architect*  
🌟 *Desenvolvedor Go Apaixonado & Arquiteto de Microserviços*

---

## 💡 Welcome Message | Mensagem de Boas-Vindas

### 🇺🇸 To Fellow Developers

Welcome to the **Validator Package**! 🎉

As developers, we know that robust validation is the foundation of reliable microservices. This package was crafted with love and attention to detail to make your validation logic clean, maintainable, and powerful.

**Why you'll love this package:**
- 🎯 **Focus on Business Logic**: Spend less time writing validation code, more time building features
- 🚀 **Production Ready**: Battle-tested patterns for real-world applications
- 🔧 **Developer Experience**: Clear error messages and intuitive API design
- 📈 **Scalable**: Built for high-performance microservice architectures

Remember: Great software starts with great validation. Let's build something amazing together! 💪

### 🇧🇷 Para Desenvolvedores Companheiros

Bem-vindos ao **Pacote Validator**! 🎉

Como desenvolvedores, sabemos que validação robusta é a base de microserviços confiáveis. Este pacote foi criado com amor e atenção aos detalhes para tornar sua lógica de validação limpa, mantível e poderosa.

**Por que você vai amar este pacote:**
- 🎯 **Foco na Lógica de Negócio**: Gaste menos tempo escrevendo código de validação, mais tempo construindo funcionalidades
- 🚀 **Pronto para Produção**: Padrões testados em batalha para aplicações do mundo real
- 🔧 **Experiência do Desenvolvedor**: Mensagens de erro claras e design de API intuitivo
- 📈 **Escalável**: Construído para arquiteturas de microserviços de alta performance

Lembre-se: Grandes softwares começam com grande validação. Vamos construir algo incrível juntos! 💪

---

*Happy Coding! | Codificação Feliz!* 🚀✨
