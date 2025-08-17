package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	_ "template-mobile-app-api/docs"
)

type ORSDirectionsRequest struct {
	// Required field - The waypoints to use for the route as an array of longitude/latitude pairs in WGS 84
	Coordinates        [][]float64 `json:"coordinates" validate:"required"`
	ViaBikeParking     bool        `json:"viaBikeParking"`
	AvoidBusStops      bool        `json:"avoidBusStops"`
	AvoidTrafficLights bool        `json:"avoidTrafficLights"`

	// // Optional fields
	// ID                 string                   `json:"id,omitempty"`
	// Preference         string                   `json:"preference,omitempty"`         // "fastest", "shortest", "recommended", "custom"
	// Units              string                   `json:"units,omitempty"`              // "m", "km", "mi"
	// Language           string                   `json:"language,omitempty"`           // Language code for route instructions
	// Geometry           *bool                    `json:"geometry,omitempty"`           // Specifies whether to return geometry
	// Instructions       *bool                    `json:"instructions,omitempty"`       // Specifies whether to return instructions
	// InstructionsFormat string                   `json:"instructions_format,omitempty"` // "html", "text"
	// RoundaboutExits    *bool                    `json:"roundabout_exits,omitempty"`   // Provides bearings of roundabout exits
	// Attributes         []string                 `json:"attributes,omitempty"`         // List of route attributes
	// Maneuvers          *bool                    `json:"maneuvers,omitempty"`          // Include maneuver object in steps
	// Radiuses           []float64                `json:"radiuses,omitempty"`           // Search radius for each waypoint
	// Bearings           [][]float64              `json:"bearings,omitempty"`           // Bearing constraints for waypoints
	// ContinueStraight   *bool                    `json:"continue_straight,omitempty"`  // Forces route to continue straight at waypoints
	// Elevation          *bool                    `json:"elevation,omitempty"`          // Return elevation values for points
	// ExtraInfo          []string                 `json:"extra_info,omitempty"`         // Extra info items to include
	// Options            *ORSRequestOptions       `json:"options,omitempty"`            // Advanced routing options
	// SuppressWarnings   *bool                    `json:"suppress_warnings,omitempty"`  // Suppress warning messages
	// GeometrySimplify   *bool                    `json:"geometry_simplify,omitempty"`  // Simplify the geometry
	// SkipSegments       []int                    `json:"skip_segments,omitempty"`      // Segments to skip in route calculation
	// AlternativeRoutes  *ORSRequestAlternativeRoutes `json:"alternative_routes,omitempty"` // Alternative routes parameters
	// MaximumSpeed       *float64                 `json:"maximum_speed,omitempty"`      // Maximum speed for driving profiles

	// // Public transport specific fields
	// Schedule         *bool  `json:"schedule,omitempty"`          // Return public transport schedule
	// ScheduleDuration string `json:"schedule_duration,omitempty"` // Schedule duration (e.g., "PT2H30M")
	// ScheduleRows     *int   `json:"schedule_rows,omitempty"`     // Maximum schedule entries
	// WalkingTime      string `json:"walking_time,omitempty"`      // Walking time (e.g., "PT2H30M")
	// IgnoreTransfers  *bool  `json:"ignore_transfers,omitempty"`  // Ignore transfers criterion

	// // Custom model
	// CustomModel *ORSRequestCustomModel `json:"custom_model,omitempty"` // Custom weighting model
}

// // ORSRequestOptions represents advanced routing options
// type ORSRequestOptions struct {
// 	AvoidFeatures  []string                      `json:"avoid_features,omitempty"`  // Features to avoid
// 	AvoidBorders   string                        `json:"avoid_borders,omitempty"`   // Border crossing type to avoid
// 	AvoidCountries []string                      `json:"avoid_countries,omitempty"` // Countries to exclude
// 	VehicleType    string                        `json:"vehicle_type,omitempty"`    // Vehicle type for HGV profile
// 	ProfileParams  *ORSRequestProfileParams      `json:"profile_params,omitempty"`  // Profile-specific parameters
// 	AvoidPolygons  map[string]interface{}        `json:"avoid_polygons,omitempty"`  // Areas to avoid (GeoJSON)
// 	RoundTrip      *ORSRequestRoundTripOptions   `json:"round_trip,omitempty"`      // Round trip parameters
// }

