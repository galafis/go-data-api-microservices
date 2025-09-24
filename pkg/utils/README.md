# Utils Package / Pacote de Utilitários

## English

### Purpose
The `utils` package provides a comprehensive collection of utility functions for the Go Data API Microservices project. It contains commonly used helper functions for string manipulation, data conversion, validation, HTTP handling, and various other operations that are shared across different microservices.

### Description
This package serves as a centralized utility library that eliminates code duplication across services and provides consistent, well-tested helper functions. It includes utilities for UUID generation, string transformations, JSON operations, HTTP response handling, pagination, and much more.

### Features
- **UUID Operations**: Generate and parse UUIDs
- **String Transformations**: Snake case, camel case, Pascal case conversions
- **Data Parsing**: Safe parsing with default values for int, float, bool, time
- **Random Generation**: Cryptographically secure random bytes and strings
- **JSON Operations**: Marshal and unmarshal utilities
- **HTTP Utilities**: Request/response helpers, client IP detection
- **Functional Programming**: Map, filter, contains functions for slices
- **Formatting**: Human-readable byte and duration formatting
- **Gin Integration**: Pagination and response utilities

### Expected Structure
```
utils/
├── utils.go     # Main utility functions implementation
└── README.md    # This documentation
```

### Main Function Categories

#### UUID Functions
- `NewUUID()`: Generate new UUID
- `ParseUUID(string)`: Parse string to UUID

#### String Functions
- `ToSnakeCase()`: Convert to snake_case
- `ToCamelCase()`: Convert to camelCase
- `ToPascalCase()`: Convert to PascalCase
- `Truncate()`: Truncate string with ellipsis

#### Parsing Functions
- `ParseInt()`: Safe int parsing with default
- `ParseFloat()`: Safe float parsing with default
- `ParseBool()`: Safe bool parsing with default
- `ParseTime()`: Safe time parsing with default

#### HTTP/Gin Functions
- `PaginationParams()`: Extract pagination from request
- `RespondWithJSON()`: JSON response helper
- `RespondWithError()`: Error response helper
- `GetClientIP()`: Extract client IP from request

### Usage

#### Basic Usage
```go
import "github.com/galafis/go-data-api-microservices/pkg/utils"

// UUID operations
id := utils.NewUUID()
parsedID, err := utils.ParseUUID("550e8400-e29b-41d4-a716-446655440000")

// String conversions
snakeCase := utils.ToSnakeCase("HelloWorld")     // "hello_world"
camelCase := utils.ToCamelCase("hello_world")    // "helloWorld"
pascalCase := utils.ToPascalCase("hello_world")  // "HelloWorld"

// Safe parsing
page := utils.ParseInt("10", 1)         // 10, or 1 if parsing fails
price := utils.ParseFloat("99.99", 0.0) // 99.99, or 0.0 if parsing fails
```

#### Gin Integration
```go
func GetUsers(c *gin.Context) {
    page, pageSize := utils.PaginationParams(c)
    
    users, err := getUsersFromDB(page, pageSize)
    if err != nil {
        utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
        return
    }
    
    utils.RespondWithJSON(c, http.StatusOK, users)
}
```

### Best Practices
- Use safe parsing functions to avoid panics
- Prefer utility functions over reimplementation
- Use appropriate string case functions for API consistency
- Always handle errors from parsing functions

---

## Português

### Propósito
O pacote `utils` fornece uma coleção abrangente de funções utilitárias para o projeto Go Data API Microservices. Contém funções auxiliares comumente usadas para manipulação de strings, conversão de dados, validação, tratamento HTTP e várias outras operações compartilhadas entre diferentes microsserviços.

### Descrição
Este pacote serve como uma biblioteca utilitária centralizada que elimina a duplicação de código entre serviços e fornece funções auxiliares consistentes e bem testadas. Inclui utilitários para geração de UUID, transformações de string, operações JSON, tratamento de resposta HTTP, paginação e muito mais.

### Funcionalidades
- **Operações UUID**: Gerar e analisar UUIDs
- **Transformações de String**: Conversões snake case, camel case, Pascal case
- **Análise de Dados**: Análise segura com valores padrão para int, float, bool, time
- **Geração Aleatória**: Bytes e strings aleatórios criptograficamente seguros
- **Operações JSON**: Utilitários marshal e unmarshal
- **Utilitários HTTP**: Auxiliares de request/response, detecção de IP do cliente
- **Programação Funcional**: Funções map, filter, contains para slices
- **Formatação**: Formatação legível de bytes e duração
- **Integração Gin**: Utilitários de paginação e resposta

### Estrutura Esperada
```
utils/
├── utils.go     # Implementação das funções utilitárias principais
└── README.md    # Esta documentação
```

### Como Usar no Projeto

#### Uso Básico
```go
import "github.com/galafis/go-data-api-microservices/pkg/utils"

// Operações UUID
id := utils.NewUUID()
parsedID, err := utils.ParseUUID("550e8400-e29b-41d4-a716-446655440000")

// Conversões de string
snakeCase := utils.ToSnakeCase("HelloWorld")     // "hello_world"
camelCase := utils.ToCamelCase("hello_world")    // "helloWorld"
pascalCase := utils.ToPascalCase("hello_world")  // "HelloWorld"

// Análise segura
page := utils.ParseInt("10", 1)         // 10, ou 1 se a análise falhar
price := utils.ParseFloat("99.99", 0.0) // 99.99, ou 0.0 se a análise falhar
```

#### Integração com Gin
```go
func GetUsers(c *gin.Context) {
    page, pageSize := utils.PaginationParams(c)
    
    users, err := getUsersFromDB(page, pageSize)
    if err != nil {
        utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao buscar usuários")
        return
    }
    
    utils.RespondWithJSON(c, http.StatusOK, users)
}
```

### Boas Práticas
- Use funções de análise segura para evitar pânicos
- Prefira funções utilitárias em vez de reimplementação
- Use funções apropriadas de case de string para consistência da API
- Sempre trate erros das funções de análise

---

## Author / Autor
**Gabriel Demetrios Lafis**

## Onboarding Message / Mensagem de Integração

⚡ **Don't reinvent the wheel - use these utilities!** / **Não reinvente a roda - use esses utilitários!**

This utilities package is your Swiss Army knife for common operations. These battle-tested functions will save you time and prevent bugs. Remember: consistency is key in microservices architecture!

Este pacote de utilitários é seu canivete suíço para operações comuns. Essas funções testadas em batalha economizarão seu tempo e prevenirão bugs. Lembre-se: consistência é fundamental na arquitetura de microsserviços!
