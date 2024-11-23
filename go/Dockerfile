# ベースイメージ
FROM golang:1.22.4

WORKDIR /usr/src/app

# ソースコードのコピー
COPY . .

# 依存関係のインストール
RUN go mod download

# CompileDaemon のインストール
RUN go install github.com/githubnemo/CompileDaemon@latest

# 開発環境と本番環境のビルド・実行方法を分岐
RUN if [ "$ENV_MODE" = "production" ]; then \
        go build -o main cmd/api/main.go; \
    fi

# CMD の設定
# 開発環境では CompileDaemon で自動再ビルドし、本番環境ではビルド済みのファイルを実行
CMD if [ "$ENV_MODE" = "production" ]; then \
        ./main; \
    else \
        CompileDaemon -build="go build -o main cmd/api/main.go" -command="./main"; \
    fi