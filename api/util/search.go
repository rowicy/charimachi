package util

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	_ "template-mobile-app-api/docs"

	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	// https://nominatim.openstreetmap.org/search?q={Client input}&format=json&limit=5
	// のレスポンスを構造体に
	PlaceID     int      `json:"place_id"`     // Nominatim 内での一意なID
	Licence     string   `json:"licence"`      // データのライセンス情報
	OsmType     string   `json:"osm_type"`     // OSM 要素の種類 (node, way, relation)
	OsmID       int64    `json:"osm_id"`       // OSM 内の要素ID
	Lat         string   `json:"lat"`          // 緯度
	Lon         string   `json:"lon"`          // 経度
	Class       string   `json:"class"`        // 分類（例: railway, building, amenity など）
	Type        string   `json:"type"`         // 分類のサブタイプ（例: stop, station, bus_station など）
	PlaceRank   int      `json:"place_rank"`   // 検索結果の粒度を示すランク
	Importance  float64  `json:"importance"`   // 検索結果の重要度スコア
	AddressType string   `json:"addresstype"`  // アドレスの種類（例: city, house, railway）
	Name        string   `json:"name"`         // OSM に登録されている名前
	DisplayName string   `json:"display_name"` // 住所や施設名などの連結
	BoundingBox []string `json:"boundingbox"`  // 範囲 [南緯, 北緯, 西経, 東経]
}

//TODO 定義中

// getSearch godoc
// @Summary 目的地候補検索
// @Description Nominatimのレスポンスを横流し(https://nominatim.openstreetmap.org/search?q={Client input}&format=json&limit=5)
// @Tags map
// @Accept json
// @Produce json
// @Param q query string true "検索文字列"
// @Success 200 {object} []SearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /search [get]
func GetSearch(c *gin.Context) {
	//Client inputを取得 パラメータ: q
	//https://nominatim.openstreetmap.org/search?q={Client input}&accept-language=ja&format=json&limit=5
	//nominatimレスポンスをそのまま返す
	query := c.Query("q")

	resp := GetSearchBase(query)

	if _, ok := resp.([]SearchResponse); ok {
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusInternalServerError, resp)
	}
}

func GetSearchBase(query string) (res any) {
	//Client inputを取得 パラメータ: q
	//https://nominatim.openstreetmap.org/search?q={Client input}&format=json&limit=5
	//nominatimレスポンスをそのまま返す
	encodedQuery := url.QueryEscape(query)
	resp, err := http.Get("https://nominatim.openstreetmap.org/search?q=" + encodedQuery + "&accept-language=ja&format=json&limit=5")

	if err != nil {
		response := ErrorResponse{
			Error:   "Failed to fetch data",
			Message: err.Error(),
		}
		return response
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response := ErrorResponse{
			Error:   "Failed to read response",
			Message: err.Error(),
		}
		return response
	}

	var searchResponse []SearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		response := ErrorResponse{
			Error:   "Failed to parse response",
			Message: err.Error(),
		}

		return response
	}

	return searchResponse
}
