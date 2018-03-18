package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/yo-kondo/FmtBookDirStruct/logic"
)

// エントリポイント
func main() {
	var conf logic.Config
	err := logic.GetConfig(&conf)
	if err != nil {
		log.Fatal(err)
		return
	}

	readmeList, err := loadReadme(conf.RepositoryPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, v := range readmeList {
		err2 := v.MoveFile(conf.RepositoryPath)
		if err2 != nil {
			log.Fatal(err2)
			return
		}
	}
	// TODO:READMEのリンクを修正
	// TODO:コピー元のファイルを削除

	// TODO:debug
	for _, v := range readmeList {
		fmt.Println(v.String())
	}
}

// READMEファイルを読み込み、構造体を返す。
func loadReadme(repositoryPath string) ([]logic.ReadmeIndexList, error) {
	rtnList := make([]logic.ReadmeIndexList, 0)

	file, err := os.Open(filepath.Clean(repositoryPath + "/README.md"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	// "1. yyyy/mm/dd - ["の行を取得対象とする
	reg1 := regexp.MustCompile(`1. 20[0-9]{2}/[0-2]{2}/[0-9]{2} - \[`)
	reg2 := regexp.MustCompile(`\[.*]`)
	reg3 := regexp.MustCompile(`\(.*\)`)
	for sc.Scan() {
		line := sc.Text()
		if !reg1.MatchString(line) {
			continue
		}

		// 分割して構造体に格納
		sp := strings.Split(line, " - ")
		runes := []rune(sp[0])
		dateStr := string(runes[3:])
		// yyyy/mm/ddをTime構造体に変換
		t, err := time.Parse("2006/01/02", dateStr)
		if err != nil {
			log.Printf("README内の日付が正しくありません。%s\r\n", line)
		}

		title := []rune(string(reg2.Find([]byte(sp[1]))))
		link := []rune(string(reg3.Find([]byte(sp[1]))))

		readme := logic.ReadmeIndexList{
			ReadingDate:  t,
			BookTitle:    string(title[1 : len(title)-1]),
			LinkMarkdown: string(link[1 : len(link)-1]),
			Author:       sp[2],
		}
		// ISBN取得
		isbn, err := readme.SetIsbn(repositoryPath)
		if err != nil {
			return nil, err
		}
		readme.Isbn = isbn

		rtnList = append(rtnList, readme)
	}

	return rtnList, nil
}
