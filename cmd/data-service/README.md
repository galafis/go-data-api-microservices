# Data Service

The Data Service is the core service responsible for data storage, retrieval, and manipulation in the Go Data API Microservices system.

## Overview

This service handles all data-related operations including:

- Dataset CRUD operations
- Data querying with filtering, sorting, and pagination
- Data transformation and aggregation
- Data import/export in various formats (CSV, JSON, Parquet)

## Features

- **Dataset Management**: Create, read, update, and delete datasets
- **Query Engine**: Advanced querying capabilities with filters and sorting
- **Data Processing**: Transform and aggregate data on-the-fly
- **Multi-format Support**: Import/export data in CSV, JSON, and Parquet formats
- **Performance Optimized**: Efficient data processing and caching

## Building and Running

```bash
# Build the service
go build -o data-service ./cmd/data-service

# Run the service
./data-service
```

## Configuration

The service uses environment variables for configuration. See the root `.env.example` file for available options.

## API Endpoints

Refer to the main project README for detailed API documentation.
