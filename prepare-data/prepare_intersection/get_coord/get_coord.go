package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Coordinate represents a latitude/longitude point
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Way represents a road way from OpenStreetMap
type Way struct {
	Type     string       `json:"type"`
	ID       int64        `json:"id"`
	Geometry []Coordinate `json:"geometry"`
	Tags     struct {
		Name string `json:"name"`
	} `json:"tags"`
}

// OverpassResponse represents the response from Overpass API
type OverpassResponse struct {
	Elements []Way `json:"elements"`
}

// IntersectionResult holds the intersection coordinates between two roads
type IntersectionResult struct {
	Road1Name     string       `json:"road1_name"`
	Road2Name     string       `json:"road2_name"`
	Intersections []Coordinate `json:"intersections"`
	Road1Segments []Way        `json:"road1_segments"`
	Road2Segments []Way        `json:"road2_segments"`
}

// Client handles Overpass API requests
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new Overpass API client
func NewClient() *Client {
	return &Client{
		baseURL: "https://overpass-api.de/api/interpreter",
		client: &http.Client{
			Transport: &http.Transport{
				DisableCompression: false, // Enable automatic decompression
			},
		},
	}
}

// GetRoadData fetches road data for a specific road name
func (c *Client) GetRoadData(roadName string, bbox [4]float64) (*OverpassResponse, error) {
	// Let's manually reconstruct the exact curl command data
	// The original was: data=%5Bout%3Ajson%5D%5Btimeout%3A25%5D%3B%0Away%5B%22highway%22%5D%5B%22name%22%3D%22%E7%AC%AC%E4%B8%80%E4%BA%AC%E6%B5%9C%22%5D(35.66%2C139.74%2C35.67%2C139.77)%3B%0Aout+geom%3B%0A
	// Let's decode this step by step:
	// %5B = [, %3A = :, %3D = =, %22 = ", %0A = newline, %2C = comma, etc.

	// The decoded query should be:
	// [out:json][timeout:25];
	// way["highway"]["name"="第一京浜"](35.66,139.74,35.67,139.77);
	// out geom;
	// (with actual newlines)

	// But the curl shows out+geom in the URL encoding, let me check if that's literally what it sends
	query := fmt.Sprintf("[out:json][timeout:25];\nway[\"highway\"][\"name\"=\"%s\"](%.2f,%.2f,%.2f,%.2f);\nout geom;\n",
		roadName, bbox[0], bbox[1], bbox[2], bbox[3])

	fmt.Printf("=== DEBUG: Overpass Query for '%s' ===\n", roadName)
	fmt.Printf("Query:\n%s", query)

	// Let's try to match the exact curl encoding
	// Create form data like curl does
	formData := url.Values{}
	formData.Set("data", query)
	requestBody := formData.Encode()

	fmt.Printf("Request Body (url.Values): %s\n", requestBody)

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBufferString(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to match the working curl command
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:142.0) Gecko/20100101 Firefox/142.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.7,en;q=0.3")
	// Remove Accept-Encoding to let Go handle compression automatically
	req.Header.Set("Origin", "https://overpass-turbo.eu")
	req.Header.Set("Referer", "https://overpass-turbo.eu/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")

	fmt.Printf("Making request to: %s\n", c.baseURL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("Response Body Length: %d bytes\n", len(body))
	fmt.Printf("Response Content-Type: %s\n", resp.Header.Get("Content-Type"))
	fmt.Printf("Response Content-Encoding: %s\n", resp.Header.Get("Content-Encoding"))

	// Check first few bytes for debugging
	if len(body) > 0 {
		fmt.Printf("First 10 bytes (hex): ")
		for i := 0; i < min(10, len(body)); i++ {
			fmt.Printf("%02x ", body[i])
		}
		fmt.Printf("\n")
	}

	// Check if response is HTML (error page)
	bodyStr := string(body)
	if len(bodyStr) > 0 && (strings.Contains(strings.ToLower(bodyStr), "<html>") || strings.Contains(strings.ToLower(bodyStr), "<!doctype html>")) {
		fmt.Printf("ERROR: Received HTML/XML error response!\n")
		fmt.Printf("Full Error Response:\n%s\n", bodyStr)
		return nil, fmt.Errorf("received HTML response instead of JSON - API error: %s", resp.Status)
	}

	// Only show response body if it looks like valid JSON (starts with '{' or '[')
	if len(bodyStr) > 0 && (bodyStr[0] == '{' || bodyStr[0] == '[') {
		fmt.Printf("Raw Response Body:\n%s\n", bodyStr)
	} else {
		fmt.Printf("Response does not appear to be JSON. Full response:\n%s\n", bodyStr)
	}

	var result OverpassResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("JSON Parse Error: %v\n", err)
		fmt.Printf("Attempting to parse first 500 chars: %s\n", bodyStr[:min(500, len(bodyStr))])
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Parsed Elements Count: %d\n", len(result.Elements))
	for i, element := range result.Elements {
		fmt.Printf("  Element %d: ID=%d, Type=%s, Name=%s, Geometry Points=%d\n",
			i, element.ID, element.Type, element.Tags.Name, len(element.Geometry))
		if len(element.Geometry) > 0 {
			fmt.Printf("    First point: Lat=%.6f, Lon=%.6f\n", element.Geometry[0].Lat, element.Geometry[0].Lon)
			fmt.Printf("    Last point:  Lat=%.6f, Lon=%.6f\n", element.Geometry[len(element.Geometry)-1].Lat, element.Geometry[len(element.Geometry)-1].Lon)
		}
	}
	fmt.Printf("=== END DEBUG for '%s' ===\n\n", roadName)

	return &result, nil
}

// distance calculates the distance between two coordinates in meters
func distance(p1, p2 Coordinate) float64 {
	const earthRadius = 6371000 // meters
	lat1Rad := p1.Lat * math.Pi / 180
	lat2Rad := p2.Lat * math.Pi / 180
	deltaLat := (p2.Lat - p1.Lat) * math.Pi / 180
	deltaLon := (p2.Lon - p1.Lon) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// lineSegmentIntersection finds the intersection point of two line segments
func lineSegmentIntersection(p1, p2, p3, p4 Coordinate) (*Coordinate, bool) {
	x1, y1 := p1.Lon, p1.Lat
	x2, y2 := p2.Lon, p2.Lat
	x3, y3 := p3.Lon, p3.Lat
	x4, y4 := p4.Lon, p4.Lat

	denom := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(denom) < 1e-10 {
		return nil, false // Lines are parallel
	}

	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / denom
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / denom

	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		// Intersection point
		x := x1 + t*(x2-x1)
		y := y1 + t*(y2-y1)
		return &Coordinate{Lat: y, Lon: x}, true
	}

	return nil, false
}

// findIntersections finds all intersection points between two sets of road ways
func findIntersections(ways1, ways2 []Way, tolerance float64) []Coordinate {
	var intersections []Coordinate
	seen := make(map[string]bool)

	fmt.Printf("=== DEBUG: Finding Intersections ===\n")
	fmt.Printf("Road 1 ways: %d, Road 2 ways: %d\n", len(ways1), len(ways2))
	fmt.Printf("Tolerance: %.2f meters\n", tolerance)

	intersectionCount := 0

	for i, way1 := range ways1 {
		fmt.Printf("Processing Way1[%d] (ID: %d) with %d geometry points\n", i, way1.ID, len(way1.Geometry))
		for segIdx := 0; segIdx < len(way1.Geometry)-1; segIdx++ {
			seg1Start := way1.Geometry[segIdx]
			seg1End := way1.Geometry[segIdx+1]

			for j, way2 := range ways2 {
				for segIdx2 := 0; segIdx2 < len(way2.Geometry)-1; segIdx2++ {
					seg2Start := way2.Geometry[segIdx2]
					seg2End := way2.Geometry[segIdx2+1]

					if intersection, found := lineSegmentIntersection(seg1Start, seg1End, seg2Start, seg2End); found {
						intersectionCount++
						fmt.Printf("  Found line intersection #%d:\n", intersectionCount)
						fmt.Printf("    Way1[%d] segment %d: (%.6f,%.6f) -> (%.6f,%.6f)\n",
							i, segIdx, seg1Start.Lat, seg1Start.Lon, seg1End.Lat, seg1End.Lon)
						fmt.Printf("    Way2[%d] segment %d: (%.6f,%.6f) -> (%.6f,%.6f)\n",
							j, segIdx2, seg2Start.Lat, seg2Start.Lon, seg2End.Lat, seg2End.Lon)
						fmt.Printf("    Intersection point: (%.6f,%.6f)\n", intersection.Lat, intersection.Lon)

						// Create a unique key for this intersection point
						key := fmt.Sprintf("%.6f,%.6f", intersection.Lat, intersection.Lon)
						if !seen[key] {
							intersections = append(intersections, *intersection)
							seen[key] = true
							fmt.Printf("    -> Added to results (unique)\n")
						} else {
							fmt.Printf("    -> Skipped (duplicate)\n")
						}
					}
				}
			}
		}
	}

	// Also check for very close endpoints (within tolerance)
	proximityCount := 0
	fmt.Printf("\nChecking for proximity intersections (tolerance: %.2fm)...\n", tolerance)
	for i, way1 := range ways1 {
		for coordIdx1, coord1 := range way1.Geometry {
			for j, way2 := range ways2 {
				for coordIdx2, coord2 := range way2.Geometry {
					dist := distance(coord1, coord2)
					if dist <= tolerance {
						proximityCount++
						fmt.Printf("  Found proximity intersection #%d:\n", proximityCount)
						fmt.Printf("    Way1[%d] point %d: (%.6f,%.6f)\n", i, coordIdx1, coord1.Lat, coord1.Lon)
						fmt.Printf("    Way2[%d] point %d: (%.6f,%.6f)\n", j, coordIdx2, coord2.Lat, coord2.Lon)
						fmt.Printf("    Distance: %.2f meters\n", dist)

						key := fmt.Sprintf("%.6f,%.6f", coord1.Lat, coord1.Lon)
						if !seen[key] {
							intersections = append(intersections, coord1)
							seen[key] = true
							fmt.Printf("    -> Added to results (unique)\n")
						} else {
							fmt.Printf("    -> Skipped (duplicate)\n")
						}
					}
				}
			}
		}
	}

	fmt.Printf("\nIntersection Summary:\n")
	fmt.Printf("  Line intersections found: %d\n", intersectionCount)
	fmt.Printf("  Proximity intersections found: %d\n", proximityCount)
	fmt.Printf("  Unique intersections: %d\n", len(intersections))
	fmt.Printf("=== END DEBUG: Finding Intersections ===\n\n")

	return intersections
}

// FindRoadIntersections makes 2 POST requests with different road names and calculates intersection coordinates
func FindRoadIntersections(road1Name, road2Name string, bbox [4]float64) (*IntersectionResult, error) {
	fmt.Printf("=== DEBUG: Starting Road Intersection Analysis ===\n")
	fmt.Printf("Making 2 POST requests to change name field:\n")
	fmt.Printf("POST Request 1: name=\"%s\"\n", road1Name)
	fmt.Printf("POST Request 2: name=\"%s\"\n", road2Name)
	fmt.Printf("BBox: [%.2f, %.2f, %.2f, %.2f]\n", bbox[0], bbox[1], bbox[2], bbox[3])

	client := NewClient()

	// POST Request 1: Get data for first road
	fmt.Printf("\n=== POST Request 1 ===\n")
	fmt.Printf("Fetching coordinate area data for Road: %s\n", road1Name)
	road1Data, err := client.GetRoadData(road1Name, bbox)
	if err != nil {
		return nil, fmt.Errorf("POST request 1 failed for road (%s): %w", road1Name, err)
	}

	// POST Request 2: Get data for second road
	fmt.Printf("\n=== POST Request 2 ===\n")
	fmt.Printf("Fetching coordinate area data for Road: %s\n", road2Name)
	road2Data, err := client.GetRoadData(road2Name, bbox)
	if err != nil {
		return nil, fmt.Errorf("POST request 2 failed for road (%s): %w", road2Name, err)
	}

	if len(road1Data.Elements) == 0 {
		fmt.Printf("ERROR: No coordinate area data found for road: %s\n", road1Name)
		return nil, fmt.Errorf("no coordinate area data found for road: %s", road1Name)
	}

	if len(road2Data.Elements) == 0 {
		fmt.Printf("ERROR: No coordinate area data found for road: %s\n", road2Name)
		return nil, fmt.Errorf("no coordinate area data found for road: %s", road2Name)
	}

	fmt.Printf("\nPOST Response Summary:\n")
	fmt.Printf("  POST Response 1 (%s): %d coordinate area segments\n", road1Name, len(road1Data.Elements))
	fmt.Printf("  POST Response 2 (%s): %d coordinate area segments\n", road2Name, len(road2Data.Elements))

	// Calculate intersection coordinates between the two coordinate-based areas
	fmt.Printf("\nCalculating intersection coordinates between coordinate areas...\n")
	intersections := findIntersections(road1Data.Elements, road2Data.Elements, 10.0)

	result := &IntersectionResult{
		Road1Name:     road1Name,
		Road2Name:     road2Name,
		Intersections: intersections,
		Road1Segments: road1Data.Elements,
		Road2Segments: road2Data.Elements,
	}

	fmt.Printf("=== DEBUG: Final Intersection Coordinates ===\n")
	fmt.Printf("Calculated intersection coordinates: %d\n", len(result.Intersections))
	for i, intersection := range result.Intersections {
		fmt.Printf("  Intersection Coordinate %d: Lat=%.6f, Lon=%.6f\n", i+1, intersection.Lat, intersection.Lon)
	}
	fmt.Printf("=== END DEBUG ===\n\n")

	return result, nil
}

func Find(road1Name string, road2Name string) ([]float64, error) {
	// Example usage - exactly as requested: change the "name" field and make 2 POST requests
	bbox := [4]float64{35.52, 138.94, 35.90, 139.92} //XXX Tokyo area from your example

	// Find intersections between the two roads from the 2 POST responses
	result, err := FindRoadIntersections(road1Name, road2Name, bbox)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	// Print results as struct output
	fmt.Printf("\n=== Road Intersection Analysis (Struct Output) ===\n")
	fmt.Printf("Road 1: %s (%d segments)\n", result.Road1Name, len(result.Road1Segments))
	fmt.Printf("Road 2: %s (%d segments)\n", result.Road2Name, len(result.Road2Segments))
	fmt.Printf("Found %d intersection coordinate(s):\n", len(result.Intersections))

	for i, intersection := range result.Intersections {
		fmt.Printf("  Intersection %d: {Lat: %.6f, Lon: %.6f}\n", i+1, intersection.Lat, intersection.Lon)
	}

	// Output the complete struct as JSON for easy consumption
	fmt.Printf("\n=== Complete Struct Output (JSON) ===\n")
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return nil, err
	}
	fmt.Printf("%s\n", string(jsonOutput))

	// Example of how to use the result struct programmatically
	fmt.Printf("\n=== Programmatic Usage Example ===\n")
	fmt.Printf("Total intersection points: %d\n", len(result.Intersections))
	fmt.Printf("Road1 has %d way segments\n", len(result.Road1Segments))
	fmt.Printf("Road2 has %d way segments\n", len(result.Road2Segments))

	if len(result.Intersections) == 0 {
		return nil, fmt.Errorf("no intersections found between roads '%s' and '%s'", road1Name, road2Name)
	}

	firstIntersection := result.Intersections[0]
	return []float64{firstIntersection.Lon, firstIntersection.Lat}, nil

}

