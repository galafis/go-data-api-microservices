# Config Package | Pacote de Configuração

⚙️ Go Data API Microservices - Configuration Package ⚙️

## 🇺🇸 English

### 📖 Overview

The Config Package is a comprehensive configuration management library designed specifically for Go microservices. It provides centralized configuration loading from multiple sources including environment variables, configuration files, and default values, making it perfect for managing complex microservice configurations across different environments.

### ✨ Key Features

• **Multi-Source Loading**: Load configuration from env vars, files, and defaults
• **Environment Management**: Seamless configuration across development, staging, and production
• **Type Safety**: Strong typing with automatic validation
• **Hot Reload**: Runtime configuration updates without restart
• **Validation**: Built-in validation for required fields and formats
• **Microservice Ready**: Optimized for distributed architecture patterns

### 🏗️ File Structure

```
internal/config/
├── config.go    # Main configuration management and loading logic
└── README.md    # This documentation file
```

### 🚀 Quick Start

#### Basic Usage

```go
import "github.com/galafis/go-data-api-microservices/internal/config"

// Load configuration
cfg, err := config.Load()
if err != nil {
    log.Fatal("Failed to load config:", err)
}

// Access configuration values
dbHost := cfg.Database.Host
serverPort := cfg.Server.Port
logLevel := cfg.Logging.Level

// Environment-specific configurations
if cfg.Environment == "production" {
    // Production-specific setup
}
```

#### Configuration Structure

```go
type Config struct {
    Environment string         `env:"ENVIRONMENT" default:"development"`
    Server      ServerConfig   `json:"server"`
    Database    DatabaseConfig `json:"database"`
    Auth        AuthConfig     `json:"auth"`
    Logging     LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
    Port         int    `env:"SERVER_PORT" default:"8080"`
    Host         string `env:"SERVER_HOST" default:"localhost"`
    ReadTimeout  int    `env:"READ_TIMEOUT" default:"30"`
    WriteTimeout int    `env:"WRITE_TIMEOUT" default:"30"`
}
```

#### Environment Variables

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=microservices_db
DB_USER=postgres
DB_PASSWORD=secret

# Authentication
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### 🎯 Best Practices

1. **Environment Separation**: Use different config files for each environment
2. **Secret Management**: Never commit secrets to version control
3. **Validation**: Always validate critical configuration values
4. **Defaults**: Provide sensible defaults for all configuration values
5. **Documentation**: Document all configuration options and their effects

## 🇧🇷 Português

### 📖 Visão Geral

O Pacote Config é uma biblioteca abrangente de gerenciamento de configuração projetada especificamente para microserviços Go. Fornece carregamento centralizado de configuração de múltiplas fontes incluindo variáveis de ambiente, arquivos de configuração e valores padrão, sendo perfeita para gerenciar configurações complexas de microserviços em diferentes ambientes.

### ✨ Funcionalidades Principais

• **Carregamento Multi-Fonte**: Carrega configuração de env vars, arquivos e padrões
• **Gerenciamento de Ambiente**: Configuração perfeita entre desenvolvimento, staging e produção
• **Segurança de Tipos**: Tipagem forte com validação automática
• **Hot Reload**: Atualizações de configuração em tempo de execução sem reinicialização
• **Validação**: Validação integrada para campos obrigatórios e formatos
• **Pronto para Microserviços**: Otimizado para padrões de arquitetura distribuída

### 🏗️ Estrutura de Arquivos

```
internal/config/
├── config.go    # Lógica principal de gerenciamento e carregamento de configuração
└── README.md    # Este arquivo de documentação
```

### 🚀 Início Rápido

#### Uso Básico

```go
import "github.com/galafis/go-data-api-microservices/internal/config"

// Carregar configuração
cfg, err := config.Load()
if err != nil {
    log.Fatal("Falha ao carregar config:", err)
}

// Acessar valores de configuração
dbHost := cfg.Database.Host
serverPort := cfg.Server.Port
logLevel := cfg.Logging.Level

// Configurações específicas do ambiente
if cfg.Environment == "production" {
    // Configuração específica de produção
}
```

