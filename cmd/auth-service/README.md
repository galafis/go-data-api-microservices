# Auth Service

The Auth Service handles user authentication and authorization for the Go Data API Microservices system.

## Overview

This service is responsible for all authentication and authorization operations including:

- User registration and authentication
- JWT token generation and validation
- Role-based access control
- Password management and security

## Features

- **User Management**: Register and authenticate users securely
- **JWT Tokens**: Generate and validate JSON Web Tokens
- **Role-Based Access**: Implement fine-grained access control
- **Password Security**: Secure password hashing and validation
- **Session Management**: Handle user sessions and token refresh
- **Security Middleware**: Protect API endpoints with authentication

## Authentication Flow

1. User registers with email/password
2. User logs in with credentials
3. Service generates JWT access and refresh tokens
4. Client uses access token for authenticated requests
5. Refresh token used to obtain new access tokens

## Building and Running

```bash
# Build the service
go build -o auth-service ./cmd/auth-service

# Run the service
./auth-service
```

## Configuration

The service uses environment variables for configuration. Key variables include:

- `JWT_SECRET`: Secret key for token signing
- `ACCESS_TOKEN_EXPIRY`: Access token expiration time
- `REFRESH_TOKEN_EXPIRY`: Refresh token expiration time
- `PASSWORD_HASH_COST`: BCrypt hash cost for passwords

## API Endpoints

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh
- `POST /api/v1/auth/logout` - User logout

Refer to the main project README for detailed API documentation.
