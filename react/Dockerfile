# ベースイメージ
FROM node:18

WORKDIR /var/www/vite-project

# プロジェクトファイルをコピー
COPY ./project/ /var/www/vite-project/
COPY ./startup.sh /startup.sh

# 依存パッケージをインストール
RUN npm install
RUN chmod 755 /startup.sh

# 起動スクリプトの実行, 本番環境ではビルド後の静的ファイルを提供する
CMD ["/startup.sh"]