package util

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	_ "template-mobile-app-api/docs"
)

// =================レスポンス=================
type DirectionsResponse struct {
	// https://openrouteservice.org/dev/#/api-docs/v2/directions/{profile}/geojson/get
	// のレスポンスを構造体に
	Type          string         `json:"type"`
	BBox          []float64      `json:"bbox"`
	Features      []ORSFeature   `json:"features"`
	Metadata      ORSMetadata    `json:"metadata"`
	WarningPoints []WarningPoint `json:"warning_points"` //XXX 追加項目
	ComfortScore  int            `json:"comfort_score"`  //XXX 追加項目, 0-100のスコア
	SessoinID     string         `json:"session_id"`     //XXX 追加項目, セッションID
}

// ORSFeature represents a feature in the GeoJSON response
type ORSFeature struct {
	BBox       []float64            `json:"bbox"`
	Type       string               `json:"type"`
	Properties ORSFeatureProperties `json:"properties"`
	Geometry   ORSGeometry          `json:"geometry"`
}

// ORSFeatureProperties contains the properties of a route feature
type ORSFeatureProperties struct {
	Segments  []ORSSegment `json:"segments"`
	Summary   ORSSummary   `json:"summary"`
	WayPoints []int        `json:"way_points"`
}

// ORSGeometry represents the geometry of the route
type ORSGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// ORSSegment represents a segment of the route
type ORSSegment struct {
	Distance float64   `json:"distance"`
	Duration float64   `json:"duration"`
	Steps    []ORSStep `json:"steps"`
}

// ORSStep represents a step within a segment
type ORSStep struct {
	Distance    float64 `json:"distance"`
	Duration    float64 `json:"duration"`
	Type        int     `json:"type"`
	Instruction string  `json:"instruction"`
	Name        string  `json:"name"`
	WayPoints   []int   `json:"way_points"`
}

// ORSManeuver represents maneuver information for a step
type ORSManeuver struct {
	Location      []float64 `json:"location"`
	BearingBefore int       `json:"bearing_before"`
	BearingAfter  int       `json:"bearing_after"`
}

// ORSSummary contains summary information about the route
type ORSSummary struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

// ORSWarning represents a warning message
type ORSWarning struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ORSMetadata contains metadata about the request and response
type ORSMetadata struct {
	Attribution string    `json:"attribution"`
	Service     string    `json:"service"`
	Timestamp   int64     `json:"timestamp"`
	Query       ORSQuery  `json:"query"`
	Engine      ORSEngine `json:"engine"`
}

// ORSQuery represents the original query parameters
type ORSQuery struct {
	Coordinates [][]float64 `json:"coordinates"`
	Profile     string      `json:"profile"`
	ProfileName string      `json:"profileName"`
	Format      string      `json:"format"`
}

// ORSEngine contains information about the OpenRouteService engine
type ORSEngine struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	GraphDate string `json:"graph_date"`
	OSMDate   string `json:"osm_date"`
}

// ORSRouteOptions represents advanced routing options
type ORSRouteOptions struct {
	AvoidFeatures  []string               `json:"avoid_features,omitempty"`
	AvoidBorders   string                 `json:"avoid_borders,omitempty"`
	AvoidCountries []string               `json:"avoid_countries,omitempty"`
	VehicleType    string                 `json:"vehicle_type,omitempty"`
	ProfileParams  *ORSProfileParams      `json:"profile_params,omitempty"`
	AvoidPolygons  map[string]interface{} `json:"avoid_polygons,omitempty"`
	RoundTrip      *ORSRoundTripOptions   `json:"round_trip,omitempty"`
}

// ORSProfileParams represents profile-specific parameters
type ORSProfileParams struct {
	Weightings          *ORSWeightings   `json:"weightings,omitempty"`
	Restrictions        *ORSRestrictions `json:"restrictions,omitempty"`
	SurfaceQualityKnown *bool            `json:"surface_quality_known,omitempty"`
	AllowUnsuitable     *bool            `json:"allow_unsuitable,omitempty"`
}

// ORSWeightings represents weightings for different routing factors
type ORSWeightings struct {
	SteepnessDifficulty *int     `json:"steepness_difficulty,omitempty"`
	Green               *float32 `json:"green,omitempty"`
	Quiet               *float32 `json:"quiet,omitempty"`
	Shadow              *float32 `json:"shadow,omitempty"`
}

