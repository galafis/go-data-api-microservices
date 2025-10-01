# Database Module / Módulo de Banco de Dados

*A comprehensive database layer for Go microservices / Uma camada abrangente de banco de dados para microserviços Go*

---

## 🌍 Language / Idioma
- [🇺🇸 English](#english)
- [🇧🇷 Português](#português)

---

## English

### 🎯 Purpose

This module provides a robust and standardized database layer for Go microservices, handling connections, operations, and best practices for database management in distributed systems.

### ✨ Features

- **Connection Management**: Efficient database connection pooling and lifecycle management
- **Migration Helpers**: Automated database schema migration utilities
- **Transaction Pattern**: Standardized transaction handling with proper rollback mechanisms
- **Context Support**: Full context propagation for timeouts and cancellation
- **Error Handling**: Comprehensive error wrapping and logging
- **Connection Pool**: Optimized connection pool configuration
- **Health Checks**: Database connectivity monitoring

### 📁 Expected Structure

```
internal/database/
├── connection.go      # Database connection management
├── migrations/        # Database migration files
│   ├── migrate.go     # Migration runner
│   └── sql/          # SQL migration files
├── models/           # Database models and schemas
├── repositories/     # Repository pattern implementations
├── transactions.go   # Transaction management utilities
├── config.go        # Database configuration
└── health.go        # Health check implementations
```

### 🚀 Basic Usage Example

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/galafis/go-data-api-microservices/internal/database"
    _ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
    // Initialize database connection
    config := database.Config{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        Username: "user",
        Password: "password",
        MaxConns: 25,
        MaxIdle:  5,
    }
    
    db, err := database.NewConnection(config)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Execute query with context
    var count int
    err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
    if err != nil {
        log.Printf("Query failed: %v", err)
        return
    }
    
    log.Printf("Total users: %d", count)
    
    // Transaction example
    err = database.WithTransaction(ctx, db, func(tx *sql.Tx) error {
        _, err := tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "John", "john@example.com")
        if err != nil {
            return err
        }
        
        _, err = tx.ExecContext(ctx, "UPDATE user_stats SET count = count + 1")
        return err
    })
    
    if err != nil {
        log.Printf("Transaction failed: %v", err)
    }
}
```

### 📋 Best Practices

#### Connection Pool Management
- Configure appropriate `MaxOpenConns` and `MaxIdleConns`
- Set reasonable `ConnMaxLifetime` to prevent stale connections
- Monitor connection pool metrics

#### Error Handling
- Always wrap database errors with context
- Implement proper retry logic for transient failures
- Log errors with sufficient detail for debugging

#### Security & Compliance
- Use parameterized queries to prevent SQL injection
- Implement proper access controls and audit logging
- Encrypt sensitive data at rest and in transit
- Follow GDPR/LGPD compliance for personal data handling

#### Performance
- Use appropriate indexes for frequent queries
- Implement query timeouts via context
- Consider read replicas for read-heavy workloads
- Monitor and optimize slow queries

---

## Português

### 🎯 Propósito

Este módulo fornece uma camada de banco de dados robusta e padronizada para microserviços Go, gerenciando conexões, operações e melhores práticas para gerenciamento de banco de dados em sistemas distribuídos.

### ✨ Funcionalidades

- **Gerenciamento de Conexões**: Pool de conexões eficiente e gerenciamento do ciclo de vida
- **Auxiliares de Migração**: Utilitários automatizados de migração de esquema de banco de dados
- **Padrão de Transação**: Manipulação padronizada de transações com mecanismos adequados de rollback
- **Suporte a Contexto**: Propagação completa de contexto para timeouts e cancelamento
- **Tratamento de Erros**: Encapsulamento abrangente de erros e logging
- **Pool de Conexões**: Configuração otimizada do pool de conexões
- **Verificações de Saúde**: Monitoramento de conectividade do banco de dados

### 📁 Estrutura Esperada

```
internal/database/
├── connection.go      # Gerenciamento de conexão do banco
├── migrations/        # Arquivos de migração do banco
│   ├── migrate.go     # Executor de migrações
│   └── sql/          # Arquivos SQL de migração
├── models/           # Modelos e esquemas do banco
├── repositories/     # Implementações do padrão repository
├── transactions.go   # Utilitários de gerenciamento de transações
├── config.go        # Configuração do banco de dados
└── health.go        # Implementações de verificação de saúde
```

### 🚀 Exemplo de Uso Básico

```go
package main

import (
    "context"
    "database/sql"
    "log"
    "time"
    
    "github.com/galafis/go-data-api-microservices/internal/database"
    _ "github.com/lib/pq" // Driver PostgreSQL
)

func main() {
    // Inicializar conexão com banco de dados
    config := database.Config{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        Username: "user",
        Password: "password",
        MaxConns: 25,
        MaxIdle:  5,
    }
    
    db, err := database.NewConnection(config)
    if err != nil {
        log.Fatal("Falha ao conectar com o banco:", err)
    }
    defer db.Close()
    
    // Criar contexto com timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Executar query com contexto
    var count int
    err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
    if err != nil {
        log.Printf("Query falhou: %v", err)
        return
    }
    
    log.Printf("Total de usuários: %d", count)
    
    // Exemplo de transação
    err = database.WithTransaction(ctx, db, func(tx *sql.Tx) error {
        _, err := tx.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", "João", "joao@exemplo.com")
        if err != nil {
            return err
        }
        
        _, err = tx.ExecContext(ctx, "UPDATE user_stats SET count = count + 1")
        return err
    })
    
    if err != nil {
        log.Printf("Transação falhou: %v", err)
    }
}
```

### 📋 Melhores Práticas

#### Gerenciamento do Pool de Conexões
- Configure `MaxOpenConns` e `MaxIdleConns` apropriados
- Defina `ConnMaxLifetime` razoável para evitar conexões obsoletas
- Monitore métricas do pool de conexões

#### Tratamento de Erros
- Sempre encapsule erros de banco com contexto
- Implemente lógica de retry adequada para falhas transitórias
- Registre erros com detalhes suficientes para debugging

#### Segurança e Compliance
- Use queries parametrizadas para prevenir injeção SQL
- Implemente controles de acesso adequados e audit logging
- Criptografe dados sensíveis em repouso e em trânsito
- Siga compliance GDPR/LGPD para tratamento de dados pessoais

#### Performance
- Use índices apropriados para queries frequentes
- Implemente timeouts de query via contexto
- Considere réplicas de leitura para cargas intensivas de leitura
- Monitore e otimize queries lentas

---

## 👨‍💻 Author / Autor

**Gabriel Demetrios Lafis**

*Building robust and scalable microservices, one module at a time.*

*Construindo microserviços robustos e escaláveis, um módulo por vez.*

---

## 🚀 Welcome to the Journey!

### English
Welcome to the database module! This is where data meets reliability. Every query you write, every transaction you manage, and every connection you establish is a building block toward creating something extraordinary. Remember: great software isn't just about code – it's about solving real problems and making a positive impact. Keep coding, keep learning, and keep pushing the boundaries of what's possible!

### Português
Bem-vindo ao módulo de banco de dados! Este é o lugar onde os dados encontram a confiabilidade. Cada query que você escreve, cada transação que gerencia e cada conexão que estabelece é um bloco de construção para criar algo extraordinário. Lembre-se: um ótimo software não é apenas sobre código – é sobre resolver problemas reais e causar um impacto positivo. Continue codificando, continue aprendendo e continue expandindo os limites do que é possível!

---

*Happy coding! / Bom código!* 🎉
