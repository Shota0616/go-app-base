#!/bin/bash

# 環境モードに基づく処理分岐
if [ "$NODE_ENV" = "production" ]; then
    # 本番環境: ビルドして、Nginx などのサーバで配信するための静的ファイルを生成
    echo "Running production build..."
    npm run build
    # ビルドされたファイルはvolumeでマウントされたディレクトリに出力されるため、別途コピーする必要はない
else
    # 開発環境: 開発サーバを起動して変更をリアルタイムに反映
    echo "Starting development server..."
    npm run dev
fi