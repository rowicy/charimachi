package util

import (
	"math"
	"math/rand"
	"time"

	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
)

type ViolationRate struct {
	Type           string    `json:"type,omitempty"`
	Name           string    `json:"name,omitempty"`
	ViolationRate  float64   `json:"violation_rate"`            //違反率スコア
	ViolationCount int       `json:"violation_count,omitempty"` //違反件数
	Coordinate     []float64 `json:"coordinate"`
	Message        string    `json:"message,omitempty"`
}

// GetViolationRates godoc
// @Summary 違反箇所の交差点
// @Description /directions/bicycleでの経路検索結果に対する違反箇所の交差点を返す
// @Param session_id query string false "/directions/bicycleのレスポンス内のsession_id" example(session_id=53238c22-12ac-8ff8-fcdd-738d30a780df)
// @Tags map
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]ViolationRate "List of violation rates for given coordinates"
// @Router /violation_rates [get]
func GetViolationRates(c *gin.Context) {
	session_id := c.Query("session_id")
	if session_id == "" {
		c.JSON(200, gin.H{
			"violation_rates": nil,
		})
		return
	}
	// // Sample data - replace with actual logic to fetch violation rates
	// violationRates := []ViolationRate{
	// 	{
	// 		Type:           "intersection",
	// 		ViolationRate:  0.1,
	// 		ViolationCount: 5,
	// 		Coordinate:     []float64{139.701961, 35.689727},
	// 		Message:        "違反が多い箇所です",
	// 	},
	// 	{
	// 		Type:           "intersection",
	// 		ViolationRate:  0.2,
	// 		ViolationCount: 8,
	// 		Coordinate:     []float64{139.760530, 35.707620},
	// 		Message:        "違反が多い箇所です",
	// 	},
	// }

	featureORSGeometry := SessionIDResponse[session_id]
	filteredRates := FilterViolationRates(featureORSGeometry, violationRates)

	c.JSON(200, gin.H{
		"violation_rates": filteredRates,
	})
}

var warningMessages = []string{
	"歩道ではなく車道を走行しましょう。",
	"歩道走行は禁止されています。安全な車道を走りましょう。",
	"歩行者を守るため、歩道では走らないでください。",
	"自転車は原則、車道を通行します。",
	"歩道では降りて押して歩きましょう。",
}

// 東京付近の1mを度に換算（経度方向のみ簡易）
const metersPerDegree = 111319.0 // 赤道での1° ≈ 111.3 km

// 距離の閾値（20mを度換算）
var threshold float64 = 20.0 / metersPerDegree

// 距離（度ベース）を計算（単純なユークリッド距離）
func distance(coord1, coord2 []float64) float64 {
	dLon := coord1[0] - coord2[0]
	dLat := coord1[1] - coord2[1]
	return math.Sqrt(dLon*dLon + dLat*dLat)
}

func FilterViolationRates(geometry ORSGeometry, violations []ViolationRate) []ViolationRate {
	var result []ViolationRate
	for _, v := range violations {
		for _, c := range geometry.Coordinates {
			if distance(v.Coordinate, c) < threshold {
				var title string
				if v.ViolationRate < 0.3 { // 低リスク
					title = "注意 交差点"
				} else if v.ViolationRate < 0.6 { // 中リスク
					title = "警告 交差点"
				} else { // 高リスク
					title = "違反多発 交差点"
				}
				rand.Seed(time.Now().UnixNano())
				violationRate := ViolationRate{
					Type:           "intersection",
					Name:           title,
					ViolationRate:  math.Floor(v.ViolationRate*100) / 100, // 小数点以下2桁に丸める
					ViolationCount: v.ViolationCount,
					Coordinate:     c,
					Message:        warningMessages[rand.Intn(len(warningMessages))]}

				result = append(result, violationRate)
				break
			}
		}
	}
	return result
}

func LoadViolationRates(filePath string) ([]ViolationRate, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var violations []ViolationRate
	if err := json.Unmarshal(data, &violations); err != nil {
		panic(err)
	}
	return violations, nil
}
