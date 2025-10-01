# Go Data API Microservices

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.18+-00ADD8.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)
![Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)

A high-performance microservices system for data APIs, built with Go, providing robust data processing, analysis, and visualization capabilities.

## 🚀 Features

- **RESTful API**: Clean and consistent API design following REST principles
- **Microservices Architecture**: Modular services for scalability and resilience
- **Data Processing**: Powerful query, transform, and aggregation capabilities
- **Analytics**: Statistical analysis, correlation, time series, and forecasting
- **Authentication & Authorization**: Secure JWT-based authentication
- **Documentation**: Auto-generated Swagger/OpenAPI documentation
- **Monitoring**: Prometheus metrics and structured logging
- **Containerization**: Docker and Kubernetes ready

## 📋 Table of Contents

- [Architecture](#architecture)
- [Services](#services)
- [Installation](#installation)
- [Usage](#usage)
- [API Reference](#api-reference)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## 🏗️ Architecture

The system follows a microservices architecture with the following components:

- **API Gateway**: Entry point for all client requests, handles routing and authentication
- **Data Service**: Core service for data storage, retrieval, and manipulation
- **Auth Service**: Handles user authentication and authorization
- **Analytics Service**: Provides data analysis and visualization capabilities

Each service is independently deployable and communicates via gRPC for internal communication and REST for external clients.

## 🧩 Services

### API Gateway

The API Gateway serves as the entry point for all client requests. It handles:

- Request routing to appropriate services
- Authentication and authorization
- Rate limiting and throttling
- Request/response transformation
- API documentation (Swagger)

### Data Service

The Data Service manages data operations:

- Dataset CRUD operations
- Data querying with filtering, sorting, and pagination
- Data transformation and aggregation
- Data import/export in various formats (CSV, JSON, Parquet)

### Auth Service

The Auth Service handles user management and security:

- User registration and authentication
- JWT token generation and validation
- Role-based access control
- Password management and security

### Analytics Service

The Analytics Service provides data analysis capabilities:

- Statistical analysis (mean, median, standard deviation, etc.)
- Correlation analysis
- Time series analysis
- Forecasting and prediction

## 📦 Installation

### Prerequisites

- Go 1.18 or higher
- PostgreSQL 13 or higher
- Docker (optional)
- Kubernetes (optional)

### From Source

```bash
# Clone the repository
git clone https://github.com/galafis/go-data-api-microservices.git
cd go-data-api-microservices

# Install dependencies
go mod download

# Build the services
make build

# Run the services
make run
```

### Using Docker

```bash
# Build Docker images
docker-compose build

# Run the services
docker-compose up -d
```

### Using Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/kubernetes/
```

## 🔧 Usage

### Starting the Services

```bash
# Start all services
make run

# Start a specific service
make run-api-gateway
make run-data-service
make run-auth-service
make run-analytics-service
```

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Server
ENVIRONMENT=development
SERVER_PORT=8080

# Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=data_api
DB_SSL_MODE=disable

# Authentication
JWT_SECRET=your-secret-key
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
PASSWORD_HASH_COST=10

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

## 📚 API Reference

The API documentation is available at `/swagger/index.html` when the API Gateway is running.

### Authentication

```
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### Datasets

```
GET /api/v1/data/datasets
POST /api/v1/data/datasets
GET /api/v1/data/datasets/{id}
PUT /api/v1/data/datasets/{id}
DELETE /api/v1/data/datasets/{id}
```

### Data Operations

```
POST /api/v1/data/query
POST /api/v1/data/transform
POST /api/v1/data/aggregate
POST /api/v1/data/join
```

### Analytics

```
GET /api/v1/analytics/summary
POST /api/v1/analytics/statistics
POST /api/v1/analytics/correlation
POST /api/v1/analytics/timeseries
POST /api/v1/analytics/forecast
```

### Users

```
GET /api/v1/users/me
PUT /api/v1/users/me
DELETE /api/v1/users/me
```

## 💻 Development

### Project Structure

```
.
├── cmd/                    # Service entry points
│   ├── api-gateway/        # API Gateway service
│   ├── data-service/       # Data service
│   ├── auth-service/       # Auth service
│   └── analytics-service/  # Analytics service
├── internal/               # Internal packages
│   ├── auth/               # Authentication logic
│   ├── config/             # Configuration
│   ├── database/           # Database connections
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   └── models/             # Data models
├── pkg/                    # Public packages
│   ├── logger/             # Logging utilities
│   ├── validator/          # Validation utilities
│   └── utils/              # General utilities
├── api/                    # API definitions
│   └── v1/                 # API v1
├── deployments/            # Deployment configurations
│   ├── docker/             # Docker configurations
│   └── kubernetes/         # Kubernetes manifests
├── docs/                   # Documentation
├── scripts/                # Scripts
├── .env.example            # Example environment variables
├── Dockerfile              # Dockerfile
├── docker-compose.yml      # Docker Compose configuration
├── go.mod                  # Go modules
├── go.sum                  # Go modules checksums
├── Makefile                # Makefile
└── README.md               # README
```

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for a specific package
go test -v ./internal/auth/...
```

### Benchmarks

```bash
# Run benchmarks
make benchmark
```

## 🚢 Deployment

### Docker

```bash
# Build Docker images
docker-compose build

# Run the services
docker-compose up -d
```

### Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f deployments/kubernetes/
```

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Created by Gabriel Demetrios Lafis

