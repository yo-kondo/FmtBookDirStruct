package lib

import (
	"fmt"
	"time"
)

// READMEのINDEXリスト
type ReadmeIndexList struct {
	// 読了日
	ReadingDate time.Time
	// 読了日（年）
	ReadingYear string
	// ISBN13 > ISBN10 > ASIN
	Isbn string
	// 本のタイトル
	BookTitle string
	// Markdownへのリンク（旧）
	OldLinkMarkdown string
	// Markdownへのリンク（新）
	NewLinkMarkdown string
	// 著者
	Author string
}

// 日時のフォーマット（Goでは"yyyy/MM/dd HH:mm:ss"ではなく、"2016/01/02 15:04:05"と書く。
const datetimeLayout = "2006/01/02 15:04:05"

// readmeIndexListの文字列を返す。
func (r *ReadmeIndexList) String() string {
	return fmt.Sprintf(
		"ReadingDate=[%s], ReadingYear=[%s], Isbn=[%s], BookTitle=[%s], OldLinkMarkdown=[%s], NewLinkMarkdown=[%s], Author=[%s]",
		r.ReadingDate.Format(datetimeLayout), r.ReadingYear, r.Isbn, r.BookTitle, r.OldLinkMarkdown, r.NewLinkMarkdown, r.Author)
}

// ReadmeIndexListにISBNをセットする。
func SetIsbn(r ReadmeIndexList, isbn string) ReadmeIndexList {
	return ReadmeIndexList{
		ReadingDate:     r.ReadingDate,
		ReadingYear:     r.ReadingYear,
		Isbn:            isbn,
		BookTitle:       r.BookTitle,
		OldLinkMarkdown: r.OldLinkMarkdown,
		NewLinkMarkdown: r.NewLinkMarkdown,
		Author:          r.Author,
	}
}

// ReadmeIndexListにMarkdownへのリンク（新）をセットする。
func SetNewLinkMarkdown(r ReadmeIndexList, newLink string) ReadmeIndexList {
	return ReadmeIndexList{
		ReadingDate:     r.ReadingDate,
		ReadingYear:     r.ReadingYear,
		Isbn:            r.Isbn,
		BookTitle:       r.BookTitle,
		OldLinkMarkdown: r.OldLinkMarkdown,
		NewLinkMarkdown: newLink,
		Author:          r.Author,
	}
}
