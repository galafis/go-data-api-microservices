# Analytics Service Scripts

Automation and utility scripts for the Analytics Service development, testing, deployment, and maintenance workflows.

## 🔗 Quick Links

- **[Analytics Service](../README.md)** - Main service documentation
- **[API Documentation](../api/v1/README.md)** - REST API reference
- **[Deployment Guides](../deployments/)** - Docker and Kubernetes deployment
- **[Technical Documentation](../docs/README.md)** - Comprehensive technical docs

## 📋 Script Directory

### 🔨 Build & Development
- **`build.sh`** - Build and compile the Analytics Service
- **`dev.sh`** - Start development server with hot reload
- **`clean.sh`** - Clean build artifacts and temporary files
- **`deps.sh`** - Install and update dependencies

### 🧪 Testing & Quality
- **`test.sh`** - Run all tests with coverage reports
- **`test-unit.sh`** - Run unit tests only
- **`test-integration.sh`** - Run integration tests
- **`lint.sh`** - Code linting and formatting
- **`security.sh`** - Security vulnerability scanning

### 🚀 Deployment
- **`deploy.sh`** - Complete deployment automation
- **`docker-build.sh`** - Build Docker images
- **`k8s-deploy.sh`** - Kubernetes deployment
- **`rollback.sh`** - Rollback to previous version

### 💾 Database & Data
- **`migrate.sh`** - Database schema migrations
- **`seed.sh`** - Seed database with test data
- **`backup.sh`** - Database backup utilities
- **`restore.sh`** - Database restore utilities

### 📈 Monitoring & Maintenance
- **`health-check.sh`** - Service health verification
- **`logs.sh`** - Log collection and analysis
- **`metrics.sh`** - Performance metrics collection
- **`cleanup.sh`** - Maintenance cleanup tasks

## 🚀 Quick Start

### Development Workflow

```bash
# 1. Set up development environment
./scripts/deps.sh

# 2. Run tests
./scripts/test.sh

# 3. Start development server
./scripts/dev.sh

# 4. Build for production
./scripts/build.sh
```

### Deployment Workflow

```bash
# 1. Run full test suite
./scripts/test.sh

# 2. Build Docker image
./scripts/docker-build.sh

# 3. Deploy to staging
./scripts/deploy.sh staging

# 4. Deploy to production
./scripts/deploy.sh production
```

## 📄 Detailed Script Usage

### build.sh
Builds the Analytics Service binary with optimization flags.

```bash
# Basic build
./scripts/build.sh

# Build with specific version
./scripts/build.sh --version=v1.2.3

# Build for different platform
./scripts/build.sh --os=linux --arch=amd64
```

**Options:**
- `--version`: Set build version
- `--os`: Target operating system
- `--arch`: Target architecture
- `--debug`: Build with debug symbols

### test.sh
Executes comprehensive test suite with coverage reporting.

```bash
# Run all tests
./scripts/test.sh

# Run with verbose output
./scripts/test.sh --verbose

# Generate HTML coverage report
./scripts/test.sh --html-coverage

# Run specific test package
./scripts/test.sh --package=./internal/analytics
```

**Features:**
- Unit and integration test execution
- Coverage report generation
- Benchmark testing
- Race condition detection

### deploy.sh
Automated deployment to various environments.

```bash
# Deploy to development
./scripts/deploy.sh dev

# Deploy to staging
./scripts/deploy.sh staging

# Deploy to production with confirmation
./scripts/deploy.sh prod --confirm

# Dry run deployment
./scripts/deploy.sh prod --dry-run
```

**Environments:**
- `dev` - Development environment
- `staging` - Staging environment
- `prod` - Production environment

**Options:**
- `--confirm`: Require manual confirmation
- `--dry-run`: Simulate deployment
- `--rollback`: Rollback to previous version

### migrate.sh
Database migration management.

```bash
# Apply all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Show migration status
./scripts/migrate.sh status

# Create new migration
./scripts/migrate.sh create add_user_preferences
```

