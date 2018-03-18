package logic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// READMEのINDEXリスト
type ReadmeIndexList struct {
	// 読了日
	ReadingDate time.Time
	// ISBN13 > ISBN10 > ASIN
	Isbn string
	// 本のタイトル
	BookTitle string
	// Markdownへのリンク
	LinkMarkdown string
	// 著者
	Author string
}

// 日時のフォーマット（Goでは"yyyy/MM/dd HH:mm:ss"ではなく、"2016/01/02 15:04:05"と書く。
const datetimeLayout = "2006/01/02 15:04:05"

// readmeIndexListの文字列を返す。
func (r *ReadmeIndexList) String() string {
	return fmt.Sprintf("ReadingDate=[%s], Isbn=[%s] BookTitle=[%s], LinkMarkdown=[%s], Author=[%s]",
		r.ReadingDate.Format(datetimeLayout), r.Isbn, r.BookTitle, r.LinkMarkdown, r.Author)
}

// リンクからMarkdownを探索し、Markdownの中からISBNを取得する。
// 取得順 ISBN13 > ISBN10 > ASIN
func (r *ReadmeIndexList) SetIsbn(repositoryPath string) (string, error) {

	path := filepath.Clean(repositoryPath + "/" + r.LinkMarkdown)
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

// Markdownファイルをコピーする。
func (r *ReadmeIndexList) MoveFile(repositoryPath string) error {
	// リンクからディレクトリ名を取得
	linkSp := strings.Split(r.LinkMarkdown, "/")
	linkPath := ""
	for _, v := range linkSp[:len(linkSp)-1] {
		linkPath += "/" + v
	}

	oldPath := filepath.Clean(repositoryPath + linkPath)
	newPath := filepath.Clean(repositoryPath + "/md/" + strconv.Itoa(r.ReadingDate.Year()) + "/" + r.Isbn)
	if err := copyDir(oldPath, newPath); err != nil {
		return err
	}
	return nil
}

// ディレクトリ単位でコピーする
// https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
func copyDir(oldPath, newPath string) (err error) {

	// コピー元のディレクトリ情報を取得
	oldFileInfo, err := os.Stat(oldPath)
	if err != nil {
		return err
	}

	// コピー先のディレクトリを作成
	err = os.MkdirAll(newPath, oldFileInfo.Mode())
	if err != nil {
		return err
	}

	// ディレクトリのFile構造体を取得
	dirInfo, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	// ディレクトリ内のファイルを取得
	fileInfos, err := dirInfo.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		oldFile := oldPath + "/" + fileInfo.Name()
		newFile := newPath + "/" + fileInfo.Name()

		if fileInfo.IsDir() {
			// サブディレクトリを再帰で実行
			err = copyDir(oldFile, newFile)
			if err != nil {
				return err
			}
		} else {
			// ファイルのコピー
			err = copyFile(oldFile, newFile)
			if err != nil {
				return err
			}
		}
	}
	return
}

// ファイルをコピーする。
// https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
func copyFile(oldFilePath string, newFilePath string) (err error) {

	oldFile, err := os.Open(oldFilePath)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	newFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// WriterとReaderを渡して、ファイルの内容をコピーする
	_, err = io.Copy(newFile, oldFile)
	if err == nil {
		// パーミッションをコピー元とあわせる。
		oldFileInfo, err := os.Stat(oldFilePath)
		if err != nil {
			err = os.Chmod(newFilePath, oldFileInfo.Mode())
		}
	}
	return
}
