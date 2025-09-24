# Models Module / Módulo de Models

Persistent data structures and types definition for Go microservices / Definição de estruturas e tipos de dados persistentes para microserviços Go

## 🌍 Language / Idioma

• [🇺🇸 English](#english)
• [🇧🇷 Português](#português)

## English

### 🎯 Purpose

This module defines persistent data structures and types (ORM, schema) for Go microservices, providing a centralized repository for all data models, database mappings, and type definitions used across the application.

### ✨ Key Features

• **Struct Definitions**: Clean, well-documented Go structs representing business entities
• **Database Tags**: GORM tags for ORM mapping and database schema generation
• **JSON Serialization**: Proper JSON tags for API responses and requests
• **Schema Versioning**: Migration-friendly structures with version control
• **Type Safety**: Strong typing with validation and constraints
• **Relationships**: Foreign keys and associations between entities
• **Auditing**: Created/updated timestamps and soft delete support

### 📁 Typical Structure

```
internal/models/
├── user.go              # User entity and related structures
├── dataset.go           # Dataset management models
├── query.go             # Query execution and result models
├── analytics.go         # Analytics and metrics models
├── base.go              # Common base model with timestamps
├── dto/                 # Data Transfer Objects
│   ├── user_dto.go      # User API request/response DTOs
│   └── dataset_dto.go   # Dataset API DTOs
└── migrations/          # Database schema migrations
    ├── 001_initial.sql
    └── 002_add_indexes.sql
```

### 🚀 Basic Usage Example

#### Model Definition with GORM

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Base model with common fields
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a system user
type User struct {
    BaseModel
    Username string    `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
    Email    string    `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
    FullName string    `gorm:"size:100" json:"full_name"`
    Active   bool      `gorm:"default:true" json:"active"`
    Datasets []Dataset `gorm:"foreignKey:UserID" json:"datasets,omitempty"`
}

// Dataset represents a data collection
type Dataset struct {
    BaseModel
    Name        string `gorm:"size:100;not null" json:"name" validate:"required"`
    Description string `gorm:"type:text" json:"description"`
    UserID      uint   `gorm:"not null;index" json:"user_id"`
    User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
    Status      string `gorm:"size:20;default:'active'" json:"status"`
}
```

#### DTO Separation

```go
package dto

// UserCreateRequest for user creation API
type UserCreateRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    FullName string `json:"full_name" validate:"max=100"`
}

// UserResponse for API responses
type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    FullName string `json:"full_name"`
    Active   bool   `json:"active"`
}

// ToResponse converts User model to UserResponse DTO
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
        FullName: u.FullName,
        Active:   u.Active,
    }
}
```

### 📋 Best Practices

#### Model/DTO Separation

• **Clear Boundaries**: Keep domain models separate from API DTOs
• **Conversion Methods**: Implement ToDTO/FromDTO methods for clean transformation
• **Validation**: Use struct tags for input validation in DTOs

#### Safe Type Usage

• **Pointer Fields**: Use pointers for optional fields to distinguish between zero values and nil
• **Constraints**: Leverage database constraints and Go validation tags
• **Enums**: Define constants for status fields and use custom types when needed

#### Semantic Comments

• **Purpose Documentation**: Clearly describe what each model represents
• **Field Descriptions**: Document business rules and constraints
• **Relationship Mapping**: Explain foreign key relationships and their purpose

## Português

### 🎯 Propósito

Este módulo define estruturas e tipos de dados persistentes (ORM, schema) para microserviços Go, fornecendo um repositório centralizado para todos os modelos de dados, mapeamentos de banco de dados e definições de tipos usados na aplicação.

### ✨ Funcionalidades Principais

• **Definições de Struct**: Structs Go limpos e bem documentados representando entidades de negócio
• **Tags de Banco de Dados**: Tags GORM para mapeamento ORM e geração de schema de banco
• **Serialização JSON**: Tags JSON adequadas para respostas e requisições de API
• **Versionamento de Schema**: Estruturas amigáveis à migração com controle de versão
• **Segurança de Tipos**: Tipagem forte com validação e restrições
• **Relacionamentos**: Chaves estrangeiras e associações entre entidades
• **Auditoria**: Timestamps de criação/atualização e suporte a soft delete

### 📁 Estrutura Típica

```
internal/models/
├── user.go              # Entidade de usuário e estruturas relacionadas
├── dataset.go           # Modelos de gerenciamento de dataset
├── query.go             # Modelos de execução de query e resultados
├── analytics.go         # Modelos de analytics e métricas
├── base.go              # Modelo base comum com timestamps
├── dto/                 # Objetos de Transferência de Dados
│   ├── user_dto.go      # DTOs de requisição/resposta de usuário
│   └── dataset_dto.go   # DTOs de dataset
└── migrations/          # Migrações de schema de banco de dados
    ├── 001_initial.sql
    └── 002_add_indexes.sql
