package main

import (
	"fmt"
	"log"

	"github.com/yo-kondo/FmtBookDirStruct/model"
)

// エントリポイント
func main() {
	var conf model.Config
	err := model.GetConfig(&conf)
	if err != nil {
		log.Fatal(err)
		return
	}

	readmeList, err := model.LoadReadme(conf.RepositoryPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// TODO:ISBN-13、読了日のディレクトリを作成
	// TODO:ファイルを移動
	// TODO:READMEのリンクを修正

	// TODO:debug
	for _, v := range readmeList {
		fmt.Println(v.String())
	}
}
