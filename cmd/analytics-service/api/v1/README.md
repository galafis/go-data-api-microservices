# Analytics Service API v1

RESTful API endpoints for the Analytics Service, providing comprehensive data analysis capabilities through HTTP requests.

## 🔗 Quick Links

- **[Main Service Documentation](../../README.md)** - Complete Analytics Service overview
- **[Deployment Guides](../../deployments/)** - Docker and Kubernetes deployment
- **[Technical Documentation](../../docs/README.md)** - Detailed API specifications
- **[Build Scripts](../../scripts/README.md)** - Development and deployment scripts

## 📋 Table of Contents

- [Overview](#overview)
- [API Architecture](#api-architecture)
- [Authentication](#authentication)
- [Rate Limiting](#rate-limiting)
- [Endpoint Categories](#endpoint-categories)
- [Request/Response Format](#requestresponse-format)
- [Error Handling](#error-handling)
- [Code Examples](#code-examples)
- [Testing](#testing)
- [Development](#development)

## Overview

The Analytics Service API v1 provides RESTful endpoints for advanced statistical analysis, machine learning, and data visualization. All endpoints follow RESTful conventions and return JSON responses.

### Base URL
```
http://localhost:8083/api/v1
```

### Supported Content Types
- `application/json` (primary)
- `application/x-www-form-urlencoded`
- `multipart/form-data` (for file uploads)

## API Architecture

```
api/v1/
├── handlers/           # HTTP request handlers
│   ├── analytics.go   # Statistical analysis endpoints
│   ├── models.go      # Machine learning endpoints
│   ├── timeseries.go  # Time series analysis endpoints
│   ├── correlation.go # Correlation analysis endpoints
│   └── anomaly.go     # Anomaly detection endpoints
├── middleware/        # HTTP middleware
│   ├── auth.go       # Authentication middleware
│   ├── validate.go   # Request validation
│   ├── ratelimit.go  # Rate limiting
│   └── logging.go    # Request/response logging
├── models/           # Data structures
│   ├── requests.go   # API request models
│   ├── responses.go  # API response models
│   └── errors.go     # Error response models
└── routes/           # Route configuration
    ├── router.go     # Main router setup
    └── routes.go     # Endpoint definitions
```

## Authentication

All API endpoints require authentication via JWT tokens:

```http
Authorization: Bearer <jwt_token>
```

### Obtaining a Token
```bash
# Get authentication token from Auth Service
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}'
```

## Rate Limiting

- **Default Limit**: 1000 requests per hour per API key
- **Burst Limit**: 100 requests per minute
- **Headers**: Rate limit information included in response headers

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1633027200
```

## Endpoint Categories

### 📊 Statistical Analysis
- `POST /analytics/statistics` - Calculate descriptive statistics
- `POST /analytics/hypothesis-test` - Perform statistical hypothesis tests
- `POST /analytics/distribution` - Analyze data distributions
- `POST /analytics/summary` - Generate comprehensive data summaries

### 🔗 Correlation Analysis
- `POST /analytics/correlation` - Calculate correlation matrices
- `POST /analytics/partial-correlation` - Partial correlation analysis
- `POST /analytics/covariance` - Covariance matrix calculation
- `POST /analytics/regression` - Linear and non-linear regression

### ⏰ Time Series Analysis
- `POST /analytics/timeseries/forecast` - Generate time series forecasts
- `POST /analytics/timeseries/decompose` - Seasonal decomposition
- `POST /analytics/timeseries/trend` - Trend analysis and detection
- `POST /analytics/timeseries/anomalies` - Time series anomaly detection

### 🤖 Machine Learning Models
- `POST /analytics/models/train` - Train predictive models
- `POST /analytics/models/predict` - Make predictions with trained models
- `GET /analytics/models/{id}` - Get model information
- `GET /analytics/models/{id}/metrics` - Get model performance metrics
- `DELETE /analytics/models/{id}` - Delete a trained model

### 🔍 Anomaly Detection
- `POST /analytics/anomalies/detect` - Detect outliers and anomalies
- `POST /analytics/anomalies/threshold` - Configure detection thresholds
- `GET /analytics/anomalies/rules` - List detection rules
- `PUT /analytics/anomalies/rules/{id}` - Update detection rule

### 📈 Data Visualization
- `POST /analytics/visualize/chart` - Generate charts and graphs
- `POST /analytics/visualize/heatmap` - Create correlation heatmaps
- `POST /analytics/visualize/distribution` - Plot data distributions
- `POST /analytics/visualize/timeseries` - Create time series plots

## Request/Response Format

### Standard Request Structure
```json
{
  "data": {
    "values": [1, 2, 3, 4, 5],
    "columns": ["feature1", "feature2"],
    "metadata": {
      "source": "dataset_name",
      "timestamp": "2025-09-23T20:42:00Z"
    }
  },
  "options": {
    "method": "pearson",
    "confidence_level": 0.95,
    "include_plots": true
  }
}
```

### Standard Response Structure
```json
{
  "success": true,
  "data": {
    "result": {
      "correlation_matrix": [[1.0, 0.85], [0.85, 1.0]],
      "p_values": [[0.0, 0.001], [0.001, 0.0]]
    },
    "metadata": {
      "execution_time": "150ms",
      "method": "pearson",
      "sample_size": 1000
    }
  },
  "timestamp": "2025-09-23T20:42:05Z",
  "request_id": "req_abc123def456"
}
```

### Error Response Structure
```json
{
  "success": false,
  "error": {
    "code": "INVALID_INPUT",
    "message": "Data must contain at least 2 observations",
    "details": {
      "field": "data.values",
      "received_count": 1,
      "minimum_required": 2
    }
  },
  "timestamp": "2025-09-23T20:42:05Z",
  "request_id": "req_abc123def456"
}
```

## Error Handling

### HTTP Status Codes
- `200 OK` - Request successful
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Endpoint or resource not found
- `422 Unprocessable Entity` - Valid request but processing failed
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

### Common Error Codes
- `INVALID_INPUT` - Request data validation failed
- `INSUFFICIENT_DATA` - Not enough data for analysis
- `MODEL_NOT_FOUND` - Specified model doesn't exist
- `ANALYSIS_FAILED` - Statistical analysis couldn't complete
- `TIMEOUT_EXCEEDED` - Request processing timeout

## Code Examples

### Statistical Analysis
```bash
# Calculate correlation matrix
curl -X POST http://localhost:8083/api/v1/analytics/correlation \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "values": [[1,2,3], [4,5,6], [7,8,9]],
      "columns": ["x", "y", "z"]
    },
    "options": {
      "method": "pearson",
      "return_p_values": true
    }
  }'
```

### Time Series Forecasting
```bash
# Generate forecast
curl -X POST http://localhost:8083/api/v1/analytics/timeseries/forecast \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "timestamps": ["2025-01-01", "2025-01-02", "2025-01-03"],
      "values": [100, 105, 98]
    },
    "options": {
      "method": "arima",
      "forecast_periods": 7,
      "confidence_intervals": [0.8, 0.95]
    }
  }'
```

### Model Training
```bash
# Train a regression model
curl -X POST http://localhost:8083/api/v1/analytics/models/train \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {
      "features": [[1,2], [3,4], [5,6]],
      "target": [10, 20, 30]
    },
    "options": {
      "model_type": "linear_regression",
      "validation_split": 0.2,
      "cross_validation_folds": 5
    }
  }'
```

## Testing

### Running Tests
```bash
# Unit tests
go test ./api/v1/handlers/...

# Integration tests
go test -tags=integration ./api/v1/...

# API endpoint tests
go test ./api/v1/routes/...
```

### Test Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./api/v1/...
go tool cover -html=coverage.out
```

## Development

### Adding New Endpoints

1. **Create Handler**: Add handler function in `handlers/`
2. **Define Models**: Add request/response models in `models/`
3. **Add Route**: Register endpoint in `routes/routes.go`
4. **Add Middleware**: Configure authentication and validation
5. **Write Tests**: Add comprehensive test coverage
6. **Update Documentation**: Update this README and API docs

### Code Style Guidelines

- Follow Go naming conventions
- Use descriptive variable names
- Add comprehensive error handling
- Include detailed logging
- Write unit tests for all functions
- Document all public functions

### Environment Setup

```bash
# Install dependencies
go mod download

# Run development server
go run ../../main.go

# Run with hot reload (using air)
air -c .air.toml
```

## 🔍 Related Documentation

- **[Analytics Service](../../README.md)** - Main service documentation
- **[Deployment](../../deployments/)** - Docker and Kubernetes setup
- **[Technical Docs](../../docs/README.md)** - Detailed specifications
- **[Scripts](../../scripts/README.md)** - Development utilities
- **[Main Project](../../../../README.md)** - Complete system overview

---

**Part of the [Analytics Service](../../README.md) | [Go Data API Microservices](../../../../README.md) ecosystem**