// ORSRestrictions represents restrictions for routing
type ORSRestrictions struct {
	Length            *float32 `json:"length,omitempty"`
	Width             *float32 `json:"width,omitempty"`
	Height            *float32 `json:"height,omitempty"`
	Axleload          *float32 `json:"axleload,omitempty"`
	Weight            *float32 `json:"weight,omitempty"`
	Hazmat            *bool    `json:"hazmat,omitempty"`
	SurfaceType       string   `json:"surface_type,omitempty"`
	TrackType         string   `json:"track_type,omitempty"`
	SmoothnessType    string   `json:"smoothness_type,omitempty"`
	MaximumSlopedKerb *float32 `json:"maximum_sloped_kerb,omitempty"`
	MaximumIncline    *int     `json:"maximum_incline,omitempty"`
	MinimumWidth      *float32 `json:"minimum_width,omitempty"`
}

// ORSRoundTripOptions represents options for round trip routing
type ORSRoundTripOptions struct {
	Length *float32 `json:"length,omitempty"`
	Points *int     `json:"points,omitempty"`
	Seed   *int64   `json:"seed,omitempty"`
}

// ORSAlternativeRoutes represents alternative route options
type ORSAlternativeRoutes struct {
	TargetCount  *int     `json:"target_count,omitempty"`
	WeightFactor *float64 `json:"weight_factor,omitempty"`
	ShareFactor  *float64 `json:"share_factor,omitempty"`
}

// ORSCustomModel represents custom model for weighting
type ORSCustomModel struct {
	DistanceInfluence *float64               `json:"distance_influence,omitempty"`
	Speed             []ORSCustomRule        `json:"speed,omitempty"`
	Priority          []ORSCustomRule        `json:"priority,omitempty"`
	Areas             map[string]interface{} `json:"areas,omitempty"`
}

// ORSCustomRule represents a rule in the custom model
type ORSCustomRule struct {
	Keyword   string   `json:"keyword,omitempty"`
	Condition string   `json:"condition,omitempty"`
	Operation string   `json:"operation,omitempty"`
	Value     *float64 `json:"value,omitempty"`
}

// ORSErrorResponse represents error response structure
type ORSErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// getDirections godoc
// @Summary 自転車ルート検索
// @Description 出発地点と目的地の座標をクエリパラメータで受け取り、OpenRouteServiceのAPIを呼び出してルート情報を取得する
// @Tags map
// @Accept json
// @Produce json
// @Param start query string true "出発地点の座標 (経度,緯度)" example:"139.745494,35.659071"
// @Param end query string true "目的地の座標 (経度,緯度)" example:"139.808617,35.709907"
// @Param via_bike_parking query boolean false "自転車駐輪場経由ルート" default(true)
// @Param avoid_bus_stops query boolean false "バス停回避" default(true)
// @Param avoid_traffic_lights query boolean false "信号回避" default(true)
// @Success 200 {object} DirectionsResponse "GeoJson形式のルート情報"
// @Failure 400 {object} ORSErrorResponse "リクエストパラメータ不正"
// @Failure 404 {object} ORSErrorResponse "ルートが見つからない"
// @Failure 500 {object} ORSErrorResponse "サーバー内部エラー"
// @Router /directions/bicycle [get]
func GetDirections(c *gin.Context) {
	// クエリパラメータの取得
	start := c.Query("start")
	end := c.Query("end")
	viaBikeParking := c.DefaultQuery("via_bike_parking", "false")
	avoidBusStops := c.DefaultQuery("avoid_bus_stops", "false")
	avoidTrafficLights := c.DefaultQuery("avoid_traffic_lights", "false")

	status, orsResp := GetDirectionsBase(start, end)
	directionsResponse, ok := orsResp.(DirectionsResponse)
	if ok {
		directionsResponse.Features[0].Geometry.Coordinates = GetBicycleParkingDirection(directionsResponse.Features[0].Geometry.Coordinates, [][]float64{})

		//TODO ComfortScoreはサンプル
		directionsResponse.ComfortScore = 49 // https://github.com/rowicy/charimachi/issues/34
		var comfortLevel = 0
		if viaBikeParking == "true" {
			comfortLevel++
		}
		if avoidBusStops == "true" {
			comfortLevel++
		}
		if avoidTrafficLights == "true" {
			comfortLevel++
		}
		switch comfortLevel {
		case 1:
			directionsResponse.ComfortScore = 79
		case 2:
			directionsResponse.ComfortScore = 89
		case 3:
			directionsResponse.ComfortScore = 100
		}
		directionsResponse.SessoinID = GenerateSessionID()
		SessionIDResponse[directionsResponse.SessoinID] = directionsResponse.Features[0].Geometry
	}
	c.JSON(status, directionsResponse)
}

