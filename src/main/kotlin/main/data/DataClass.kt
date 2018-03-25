package main.data

import java.time.LocalDate

/**
 * 設定ファイルのクラス
 */
internal data class ConfigData(
        /** リポジトリのディレクトリパス */
        val repositoryPath: String
)

/**
 * READMEのINDEXクラス
 */
internal data class ReadmeIndexData(
        /** 読了日 */
        var readingData: LocalDate,
        /** 読了日（年） */
        var readingYear: String,
        /** ISBN */
        var isbn: String,
        /** 本のタイトル */
        var bookTitle: String,
        /** Markdownへのリンク（旧） */
        var oldLinkMarkdown: String,
        /** Markdownへのリンク（新） */
        var newLinkMarkdown: String,
        /** 著者 */
        var author: String
)
