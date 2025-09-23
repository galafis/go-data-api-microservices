# Analytics Service

The Analytics Service provides advanced statistical analysis, correlation studies, time series analysis, and predictive modeling for the Go Data API Microservices system.

## Overview

This service is responsible for all analytical operations including:

• **Statistical Analysis**: Descriptive statistics, distributions, and hypothesis testing
• **Correlation Analysis**: Pearson, Spearman, and partial correlation analysis
• **Time Series Analysis**: Trend analysis, seasonality detection, and forecasting
• **Predictive Modeling**: Machine learning models for forecasting and classification
• **Data Mining**: Pattern recognition and anomaly detection
• **Visualization**: Statistical charts and analytical dashboards

## Features

• **Statistical Computing**: Calculate descriptive statistics, confidence intervals, and hypothesis tests
• **Correlation Matrix**: Generate correlation matrices with various correlation coefficients
• **Time Series Forecasting**: ARIMA, exponential smoothing, and seasonal decomposition
• **Predictive Analytics**: Linear regression, logistic regression, and ensemble methods
• **Anomaly Detection**: Statistical and machine learning-based outlier detection
• **Data Transformation**: Normalization, standardization, and feature engineering
• **Interactive Dashboards**: Real-time analytical visualizations and reports

## Analytics Workflow

1. **Data Ingestion**: Receive data from data-service or external sources
2. **Data Preprocessing**: Clean, transform, and prepare data for analysis
3. **Statistical Analysis**: Perform descriptive and inferential statistics
4. **Correlation Studies**: Analyze relationships between variables
5. **Time Series Processing**: Extract trends, patterns, and seasonality
6. **Model Training**: Train predictive models on historical data
7. **Forecasting**: Generate predictions and confidence intervals
8. **Visualization**: Create charts, reports, and interactive dashboards
9. **Results Export**: Deliver insights via API endpoints or notifications

## Building and Running

```bash
# Build the service
go build -o analytics-service ./cmd/analytics-service

# Run the service
./analytics-service
```

## Configuration

The service uses environment variables for configuration. Key variables include:

• **ANALYTICS_PORT**: Service port (default: 8083)
• **DATABASE_URL**: Connection string for analytical database
• **CACHE_TTL**: Cache expiration time for computed results
• **MODEL_UPDATE_INTERVAL**: Frequency for retraining predictive models
• **MAX_PARALLEL_JOBS**: Maximum concurrent analytical jobs
• **PREDICTION_HORIZON**: Default forecasting time horizon
• **CONFIDENCE_LEVEL**: Default confidence level for statistical tests (0.95)

## API Endpoints

### Statistical Analysis
• **POST /api/v1/analytics/statistics** - Calculate descriptive statistics
• **POST /api/v1/analytics/hypothesis-test** - Perform statistical hypothesis tests
• **POST /api/v1/analytics/distribution** - Analyze data distributions

### Correlation Analysis
• **POST /api/v1/analytics/correlation** - Calculate correlation matrices
• **POST /api/v1/analytics/partial-correlation** - Partial correlation analysis
• **POST /api/v1/analytics/covariance** - Covariance matrix calculation

### Time Series
• **POST /api/v1/analytics/timeseries/forecast** - Generate forecasts
• **POST /api/v1/analytics/timeseries/decompose** - Seasonal decomposition
• **POST /api/v1/analytics/timeseries/trend** - Trend analysis

### Predictive Modeling
• **POST /api/v1/analytics/models/train** - Train predictive models
• **POST /api/v1/analytics/models/predict** - Make predictions
• **GET /api/v1/analytics/models/{id}/metrics** - Model performance metrics

### Anomaly Detection
• **POST /api/v1/analytics/anomalies/detect** - Detect outliers and anomalies
• **POST /api/v1/analytics/anomalies/threshold** - Set anomaly detection thresholds

## Supported Algorithms

### Statistical Methods
• T-tests, Chi-square tests, ANOVA
• Bootstrap sampling and confidence intervals
• Non-parametric tests (Mann-Whitney, Kruskal-Wallis)

### Time Series Methods
• ARIMA (AutoRegressive Integrated Moving Average)
• Exponential Smoothing (Simple, Double, Triple)
• Seasonal Decomposition (STL, X-13ARIMA-SEATS)
• Fourier Transform for frequency analysis

### Machine Learning Models
• Linear and Polynomial Regression
• Logistic Regression for classification
• Decision Trees and Random Forest
• Support Vector Machines (SVM)
• K-means clustering
• Isolation Forest for anomaly detection

## Performance Considerations

• **Caching**: Frequently requested analyses are cached for improved response times
• **Parallel Processing**: CPU-intensive computations utilize multiple cores
• **Streaming Analytics**: Support for real-time data processing
• **Model Persistence**: Trained models are saved and reused
• **Memory Management**: Efficient handling of large datasets

## Dependencies

• **GoStats**: Statistical computing library
• **GoNum**: Numerical computing and linear algebra
• **TimeSeries-Go**: Time series analysis toolkit
• **GoML**: Machine learning algorithms
• **Plotly-Go**: Data visualization and charting
• **Redis**: Caching and session storage
• **PostgreSQL**: Analytical data warehouse

Refer to the main project README for detailed API documentation and usage examples.
