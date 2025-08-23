package util

// XXX カスタム構造体
type WarningPoint struct {
	Type       string    `json:"type"`       //交差点なのか直線道路なのか(このフィールドいらない)
	Name       string    `json:"name"`       //名称, 地元での言われ名など(地獄谷etc)
	Coordinate []float64 `json:"coordinate"` //座標
	Message    string    `json:"message"`    //警告のメッセージ(どんな事故が多かったか)
}
