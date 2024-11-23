package config

import (
    "encoding/json"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
    "os"
)

var Localizer *i18n.Localizer

func InitI18n() {
    // i18n バンドルを作成
    bundle := i18n.NewBundle(language.English)

    // JSON のアンマーシャル関数を登録
    bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

    // ローカルファイルをロード
    if _, err := bundle.LoadMessageFile("../locales/en.json"); err != nil {
        panic("Failed to load en.json: " + err.Error())
    }
    if _, err := bundle.LoadMessageFile("../locales/ja.json"); err != nil {
        panic("Failed to load ja.json: " + err.Error())
    }

    // 環境変数から言語設定を取得
    lang := os.Getenv("APP_LANG")
    if lang == "" {
        lang = "en" // デフォルトの言語を英語に設定
    }

    // ローカライザーを初期化
    Localizer = i18n.NewLocalizer(bundle, lang)
}
