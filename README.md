# go-app-base

```
docker-compose up --build
```


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
