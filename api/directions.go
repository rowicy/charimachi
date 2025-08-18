package main

import (
	"net/http"

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
	WarningPoints []WarningPoint `json:"warningPoints"` //XXX 追加項目
}

// XXX カスタム構造体
type WarningPoint struct {
	Type       string    `json:"type"`       //交差点なのか直線道路なのか(このフィールドいらない)
	Name       string    `json:"name"`       //名称, 地元での言われ名など(地獄谷etc)
	Coordinate []float64 `json:"coordinate"` //座標
	Message    string    `json:"message"`    //警告のメッセージ(どんな事故が多かったか)
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
func getDirections(c *gin.Context) {
	// クエリパラメータの取得
	start := c.Query("start")
	end := c.Query("end")
	viaBikeParking := c.DefaultQuery("via_bike_parking", "false")
	avoidBusStops := c.DefaultQuery("avoid_bus_stops", "false")
	avoidTrafficLights := c.DefaultQuery("avoid_traffic_lights", "false")

	// TODO: start, endの座標をパースしてOpenRouteServiceのAPIを呼び出す
	// TODO: via_bike_parking, avoid_bus_stops, avoid_traffic_lightsの処理を実装
	_ = start
	_ = end
	_ = viaBikeParking
	_ = avoidBusStops
	_ = avoidTrafficLights

	// ORSへのリクエスト
	// ルート整形処理

	directionsResponse := DirectionsResponse{
		Type: "FeatureCollection",
		BBox: []float64{139.745496, 35.658588, 139.746388, 35.659073},
		Features: []ORSFeature{
			{
				BBox: []float64{139.745496, 35.658588, 139.746388, 35.659073},
				Type: "Feature",
				Properties: ORSFeatureProperties{
					Segments: []ORSSegment{
						{
							Distance: 109.6,
							Duration: 32.9,
							Steps: []ORSStep{
								{
									Distance:    26.6,
									Duration:    15.9,
									Type:        11,
									Instruction: "Head southeast",
									Name:        "-",
									WayPoints:   []int{0, 1},
								},
								{
									Distance:    3.8,
									Duration:    2.3,
									Type:        0,
									Instruction: "Turn left",
									Name:        "-",
									WayPoints:   []int{1, 2},
								},
								{
									Distance:    63.3,
									Duration:    8.8,
									Type:        1,
									Instruction: "Turn right onto 東京タワー通り",
									Name:        "東京タワー通り",
									WayPoints:   []int{2, 3},
								},
								{
									Distance:    8.0,
									Duration:    1.1,
									Type:        1,
									Instruction: "Turn right onto 降車場",
									Name:        "降車場",
									WayPoints:   []int{3, 4},
								},
								{
									Distance:    7.9,
									Duration:    4.8,
									Type:        0,
									Instruction: "Turn left",
									Name:        "-",
									WayPoints:   []int{4, 5},
								},
								{
									Distance:    0.0,
									Duration:    0.0,
									Type:        10,
									Instruction: "Arrive at your destination, on the left",
									Name:        "-",
									WayPoints:   []int{5, 5},
								},
							},
						},
					},
					Summary: ORSSummary{
						Distance: 109.6,
						Duration: 32.9,
					},
					WayPoints: []int{0, 5},
				},
				Geometry: ORSGeometry{
					Type: "LineString",
					Coordinates: [][]float64{
						{139.745496, 35.659073},
						{139.74574, 35.65894},
						{139.745761, 35.65897},
						{139.746365, 35.65868},
						{139.746312, 35.658623},
						{139.746388, 35.658588},
					},
				},
			},
		},
		Metadata: ORSMetadata{
			Attribution: "openrouteservice.org | OpenStreetMap contributors",
			Service:     "routing",
			Timestamp:   1755538285980,
			Query: ORSQuery{
				Coordinates: [][]float64{{139.745494, 35.659071}, {139.746396, 35.6586}},
				Profile:     "cycling-road",
				ProfileName: "cycling-road",
				Format:      "json",
			},
			Engine: ORSEngine{
				Version:   "9.3.0",
				BuildDate: "2025-06-06T15:39:25Z",
				GraphDate: "2025-08-13T04:04:04Z",
				OSMDate:   "2025-08-04T00:00:00Z",
			},
		},
		WarningPoints: []WarningPoint{
			{
				Type:       "intersection",
				Name:       "地獄谷",
				Coordinate: []float64{139.745494, 35.659071},
				Message:    "過去に事故が多発した地点です。注意してください。",
			},
			{
				Type:       "straight_road",
				Name:       "直線道路",
				Coordinate: []float64{139.747726, 35.689709},
				Message:    "この道路は直線で、速度を出しやすいです。安全運転を心掛けてください。",
			},
		},
	}

	c.JSON(http.StatusOK, directionsResponse)
}
