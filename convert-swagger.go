package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Swagger2 represents a Swagger 2.0 specification
type Swagger2 struct {
	Swagger     string                 `json:"swagger"`
	Info        Info                   `json:"info"`
	Host        string                 `json:"host"`
	BasePath    string                 `json:"basePath"`
	Paths       map[string]interface{} `json:"paths"`
	Definitions map[string]interface{} `json:"definitions"`
}

// OpenAPI3 represents an OpenAPI 3.0.3 specification
type OpenAPI3 struct {
	OpenAPI    string                 `json:"openapi"`
	Info       Info                   `json:"info"`
	Servers    []Server               `json:"servers"`
	Paths      map[string]interface{} `json:"paths"`
	Components Components             `json:"components"`
}

type Info struct {
	Description string      `json:"description"`
	Title       string      `json:"title"`
	Contact     interface{} `json:"contact"`
	Version     string      `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Components struct {
	Schemas map[string]interface{} `json:"schemas"`
}

func convertSwagger2ToOpenAPI3(swagger2 Swagger2) OpenAPI3 {
	// Convert server info
	serverURL := "http://" + swagger2.Host + swagger2.BasePath
	servers := []Server{
		{
			URL:         serverURL,
			Description: "Development server",
		},
	}

	// Convert paths by replacing $ref patterns and updating structure
	convertedPaths := make(map[string]interface{})
	for path, pathItem := range swagger2.Paths {
		convertedPathItem := convertPathItem(pathItem)
		convertedPaths[path] = convertedPathItem
	}

	// Convert definitions to components/schemas
	convertedSchemas := make(map[string]interface{})
	for defName, defSchema := range swagger2.Definitions {
		// Remove the "main." prefix from schema names
		cleanName := strings.TrimPrefix(defName, "main.")
		convertedSchemas[cleanName] = defSchema
	}

	return OpenAPI3{
		OpenAPI: "3.0.3",
		Info:    swagger2.Info,
		Servers: servers,
		Paths:   convertedPaths,
		Components: Components{
			Schemas: convertedSchemas,
		},
	}
}

func convertPathItem(pathItem interface{}) interface{} {
	pathMap, ok := pathItem.(map[string]interface{})
	if !ok {
		return pathItem
	}

	convertedPath := make(map[string]interface{})
	
	for method, operation := range pathMap {
		convertedOperation := convertOperation(operation)
		convertedPath[method] = convertedOperation
	}

	return convertedPath
}

func convertOperation(operation interface{}) interface{} {
	opMap, ok := operation.(map[string]interface{})
	if !ok {
		return operation
	}

	convertedOp := make(map[string]interface{})

	// Copy basic fields
	for key, value := range opMap {
		if key == "consumes" || key == "produces" {
			// Skip consumes/produces as they are handled differently in OpenAPI 3
			continue
		} else if key == "responses" {
			convertedOp[key] = convertResponses(value)
		} else {
			convertedOp[key] = value
		}
	}

	return convertedOp
}

func convertResponses(responses interface{}) interface{} {
	respMap, ok := responses.(map[string]interface{})
	if !ok {
		return responses
	}

	convertedResponses := make(map[string]interface{})

	for statusCode, response := range respMap {
		convertedResponse := convertResponse(response)
		convertedResponses[statusCode] = convertedResponse
	}

	return convertedResponses
}

func convertResponse(response interface{}) interface{} {
	respMap, ok := response.(map[string]interface{})
	if !ok {
		return response
	}

	convertedResp := make(map[string]interface{})

	// Copy description
	if desc, exists := respMap["description"]; exists {
		convertedResp["description"] = desc
	}

	// Convert schema to content
	if schema, exists := respMap["schema"]; exists {
		content := map[string]interface{}{
			"application/json": map[string]interface{}{
				"schema": convertSchemaRef(schema),
			},
		}

		// Add examples based on schema
		if schemaMap, ok := schema.(map[string]interface{}); ok {
			if ref, exists := schemaMap["$ref"]; exists {
				content["application/json"].(map[string]interface{})["examples"] = getExampleForRef(ref.(string))
			} else if schemaType, exists := schemaMap["type"]; exists && schemaType == "array" {
				if items, exists := schemaMap["items"]; exists {
					if itemsMap, ok := items.(map[string]interface{}); ok {
						if ref, exists := itemsMap["$ref"]; exists {
							content["application/json"].(map[string]interface{})["examples"] = getArrayExampleForRef(ref.(string))
						}
					}
				}
			}
		}

		convertedResp["content"] = content
	}

	return convertedResp
}

func convertSchemaRef(schema interface{}) interface{} {
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return schema
	}

	convertedSchema := make(map[string]interface{})

	for key, value := range schemaMap {
		if key == "$ref" && value != nil {
			// Update reference path from #/definitions/ to #/components/schemas/
			ref := value.(string)
			if strings.HasPrefix(ref, "#/definitions/") {
				// Remove "main." prefix from the reference
				refName := strings.TrimPrefix(strings.TrimPrefix(ref, "#/definitions/"), "main.")
				convertedSchema["$ref"] = "#/components/schemas/" + refName
			} else {
				convertedSchema[key] = value
			}
		} else if key == "items" {
			convertedSchema[key] = convertSchemaRef(value)
		} else {
			convertedSchema[key] = value
		}
	}

	return convertedSchema
}

func getExampleForRef(ref string) map[string]interface{} {
	if strings.Contains(ref, "PostResponse") {
		return map[string]interface{}{
			"success": map[string]interface{}{
				"value": map[string]interface{}{
					"id":          1,
					"title":       "Sample Post Title",
					"body":        "Sample post body content",
					"userId":      1,
					"formattedAt": "2024-01-01T12:00:00Z",
				},
			},
		}
	} else if strings.Contains(ref, "ErrorResponse") {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"value": map[string]interface{}{
					"error":   "Internal server error",
					"message": "Failed to fetch data",
				},
			},
		}
	}
	return map[string]interface{}{}
}

func getArrayExampleForRef(ref string) map[string]interface{} {
	if strings.Contains(ref, "PostResponse") {
		return map[string]interface{}{
			"success": map[string]interface{}{
				"value": []interface{}{
					map[string]interface{}{
						"id":          1,
						"title":       "Sample Post Title",
						"body":        "Sample post body content",
						"userId":      1,
						"formattedAt": "2024-01-01T12:00:00Z",
					},
					map[string]interface{}{
						"id":          2,
						"title":       "Another Post Title",
						"body":        "Another post body content",
						"userId":      2,
						"formattedAt": "2024-01-01T12:30:00Z",
					},
				},
			},
		}
	}
	return map[string]interface{}{}
}

func main() {
	// Read Swagger 2.0 file
	swaggerData, err := ioutil.ReadFile("api/docs/swagger.json")
	if err != nil {
		log.Fatalf("Failed to read swagger.json: %v", err)
	}

	var swagger2 Swagger2
	if err := json.Unmarshal(swaggerData, &swagger2); err != nil {
		log.Fatalf("Failed to parse swagger.json: %v", err)
	}

	// Convert to OpenAPI 3.0.3
	openapi3 := convertSwagger2ToOpenAPI3(swagger2)

	// Marshal to JSON with pretty printing
	openapi3Data, err := json.MarshalIndent(openapi3, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI 3.0.3: %v", err)
	}

	// Write to output file
	if err := ioutil.WriteFile("openapi-specifications/api.swagger.json", openapi3Data, 0644); err != nil {
		log.Fatalf("Failed to write openapi-specifications/api.swagger.json: %v", err)
	}

	fmt.Println("Successfully converted Swagger 2.0 to OpenAPI 3.0.3")
	fmt.Println("Output: openapi-specifications/api.swagger.json")
}