package lib

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// ファイルを読み込み、stringのスライスにして返す。
func ReadFile(path string) ([]string, error) {
	cleanPath := filepath.Clean(path)
	lines := make([]string, 0)

	file, err := os.Open(cleanPath)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, nil
}

// ディレクトリ単位でコピーする
// https://www.socketloop.com/tutorials/golang-copy-directory-including-sub-directories-files
func CopyDir(oldPath, newPath string) (err error) {

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
			err = CopyDir(oldFile, newFile)
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
func copyFile(oldFilePath, newFilePath string) (err error) {

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
