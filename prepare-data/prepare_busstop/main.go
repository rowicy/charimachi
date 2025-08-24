package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type BusStop struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	Polygon   [][]float64 `json:"polygon,omitempty"`
}

type ODPTResponse struct {
	ID    string `json:"@id"`
	Type  string `json:"@type"`
	Title struct {
		En     string `json:"en"`
		Ja     string `json:"ja"`
		JaHrkt string `json:"ja-Hrkt"`
	} `json:"title"`
	DCDate    string  `json:"dc:date"`
	GeoLat    float64 `json:"geo:lat"`
	GeoLong   float64 `json:"geo:long"`
	DCTitle   string  `json:"dc:title"`
	ODPTKana  string  `json:"odpt:kana"`
	ODPTNote  string  `json:"odpt:note"`
	OwlSameAs string  `json:"owl:sameAs"`
}

// 緯度経度の度をメートルに変換する係数（東京付近の近似値）
const (
	// 緯度1度 ≈ 111,000m
	latToMeter = 111000.0
	// 経度1度 ≈ 91,000m (東京の緯度35度付近での近似値)
	longToMeter = 91000.0
)

// 指定された中心点から半径10mの四角形（時計回り）を生成
func createPolygon(centerLat, centerLong float64, radiusMeters float64) [][]float64 {
	// メートルを度に変換
	latOffset := radiusMeters / latToMeter
	longOffset := radiusMeters / longToMeter

	// 四角形の4つの角の座標（時計回り）
	polygon := [][]float64{
		{centerLong - longOffset, centerLat + latOffset}, // 左上
		{centerLong + longOffset, centerLat + latOffset}, // 右上
		{centerLong + longOffset, centerLat - latOffset}, // 右下
		{centerLong - longOffset, centerLat - latOffset}, // 左下
		{centerLong - longOffset, centerLat + latOffset}, // 最初の点に戻る（閉じる）
	}

	return polygon
}

func fetchBusStopData(url string) ([]ODPTResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTPエラー: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み取りエラー: %v", err)
	}

	var odpтData []ODPTResponse
	if err := json.Unmarshal(body, &odpтData); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %v", err)
	}

	return odpтData, nil
}

func convertToBusStops(odpтData []ODPTResponse) []BusStop {
	var busStops []BusStop

	for _, data := range odpтData {
		// IDを抽出（owl:sameAsから取得、なければ@idを使用）
		id := data.OwlSameAs
		if id == "" {
			id = data.ID
		}

		// 名前を取得（日本語タイトルを優先）
		name := data.Title.Ja
		if name == "" {
			name = data.DCTitle
		}
		if name == "" {
			name = data.Title.En
		}

		// 座標が有効かチェック
		if data.GeoLat == 0 || data.GeoLong == 0 {
			continue
		}

		// 半径10mの四角形ポリゴンを生成
		polygon := createPolygon(data.GeoLat, data.GeoLong, 10.0)

		busStop := BusStop{
			ID:        id,
			Name:      name,
			Latitude:  data.GeoLat,
			Longitude: data.GeoLong,
			Polygon:   polygon,
		}

		busStops = append(busStops, busStop)
	}

	return busStops
}

func main() {
	outdir := flag.String("outdir", ".", "Output dir")
	flag.Parse()

	url := "https://api-public.odpt.org/api/v4/odpt:BusstopPole?odpt:operator=odpt.Operator:Toei"

	fmt.Println("都営バスのバス停データを取得中...")
	odpтData, err := fetchBusStopData(url)
	if err != nil {
		log.Fatalf("データ取得エラー: %v", err)
	}

	fmt.Printf("取得したバス停数: %d\n", len(odpтData))

	// BusStop構造体に変換
	busStops := convertToBusStops(odpтData)

	fmt.Printf("変換後のバス停数: %d\n", len(busStops))

	// JSONファイルに出力
	outputFile := filepath.Join(*outdir, "bus_stopss.json")
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("ファイル作成エラー: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(busStops); err != nil {
		log.Fatalf("JSON書き込みエラー: %v", err)
	}

	fmt.Printf("バス停データを %s に出力しました\n", outputFile)

	// コンソールにも最初の3件を表示
	fmt.Println("\n=== 最初の3件の詳細（プレビュー） ===")
	for i, busStop := range busStops {
		if i >= 3 {
			break
		}
		fmt.Printf("\nバス停 %d:\n", i+1)
		fmt.Printf("  ID: %s\n", busStop.ID)
		fmt.Printf("  名前: %s\n", busStop.Name)
		fmt.Printf("  緯度: %.6f\n", busStop.Latitude)
		fmt.Printf("  経度: %.6f\n", busStop.Longitude)
		fmt.Printf("  ポリゴン（最初の2点）:\n")
		for j, point := range busStop.Polygon {
			if j >= 2 {
				fmt.Printf("    ...\n")
				break
			}
			fmt.Printf("    [%.7f, %.8f]\n", point[0], point[1])
		}
	}
}