#### Estrutura de Configuração

```go
type Config struct {
    Environment string         `env:"ENVIRONMENT" default:"development"`
    Server      ServerConfig   `json:"server"`
    Database    DatabaseConfig `json:"database"`
    Auth        AuthConfig     `json:"auth"`
    Logging     LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
    Port         int    `env:"SERVER_PORT" default:"8080"`
    Host         string `env:"SERVER_HOST" default:"localhost"`
    ReadTimeout  int    `env:"READ_TIMEOUT" default:"30"`
    WriteTimeout int    `env:"WRITE_TIMEOUT" default:"30"`
}
```

#### Variáveis de Ambiente

```bash
# Configuração do Servidor
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Configuração do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_NAME=microservices_db
DB_USER=postgres
DB_PASSWORD=secret

# Autenticação
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### 🎯 Melhores Práticas

1. **Separação de Ambiente**: Use arquivos de configuração diferentes para cada ambiente
2. **Gerenciamento de Segredos**: Nunca faça commit de segredos no controle de versão
3. **Validação**: Sempre valide valores críticos de configuração
4. **Padrões**: Forneça padrões sensatos para todos os valores de configuração
5. **Documentação**: Documente todas as opções de configuração e seus efeitos

## 🔧 Advanced Usage | Uso Avançado

### Configuration Validation | Validação de Configuração

```go
// Custom validation functions | Funções de validação personalizadas
func (c *Config) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return errors.New("invalid server port")
    }
    
    if c.Database.Host == "" {
        return errors.New("database host is required")
    }
    
    return nil
}
```

### Environment-Specific Loading | Carregamento Específico do Ambiente

```go
// Load config based on environment | Carregar config baseado no ambiente
func LoadForEnvironment(env string) (*Config, error) {
    configFile := fmt.Sprintf("config.%s.yaml", env)
    return config.LoadFromFile(configFile)
}
```

### Hot Reload Implementation | Implementação de Hot Reload

```go
// Watch for configuration changes | Observar mudanças de configuração
func (c *Config) WatchChanges(callback func(*Config)) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    
    // Implementation for file watching | Implementação para observar arquivos
    // ...
}
```

## 👨‍💻 Author | Autor

**Gabriel Demetrios Lafis**
🚀 Passionate Go Developer & Microservices Architect  
🌟 Desenvolvedor Go Apaixonado & Arquiteto de Microserviços

## 💡 Welcome Message | Mensagem de Boas-Vindas

### 🇺🇸 To Fellow Developers

Welcome to the Config Package! 🎉

Configuration is the backbone of flexible, maintainable software. This package was designed with the 12-factor app principles in mind, ensuring your microservices can adapt to any environment seamlessly.

Why you'll love this package:
• ⚙️ **Environment Agnostic**: Deploy the same code across all environments
• 🔒 **Secure by Design**: Built-in support for secret management
• 🚀 **Production Ready**: Battle-tested configuration patterns
• 🔧 **Developer Friendly**: Clear configuration structure and validation

Remember: Good configuration management is the foundation of reliable deployments. Let's build configurable systems together! 💪

### 🇧🇷 Para Desenvolvedores Companheiros

Bem-vindos ao Pacote Config! 🎉

Configuração é a espinha dorsal de software flexível e mantível. Este pacote foi projetado com os princípios da aplicação de 12 fatores em mente, garantindo que seus microserviços possam se adaptar a qualquer ambiente perfeitamente.

Por que você vai amar este pacote:
• ⚙️ **Agnóstico ao Ambiente**: Implante o mesmo código em todos os ambientes
• 🔒 **Seguro por Design**: Suporte integrado para gerenciamento de segredos
• 🚀 **Pronto para Produção**: Padrões de configuração testados em batalha
• 🔧 **Amigável ao Desenvolvedor**: Estrutura clara de configuração e validação

Lembre-se: Bom gerenciamento de configuração é a base de implantações confiáveis. Vamos construir sistemas configuráveis juntos! 💪

Happy Coding! | Codificação Feliz! 🚀✨
