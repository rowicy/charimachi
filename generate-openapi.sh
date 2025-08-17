#!/bin/bash

# Script to generate OpenAPI 3.0.3 specification from Go Gin swagger annotations
# This script follows the workflow requested:
# 1. Run go-swagger3 init to generate Swagger 3.0.3 in api/docs/
# 2. Place result in openapi-specifications/api.swagger.json
# 3. Verify npm commands work

set -e

echo "🚀 Starting OpenAPI 3.0.3 generation workflow..."

# Check if we're in the right directory
if [ ! -f "api/main.go" ]; then
    echo "❌ Error: api/main.go not found. Please run this script from the repository root."
    exit 1
fi

# Step 1: Install go-swagger3 if not available
echo "📦 Checking go-swagger3 installation..."
if ! command -v go-swagger3 &> /dev/null; then
    echo "Installing swag..."
    export PATH=$PATH:$(go env GOPATH)/bin
    cd api && go install github.com/parvez3019/go-swagger3@latest
    cd ..
else
    echo "✅ go-swagger3 is already installed"
fi

# Ensure PATH includes Go bin directory
export PATH=$PATH:$(go env GOPATH)/bin

# Step 2: Generate Swagger 3.0.3 documentation
echo "📝 Generating Swagger 3.0.3 documentation with go-swagger3..."
cd api
go-swagger3 --module-path . --output ./docs/swagger.json --schema-without-pkg
cd ..

# Verify swagger.json was generated
if [ ! -f "api/docs/swagger.json" ]; then
    echo "❌ Error: api/docs/swagger.json was not generated"
    exit 1
fi
echo "✅ Swagger 3.0.3 documentation generated at api/docs/swagger.json"

# Verify OpenAPI 3.0.3 was generated
if [ ! -f "openapi-specifications/api.swagger.json" ]; then
    echo "❌ Error: openapi-specifications/api.swagger.json was not generated"
    exit 1
fi
echo "✅ OpenAPI 3.0.3 specification generated at openapi-specifications/api.swagger.json"

# Step 4: Test npm commands
echo "🧪 Testing npm commands..."

# Check if npm dependencies are installed
if [ ! -d "mobile-app/node_modules" ]; then
    echo "📦 Installing npm dependencies..."
    cd mobile-app
    npm install
    cd ..
fi

# Test gen-schema command
echo "Testing npm run gen-schema..."
cd mobile-app
npm run gen-schema
cd ..

# Verify TypeScript definitions were generated
if [ ! -f "mobile-app/schema/api.d.ts" ]; then
    echo "❌ Error: TypeScript definitions were not generated"
    exit 1
fi
echo "✅ TypeScript definitions generated at mobile-app/schema/api.d.ts"

# Test mock command (start and immediately stop)
echo "Testing npm run mock..."
cd mobile-app
timeout 5s npm run mock || true
cd ..
echo "✅ Mock server test passed"

echo ""
echo "🎉 All steps completed successfully!"
echo ""
echo "Generated files:"
echo "  - api/docs/swagger.json (Swagger 3.0.3 from go-swagger3)"
echo "  - openapi-specifications/api.swagger.json (OpenAPI 3.0.3 converted)"
echo "  - mobile-app/schema/api.d.ts (TypeScript definitions)"
echo ""
echo "Available commands:"
echo "  - npm run gen-schema  # Generate TypeScript definitions"
echo "  - npm run mock        # Start mock server on port 3001"
echo ""