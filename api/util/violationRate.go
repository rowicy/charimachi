package util

import "github.com/gin-gonic/gin"

type ViolationRate struct {
	Type           string    `json:"type,omitempty"`
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
	// Sample data - replace with actual logic to fetch violation rates
	violationRates := []ViolationRate{
		{
			Type:           "intersection",
			ViolationRate:  0.1,
			ViolationCount: 5,
			Coordinate:     []float64{139.701961, 35.689727},
			Message:        "違反が多い箇所です",
		},
		{
			Type:           "intersection",
			ViolationRate:  0.2,
			ViolationCount: 8,
			Coordinate:     []float64{139.760530, 35.707620},
			Message:        "違反が多い箇所です",
		},
	}

	c.JSON(200, gin.H{
		"violation_rates": violationRates,
	})
}
