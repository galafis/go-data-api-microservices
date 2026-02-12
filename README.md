# 📊 Go Data Api Microservices

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8.svg)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

[English](#english) | [Português](#português)

---

## English

### 🎯 Overview

**Go Data Api Microservices** — Data Science project - go-data-api-microservices

Total source lines: **6,747** across **29** files in **2** languages.

### ✨ Key Features

- **Production-Ready Architecture**: Modular, well-documented, and following best practices
- **Comprehensive Implementation**: Complete solution with all core functionality
- **Clean Code**: Type-safe, well-tested, and maintainable codebase
- **Easy Deployment**: Docker support for quick setup and deployment

### 🚀 Quick Start

#### Prerequisites
- Go 1.22+
- Docker and Docker Compose (optional)

#### Installation

1. **Clone the repository**
```bash
git clone https://github.com/galafis/go-data-api-microservices.git
cd go-data-api-microservices
```

2. **Install dependencies**
```bash
go mod download
```

#### Running

```bash
go run ./...
```

## 🐳 Docker

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### 🧪 Testing

```bash
go test ./...
```

### 📁 Project Structure

```
go-data-api-microservices/
├── bin/
├── docs/
│   ├── README.en-us.md
│   └── README.pt-br.md
├── src/
│   ├── cmd/
│   │   ├── analytics-service/
│   │   ├── api-gateway/
│   │   ├── auth-service/
│   │   └── data-service/
│   ├── internal/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── models/
│   ├── pkg/
│   │   ├── logger/
│   │   ├── utils/
│   │   └── validator/
│   └── docker-compose.yml
├── CONTRIBUTING.md
└── README.md
```

### 🛠️ Tech Stack

| Technology | Usage |
|------------|-------|
| Go | 23 files |
| Shell | 6 files |

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### 👤 Author

**Gabriel Demetrios Lafis**

- GitHub: [@galafis](https://github.com/galafis)
- LinkedIn: [Gabriel Demetrios Lafis](https://linkedin.com/in/gabriel-demetrios-lafis)

---

## Português

### 🎯 Visão Geral

**Go Data Api Microservices** — Data Science project - go-data-api-microservices

Total de linhas de código: **6,747** em **29** arquivos em **2** linguagens.

### ✨ Funcionalidades Principais

- **Arquitetura Pronta para Produção**: Modular, bem documentada e seguindo boas práticas
- **Implementação Completa**: Solução completa com todas as funcionalidades principais
- **Código Limpo**: Type-safe, bem testado e manutenível
- **Fácil Implantação**: Suporte Docker para configuração e implantação rápidas

### 🚀 Início Rápido

#### Pré-requisitos
- Go 1.22+
- Docker e Docker Compose (opcional)

#### Instalação

1. **Clone the repository**
```bash
git clone https://github.com/galafis/go-data-api-microservices.git
cd go-data-api-microservices
```

2. **Install dependencies**
```bash
go mod download
```

#### Execução

```bash
go run ./...
```

### 🧪 Testes

```bash
go test ./...
```

### 📁 Estrutura do Projeto

```
go-data-api-microservices/
├── bin/
├── docs/
│   ├── README.en-us.md
│   └── README.pt-br.md
├── src/
│   ├── cmd/
│   │   ├── analytics-service/
│   │   ├── api-gateway/
│   │   ├── auth-service/
│   │   └── data-service/
│   ├── internal/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── models/
│   ├── pkg/
│   │   ├── logger/
│   │   ├── utils/
│   │   └── validator/
│   └── docker-compose.yml
├── CONTRIBUTING.md
└── README.md
```

### 🛠️ Stack Tecnológica

| Tecnologia | Uso |
|------------|-----|
| Go | 23 files |
| Shell | 6 files |

### 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

### 👤 Autor

**Gabriel Demetrios Lafis**

- GitHub: [@galafis](https://github.com/galafis)
- LinkedIn: [Gabriel Demetrios Lafis](https://linkedin.com/in/gabriel-demetrios-lafis)
