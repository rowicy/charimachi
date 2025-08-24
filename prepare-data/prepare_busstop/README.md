# バス停領域の座標変換

[公共交通機関オープンデータ バス停情報](https://ckan.odpt.org/dataset/b_busstop-toei/resource/f340278d-aefe-47ea-bc8f-15ebe48c286d)
からバス停座標取得　回避エリア群の座標ポリゴン生成

バッチ処理は手動

1. 	バス停座標取得 & それを囲むポリゴン作成

    ```
    go run . -outdir ../../api/data
    ```