# Job Applicants Tracker API

A Go microservice for tracking job applicants with gRPC + REST APIs, PostgreSQL database, and comprehensive validation.

Built with production-ready patterns: type-safe queries (sqlc), protocol buffers for service definitions, structured logging, and proper error handling.

## Quick Start

### Prerequisites

- **Docker & Docker Compose** (recommended for easiest setup)
- **Go 1.25+** (if running locally)

### Start the Application (Docker)

```bash
# Start database, run migrations, seed data, and start API
make docker-up

# The API is now running at:
# - REST API: http://localhost:8080
# - gRPC API: localhost:9090
# - API Documentation: http://localhost:8080/docs/
```

That's it! The service is ready to use.

### Test the API with Swagger UI

1. Open your browser to **http://localhost:8080/docs/**
2. You'll see the interactive Swagger UI with all available endpoints
3. Click on any endpoint to expand it
4. Click "Try it out" to test the endpoint
5. Fill in the parameters and click "Execute"

**Try these first:**
- `GET /v1/applicants` - List all applicants
- `GET /v1/applicants/best` - Get the top-rated applicant
- `GET /v1/applicants/{id}` - Get a specific applicant by ID

### Example API Calls (curl)

#### Get All Applicants
```bash
curl http://localhost:8080/v1/applicants
```

#### Get Specific Applicant
```bash
curl http://localhost:8080/v1/applicants/1
```

#### Get Best Applicant
```bash
curl http://localhost:8080/v1/applicants/best
```

#### Create New Applicant
```bash
curl -X POST http://localhost:8080/v1/applicants \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Developer",
    "email": "jane@example.com",
    "position": "Senior Golang Developer",
    "yearsExperience": 5,
    "skills": ["Go", "Kubernetes", "PostgreSQL"],
    "githubStars": 150,
    "canExitVim": true,
    "knowsGo": true,
    "debugsInProduction": false,
    "interviewScore": 85.0,
    "culturalFitScore": 87.0,
    "technicalScore": 88.0,
    "funFact": "Writes Go tests before implementation",
    "availability": "2 weeks",
    "salaryExpectation": "Competitive"
  }'
```

#### Update Applicant
```bash
curl -X PUT http://localhost:8080/v1/applicants/2 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Senior Developer",
    "email": "jane@example.com",
    "position": "Lead Golang Developer",
    "yearsExperience": 6,
    "skills": ["Go", "Kubernetes", "PostgreSQL", "gRPC"],
    "githubStars": 200,
    "canExitVim": true,
    "knowsGo": true,
    "debugsInProduction": true,
    "interviewScore": 90.0,
    "culturalFitScore": 88.0,
    "technicalScore": 92.0,
    "status": "HIRED",
    "funFact": "Promoted after excellent performance",
    "availability": "Immediate",
    "salaryExpectation": "Negotiated"
  }'
```

#### Delete Applicant
```bash
curl -X DELETE http://localhost:8080/v1/applicants/3
```

#### Health Check
```bash
# Check if service and database are healthy
curl http://localhost:8080/health
```

## Running Locally (Without Docker)

If you prefer to run without Docker:

```bash
# 1. Start PostgreSQL
docker-compose up -d postgres

# 2. Run migrations
make migrate-up

# 3. Seed sample data
make seed

# 4. Run the server
make run
```

The API will be available at:
- REST: http://localhost:8080
- gRPC: localhost:9090
- Docs: http://localhost:8080/docs/

## Stopping the Service

```bash
# Stop all Docker services
make docker-down
```

## Available Make Commands

```bash
make help              # Show all commands
make docker-up         # Start all services with Docker
make docker-down       # Stop all services
make test              # Run tests
make test-coverage     # Run tests with coverage
make build             # Build binaries
make run               # Run server locally
make migrate-up        # Apply database migrations
make migrate-down      # Rollback migrations
make seed              # Seed database
make reset-db          # Reset database (down, up, seed)
```

## Configuration

Configuration via environment variables (see `.env.example`):

```bash
DATABASE_URL=postgres://user:pass@localhost:5432/applicants?sslmode=disable
SERVER_PORT=8080
GRPC_PORT=9090
LOG_LEVEL=debug
CORS_ORIGINS=*
```