package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// READMEのINDEXリスト
type readmeIndexList struct {
	// 読了日
	readingDate time.Time
	// 本のタイトル
	bookTitle string
	// Markdownへのリンク
	linkMarkdown string
	// 著者
	author string
}

// 日時のフォーマット（Goでは"yyyy/MM/dd HH:mm:ss"ではなく、"2016/01/02 15:04:05"と書く。
const datetimeLayout = "2006/01/02 15:04:05"

// readmeIndexListの文字列を返す。
func (r *readmeIndexList) String() string {
	return fmt.Sprintf("readingDate=[%s], bookTitle=[%s], linkMarkdown=[%s], author=[%s]",
		r.readingDate.Format(datetimeLayout), r.bookTitle, r.linkMarkdown, r.author)
}

// READMEファイルを読み込み、構造体を返す。
func loadReadme(repositoryPath string) ([]readmeIndexList, error) {
	rtnList := make([]readmeIndexList, 0)

	file, err := os.Open(repositoryPath + "\\README.md")
	if err != nil {
		return nil, err
	}

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

		title := string(reg2.Find([]byte(sp[1])))
		rTitle := []rune(title)
		link := string(reg3.Find([]byte(sp[1])))
		rLink := []rune(link)

		readme := readmeIndexList{
			t,
			string(rTitle[1 : len(rTitle)-1]),
			string(rLink[1 : len(rLink)-1]),
			sp[2],
		}
		rtnList = append(rtnList, readme)
	}
	return rtnList, nil
}