**Commands:**
- `up` - Apply pending migrations
- `down` - Rollback migrations
- `status` - Show migration status
- `create` - Create new migration file

## 🔧 Configuration

### Environment Variables
Scripts use these environment variables for configuration:

```bash
# Service Configuration
export ANALYTICS_ENV="development"
export ANALYTICS_PORT="8083"
export LOG_LEVEL="info"

# Database Configuration
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_NAME="analytics_dev"
export DB_USER="analytics_user"

# Docker Configuration
export DOCKER_REGISTRY="your-registry.com"
export DOCKER_TAG="latest"

# Kubernetes Configuration
export KUBECONFIG="~/.kube/config"
export K8S_NAMESPACE="analytics"
```

### Script Configuration Files
- **`.env.scripts`** - Script-specific environment variables
- **`config/build.yaml`** - Build configuration
- **`config/deploy.yaml`** - Deployment configuration
- **`config/test.yaml`** - Testing configuration

## 🤖 CI/CD Integration

### GitHub Actions
```yaml
# .github/workflows/analytics-service.yml
name: Analytics Service CI/CD

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Tests
        run: ./cmd/analytics-service/scripts/test.sh
      
  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Service
        run: ./cmd/analytics-service/scripts/build.sh
      - name: Build Docker Image
        run: ./cmd/analytics-service/scripts/docker-build.sh
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    
    stages {
        stage('Test') {
            steps {
                sh './cmd/analytics-service/scripts/test.sh'
            }
        }
        
        stage('Build') {
            steps {
                sh './cmd/analytics-service/scripts/build.sh'
                sh './cmd/analytics-service/scripts/docker-build.sh'
            }
        }
        
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh './cmd/analytics-service/scripts/deploy.sh prod'
            }
        }
    }
}
```

## 🐛 Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Or make all at once
find scripts/ -name "*.sh" -exec chmod +x {} \;
```

#### Missing Dependencies
```bash
# Install required tools
./scripts/deps.sh

# Check system requirements
./scripts/health-check.sh --prereq
```

#### Build Failures
```bash
# Clean build cache
./scripts/clean.sh

# Rebuild dependencies
go mod download
go mod tidy

# Run build with debug info
./scripts/build.sh --debug
```

### Debug Mode
Most scripts support debug mode for troubleshooting:

```bash
# Enable debug output
DEBUG=1 ./scripts/build.sh

# Or use the debug flag
./scripts/build.sh --debug
```

## 📚 Best Practices

### Script Development
1. **Use set -e**: Exit on error
2. **Use set -u**: Exit on undefined variables
3. **Use set -o pipefail**: Exit on pipe failures
4. **Add help documentation**: Include --help flag
5. **Validate prerequisites**: Check required tools and permissions
6. **Use logging**: Provide clear output and error messages

### Security Considerations
1. **Validate inputs**: Sanitize all user inputs
2. **Use secure defaults**: Fail securely when possible
3. **Avoid hardcoded secrets**: Use environment variables
4. **Audit script permissions**: Minimal required permissions
5. **Log security events**: Track sensitive operations

### Maintenance
1. **Version your scripts**: Use semantic versioning
2. **Test regularly**: Include scripts in CI/CD testing
3. **Document changes**: Maintain changelog
4. **Monitor usage**: Track script performance and errors
5. **Update dependencies**: Keep tools and dependencies current

## 🔍 Related Resources

- **[Analytics Service](../README.md)** - Main service documentation
- **[API Reference](../api/v1/README.md)** - REST API endpoints
- **[Deployment Guides](../deployments/)** - Docker and Kubernetes setup
- **[Technical Documentation](../docs/README.md)** - Comprehensive technical reference
- **[Main Project](../../../README.md)** - Complete system overview

---

**Part of the [Analytics Service](../README.md) | [Go Data API Microservices](../../../README.md) ecosystem**
