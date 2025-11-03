# go-app-base

使用方法

1. 別途`.env`ファイルをルートディレクトリに作成してください。
.env
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

2. コンテナ作成
```
docker-compose up --build
```
3. url確認
http://localhost:8000



## app

goのソースコードが配置されている。

## db

mysdql関連のファイルが格納されている

## redis

redis関連のファイルが格納されている。

## web

nginx関連のファイルが格納されている。
