package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yo-kondo/FmtBookDirStruct/lib"
)

// エントリポイント
func main() {
	// 設定ファイル取得
	var conf lib.Config
	err := lib.GetConfig(&conf)
	if err != nil {
		log.Print("設定ファイルの読み込みに失敗しました。")
		log.Fatal(err)
		return
	}

	// README読み込み
	readmeFile, err := lib.ReadFile(conf.RepositoryPath + "/README.md")
	if err != nil {
		log.Print("READMEファイルの読み込みに失敗しました。")
		log.Fatal(err)
		return
	}

	// "1. yyyy/mm/dd - ["の行を取得対象とする
	reg1 := regexp.MustCompile(`1. 20[0-9]{2}/[0-2]{2}/[0-9]{2} - \[`)
	reg2 := regexp.MustCompile(`\[.*]`)
	reg3 := regexp.MustCompile(`\(.*\)`)

	// READMEを構造体に格納
	readmeList := make([]lib.ReadmeIndexList, 0)
	for _, v := range readmeFile {
		if !reg1.MatchString(v) {
			continue
		}

		// 分割して構造体に格納
		sp := strings.Split(v, " - ")
		runes := []rune(sp[0])
		dateStr := string(runes[3:])
		// yyyy/mm/ddをTime構造体に変換
		t, err := time.Parse("2006/01/02", dateStr)
		if err != nil {
			log.Printf("README内の日付が正しくありません。%s\r\n", v)
			log.Fatal(err)
			return
		}

		title := []rune(string(reg2.Find([]byte(sp[1]))))
		link := []rune(string(reg3.Find([]byte(sp[1]))))

		readme := lib.ReadmeIndexList{
			ReadingDate:     t,
			ReadingYear:     strconv.Itoa(t.Year()),
			Isbn:            "",
			BookTitle:       string(title[1 : len(title)-1]),
			OldLinkMarkdown: string(link[1 : len(link)-1]),
			NewLinkMarkdown: "",
			Author:          sp[2],
		}
		readmeList = append(readmeList, readme)
	}

	// ISBNを設定
	tempReadmeList := make([]lib.ReadmeIndexList, 0)
	for _, v := range readmeList {
		markdownFile, err := lib.ReadFile(conf.RepositoryPath + "/" + v.OldLinkMarkdown)
		if err != nil {
			log.Print("Markdownファイルの読み込みに失敗しました。")
			log.Fatal(err)
			return
		}

		// 構造体のスライスをfor rangeでループした場合、その中では値を変更することはできない。
		// rangeの戻り値は同じ参照先が使用されているため。
		// https://qiita.com/taileagler17/items/008e2b304f27b7fb168a
		// for rangeの中でポインタ先の値を変更したい場合は、一時変数を用意して、一時変数にコピーする。

		tempIsbn := ""
		for _, line := range markdownFile {
			if strings.Contains(line, "ISBN-13") {
				tempIsbn = strings.Split(line, "|")[2]
			}
			if strings.Contains(line, "ISBN-10") {
				tempIsbn = strings.Split(line, "|")[2]
			}
			if strings.Contains(line, "ASIN") {
				tempIsbn = strings.Split(line, "|")[2]
			}
		}

		// 一時変数にコピー
		tempReadmeList = append(tempReadmeList, lib.SetIsbn(v, tempIsbn))
	}
	readmeList = tempReadmeList

	// Markdownファイルをコピー
	for _, v := range readmeList {
		linkSp := strings.Split(v.OldLinkMarkdown, "/")
		linkPath := "/" + strings.Join(linkSp[:len(linkSp)-1], "/")

		oldPath := conf.RepositoryPath + linkPath
		newPath := conf.RepositoryPath + "/md/" + v.ReadingYear + "/" + v.Isbn
		err = lib.CopyDir(oldPath, newPath)
		if err != nil {
			log.Print("Markdownファイルのコピーに失敗しました。")
			log.Fatal(err)
			return
		}
	}

	// READMEのリンクを修正
	tempReadmeList2 := make([]lib.ReadmeIndexList, 0)
	for _, v := range readmeList {
		sp := strings.Split(v.OldLinkMarkdown, "/")
		newLink := sp[0] + "/" + v.ReadingYear + "/" + v.Isbn + "/" + sp[3]

		// 一時変数にコピー
		tempReadmeList2 = append(tempReadmeList2, lib.SetNewLinkMarkdown(v, newLink))
	}
	readmeList = tempReadmeList2

	// TODO:READMEのINDEXだけ書いたファイルを出力
	// TODO:コピー元のファイルを削除

	// TODO:debug
	for _, v := range readmeList {
		fmt.Println(v.String())
	}
}
