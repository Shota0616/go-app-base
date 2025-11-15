package config

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"log" // Added for logging
)

var Localizer *i18n.Localizer

func InitI18n() {
	// i18n バンドルを作成
	bundle := i18n.NewBundle(language.English)

	// JSON のアンマーシャル関数を登録
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	appRoot := os.Getenv("APP_ROOT")
	if appRoot == "" {
		appRoot = "/usr/src" // Default to /usr/src if APP_ROOT is not set, assuming locales is at /usr/src/locales
	}

	// ローカルファイルをロード
	enPath := filepath.Join(appRoot, "locales", "en.json")
	jaPath := filepath.Join(appRoot, "locales", "ja.json")

	log.Printf("Loading en.json from: %s", enPath)
	if _, err := bundle.LoadMessageFile(enPath); err != nil {
		panic("Failed to load en.json: " + err.Error())
	}
	log.Printf("Loading ja.json from: %s", jaPath)
	if _, err := bundle.LoadMessageFile(jaPath); err != nil {
		panic("Failed to load ja.json: " + err.Error())
	}

	// 環境変数から言語設定を取得
	lang := os.Getenv("APP_LANG")
	if lang == "" {
		lang = "en" // デフォルトの言語を英語に設定
	}

	// ローカライザーを初期化
	Localizer = i18n.NewLocalizer(bundle, lang)
	log.Printf("i18n Localizer initialized for language: %s", lang)
}
