# API

This is a template API implementation using Go and the Gin framework, with automatic Swagger documentation generation.

## Technology Stack

- **[Go](https://go.dev/)** - Programming language
- **[Gin](https://gin-gonic.com/)** - Web framework
- **[gin-swagger](https://github.com/swaggo/gin-swagger)** - Swagger/OpenAPI documentation generation

## Features

- RESTful API endpoints
- Automatic Swagger/OpenAPI documentation
- **Automated Swagger 2.0 to OpenAPI 3.0.3 conversion** (pure Go implementation)
- CORS support
- Health check endpoint
- Example integration with external API (JSONPlaceholder)
- Proper error handling and response formatting

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Internet connection (for external API calls)

### Installation

1. Navigate to the API directory:
   ```bash
   cd api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install swag tool (for Swagger generation):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

4. Ensure swag is in your PATH (add this to your shell profile if needed):
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

### Running the API

1. Generate OpenAPI 3.0.3 documentation:
   ```bash
   # From repository root directory
   ./generate-openapi.sh
   ```

2. Start the server:
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

### API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/posts` - Get all posts from JSONPlaceholder
- `GET /api/v1/posts/{id}` - Get a specific post by ID
- `GET /swagger/index.html` - Swagger UI documentation

### Swagger Documentation

The API automatically generates OpenAPI/Swagger documentation through the following workflow:

1. **Swagger 2.0 Generation**: `swag init` generates Swagger 2.0 format in `docs/swagger.json`
2. **Conversion to OpenAPI 3.0.3**: A Go converter transforms it to OpenAPI 3.0.3 format
3. **Output**: Final specification is placed in `../openapi-specifications/api.swagger.json`

The documentation is:

- Served at `/swagger/index.html` when the server is running (Swagger 2.0 format)
- Available as OpenAPI 3.0.3 in `../openapi-specifications/api.swagger.json` for tooling

### Development

To regenerate OpenAPI 3.0.3 documentation after making changes:

**Complete workflow (Recommended):**

```bash
# From repository root
./generate-openapi.sh
```

This script will:

1. Run `swag init` to generate Swagger 2.0 documentation in `api/docs/`
2. Convert Swagger 2.0 to OpenAPI 3.0.3 format using a Go converter
3. Place the result in `openapi-specifications/api.swagger.json`
4. Test that `npm run gen-schema` and `npm run mock` work correctly

**Manual steps:**

```bash
# 1. Generate Swagger 2.0
cd api && swag init && cd ..

# 2. Convert to OpenAPI 3.0.3
go run convert-swagger.go

# 3. Test React Native tooling
cd mobile-app
npm run gen-schema  # Generate TypeScript definitions
npm run mock        # Start mock server on port 3001
```
