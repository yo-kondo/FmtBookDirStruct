package main

import com.google.gson.Gson
import main.data.ConfigData
import main.data.ReadmeIndexData
import java.io.File
import java.nio.charset.StandardCharsets

/**
 * エントリポイント
 * @param args コマンドライン引数
 */
fun main(args: Array<String>) {
    // 設定ファイルを読み込む
    // https://qiita.com/devneko/items/93ee1212ce189f910891
    val source = File("config.json").readText(StandardCharsets.UTF_8)
    val conf = Gson().fromJson(source, ConfigData::class.java)

    // README読み込み
    val readmeFile: MutableList<String> = mutableListOf()
    // Fileクラスの拡張関数forEachLineで行単位で取得
    File(conf.repositoryPath + "/README.md")
            .forEachLine(StandardCharsets.UTF_8) {
                readmeFile.add(it)
            }

    // "1. yyyy/mm/dd - ["の行を取得対象とする
    val reg1 = Regex("""1. 20[0-9]{2}/[0-2]{2}/[0-9]{2} - \[""")

    val readmeList: MutableList<ReadmeIndexData> = mutableListOf()
    // READMEをReadmeIndexに変換
    for (line in readmeFile) {
        if (!reg1.containsMatchIn(line)) {
            continue
        }

        // 分割してReadmeIndexに格納
        val sp = line.split(" - ")

        // yyyy/mm/ddをDateに変換
        val date = sp[0].substring(3).toDate()
        if (date == null) {
            println("README内の日付が不正です。line = $line")
            return
        }

        // タイトルとリンクを分割
        val result1 = Regex("""\[.*]""").find(sp[1])
        // substringで[]を除去
        val title = result1?.let {
            it.value.substring(1, it.value.length - 1)
        } ?: ""

        val result2 = Regex("""\(.*\)""").find(sp[1])
        // substringで()を除去
        val link = result2?.let {
            it.value.substring(1, it.value.length - 1)
        } ?: ""

        readmeList.add(
                ReadmeIndexData(
                        readingData = date,
                        readingYear = date.year.toString(),
                        isbn = getIsbn(conf, link),
                        bookTitle = title,
                        oldLinkMarkdown = link,
                        newLinkMarkdown = "",
                        author = sp[2]))
    }

    // TODO:Markdownファイルをコピー
    // TODO:READMEのリンクを修正
    // TODO:READMEのINDEXだけ書いたファイルを出力
    // TODO:コピー元のファイルを削除

    // TODO:debug
    println(conf)
    println(readmeList.joinToString("\n"))
}

/**
 * ISBNを取得します。
 * @param conf 設定ファイル
 * @param markdownLink ISBN取得対象のMarkdownへのリンク
 * @return ISBN。優先度は、ISBN13 > ISBN10 > ASIN
 */
fun getIsbn(conf: ConfigData, markdownLink: String): String {

    val markdownFile: MutableList<String> = mutableListOf()
    File(conf.repositoryPath + "/" + markdownLink)
            .forEachLine(StandardCharsets.UTF_8) {
                markdownFile.add(it)
            }

    var isbn13 = ""
    var isbn10 = ""
    var asin = ""
    for (line in markdownFile) {
        if (line.contains("ISBN-13")) {
            isbn13 = line.split("|")[2]
        }
        if (line.contains("ISBN-10")) {
            isbn10 = line.split("|")[2]
        }
        if (line.contains("ASIN")) {
            asin = line.split("|")[2]
        }
    }

    if (!isbn13.isEmpty() || isbn13 == "－") {
        return isbn13
    }
    if (!isbn10.isEmpty() || isbn13 == "－") {
        return isbn10
    }
    return asin
}