func GetBicycleParkingDirection(searchCoordinates [][]float64, coordinates [][]float64) (returnCoordinates [][]float64) {
	var distance = 0.0
	var prePosition = searchCoordinates[0]
	for _, v := range searchCoordinates {
		const metersPerDegree = 111320.0 // 緯度1度あたりのメートル
		var lon = prePosition[0] - v[0]
		var lat = prePosition[1] - v[1]
		distance += math.Abs(lon) * metersPerDegree
		distance += math.Abs(lat) * metersPerDegree

		prePosition = v
		fmt.Println("distance:", distance)

		var query = "[bicycle_parking]&viewbox="
		query += strconv.FormatFloat(v[0]-0.011, 'f', -1, 64)
		query += ","
		query += strconv.FormatFloat(v[1]-0.011, 'f', -1, 64)
		query += strconv.FormatFloat(v[0]+0.011, 'f', -1, 64)
		query += ","
		query += strconv.FormatFloat(v[1]+0.011, 'f', -1, 64)
		searchResp := GetSearchBase(query)
		searchResponse, ok := searchResp.([]SearchResponse)
		if ok && (distance > 500) && len(searchResponse) != 0 {
			_, orsResp := GetDirectionsBase(searchResponse[0].Lon, searchResponse[0].Lat)
			if directionsResponse, ok := orsResp.(DirectionsResponse); ok {
				coordinates = append(coordinates, v)
				coordinates = GetBicycleParkingDirection(directionsResponse.Features[0].Geometry.Coordinates, coordinates)
				return coordinates
			}
		} else {
			coordinates = append(coordinates, v)
		}
	}

	return coordinates
}

func GetDirectionsBase(
	start string,
	end string) (status int, res any) {

	// バリデーション
	if start == "" || end == "" {
		var er ORSErrorResponse
		er.Error.Code = http.StatusBadRequest
		er.Error.Message = "start and end query parameters are required"
		return http.StatusBadRequest, er
	}

	// APIキー取得
	apiKey := os.Getenv("OPEN_ROUTE_SERVICE_API_KEY")
	if apiKey == "" {
		var er ORSErrorResponse
		er.Error.Code = http.StatusInternalServerError
		er.Error.Message = "OPEN_ROUTE_SERVICE_API_KEY is not set"
		return http.StatusInternalServerError, er
	}

	// ORSへリクエスト
	base := "https://api.openrouteservice.org/v2/directions/cycling-road"
	u, err := url.Parse(base)
	if err != nil {
		var er ORSErrorResponse
		er.Error.Code = http.StatusInternalServerError
		er.Error.Message = fmt.Sprintf("invalid base url: %v", err)
		return http.StatusInternalServerError, er
	}

	q := u.Query()
	q.Set("api_key", apiKey)
	q.Set("start", start)
	q.Set("end", end)
	u.RawQuery = q.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		var er ORSErrorResponse
		er.Error.Code = http.StatusBadGateway
		er.Error.Message = err.Error()
		return http.StatusBadGateway, er
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		var er ORSErrorResponse
		er.Error.Code = http.StatusBadGateway
		er.Error.Message = fmt.Sprintf("failed to read upstream response: %v", err)
		return http.StatusBadGateway, er
	}

	if resp.StatusCode != http.StatusOK {
		var upstreamErr ORSErrorResponse
		_ = json.Unmarshal(body, &upstreamErr)
		if upstreamErr.Error.Message == "" {
			upstreamErr.Error.Code = resp.StatusCode
			upstreamErr.Error.Message = fmt.Sprintf("upstream returned status %d", resp.StatusCode)
		}
		return http.StatusBadGateway, upstreamErr
	}

	//TODO OpenRouteServiceのレスポンスをそのままレスポンスにパースしているのでorsRespを加工する処理をかく
	var orsResp DirectionsResponse
	if err := json.Unmarshal(body, &orsResp); err != nil {
		var er ORSErrorResponse
		er.Error.Code = http.StatusInternalServerError
		er.Error.Message = fmt.Sprintf("failed to parse upstream response: %v", err)
		return http.StatusInternalServerError, er
	}

	//TODO WarningPointsはサンプル
	// if orsResp.WarningPoints == nil {
	// 	coodinates := orsResp.Features[0].Geometry.Coordinates
	// 	orsResp.WarningPoints = []WarningPoint{
	// 		{
	// 			Type:       "intersection",
	// 			Name:       "地獄谷",
	// 			Coordinate: coodinates[1],
	// 			Message:    "過去に事故が多発した地点です。注意してください。",
	// 		},
	// 		{
	// 			Type:       "straight_road",
	// 			Name:       "無限道路",
	// 			Coordinate: coodinates[len(coodinates)-2],
	// 			Message:    "この道路は直線で、速度を出しやすいです。安全運転を心掛けてください。",
	// 		},
	// 	}
	// }

	return http.StatusOK, orsResp
}
