package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// 引数
	depth := flag.Int("r", 1, "search depth")
	outdirPath := flag.String("outdir", ".", "output dir path")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: go run . <dir_path> -r <depth> -o <output_dir_path>")
		os.Exit(1)
	}
	dirPath := flag.Arg(0)

	// ディレクトリの存在確認
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' does not exist\n", dirPath)
		os.Exit(1)
	}

	// 正規表現
	re := regexp.MustCompile(`^ ?[0-9]+,[^,]+,[^,]+Ｘ[^,]+,[^,]+,[^,]+,[^,]+,全方向計,12時間計,[0-9]+,[0-9]+,[0-9\-]+`)

	// 出力ファイル
	name := filepath.Base(dirPath)
	filePath := filepath.Join(*outdirPath, name+"_filtered.csv")
	// fmt.Printf("Output dir: %s\n", *outdirPath)
	// fmt.Printf("Output file: %s\n", filePath)
	outFile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error: Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// 処理統計
	var filesProcessed, totalMatches int

	// ディレクトリ探索
	err = walkDir(dirPath, *depth, func(path string) {
		if strings.Contains(filepath.Base(path), "jitensya") && strings.HasSuffix(path, ".csv") {
			fmt.Printf("Processing file: %s\n", path)
			matches := processFile(path, re, writer)
			filesProcessed++
			totalMatches += matches
			fmt.Printf("  -> Found %d matching lines\n", matches)
		}
	})

	if err != nil {
		fmt.Printf("Error during directory walk: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("Files processed: %d\n", filesProcessed)
	fmt.Printf("Total matches found: %d\n", totalMatches)
	fmt.Printf("Output written to: filtered.csv\n")
}

// ディレクトリを深さ付きで探索
func walkDir(root string, maxDepth int, fn func(path string)) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Warning: Error accessing %s: %v\n", path, err)
			return nil // エラーをスキップして続行
		}

		// 深さ制御の修正
		if d.IsDir() && path != root {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			depth := strings.Count(rel, string(os.PathSeparator)) + 1
			if depth > maxDepth {
				return filepath.SkipDir
			}
		}

		if !d.IsDir() {
			fn(path)
		}
		return nil
	})
}

// CSVファイルを読み、正規表現に一致する行を書き出す
func processFile(path string, re *regexp.Regexp, writer *csv.Writer) int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("  Error: Failed to open %s: %v\n", path, err)
		return 0
	}
	defer file.Close()

	fmt.Printf("  Successfully opened: %s\n", path)

	scanner := bufio.NewScanner(file)
	matchCount := 0
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if re.MatchString(line) {
			matchCount++
			fmt.Printf("    Hit on line %d: %.100s...\n", lineNumber, line)

			// CSVパーサーを使って正しく分割
			reader := csv.NewReader(strings.NewReader(line))
			records, err := reader.Read()
			if err != nil {
				// CSVパースに失敗した場合は単純分割にフォールバック
				fmt.Printf("    Warning: CSV parse failed for line %d, using simple split\n", lineNumber)
				records = strings.Split(line, ",")
			}

			err = writer.Write(records)
			if err != nil {
				fmt.Printf("    Error: Failed to write line %d: %v\n", lineNumber, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("  Error: Failed to scan %s: %v\n", path, err)
	}

	return matchCount
}
