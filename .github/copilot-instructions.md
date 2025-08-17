# Mobile App Template Development Guide

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Environment Setup
- **Node.js**: MUST use version 22.14.0 (specified in `mobile-app/.nvmrc`)
  - Install: `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash`
  - Setup: `export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"`
  - Install version: `nvm install 22.14.0 && nvm use 22.14.0`
- **Go**: Version 1.24.5 (already available in environment)
- **swag tool**: Required for OpenAPI generation
  - Install: `cd api && go install github.com/swaggo/swag/cmd/swag@latest`
  - Ensure PATH: `export PATH=$PATH:$(go env GOPATH)/bin`

### Bootstrap and Build Process

#### Complete OpenAPI Generation Workflow
- **NEVER CANCEL**: Full workflow takes ~3 minutes. Set timeout to 5+ minutes.
- Command: `./generate-openapi.sh` (from repository root)
- This script orchestrates:
  1. Go dependencies installation (~20 seconds)
  2. Swagger 2.0 generation via `swag init` (~1 second)
  3. OpenAPI 3.0.3 conversion (~1 second)
  4. npm install for mobile app (~2 minutes)
  5. TypeScript schema generation (~1 second each)
  6. Mock server validation (~1 second)

#### Go API Development
- **Dependencies**: `cd api && go mod tidy` (~20 seconds)
- **Generate docs**: `cd api && swag init` (~1 second) 
- **Build**: `cd api && go build -o api-server main.go` (~8 seconds)
- **Run**: `cd api && ./api-server` (starts on port 8080)
- **Test**: Available endpoints:
  - `curl http://localhost:8080/api/v1/health` (always works)
  - `curl http://localhost:8080/swagger/index.html` (Swagger UI)
  - Posts endpoints may fail in sandboxed environments (expected)

#### Mobile App Development
- **ALWAYS use Node.js 22.14.0**: `nvm use 22.14.0` before any npm commands
- **Dependencies**: `cd mobile-app && npm ci` (~15 seconds)
- **Development server**: `cd mobile-app && npx expo start` (~5 seconds to start)
  - QR code for mobile testing
  - Web development available but has known limitation with react-native-maps
- **Schema generation**: 
  - Example: `npm run gen-schema:example` (~1 second)
  - API: `npm run gen-schema` (~1 second)

### Validation and Testing

#### Linting and Formatting - NEVER CANCEL: Run all before committing
- **Mobile app directory linting**: `cd mobile-app && npm run ls-lint` (~1 second)
- **ESLint**: `cd mobile-app && npm run lint` (~4 seconds)  
- **Prettier check**: `cd mobile-app && npm run format:check` (~1 second)
- **TypeScript checking**: `cd mobile-app && npm run test:ts` (~4 seconds)
- **Root directory linting**: `npx @ls-lint/ls-lint` (~1 second)

#### Mock Server Testing
- **Start mock**: `cd mobile-app && npm run mock` (port 8080)
- **Test endpoints**:
  - `curl http://localhost:8080/health` (returns {})
  - `curl http://localhost:8080/posts` (returns sample data array)
  - `curl http://localhost:8080/posts/1` (returns single sample object)

### Known Limitations and Workarounds

#### Web Build Limitation
- **Issue**: `expo export --platform web` fails due to react-native-maps importing native modules
- **Workaround**: This is expected behavior. App has platform-specific files:
  - `maps.tsx` - Web fallback (displays text message)
  - `maps.native.tsx` - Native implementation (full map functionality)
- **Development**: Use `npx expo start` for mobile development (works perfectly)

#### External API Dependencies
- Posts endpoints (`/api/v1/posts`) depend on jsonplaceholder.typicode.com
- In sandboxed environments, these will return error responses (expected behavior)
- Health endpoint always works for testing API functionality

### Timeout Values and Timing Expectations

#### NEVER CANCEL - Critical Build Operations
- **Complete OpenAPI workflow**: 5+ minutes timeout (includes npm install)
- **npm ci**: 3+ minutes timeout (includes dependency resolution)
- **go mod tidy**: 2+ minutes timeout (first time, includes downloads)

#### Quick Operations (under 30 seconds)
- Individual linting/formatting checks: 30 seconds timeout
- Schema generation: 30 seconds timeout
- Go builds: 30 seconds timeout
- swag init: 30 seconds timeout

### Manual Validation Scenarios

#### Complete API Workflow
1. Generate OpenAPI specs: `./generate-openapi.sh`
2. Start API server: `cd api && ./api-server`
3. Test health endpoint: `curl http://localhost:8080/api/v1/health`
4. Verify Swagger UI loads: `curl http://localhost:8080/swagger/index.html`
5. Stop server and start mock: `cd mobile-app && npm run mock`
6. Test mock endpoints: `curl http://localhost:8080/posts`

#### Complete Mobile App Workflow  
1. Ensure Node.js 22.14.0: `nvm use 22.14.0`
2. Install dependencies: `cd mobile-app && npm ci`
3. Generate schemas: `npm run gen-schema:example && npm run gen-schema`
4. Run all validation: `npm run ls-lint && npm run lint && npm run format:check && npm run test:ts`
5. Start development server: `npx expo start`
6. Verify QR code and mobile development options appear

### CI/CD Integration
- GitHub Actions configured with 60-minute timeouts (appropriate)
- Mobile app workflow: schema generation, linting, formatting, type checking
- Project workflow: OpenAPI generation and validation
- All timeouts already properly configured in workflows

### Important File Locations

#### Configuration Files
- `mobile-app/.nvmrc` - Required Node.js version (22.14.0)
- `mobile-app/package.json` - Mobile app dependencies and scripts
- `api/go.mod` - Go dependencies
- `.ls-lint.yml` - Directory structure linting rules
- `generate-openapi.sh` - Complete OpenAPI workflow script

#### Generated Files (Do Not Edit)
- `api/docs/swagger.json` - Generated Swagger 2.0 specification
- `openapi-specifications/api.swagger.json` - Generated OpenAPI 3.0.3 specification  
- `mobile-app/schema/api.d.ts` - Generated TypeScript definitions
- `mobile-app/schema/example.d.ts` - Generated example TypeScript definitions

#### Key Source Directories
- `api/` - Go API server with Gin framework
- `mobile-app/app/` - React Native app with Expo Router
- `mobile-app/components/` - Reusable UI components with Gluestack-UI
- `mobile-app/api-client/` - API client implementations

### Troubleshooting Common Issues

#### "Module not found" errors in mobile app
- Ensure Node.js 22.14.0: `nvm use 22.14.0`
- Regenerate schemas: `npm run gen-schema:example && npm run gen-schema`
- Clear cache: `npx expo start --clear`

#### Go build failures
- Generate swagger docs first: `cd api && swag init`
- Install swag tool: `go install github.com/swaggo/swag/cmd/swag@latest`
- Ensure PATH includes Go bin: `export PATH=$PATH:$(go env GOPATH)/bin`

#### Web bundle failures with react-native-maps
- This is expected behavior due to native module imports
- Use mobile development workflow instead: `npx expo start` (without web)
- App automatically uses platform-specific map implementations