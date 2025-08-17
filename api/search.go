package main

import (
	"net/http"
	_ "template-mobile-app-api/docs"

	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	// https://nominatim.openstreetmap.org/search?q={Client input}&format=json&limit=5
	// のレスポンスを構造体に
}

//TODO 定義中

// getSearch godoc
// @Title 目的地候補検索
// @Description Nominatimのレスポンスを横流し(https://nominatim.openstreetmap.org/search?q={Client input}&format=json&limit=5)
// @Tags map
// @Accept json
// @Produce json
// @Param q query string true "検索文字列"
// @Success 200 {object} SearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /search [get]
func getSearch(c *gin.Context) {
	//Client inputを取得 パラメータ: q
	//https://nominatim.openstreetmap.org/search?q={Client input}&format=json&limit=5
	//nominatimレスポンスをそのまま返す
	var searchResponse []SearchResponse

	c.JSON(http.StatusOK, searchResponse)
}
