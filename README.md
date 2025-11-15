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