```

### 🚀 Exemplo de Uso Básico

#### Definição de Model com GORM

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// Modelo base com campos comuns
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User representa um usuário do sistema
type User struct {
    BaseModel
    Username string    `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,min=3,max=50"`
    Email    string    `gorm:"uniqueIndex;size:100;not null" json:"email" validate:"required,email"`
    FullName string    `gorm:"size:100" json:"full_name"`
    Active   bool      `gorm:"default:true" json:"active"`
    Datasets []Dataset `gorm:"foreignKey:UserID" json:"datasets,omitempty"`
}

// Dataset representa uma coleção de dados
type Dataset struct {
    BaseModel
    Name        string `gorm:"size:100;not null" json:"name" validate:"required"`
    Description string `gorm:"type:text" json:"description"`
    UserID      uint   `gorm:"not null;index" json:"user_id"`
    User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
    Status      string `gorm:"size:20;default:'active'" json:"status"`
}
```

#### Separação de DTO

```go
package dto

// UserCreateRequest para API de criação de usuário
type UserCreateRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    FullName string `json:"full_name" validate:"max=100"`
}

// UserResponse para respostas de API
type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    FullName string `json:"full_name"`
    Active   bool   `json:"active"`
}

// ToResponse converte User model para UserResponse DTO
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:       u.ID,
        Username: u.Username,
        Email:    u.Email,
        FullName: u.FullName,
        Active:   u.Active,
    }
}
```

### 📋 Melhores Práticas

#### Separação Model/DTO

• **Fronteiras Claras**: Mantenha modelos de domínio separados dos DTOs de API
• **Métodos de Conversão**: Implemente métodos ToDTO/FromDTO para transformação limpa
• **Validação**: Use struct tags para validação de entrada em DTOs

#### Uso Seguro de Tipos

• **Campos Pointer**: Use pointers para campos opcionais para distinguir entre valores zero e nil
• **Restrições**: Aproveite restrições de banco de dados e tags de validação Go
• **Enums**: Defina constantes para campos de status e use tipos customizados quando necessário

#### Comentários Semânticos

• **Documentação de Propósito**: Descreva claramente o que cada modelo representa
• **Descrições de Campo**: Documente regras de negócio e restrições
• **Mapeamento de Relacionamentos**: Explique relacionamentos de chave estrangeira e seu propósito

## 👨‍💻 Author / Autor

**Gabriel Demetrios Lafis**

Building robust and scalable microservices, one module at a time.
Construindo microserviços robustos e escaláveis, um módulo por vez.

## 🚀 Welcome to the Data Foundation!

### English

Welcome to the models module - the foundation of your data architecture! Every struct you define here becomes the blueprint for your application's data flow. Think of models as the DNA of your microservices - they carry the essential information that brings your application to life. Your careful attention to relationships, constraints, and type safety here will pay dividends throughout the entire application lifecycle. Remember: well-designed models lead to robust APIs, efficient queries, and maintainable code. Keep modeling, keep structuring, and keep building the data foundation that empowers everything else!

### Português

Bem-vindo ao módulo de models - a fundação da sua arquitetura de dados! Cada struct que você define aqui se torna o blueprint para o fluxo de dados da sua aplicação. Pense nos models como o DNA dos seus microserviços - eles carregam as informações essenciais que dão vida à sua aplicação. Sua atenção cuidadosa aos relacionamentos, restrições e segurança de tipos aqui trará dividendos durante todo o ciclo de vida da aplicação. Lembre-se: models bem projetados levam a APIs robustas, queries eficientes e código maintível. Continue modelando, continue estruturando e continue construindo a fundação de dados que empodera tudo mais!

Happy modeling! / Boa modelagem! 🎯
