package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	util "template-mobile-app-api/util"
)

type WorningIntersectionResponse struct {
	Total    int      `json:"total"`    // 総件数
	Subtotal int      `json:"subtotal"` // サブ合計
	Limit    int      `json:"limit"`    // 1ページあたりの件数
	Offset   int      `json:"offset"`   // オフセット
	Metadata Metadata `json:"metadata"` // メタデータ
	Hits     []Hit    `json:"hits"`     // データ配列
}

// metadata 部分
type Metadata struct {
	ApiId        string `json:"apiId"`
	Title        string `json:"title"`
	DatasetId    string `json:"datasetId"`
	DatasetTitle string `json:"datasetTitle"`
	DatasetDesc  string `json:"datasetDesc"`
	DataTitle    string `json:"dataTitle"`
	DataDesc     string `json:"dataDesc"`
	Sheetname    string `json:"sheetname"`
	Version      string `json:"version"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
}

// hits 配列の各要素（英語フィールド名）
type Hit struct {
	Row          int    `json:"row"`
	Ward         string `json:"行政区"`     // 行政区
	Affiliation  string `json:"所属"`      // 所属
	Location     string `json:"実施場所"`    // 実施場所
	Intersection string `json:"交差点・路線名"` // 交差点・路線名
	Reason       string `json:"取締理由"`    // 取締理由
}

//var worningIntersectionPoints []WarningPoint

// 取締強化交差点注意（オープンデータ）から全データキャッシュ
func main() {
	url := "https://service.api.metro.tokyo.lg.jp/api/t000022d1700000024-29a128f7bb366ba2832927fac7feeaa4-0/json?limit=1000"

	reqBody := map[string]interface{}{
		"sortOrder": nil,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Request creation error:", err)
		return
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	worningIntersectionPoints := []util.WarningPoint{}
	worningIntersectionResponse := WorningIntersectionResponse{}
	json.Unmarshal(body, &worningIntersectionResponse)
	for i, v := range worningIntersectionResponse.Hits {
		//取締り強化交差点データのLocationには「〇〇付近」とあり、検索の邪魔なので消す。
		Location := strings.Replace(v.Location, "付近", "", -1)
		searchResp := util.GetSearchBase(Location)
		fmt.Println(v.Location)
		if searchResponse, ok := searchResp.([]util.SearchResponse); ok {
			worningIntersectionPoints = append(worningIntersectionPoints, util.WarningPoint{})

			if len(searchResponse) <= 0 {
				continue
			}
			lon, _ := strconv.ParseFloat(searchResponse[0].Lon, 64)
			lat, _ := strconv.ParseFloat(searchResponse[0].Lat, 64)
			worningIntersectionPoints[i].Coordinate = []float64{lon, lat}
			worningIntersectionPoints[i].Name = v.Location
			worningIntersectionPoints[i].Message = v.Reason

		}
	}

	defer resp.Body.Close()
}
