# Hospital Middleware

A middleware service built with Go and Gin framework that serves as an intermediary between hospital systems and client applications.

## Overview

Hospital Middleware is a RESTful API service that provides:
- Staff authentication and management
- Patient information retrieval
- Secure communication with core hospital systems


## API Endpoints

The service provides the following API endpoints:

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|--------------|
| POST | `/api/staff/create` | Create a new staff account | No |
| POST | `/api/staff/login` | Authenticate and receive JWT token | No |
| GET | `/api/patient/search?national_id=12345` | Search for patients by national ID | Yes |

## Requirements

- Go 1.18 or higher
- PostgreSQL 13 or higher
- Docker and Docker Compose (for containerized deployment)

## Setup and Installation

### Local Development

1. Clone the repository:
```bash
git clone git@github.com:roasted99/hospital-middleware.git
cd hospital-middleware
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run cmd/server/main.go
```

### Docker Deployment

1. Clone the repository:
```bash
git clone git@github.com:roasted99/hospital-middleware.git
cd hospital-middleware
```

2. Configure environment variables:
```bash
cp .env
# Edit .env with your configuration
```

3. Create required directories:
```bash
mkdir -p nginx/conf.d nginx/logs
```

4. Start with Docker Compose:
```bash
docker compose up -d
```

The API will be available at http://localhost/api/

## Environment Variables

Key environment variables:

```
# API Server Configuration
PORT=8080

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=hospital_user
DB_PASSWORD=hospital_password
DB_NAME=hospital_db

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

HOSPITAL_A_BASE_URL=https://hospital-a.api.co.th
```

## Authentication

The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <token>
```

To obtain a token, use the `/api/staff/login` endpoint.

## Database Migrations

Migrations are located in the `internal/db/migrations` directory and are run automatically when the application starts.

## Testing

Run tests with:

```bash
go test ./...
```

## Building for Production

1. Build the binary:
```bash
go build -o hospital-api cmd/server/main.go
```

