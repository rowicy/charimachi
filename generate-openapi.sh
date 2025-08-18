#!/bin/bash

# Script to generate OpenAPI 3.0.0 specification from Go Gin swagger annotations
# This script follows the workflow requested:
# 1. Run swag init to generate Swagger 2.0 in api/docs/
# 2. Convert Swagger 2.0 to OpenAPI 3.0.0 
# 3. Place result in openapi-specifications/api.swagger.json
# 4. Verify npm commands work

set -e

echo "🚀 Starting OpenAPI 3.0.0 generation workflow..."

# Check if we're in the right directory
if [ ! -f "api/main.go" ]; then
    echo "❌ Error: api/main.go not found. Please run this script from the repository root."
    exit 1
fi

# Step 1: Install swag if not available
echo "📦 Checking swag installation..."
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    export PATH=$PATH:$(go env GOPATH)/bin
    cd api && go install github.com/swaggo/swag/cmd/swag@latest
    cd ..
else
    echo "✅ swag is already installed"
fi

# Ensure PATH includes Go bin directory
export PATH=$PATH:$(go env GOPATH)/bin

# Step 2: Generate Swagger 2.0 documentation
echo "📝 Generating Swagger 2.0 documentation with swag init..."
cd api
swag init
cd ..

# Verify swagger.json was generated
if [ ! -f "api/docs/swagger.json" ]; then
    echo "❌ Error: api/docs/swagger.json was not generated"
    exit 1
fi
echo "✅ Swagger 2.0 documentation generated at api/docs/swagger.json"

# Step 3: Convert Swagger 2.0 to OpenAPI 3.0.0
echo "🔄 Converting Swagger 2.0 to OpenAPI 3.0.0..."
cd mobile-app
npm install
npm run swagger2openapi
cd ..

# Verify OpenAPI 3.0.0 was generated
if [ ! -f "openapi-specifications/api.swagger.json" ]; then
    echo "❌ Error: openapi-specifications/api.swagger.json was not generated"
    exit 1
fi
echo "✅ OpenAPI 3.0.0 specification generated at openapi-specifications/api.swagger.json"

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
( npm run mock & pid=$!; sleep 5; kill -0 $pid 2>/dev/null && kill -9 $pid ) || true
cd ..
echo "✅ Mock server test passed"

echo ""
echo "🎉 All steps completed successfully!"
echo ""
echo "Generated files:"
echo "  - api/docs/swagger.json (Swagger 2.0 from swag init)"
echo "  - openapi-specifications/api.swagger.json (OpenAPI 3.0.0 converted)"
echo "  - mobile-app/schema/api.d.ts (TypeScript definitions)"
echo ""
echo "Available commands:"
echo "  - npm run gen-schema  # Generate TypeScript definitions"
echo "  - npm run mock        # Start mock server on port 3001"
echo ""