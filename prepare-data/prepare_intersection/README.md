# 交差点座標データ作成

[交通量統計表](https://catalog.data.metro.tokyo.lg.jp/dataset/t000022d0000000035)から交差点座標(json)を推測

第3カラム "外堀通りＸ第一京浜" から交差座標を割り出し必要データに整形

バッチ処理は手動

1. データ配置

    `../open-data/intersection` にCSVオープンデータを配置(utf8保存を確認) ([utf8変換](https://github.com/riiim400th/shitfjis2utf8))

2. 行抽出

    ```bash
    go run ./extract -r 3 -outdir ./predata ../../open-data/intersection/even_year/02_kousatenkubu_csv
    go run ./extract -r 3 -outdir ./predata ../../open-data/intersection/even_year/02_kousatentamabu_csv
    ```

3. 抽出csvを元に座標変換

    ```bash
    go run ./get_coord -files ./predata/*_filtered.csv -outdir ../../api/data
    ```