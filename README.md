# template-mobile-app

Templates for mobile app development.

## [mobile-app](./mobile-app)

This is a template for mobile app development using:

- [Expo](https://expo.dev/)
- [React Native](https://reactnative.dev/)
- [Node.js](https://nodejs.org/) v22.14.0

## [api](./api)

This is a template for building APIs using:

- [Go](https://go.dev/) with [Gin](https://gin-gonic.com/) framework
- [gin-swagger](https://github.com/swaggo/gin-swagger) for OpenAPI/Swagger generation

## [openapi-specifications](./openapi-specifications)

This directory contains OpenAPI specifications swagger files.
OpenAPI version 3.0 is used for the specifications, and the files are in JSON format.

To generate the OpenAPI definitions, run the following command from the repository root:

```bash
./generate-openapi.sh
```

This script will:

1. Generate Swagger 2.0 documentation using `swag init` in the API directory
2. Convert the Swagger 2.0 specification to OpenAPI 3.0.3 format
3. Place the result in `openapi-specifications/api.swagger.json`
4. Verify that React Native type generation and mock server commands work correctly

- [OpenAPI](https://www.openapis.org/)
- [Swagger](https://swagger.io/)

## [ls-lint](./.ls-lint.yml)

This is a configuration file for [ls-lint](https://ls-lint.org/), a linter for directory structures.  
GitHub Actions are set up to run ls-lint on pull requests to ensure that the directory structure adheres to the defined rules.  
For [api](./api) and [mobile-app](./mobile-app) directories, add settings for each directory, as they are ignored except for markdown files.
