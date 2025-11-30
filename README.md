# go-app-base

使用方法

別途.envファイルをルートディレクトリに作成してください。 
```env
TZ=Asia/Tokyo

# DB
MYSQL_ROOT_PASSWORD=password
MYSQL_DATABASE=sns
MYSQL_USER=user
MYSQL_PASSWORD=password

# メール（googleのsmtpを使用するとき）
EMAIL_ADDRESS=xxxxxxxxx@xxxxxx
EMAIL_PASSWORD="xxxxxxxxxxxxxxxx"

# URL
#GO_API_URL=http://192.168.111.102:8080/api
APP_URL=http://localhost:8000

# go jwt
JWT_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxx"
JWT_REFRESH_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxx"

# 環境に応じて "development" または "production" を設定
ENV_MODE="development"

# 言語設定（en or ja）
APP_LANG="en"
```
SECRETの生成
```
openssl rand -base64 32
```
コンテナ作成
```
docker-compose up --build
```
url確認 http://localhost:8000


## app

goのソースコードが配置されている。

## db

mysdql関連のファイルが格納されている

## elasticsearch

elasticsearch関連のファイルが格納されている。single node設定

## redis

redis関連のファイルが格納されている。

## web

nginx関連のファイルが格納されている。

## API Endpoints

### GET /api/ping

A simple health check endpoint to verify the application is running.

**Response:**
```json
{
  "message": "Hello World! Pong"
}
```

### GET /api/db-check



Checks the connection status to the MySQL database.



**Response (Success):**

```json

{

  "status": "ok",

  "message": "Database connection is healthy"

}

```



**Response (Failure):**

```json

{

  "status": "error",

  "message": "Failed to connect to database",

  "details": "..."

}

```



## 開発環境でのメール送信について



このプロジェクトは、開発環境で送信されるすべてのメールを捕捉するために [Mailpit](https://github.com/axllent/mailpit) を使用しています。



ユーザー登録時の認証メールやパスワードリセットのメールは、**実際のメールアドレスには送信されません。**



送信されたメールの内容を確認するには、以下のMailpitのWebインターフェースにアクセスしてください。



- **Mailpit Web UI:** [http://localhost:8025](http://localhost:8025)

## テスト

このプロジェクトには包括的なテストスイートが含まれています。

### テストの実行

```bash
# すべてのテストを実行
make test

# カバレッジ付きでテストを実行
make test-coverage

# 詳細出力でテストを実行
make test-verbose

# Docker環境内で実行
docker exec go-app-base-go-1 sh -c "cd /usr/src/app && go test -v ./..."
```

### テストカバレッジ

**合計: 38テスト（100%成功）**

#### 認証・ユーザー管理（29テスト）
- ユーザー登録（成功、重複メール、無効な入力）
- メール認証（登録時、メールアドレス変更時）
- ログイン（成功、無効な認証情報、非アクティブアカウント）
- パスワードリセット（リクエスト、実行、無効なトークン）
- 認証コード再送（成功、制限到達）
- ユーザー情報取得・更新（ユーザー名、メールアドレス、パスワード）
- ユーザー削除
- トークンリフレッシュ

#### ミドルウェア（3テスト）
- 認証ミドルウェア（有効なトークン、トークンなし、無効なトークン）

#### 認証ユーティリティ（6テスト）
- JWT生成・検証
- 期限切れトークンの検証
- パスワードハッシュ化・検証

#### ヘルスチェック（2テスト）
- Pingエンドポイント
- データベース接続確認

詳細は [TESTING.md](TESTING.md) を参照してください。

