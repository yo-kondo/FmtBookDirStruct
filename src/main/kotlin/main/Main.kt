package main

import com.google.gson.Gson
import main.data.ConfigData
import main.data.ReadmeIndexData
import java.io.File

/**
 * エントリポイント
 * @param args コマンドライン引数
 */
fun main(args: Array<String>) {
    val conf = readConfig()
    val readmeList = readReadme(conf.repositoryPath)
    copyMarkdown(conf.repositoryPath, readmeList)
    addNewLinkMarkdown(readmeList)
    writeNewReadmeIndex(conf.repositoryPath, readmeList)
    deleteOldMarkdown(conf.repositoryPath, readmeList)
    println("完了")
}

/**
 * 設定ファイルを読み込んで返します。
 * @return 設定ファイルのデータクラス
 */
private fun readConfig(): ConfigData {
    // https://qiita.com/devneko/items/93ee1212ce189f910891
    val source = File("config.json").readText(Charsets.UTF_8)
    return Gson().fromJson(source, ConfigData::class.java)
}

/**
 * READMEファイルを読み込んでReadmeIndexDataのリストを返します。
 * @param repPath リポジトリのディレクトリパス
 * @return ReadmeIndexDataのリスト
 */
private fun readReadme(repPath: String): MutableList<ReadmeIndexData> {
    val readmeFile: MutableList<String> = mutableListOf()

    // Fileクラスの拡張関数forEachLineで行単位で取得
    File("$repPath/README.md")
            .forEachLine(Charsets.UTF_8) {
                readmeFile.add(it)
            }

    // "1. yyyy/mm/dd - ["の行を取得対象とする
    val reg1 = Regex("""1. 20[0-9]{2}/[0-9]{2}/[0-9]{2} - \[""")

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
            // 空のリストを返す
            return mutableListOf()
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
                        isbn = getIsbn(repPath, link),
                        bookTitle = title,
                        oldLinkMarkdown = link,
                        newLinkMarkdown = "",
                        author = sp[2]))
    }
    return readmeList
}

/**
 * ISBNを取得します。
 * @param repPath 設定ファイル
 * @param markdownLink ISBN取得対象のMarkdownへのリンク
 * @return ISBN。優先度は、ISBN13 > ISBN10 > ASIN
 */
private fun getIsbn(repPath: String, markdownLink: String): String {

    val markdownFile: MutableList<String> = mutableListOf()
    File("$repPath/$markdownLink")
            .forEachLine(Charsets.UTF_8) {
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

    if (!isbn13.isEmpty() && isbn13 != "－") {
        return isbn13
    }
    if (!isbn10.isEmpty() && isbn13 != "－") {
        return isbn10
    }
    return asin
}

/**
 * Markdownファイルをコピーします。
 * @param repPath リポジトリのディレクトリパス
 * @param readmeList ReadmeIndexDataのリスト
 */
private fun copyMarkdown(repPath: String, readmeList: MutableList<ReadmeIndexData>) {
    for (m in readmeList) {
        val linkSp = m.oldLinkMarkdown.split("/")
        val sourcePath =
                "$repPath/" + linkSp
                        .take(linkSp.size - 1)
                        .joinToString("/")
        val targetPath = "$repPath/md/${m.readingYear}/${m.isbn}"
        File(sourcePath).copyRecursively(File(targetPath), true)
    }
}

/**
 * ReadmeIndexDataのnewLinkMarkdownを追加します。
 * @param readmeList ReadmeIndexDataのリスト
 */
private fun addNewLinkMarkdown(readmeList: MutableList<ReadmeIndexData>) {
    for (m in readmeList) {
        val fileName = m.oldLinkMarkdown.split("/").last()
        m.newLinkMarkdown = "md/${m.readingYear}/${m.isbn}/$fileName"
    }
}

/**
 * 新しいリンクに変更したREADMEファイルを作成します。
 * READMEファイルのINDEXのみ記載したファイルを作成します。
 * @param repPath リポジトリのディレクトリパス
 * @return ReadmeIndexDataのリスト
 */
private fun writeNewReadmeIndex(repPath: String, readmeList: MutableList<ReadmeIndexData>) {
    val inFile = File("$repPath/newREADME_IndexOnly.md")

    // 先にファイルを消す
    inFile.delete()

    for (m in readmeList) {
        val fileName = m.oldLinkMarkdown.split("/").last()
        inFile.appendText("1. ${m.readingData.toStringEx()} - " +
                "[${m.bookTitle}](md/${m.readingData.year}/${m.isbn}/$fileName) - " +
                "${m.author}\r\n", Charsets.UTF_8)
    }
}

/**
 * 古いMarkdownファイルを削除します。
 * @param repPath リポジトリのディレクトリパス
 * @return ReadmeIndexDataのリスト
 */
private fun deleteOldMarkdown(repPath: String, readmeList: MutableList<ReadmeIndexData>) {
    for (m in readmeList) {
        val sp = m.oldLinkMarkdown.split("/")
        val delFile = File("$repPath/${sp[0]}/${sp[1]}")
        delFile.deleteRecursively()
    }
}
