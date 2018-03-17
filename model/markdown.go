package model

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
type ReadmeIndexList struct {
	// 読了日
	readingDate time.Time
	// ISBN13 > ISBN10 > ASIN
	isbn string
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
func (r *ReadmeIndexList) String() string {
	return fmt.Sprintf("readingDate=[%s], isbn=[%s] bookTitle=[%s], linkMarkdown=[%s], author=[%s]",
		r.readingDate.Format(datetimeLayout), r.isbn, r.bookTitle, r.linkMarkdown, r.author)
}

// READMEファイルを読み込み、構造体を返す。
func LoadReadme(repositoryPath string) ([]ReadmeIndexList, error) {
	rtnList := make([]ReadmeIndexList, 0)

	file, err := os.Open(repositoryPath + "\\README.md")
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

		title := string(reg2.Find([]byte(sp[1])))
		rTitle := []rune(title)
		link := string(reg3.Find([]byte(sp[1])))
		rLink := []rune(link)

		readme := ReadmeIndexList{
			t,
			"",
			string(rTitle[1 : len(rTitle)-1]),
			string(rLink[1 : len(rLink)-1]),
			sp[2],
		}
		// ISBN取得
		isbn, err := setIsbn(readme, repositoryPath)
		if err != nil {
			return nil, err
		}
		readme.isbn = isbn

		rtnList = append(rtnList, readme)
	}

	return rtnList, nil
}

// リンクからMarkdownを探索し、Markdownの中からISBNを取得する。
// 取得順 ISBN13 > ISBN10 > ASIN
func setIsbn(r ReadmeIndexList, repositoryPath string) (string, error) {

	path := repositoryPath + "\\" + strings.Replace(r.linkMarkdown, "/", "\\", -1)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	isbn13 := ""
	isbn10 := ""
	asin := ""
	for sc.Scan() {
		line := sc.Text()

		if strings.Contains(line, "ISBN-13") {
			sp := strings.Split(line, "|")
			isbn13 = sp[2]
		}
		if strings.Contains(line, "ISBN-10") {
			isbn10 = strings.Split(line, "|")[2]
		}
		if strings.Contains(line, "ASIN") {
			asin = strings.Split(line, "|")[2]
		}
	}

	if isbn13 != "" && isbn13 != "－" {
		return isbn13, nil
	}
	if isbn10 != "" && isbn10 != "－" {
		return isbn10, nil
	}
	return asin, nil
}
