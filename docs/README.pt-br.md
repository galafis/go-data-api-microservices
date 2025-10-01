# Microsserviços de API de Dados em Go

![Licença](https://img.shields.io/badge/license-MIT-blue.svg)
![Versão Go](https://img.shields.io/badge/go-1.18+-00ADD8.svg)
![Status da Build](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Cobertura](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)

Um sistema de microsserviços de alta performance para APIs de dados, construído com Go, fornecendo capacidades robustas de processamento, análise e visualização de dados.

## 🚀 Funcionalidades

- **API RESTful**: Design de API limpo e consistente seguindo os princípios REST
- **Arquitetura de Microsserviços**: Serviços modulares para escalabilidade e resiliência
- **Processamento de Dados**: Poderosas capacidades de consulta, transformação e agregação
- **Análise**: Análise estatística, correlação, séries temporais e previsão
- **Autenticação e Autorização**: Autenticação segura baseada em JWT
- **Documentação**: Documentação Swagger/OpenAPI gerada automaticamente
- **Monitoramento**: Métricas Prometheus e log estruturado
- **Containerização**: Pronto para Docker e Kubernetes

## 📋 Sumário

- [Arquitetura](#arquitetura)
- [Serviços](#serviços)
- [Instalação](#instalação)
- [Uso](#uso)
- [Referência da API](#referência-da-api)
- [Desenvolvimento](#desenvolvimento)
- [Testes](#testes)
- [Implantação](#implantação)
- [Contribuindo](#contribuindo)
- [Licença](#licença)

## 🏗️ Arquitetura

O sistema segue uma arquitetura de microsserviços com os seguintes componentes:

- **API Gateway**: Ponto de entrada para todas as requisições do cliente, lida com roteamento e autenticação
- **Serviço de Dados**: Serviço principal para armazenamento, recuperação e manipulação de dados
- **Serviço de Autenticação**: Lida com a autenticação e autorização do usuário
- **Serviço de Análise**: Fornece capacidades de análise e visualização de dados

Cada serviço é implantável independentemente e se comunica via gRPC para comunicação interna e REST para clientes externos.

## 🧩 Serviços

### API Gateway

O API Gateway serve como ponto de entrada para todas as requisições do cliente. Ele lida com:

- Roteamento de requisições para os serviços apropriados
- Autenticação e autorização
- Limitação de taxa e controle de fluxo
- Transformação de requisição/resposta
- Documentação da API (Swagger)

### Serviço de Dados

O Serviço de Dados gerencia as operações de dados:

- Operações CRUD de conjuntos de dados
- Consulta de dados com filtragem, ordenação e paginação
- Transformação e agregação de dados
- Importação/exportação de dados em vários formatos (CSV, JSON, Parquet)

### Serviço de Autenticação

O Serviço de Autenticação lida com o gerenciamento de usuários e segurança:

- Registro e autenticação de usuários
- Geração e validação de tokens JWT
- Controle de acesso baseado em função
- Gerenciamento de senhas e segurança

### Serviço de Análise

O Serviço de Análise fornece capacidades de análise de dados:

- Análise estatística (média, mediana, desvio padrão, etc.)
- Análise de correlação
- Análise de séries temporais
- Previsão e predição

## 📦 Instalação

### Pré-requisitos

- Go 1.18 ou superior
- PostgreSQL 13 ou superior
- Docker (opcional)
- Kubernetes (opcional)

### Do Código Fonte

```bash
# Clone o repositório
git clone https://github.com/galafis/go-data-api-microservices.git
cd go-data-api-microservices

# Instale as dependências
go mod download

# Compile os serviços
make build

# Execute os serviços
make run
```

### Usando Docker

```bash
# Construa as imagens Docker
docker-compose build

# Execute os serviços
docker-compose up -d
```

### Usando Kubernetes

```bash
# Aplique os manifestos do Kubernetes
kubectl apply -f deployments/kubernetes/
```

## 🔧 Uso

### Iniciando os Serviços

```bash
# Inicie todos os serviços
make run

# Inicie um serviço específico
make run-api-gateway
make run-data-service
make run-auth-service
make run-analytics-service
```

### Variáveis de Ambiente

Crie um arquivo `.env` no diretório raiz com as seguintes variáveis:

```
# Servidor
ENVIRONMENT=development
SERVER_PORT=8080

# Banco de Dados
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=data_api
DB_SSL_MODE=disable

# Autenticação
JWT_SECRET=your-secret-key
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
PASSWORD_HASH_COST=10

# Log
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

## 📚 Referência da API

A documentação da API está disponível em `/swagger/index.html` quando o API Gateway estiver em execução.

### Autenticação

```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### Conjuntos de Dados

```
GET /api/v1/data/datasets
POST /api/v1/data/datasets
GET /api/v1/data/datasets/{id}
PUT /api/v1/data/datasets/{id}
DELETE /api/v1/data/datasets/{id}
```

### Operações de Dados

```
POST /api/v1/data/query
POST /api/v1/data/transform
POST /api/v1/data/aggregate
POST /api/v1/data/join
```

### Análise

```
GET /api/v1/analytics/summary
POST /api/v1/analytics/statistics
POST /api/v1/analytics/correlation
POST /api/v1/analytics/timeseries
POST /api/v1/analytics/forecast
```

### Usuários

```
GET /api/v1/users/me
PUT /api/v1/users/me
DELETE /api/v1/users/me
```

## 💻 Desenvolvimento

### Estrutura do Projeto

```
.
├── cmd/                    # Pontos de entrada dos serviços
│   ├── api-gateway/        # Serviço API Gateway
│   ├── data-service/       # Serviço de dados
│   ├── auth-service/       # Serviço de autenticação
│   └── analytics-service/  # Serviço de análise
├── internal/               # Pacotes internos
│   ├── auth/               # Lógica de autenticação
│   ├── config/             # Configuração
│   ├── database/           # Conexões de banco de dados
│   ├── handlers/           # Handlers HTTP
│   ├── middleware/         # Middleware HTTP
│   └── models/             # Modelos de dados
├── pkg/                    # Pacotes públicos
│   ├── logger/             # Utilitários de log
│   ├── validator/          # Utilitários de validação
│   └── utils/              # Utilitários gerais
├── api/                    # Definições de API
│   └── v1/                 # API v1
├── deployments/            # Configurações de implantação
│   ├── docker/             # Configurações Docker
│   └── kubernetes/         # Manifestos Kubernetes
├── docs/                   # Documentação
├── scripts/                # Scripts
├── .env.example            # Exemplo de variáveis de ambiente
├── Dockerfile              # Dockerfile
├── docker-compose.yml      # Configuração Docker Compose
├── go.mod                  # Módulos Go
├── go.sum                  # Checksums dos módulos Go
├── Makefile                # Makefile
└── README.md               # README
```

### Fluxo de Desenvolvimento

1. Faça um fork do repositório
2. Crie uma branch de funcionalidade (`git checkout -b feature/minha-funcionalidade-incrivel`)
3. Faça suas alterações
4. Execute os testes (`make test`)
5. Faça commit de suas alterações (`git commit -m 'Adiciona uma funcionalidade incrível'`)
6. Envie para a branch (`git push origin feature/minha-funcionalidade-incrivel`)
7. Abra um Pull Request

## 🧪 Testes

### Executando Testes

```bash
# Execute todos os testes
make test

# Execute testes com cobertura
make test-coverage

# Execute testes para um pacote específico
go test -v ./internal/auth/...
```

### Benchmarks

```bash
# Execute os benchmarks
make benchmark
```

## 🚢 Implantação

### Docker

```bash
# Construa as imagens Docker
docker-compose build

# Execute os serviços
docker-compose up -d
```

### Kubernetes

```bash
# Aplique os manifestos do Kubernetes
kubectl apply -f deployments/kubernetes/
```

## 👥 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para enviar um Pull Request.

1. Faça um fork do repositório
2. Crie sua branch de funcionalidade (`git checkout -b feature/minha-funcionalidade-incrivel`)
3. Faça commit de suas alterações (`git commit -m 'Adiciona uma funcionalidade incrível'`)
4. Envie para a branch (`git push origin feature/minha-funcionalidade-incrivel`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo LICENSE para detalhes.

---

Criado por Gabriel Demetrios Lafis