// // ORSRequestProfileParams represents profile-specific parameters
// type ORSRequestProfileParams struct {
// 	Weightings          *ORSRequestWeightings    `json:"weightings,omitempty"`
// 	Restrictions        *ORSRequestRestrictions  `json:"restrictions,omitempty"`
// 	SurfaceQualityKnown *bool                    `json:"surface_quality_known,omitempty"` // For wheelchair profile
// 	AllowUnsuitable     *bool                    `json:"allow_unsuitable,omitempty"`      // For wheelchair profile
// }

// // ORSRequestWeightings represents weightings for different routing factors
// type ORSRequestWeightings struct {
// 	SteepnessDifficulty *int     `json:"steepness_difficulty,omitempty"` // Fitness level for cycling profiles (0-3)
// 	Green               *float32 `json:"green,omitempty"`                // Green factor for foot profiles (0-1)
// 	Quiet               *float32 `json:"quiet,omitempty"`                // Quiet factor for foot profiles (0-1)
// 	Shadow              *float32 `json:"shadow,omitempty"`               // Shadow factor for foot profiles (0-1)
// }

// // ORSRequestRestrictions represents restrictions for routing
// type ORSRequestRestrictions struct {
// 	// HGV restrictions
// 	Length   *float32 `json:"length,omitempty"`   // Length restriction in metres
// 	Width    *float32 `json:"width,omitempty"`    // Width restriction in metres
// 	Height   *float32 `json:"height,omitempty"`   // Height restriction in metres
// 	Axleload *float32 `json:"axleload,omitempty"` // Axleload restriction in tons
// 	Weight   *float32 `json:"weight,omitempty"`   // Weight restriction in tons
// 	Hazmat   *bool    `json:"hazmat,omitempty"`   // Hazardous goods routing

// 	// Wheelchair restrictions
// 	SurfaceType        string   `json:"surface_type,omitempty"`        // Minimum surface type
// 	TrackType          string   `json:"track_type,omitempty"`          // Minimum grade of route
// 	SmoothnessType     string   `json:"smoothness_type,omitempty"`     // Minimum smoothness
// 	MaximumSlopedKerb  *float32 `json:"maximum_sloped_kerb,omitempty"` // Maximum sloped curb height
// 	MaximumIncline     *int     `json:"maximum_incline,omitempty"`     // Maximum incline percentage
// 	MinimumWidth       *float32 `json:"minimum_width,omitempty"`       // Minimum footway width
// }

// // ORSRequestRoundTripOptions represents round trip routing parameters
// type ORSRequestRoundTripOptions struct {
// 	Length *float32 `json:"length,omitempty"` // Target route length in metres
// 	Points *int     `json:"points,omitempty"` // Number of points for route
// 	Seed   *int64   `json:"seed,omitempty"`   // Randomization seed
// }

// // ORSRequestAlternativeRoutes represents alternative route parameters
// type ORSRequestAlternativeRoutes struct {
// 	TargetCount  *int     `json:"target_count,omitempty"`  // Target number of alternative routes
// 	WeightFactor *float64 `json:"weight_factor,omitempty"` // Maximum weight factor divergence
// 	ShareFactor  *float64 `json:"share_factor,omitempty"`  // Maximum shared path fraction
// }

// // ORSRequestCustomModel represents custom weighting model
// type ORSRequestCustomModel struct {
// 	DistanceInfluence *float64                   `json:"distance_influence,omitempty"` // Distance influence parameter
// 	Speed             []ORSRequestCustomRule     `json:"speed,omitempty"`              // Speed rules
// 	Priority          []ORSRequestCustomRule     `json:"priority,omitempty"`           // Priority rules
// 	Areas             map[string]interface{}     `json:"areas,omitempty"`              // Named areas for rules
// }

// // ORSRequestCustomRule represents a rule in the custom model
// type ORSRequestCustomRule struct {
// 	Keyword   string   `json:"keyword,omitempty"`   // "IF", "ELSEIF", "ELSE"
// 	Condition string   `json:"condition,omitempty"` // Condition string
// 	Operation string   `json:"operation,omitempty"` // "MULTIPLY", "LIMIT"
// 	Value     *float64 `json:"value,omitempty"`     // Rule value
// }

// // Validation constants and helper functions