// 有効数字で丸める関数
func roundToSignificantFigures(x float64, sig int) float64 {
	if x == 0 {
		return 0
	}
	power := math.Floor(math.Log10(math.Abs(x))) + 1 - float64(sig)
	return math.Round(x/math.Pow(10, power)) * math.Pow(10, power)
}

func main() {
	filesPattern := flag.String("files", "", "CSV files to process (supports *)")
	outdir := flag.String("outdir", ".", "Output dir")
	flag.Parse()
	if *filesPattern == "" {
		fmt.Println("Usage: -files <pattern>")
		return
	}
	files, err := filepath.Glob(*filesPattern)
	if err != nil {
		fmt.Printf("Glob error: %v\n", err)
		return
	}
	if len(files) == 0 {
		fmt.Println("No files matched")
		return
	}

	// JSONとしてファイル出力
	outFilePath := filepath.Join(*outdir, "hoge.json")
	outFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer outFile.Close()

	var results []ViolationRate

	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		f.Close()
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			continue
		}

		for _, record := range records[1:] { // ヘッダーはスキップ
			if len(record) < 11 {
				continue
			}

			cycleTotal, err := strconv.ParseFloat(record[10], 64) // 自転車計
			if err != nil || cycleTotal == 0 {
				continue
			}
			sidewalk, err := strconv.ParseFloat(record[9], 64) // 歩道通行自転車
			if err != nil {
				continue
			}

			ratio := sidewalk / cycleTotal
			ratioRounded := roundToSignificantFigures(ratio, 2)

			roads := strings.Split(record[2], "Ｘ")
			road1 := roads[0]
			road2 := roads[1]
			// Findで座標取得
			coord, err := Find(road1, road2)
			if err != nil {
				fmt.Printf("Find error for %s × %s: %v\n", record[0], record[1], err)
				continue // エラーなら JSON 出力から除外
			}

			vr := ViolationRate{
				Type:           "intersection",
				ViolationRate:  ratioRounded,
				ViolationCount: int(sidewalk),
				Coordinate:     coord,
				Message:        "",
			}

			results = append(results, vr)
		}
	}

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

type ViolationRate struct {
	Type           string    `json:"type,omitempty"`            //常に"intersection"
	ViolationRate  float64   `json:"violation_rate"`            //歩道通行自転車/自転車計(有効数字2桁)
	ViolationCount int       `json:"violation_count,omitempty"` //歩道通行自転車数(整数)
	Coordinate     []float64 `json:"coordinate"`                //[経度, 緯度]
	Message        string    `json:"message,omitempty"`         //スコアによって用意している定数を変える
}
