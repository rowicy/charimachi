package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// XXX カスタム構造体
type WarningPoint struct {
	Type       string    `json:"type"`       //交差点なのか直線道路なのか(このフィールドいらない)
	Name       string    `json:"name"`       //名称, 地元での言われ名など(地獄谷etc)
	Coordinate []float64 `json:"coordinate"` //座標
	Message    string    `json:"message"`    //警告のメッセージ(どんな事故が多かったか)
}

var WorningIntersectionPoints []WarningPoint

// GetWarningPoints godoc
// @Summary 注意点取得
// @Description 取締強化交差点などの注意点をまとめて返す。
// @Tags map
// @Accept json
// @Produce json
// @Success 200 {object} []WarningPoint "注意地点情報"
// @Router /warning_point [get]
func GetWarningPoints(c *gin.Context) {
	//ここで結合する
	warningPoints := WorningIntersectionPoints

	c.JSON(http.StatusOK, warningPoints)
}