// // Valid values for request fields
// var (
// 	ValidPreferences = []string{"fastest", "shortest", "recommended", "custom"}
// 	ValidUnits      = []string{"m", "km", "mi"}
// 	ValidLanguages  = []string{
// 		"cs", "cs-cz", "da", "dk-da", "de", "de-de", "en", "en-us",
// 		"eo", "eo-eo", "es", "es-es", "fi", "fi-fi", "fr", "fr-fr",
// 		"gr", "gr-gr", "he", "he-il", "hu", "hu-hu", "id", "id-id",
// 		"it", "it-it", "ja", "ja-jp", "ne", "ne-np", "nl", "nl-nl",
// 		"nb", "nb-no", "pl", "pl-pl", "pt", "pt-pt", "ro", "ro-ro",
// 		"ru", "ru-ru", "tr", "tr-tr", "ua", "ua-ua", "vi", "vi-vn",
// 		"zh", "zh-cn",
// 	}
// 	ValidInstructionsFormats = []string{"html", "text"}
// 	ValidAttributes = []string{"avgspeed", "detourfactor", "percentage"}
// 	ValidExtraInfo = []string{
// 		"steepness", "suitability", "surface", "waycategory", "waytype",
// 		"tollways", "traildifficulty", "osmid", "roadaccessrestrictions",
// 		"countryinfo", "green", "noise", "csv", "shadow",
// 	}
// 	ValidAvoidFeatures = []string{"highways", "tollways", "ferries", "fords", "steps"}
// 	ValidAvoidBorders  = []string{"all", "controlled", "none"}
// 	ValidVehicleTypes  = []string{"hgv", "bus", "agricultural", "delivery", "forestry", "goods", "unknown"}
// 	ValidSmoothnessTypes = []string{
// 		"excellent", "good", "intermediate", "bad", "very_bad",
// 		"horrible", "very_horrible", "impassable",
// 	}
// 	ValidCustomRuleKeywords  = []string{"IF", "ELSEIF", "ELSE"}
// 	ValidCustomRuleOperations = []string{"MULTIPLY", "LIMIT"}
// )

// // Helper function to create a basic request
// func NewBasicDirectionsRequest(coordinates [][]float64) *ORSDirectionsRequest {
// 	return &ORSDirectionsRequest{
// 		Coordinates: coordinates,
// 	}
// }

// // Helper function to set common defaults
// func (r *ORSDirectionsRequest) SetDefaults() {
// 	if r.Preference == "" {
// 		r.Preference = "recommended"
// 	}
// 	if r.Units == "" {
// 		r.Units = "m"
// 	}
// 	if r.Language == "" {
// 		r.Language = "en"
// 	}
// 	if r.Geometry == nil {
// 		geometry := true
// 		r.Geometry = &geometry
// 	}
// 	if r.Instructions == nil {
// 		instructions := true
// 		r.Instructions = &instructions
// 	}
// 	if r.InstructionsFormat == "" {
// 		r.InstructionsFormat = "text"
// 	}
// }

// =================レスポンス=================
type DirectionsResponse struct {
	// https://openrouteservice.org/dev/#/api-docs/v2/directions/{profile}/geojson/post
	// のレスポンスを構造体に
	Type          string         `json:"type"`
	BBox          []float64      `json:"bbox,omitempty"`
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
	Type       string               `json:"type"`
	Properties ORSFeatureProperties `json:"properties"`
	Geometry   ORSGeometry          `json:"geometry"`
	BBox       []float64            `json:"bbox,omitempty"`
}

// ORSFeatureProperties contains the properties of a route feature
type ORSFeatureProperties struct {
	Segments  []ORSSegment           `json:"segments,omitempty"`
	Summary   ORSSummary             `json:"summary,omitempty"`
	WayPoints []int                  `json:"way_points,omitempty"`
	Extras    map[string]interface{} `json:"extras,omitempty"`
	Warnings  []ORSWarning           `json:"warnings,omitempty"`
}

// ORSGeometry represents the geometry of the route
type ORSGeometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// ORSSegment represents a segment of the route
type ORSSegment struct {
	Distance     float64   `json:"distance"`
	Duration     float64   `json:"duration"`
	Steps        []ORSStep `json:"steps,omitempty"`
	Detourfactor float64   `json:"detourfactor,omitempty"`
	Percentage   float64   `json:"percentage,omitempty"`
	Avgspeed     float64   `json:"avgspeed,omitempty"`
}

