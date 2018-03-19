package lib

import (
	"github.com/BurntSushi/toml"
)

// configファイルの構造体
// 外部パッケージからアクセスするため、名前を大文字にして公開にする必要がある。
type Config struct {
	// リポジトリのディレクトリパス
	RepositoryPath string
}

// 設定ファイルを取得する。
func GetConfig(conf *Config) error {
	_, err := toml.DecodeFile("config.toml", &conf)
	return err
}
