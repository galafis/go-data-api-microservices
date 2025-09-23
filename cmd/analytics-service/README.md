# Analytics Service

The Analytics Service is a critical component of the [Go Data API Microservices](../../README.md) system, providing advanced statistical analysis, correlation studies, time series analysis, and predictive modeling capabilities.

## 📋 Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Architecture](#architecture)
- [Analytics Workflow](#analytics-workflow)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Supported Algorithms](#supported-algorithms)
- [Performance Considerations](#performance-considerations)
- [Directory Structure](#directory-structure)
- [Dependencies](#dependencies)
- [Related Services](#related-services)

## Overview

The Analytics Service serves as the analytical engine of the microservices ecosystem, processing data from the [Data Service](../data-service/README.md) and providing insights through sophisticated statistical and machine learning methods.

### Core Capabilities

- **Statistical Analysis**: Descriptive statistics, distributions, and hypothesis testing
- **Correlation Analysis**: Pearson, Spearman, and partial correlation analysis
- **Time Series Analysis**: Trend analysis, seasonality detection, and forecasting
- **Predictive Modeling**: Machine learning models for forecasting and classification
- **Data Mining**: Pattern recognition and anomaly detection
- **Visualization**: Statistical charts and analytical dashboards

## Key Features

### 🧮 Statistical Computing
- Descriptive statistics (mean, median, mode, standard deviation)
- Confidence intervals and hypothesis testing
- Distribution analysis and goodness-of-fit tests
- Bootstrap sampling and resampling methods

### 🔗 Correlation & Relationships
- Correlation matrices with multiple correlation coefficients
- Partial correlation and causality analysis
- Covariance matrix calculation
- Multi-dimensional relationship mapping

### ⏱️ Time Series Forecasting
- ARIMA, SARIMA models
- Exponential smoothing (Holt-Winters)
- Seasonal decomposition and trend analysis
- Real-time streaming analytics

### 🤖 Predictive Analytics
- Linear and polynomial regression
- Logistic regression for classification
- Ensemble methods (Random Forest, Gradient Boosting)
- Neural networks for complex patterns

### 🔍 Anomaly Detection
- Statistical outlier detection
- Machine learning-based anomaly detection
- Real-time anomaly monitoring
- Custom threshold configuration

### 📊 Data Transformation
- Normalization and standardization
- Feature engineering and selection
- Data aggregation and sampling
- Missing data imputation

## Architecture

The Analytics Service integrates seamlessly with other microservices:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Service  │────│ Analytics Service│────│ Visualization   │
│                 │    │                 │    │   Dashboard     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │              ┌─────────────────┐              │
         └──────────────│  Auth Service   │──────────────┘
                        │                 │
                        └─────────────────┘
```

## Analytics Workflow

1. **Data Ingestion**: Receive datasets from [Data Service](../data-service/README.md) or external sources
2. **Data Preprocessing**: Clean, validate, and transform data for analysis
3. **Statistical Analysis**: Perform descriptive and inferential statistics
4. **Correlation Studies**: Analyze relationships between variables
5. **Time Series Processing**: Extract trends, patterns, and seasonality
6. **Model Training**: Train predictive models on historical data
7. **Forecasting**: Generate predictions with confidence intervals
8. **Visualization**: Create charts, reports, and interactive dashboards
9. **Results Export**: Deliver insights via API endpoints or notifications

## Quick Start

### Prerequisites
- Go 1.18+
- PostgreSQL 13+
- Redis (for caching)

### Building and Running

```bash
# Build the service
go build -o analytics-service ./cmd/analytics-service

# Run the service
./analytics-service

# Or using Docker
docker-compose up analytics-service

# Or using Kubernetes
kubectl apply -f deployments/kubernetes/
```

For detailed deployment instructions, see:
- [Docker Deployment](deployments/docker/README.md)
- [Kubernetes Deployment](deployments/kubernetes/README.md)
- [Build Scripts](scripts/README.md)

## Configuration

The service uses environment variables for configuration. Copy `.env.example` to `.env` and adjust values:

```bash
# Service Configuration
ANALYTICS_PORT=8083
SERVER_HOST=0.0.0.0
ENVIRONMENT=production

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/analytics_db
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# Cache Configuration
REDIS_URL=redis://localhost:6379
CACHE_TTL=3600
CACHE_PREFIX=analytics:

# Analytics Configuration
MODEL_UPDATE_INTERVAL=24h
MAX_PARALLEL_JOBS=4
PREDICTION_HORIZON=30d
CONFIDENCE_LEVEL=0.95
ANOMALY_THRESHOLD=2.0

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## API Documentation

Detailed API documentation is available at:
- **Swagger UI**: `http://localhost:8083/swagger/index.html` (when running)
- **API Reference**: [api/v1/README.md](api/v1/README.md)
- **Postman Collection**: [docs/analytics-api.postman_collection.json](docs/)

### Key Endpoints

#### Statistical Analysis
- `POST /api/v1/analytics/statistics` - Calculate descriptive statistics
- `POST /api/v1/analytics/hypothesis-test` - Perform hypothesis tests
- `POST /api/v1/analytics/distribution` - Analyze data distributions

#### Correlation Analysis
- `POST /api/v1/analytics/correlation` - Calculate correlation matrices
- `POST /api/v1/analytics/partial-correlation` - Partial correlation analysis
- `POST /api/v1/analytics/covariance` - Covariance matrix calculation

#### Time Series
- `POST /api/v1/analytics/timeseries/forecast` - Generate forecasts
- `POST /api/v1/analytics/timeseries/decompose` - Seasonal decomposition
- `POST /api/v1/analytics/timeseries/trend` - Trend analysis

#### Predictive Modeling
- `POST /api/v1/analytics/models/train` - Train predictive models
- `POST /api/v1/analytics/models/predict` - Make predictions
- `GET /api/v1/analytics/models/{id}/metrics` - Model performance metrics

#### Anomaly Detection
- `POST /api/v1/analytics/anomalies/detect` - Detect outliers and anomalies
- `POST /api/v1/analytics/anomalies/threshold` - Configure detection thresholds

## Supported Algorithms

### Statistical Methods
- **Parametric Tests**: T-tests, ANOVA, F-tests
- **Non-parametric Tests**: Mann-Whitney U, Kruskal-Wallis, Wilcoxon
- **Goodness-of-fit**: Chi-square, Kolmogorov-Smirnov
- **Resampling**: Bootstrap, Jackknife, Permutation tests

### Time Series Methods
- **ARIMA Family**: ARIMA, SARIMA, ARIMAX
- **Exponential Smoothing**: Simple, Double, Triple (Holt-Winters)
- **Decomposition**: STL, X-13ARIMA-SEATS, Classical
- **Frequency Analysis**: FFT, Wavelet Transform

### Machine Learning Models
- **Regression**: Linear, Polynomial, Ridge, Lasso, Elastic Net
- **Classification**: Logistic Regression, SVM, Decision Trees
- **Ensemble Methods**: Random Forest, Gradient Boosting, XGBoost
- **Clustering**: K-means, DBSCAN, Hierarchical
- **Anomaly Detection**: Isolation Forest, One-Class SVM, LOF

## Performance Considerations

### Caching Strategy
- **Result Caching**: Frequently requested analyses cached with Redis
- **Model Persistence**: Trained models stored and reused
- **Query Optimization**: Indexed database queries for fast data retrieval

### Parallel Processing
- **CPU-Intensive Tasks**: Multi-core processing for statistical computations
- **Concurrent Analysis**: Multiple analytical jobs processed simultaneously
- **Streaming Analytics**: Real-time data processing capabilities

### Memory Management
- **Large Dataset Handling**: Efficient processing of big data
- **Memory Pooling**: Optimized memory allocation for analytical operations
- **Garbage Collection**: Tuned GC for analytical workloads

## Directory Structure

```
analytics-service/
├── README.md                 # This file
├── api/v1/                  # API v1 handlers and routes
│   └── README.md           # API documentation
├── deployments/            # Deployment configurations
│   ├── docker/            # Docker deployment files
│   │   └── README.md      # Docker deployment guide
│   └── kubernetes/        # Kubernetes manifests
│       └── README.md      # Kubernetes deployment guide
├── docs/                   # Technical documentation
│   └── README.md          # Documentation overview
└── scripts/               # Utility scripts
    └── README.md         # Scripts documentation
```

## Dependencies

### Core Libraries
- **[GoStats](https://github.com/montanaflynn/stats)**: Statistical computing
- **[GoNum](https://gonum.org/)**: Numerical computing and linear algebra
- **[Prophet](https://facebook.github.io/prophet/)**: Time series forecasting
- **[GoML](https://github.com/cdipaolo/goml)**: Machine learning algorithms

### Infrastructure
- **[Gin](https://github.com/gin-gonic/gin)**: HTTP web framework
- **[GORM](https://gorm.io/)**: ORM for database operations
- **[Redis](https://redis.io/)**: Caching and session storage
- **[PostgreSQL](https://postgresql.org/)**: Analytical data warehouse
- **[Prometheus](https://prometheus.io/)**: Metrics and monitoring

### Visualization
- **[Plotly-Go](https://github.com/go-echarts/go-echarts)**: Interactive charts
- **[Chart.js](https://www.chartjs.org/)**: Statistical visualizations

## Related Services

This service integrates with other components of the Go Data API Microservices system:

- **[Main Project](../../README.md)**: Complete system overview and setup
- **[Data Service](../data-service/README.md)**: Data storage and retrieval
- **[Auth Service](../auth-service/README.md)**: Authentication and authorization  
- **[API Gateway](../api-gateway/README.md)**: Request routing and API management

## Support

For issues, questions, or contributions:
- 📖 **Documentation**: [docs/README.md](docs/README.md)
- 🐛 **Issues**: [GitHub Issues](https://github.com/galafis/go-data-api-microservices/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/galafis/go-data-api-microservices/discussions)

---

**Part of the [Go Data API Microservices](../../README.md) ecosystem** | Licensed under [MIT License](../../LICENSE)