// ORSStep represents a step within a segment
type ORSStep struct {
	Distance     float64      `json:"distance"`
	Duration     float64      `json:"duration"`
	Type         int          `json:"type"`
	Instruction  string       `json:"instruction"`
	Name         string       `json:"name,omitempty"`
	WayPoints    []int        `json:"way_points"`
	Mode         int          `json:"mode,omitempty"`
	Maneuver     *ORSManeuver `json:"maneuver,omitempty"`
	ExitBearings []float64    `json:"exit_bearings,omitempty"`
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
	ID            string    `json:"id,omitempty"`
	Attribution   string    `json:"attribution,omitempty"`
	Service       string    `json:"service,omitempty"`
	Timestamp     int64     `json:"timestamp,omitempty"`
	Query         ORSQuery  `json:"query"`
	Engine        ORSEngine `json:"engine,omitempty"`
	SystemMessage string    `json:"system_message,omitempty"`
}

// ORSQuery represents the original query parameters
type ORSQuery struct {
	Coordinates        [][]float64           `json:"coordinates"`
	ID                 string                `json:"id,omitempty"`
	Preference         string                `json:"preference,omitempty"`
	Units              string                `json:"units,omitempty"`
	Language           string                `json:"language,omitempty"`
	Geometry           *bool                 `json:"geometry,omitempty"`
	Instructions       *bool                 `json:"instructions,omitempty"`
	InstructionsFormat string                `json:"instructions_format,omitempty"`
	RoundaboutExits    *bool                 `json:"roundabout_exits,omitempty"`
	Attributes         []string              `json:"attributes,omitempty"`
	Maneuvers          *bool                 `json:"maneuvers,omitempty"`
	Radiuses           []float64             `json:"radiuses,omitempty"`
	Bearings           [][]float64           `json:"bearings,omitempty"`
	ContinueStraight   *bool                 `json:"continue_straight,omitempty"`
	Elevation          *bool                 `json:"elevation,omitempty"`
	ExtraInfo          []string              `json:"extra_info,omitempty"`
	Options            *ORSRouteOptions      `json:"options,omitempty"`
	SuppressWarnings   *bool                 `json:"suppress_warnings,omitempty"`
	GeometrySimplify   *bool                 `json:"geometry_simplify,omitempty"`
	SkipSegments       []int                 `json:"skip_segments,omitempty"`
	AlternativeRoutes  *ORSAlternativeRoutes `json:"alternative_routes,omitempty"`
	MaximumSpeed       *float64              `json:"maximum_speed,omitempty"`
	Schedule           *bool                 `json:"schedule,omitempty"`
	ScheduleDuration   string                `json:"schedule_duration,omitempty"`
	ScheduleRows       *int                  `json:"schedule_rows,omitempty"`
	WalkingTime        string                `json:"walking_time,omitempty"`
	IgnoreTransfers    *bool                 `json:"ignore_transfers,omitempty"`
	CustomModel        *ORSCustomModel       `json:"custom_model,omitempty"`
}

// ORSEngine contains information about the OpenRouteService engine
type ORSEngine struct {
	Version   string `json:"version,omitempty"`
	BuildDate string `json:"build_date,omitempty"`
	GraphDate string `json:"graph_date,omitempty"`
	OSMDate   string `json:"osm_date,omitempty"`
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

// TODO 定義中

// postDirections godoc
// @Title 自転車ルート検索
// @Description 出発地点の座標と目的地の座標をPOSTのbodyで受け取り、OpenRouteServiceのAPIを呼び出してルート情報を取得する
// @Tags map
// @Accept json
// @Produce json
// @Param orsDirectionsRequest body ORSDirectionsRequest true "自転車ルート検索リクエスト"
// @Success 200 {object} DirectionsResponse
// @Failure 400 {object} ORSErrorResponse
// @Failure 404 {object} ORSErrorResponse
// @Failure 500 {object} ORSErrorResponse
// @Router /directions/bicycle [post]
func postDirections(c *gin.Context) {
	// ORSへのリクエスト
	// ルート整形処理

	directionsResponse := DirectionsResponse{}

	c.JSON(http.StatusOK, directionsResponse)
}
