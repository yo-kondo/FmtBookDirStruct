package main

import (
	"fmt"
)

// エントリポイント
func main() {
	var conf Config
	getConfig(&conf)
	readmeLines := readReadme(conf.RepositoryPath)

	// debug
	for _, v := range readmeLines {
		fmt.Println(v)
	}
	// TODO:READMEのリンクからMarkdownファイルを探索
	// TODO:MarkdownからISBN-13、読了日を取得
	// TODO:ISBN-13、読了日のディレクトリを作成
	// TODO:ファイルを移動
	// TODO:READMEのリンクを修正
}
