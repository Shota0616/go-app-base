# デプロイメントガイド

## 環境の種類

このアプリケーションは2つの環境モードをサポートしています：

- **development**: 開発環境（Mailpit使用、詳細ログ）
- **production**: 本番環境（実際のSMTP、最小ログ）

## 開発環境のセットアップ

### 1. 環境変数の設定

`.env`ファイルを作成：

```bash
cp .env.example .env
```

開発環境用の設定例：

```env
ENV_MODE=development
TZ=Asia/Tokyo
APP_LANG=en

# データベース
MYSQL_ROOT_PASSWORD=password
MYSQL_DATABASE=app
MYSQL_USER=user
MYSQL_PASSWORD=password
MYSQL_HOST=mysql
MYSQL_PORT=3306

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT（openssl rand -base64 32で生成）
JWT_SECRET=your_jwt_secret_here
JWT_REFRESH_SECRET=your_jwt_refresh_secret_here

# アプリケーション
APP_URL=http://localhost:8000

# メール（開発環境ではMailpit使用）
EMAIL_FROM=noreply@example.com
SMTP_HOST=mailpit
SMTP_PORT=1025

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# ログ
LOG_LEVEL=debug
```

### 2. コンテナの起動

```bash
docker-compose up -d
```

### 3. モックデータの投入（オプション）

```bash
make seed-docker
```

### 4. アクセス

- アプリケーション: http://localhost:8000
- Mailpit（メール確認）: http://localhost:8025

## 本番環境のデプロイ

### 1. 環境変数の設定

本番環境用の`.env`ファイルを作成：

```env
ENV_MODE=production
TZ=Asia/Tokyo
APP_LANG=en

# データベース（強力なパスワードを使用）
MYSQL_ROOT_PASSWORD=your_very_secure_root_password
MYSQL_DATABASE=app
MYSQL_USER=appuser
MYSQL_PASSWORD=your_very_secure_password
MYSQL_HOST=mysql
MYSQL_PORT=3306

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT（必ず新しいシークレットを生成）
JWT_SECRET=$(openssl rand -base64 32)
JWT_REFRESH_SECRET=$(openssl rand -base64 32)

# アプリケーション（実際のドメインを設定）
APP_URL=https://yourdomain.com

# メール（実際のSMTPサーバーを使用）
EMAIL_FROM=noreply@yourdomain.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# CORS（本番環境のドメインを設定）
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# ログ
LOG_LEVEL=warn
```

### 2. JWT シークレットの生成

```bash
# JWT_SECRET
openssl rand -base64 32

# JWT_REFRESH_SECRET
openssl rand -base64 32
```

生成された値を`.env`ファイルに設定してください。

### 3. SSL証明書の設定（HTTPS用）

Let's Encryptを使用する場合：

```bash
# Certbotのインストール
sudo apt-get update
sudo apt-get install certbot

# 証明書の取得
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# 証明書をnginxディレクトリにコピー
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ./nginx/ssl/
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ./nginx/ssl/
```

### 4. 本番環境用コンテナの起動

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 5. データベースマイグレーション

初回デプロイ時：

```bash
docker exec go-app-prod sh -c "cd /usr/src/app && go run cmd/migrate/main.go"
```

### 6. ヘルスチェック

```bash
# Pingエンドポイント
curl https://yourdomain.com/api/ping

# データベース接続確認
curl https://yourdomain.com/api/db-check
```

## セキュリティチェックリスト

本番環境デプロイ前に以下を確認してください：

- [ ] すべてのパスワードを強力なものに変更
- [ ] JWT_SECRETとJWT_REFRESH_SECRETを新規生成
- [ ] ENV_MODE=productionに設定
- [ ] 実際のSMTPサーバー情報を設定
- [ ] ALLOWED_ORIGINSに本番ドメインのみを設定
- [ ] SSL証明書を設定（HTTPS）
- [ ] ファイアウォール設定（必要なポートのみ開放）
- [ ] データベースのバックアップ設定
- [ ] ログ監視の設定
- [ ] レート制限の有効化

## 環境変数一覧

### 必須

| 変数名 | 説明 | 例 |
|--------|------|-----|
| ENV_MODE | 環境モード | development / production |
| MYSQL_ROOT_PASSWORD | MySQLルートパスワード | - |
| MYSQL_DATABASE | データベース名 | app |
| MYSQL_USER | データベースユーザー | user |
| MYSQL_PASSWORD | データベースパスワード | - |
| JWT_SECRET | JWTシークレット | - |
| JWT_REFRESH_SECRET | リフレッシュトークンシークレット | - |
| APP_URL | アプリケーションURL | http://localhost:8000 |
| EMAIL_FROM | 送信元メールアドレス | noreply@example.com |
| SMTP_HOST | SMTPホスト | mailpit / smtp.gmail.com |
| SMTP_PORT | SMTPポート | 1025 / 587 |

### オプション

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| TZ | タイムゾーン | Asia/Tokyo |
| APP_LANG | 言語設定 | en |
| MYSQL_HOST | MySQLホスト | mysql |
| MYSQL_PORT | MySQLポート | 3306 |
| REDIS_HOST | Redisホスト | redis |
| REDIS_PORT | Redisポート | 6379 |
| REDIS_PASSWORD | Redisパスワード | - |
| SMTP_USER | SMTPユーザー | - |
| SMTP_PASSWORD | SMTPパスワード | - |
| ALLOWED_ORIGINS | CORS許可オリジン | localhost:3000,localhost:5173 |
| LOG_LEVEL | ログレベル | info |
| SESSION_TIMEOUT | セッションタイムアウト | 24h |
| REFRESH_TOKEN_TIMEOUT | リフレッシュトークンタイムアウト | 168h |

## トラブルシューティング

### コンテナが起動しない

```bash
# ログを確認
docker-compose logs go
docker-compose logs mysql
docker-compose logs redis
```

### データベース接続エラー

```bash
# MySQLコンテナの状態確認
docker exec mysql-prod mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "SHOW DATABASES;"
```

### メール送信エラー

開発環境：
- Mailpitコンテナが起動しているか確認
- http://localhost:8025 でメールを確認

本番環境：
- SMTP認証情報が正しいか確認
- SMTPサーバーのファイアウォール設定を確認

## バックアップとリストア

### データベースバックアップ

```bash
# バックアップ
docker exec mysql-prod mysqldump -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE} > backup.sql

# リストア
docker exec -i mysql-prod mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE} < backup.sql
```

### Redisバックアップ

```bash
# バックアップ
docker exec redis-prod redis-cli SAVE
docker cp redis-prod:/data/dump.rdb ./redis-backup.rdb

# リストア
docker cp ./redis-backup.rdb redis-prod:/data/dump.rdb
docker restart redis-prod
```

## モニタリング

### ログの確認

```bash
# アプリケーションログ
docker logs -f go-app-prod

# エラーログのみ
docker logs go-app-prod 2>&1 | grep -i error
```

### リソース使用状況

```bash
docker stats
```

## アップデート手順

```bash
# 1. 最新コードを取得
git pull origin main

# 2. コンテナを再ビルド
docker-compose -f docker-compose.prod.yml build

# 3. コンテナを再起動
docker-compose -f docker-compose.prod.yml up -d

# 4. ヘルスチェック
curl https://yourdomain.com/api/ping
```
