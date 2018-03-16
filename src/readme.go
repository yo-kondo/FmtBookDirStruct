package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// READMEのINDEXリスト
type readmeIndexList struct {
	// 読了日
	readingDate string
	// 本のタイトル
	bookTitle string
	// Markdownへのリンク
	linkMarkdown string
	// 著者
	author string
}

// readmeIndexListの文字列を返す。
func (r *readmeIndexList) String() string {
	return fmt.Sprintf("readingDate=[%s], bookTitle=[%s], linkMarkdown=[%s], author[%s]",
		r.readingDate, r.bookTitle, r.linkMarkdown, r.author)
}

// READMEファイルを読み込む。ファイル内の本mdのリンクのみを返す。
func readReadme(repositoryPath string) []readmeIndexList {
	lines := make([]string, 0)
	file, err := os.Open(repositoryPath + "\\README.md")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	scanner := bufio.NewScanner(file)

	ignore := true
	for scanner.Scan() {
		// ## INDEXまでは読み捨てる
		line := scanner.Text()
		if line != "## INDEX" && ignore {
			continue
		}
		ignore = false
		// 空行
		if line == "" {
			continue
		}
		// 以降は読み込まない
		if line == "## INDEXのルール" {
			break
		}
		lines = append(lines, line)
	}

	return lines
}
