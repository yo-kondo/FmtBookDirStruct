package main

import (
	"fmt"
	"log"
)

// エントリポイント
func main() {
	var conf Config
	err := getConfig(&conf)
	if err != nil {
		log.Fatal(err)
		return
	}

	readmeList, err := loadReadme(conf.RepositoryPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// TODO:debug
	for _, v := range readmeList {
		fmt.Println(v.String())
	}

	// TODO:READMEのリンクからMarkdownファイルを探索
	// TODO:MarkdownからISBN-13、読了日を取得
	// TODO:ISBN-13、読了日のディレクトリを作成
	// TODO:ファイルを移動
	// TODO:READMEのリンクを修正
}
