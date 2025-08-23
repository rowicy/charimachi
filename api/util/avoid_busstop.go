package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// BusStop represents a bus stop with its polygon data
type BusStop struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Polygon   [][]float64 `json:"polygon"`
}

// Coordinate represents a longitude, latitude pair
type Coordinate [2]float64

// RouteRequest represents the request body for OpenRouteService API
type RouteRequest struct {
	Coordinates []Coordinate           `json:"coordinates"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

const (
	OpenRouteServiceURL = "https://api.openrouteservice.org/v2/directions/cycling-road/geojson"
)

// AvoidBusStops reads bus stops data and gets a route avoiding all bus stop polygons
func AvoidBusStops(startCoord, endCoord Coordinate) (ORSGeometry, error) {
	// Read bus stops data
	busStopsData, err := os.ReadFile("data/bus_stops.json")
	if err != nil {
		return ORSGeometry{}, fmt.Errorf("failed to read bus stops file: %v", err)
	}

	var busStops []BusStop
	if err := json.Unmarshal(busStopsData, &busStops); err != nil {
		return ORSGeometry{}, fmt.Errorf("failed to parse bus stops JSON: %v", err)
	}

	// Convert bus stop polygons to avoid polygons
	var avoidPolygons [][][][]float64
	for _, stop := range busStops {
		if len(stop.Polygon) > 0 {
			avoidPolygons = append(avoidPolygons, [][][]float64{stop.Polygon})
		}
	}

	// Create request body
	requestBody := RouteRequest{
		Coordinates: []Coordinate{startCoord, endCoord},
		Options: map[string]interface{}{
			"avoid_polygons": map[string]interface{}{
				"type":        "MultiPolygon",
				"coordinates": avoidPolygons,
			},
		},
	}

	// Make API request
	geometry, err := makeOpenRouteServiceRequest(requestBody)
	if err != nil {
		return ORSGeometry{}, err
	}

	return *geometry, nil
}

// GetRouteAvoidingSinglePolygon gets a route avoiding a single polygon
func GetRouteAvoidingSinglePolygon(startCoord, endCoord Coordinate, avoidPolygon [][]float64) (*ORSGeometry, error) {
	requestBody := RouteRequest{
		Coordinates: []Coordinate{startCoord, endCoord},
		Options: map[string]interface{}{
			"avoid_polygons": map[string]interface{}{
				"type":        "Polygon",
				"coordinates": [][][]float64{avoidPolygon},
			},
		},
	}

	return makeOpenRouteServiceRequest(requestBody)
}

// makeOpenRouteServiceRequest makes the actual HTTP request to OpenRouteService API
func makeOpenRouteServiceRequest(requestBody RouteRequest) (*ORSGeometry, error) {
	// Marshal request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", OpenRouteServiceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", os.Getenv("OPEN_ROUTE_SERVICE_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:142.0) Gecko/20100101 Firefox/142.0")

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error! status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response and extract geometry
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// Extract geometry from features[0].geometry
	features, ok := response["features"].([]interface{})
	if !ok || len(features) == 0 {
		return nil, fmt.Errorf("no features found in response")
	}

	feature, ok := features[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid feature format")
	}

	geometryData, ok := feature["geometry"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no geometry found in feature")
	}

	// Convert to ORSGeometry
	geometryJSON, err := json.Marshal(geometryData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal geometry: %v", err)
	}

	var geometry ORSGeometry
	if err := json.Unmarshal(geometryJSON, &geometry); err != nil {
		return nil, fmt.Errorf("failed to parse geometry: %v", err)
	}

	return &geometry, nil
}
