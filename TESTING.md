# テストガイド

## 概要

このプロジェクトには、認証、ユーザー管理、ミドルウェア、ヘルスチェック、環境設定など、すべての主要機能をカバーする包括的なテストスイートが含まれています。

**テスト統計: 56テスト（100%成功）**
- 認証・ユーザー管理: 29テスト
- ミドルウェア: 3テスト
- 認証ユーティリティ: 6テスト
- ヘルスチェック: 2テスト
- 環境設定: 8テスト
- 統合テスト（開発/本番環境）: 8テスト

## テストの実行方法

### すべてのテストを実行
```bash
make test
```

### カバレッジ付きでテストを実行
```bash
make test-coverage
```

### 詳細出力でテストを実行
```bash
make test-verbose
```

### 特定のパッケージのテストを実行

コントローラーのテスト:
```bash
make test-controllers
```

認証機能のテスト:
```bash
make test-auth
```

ミドルウェアのテスト:
```bash
make test-middleware
```

環境設定のテスト:
```bash
make test-env
```

統合テスト（開発/本番環境）:
```bash
make test-integration
```

開発環境テストのみ:
```bash
make test-dev
```

本番環境テストのみ:
```bash
make test-prod
```

### Docker環境内で実行
```bash
docker exec go-app-base-go-1 sh -c "cd /usr/src/app && go test -v ./..."
```

## テストファイルの構成

```
go/
├── auth/
│   └── jwt_test.go              # JWT生成・検証、パスワードハッシュのテスト
├── middleware/
│   └── auth_test.go             # 認証ミドルウェアのテスト
└── cmd/api/controllers/
    ├── auth_test.go             # 認証関連エンドポイントのテスト（29テスト）
    ├── ping_test.go             # Pingエンドポイントのテスト
    └── db_check_test.go         # DB接続確認のテスト
```

## テストカバレッジ詳細

### 1. 認証・ユーザー管理（29テスト）

#### ユーザー登録
- ✅ `TestRegisterSuccess` - 正常な登録
- ✅ `TestRegisterDuplicateEmail` - 重複メールアドレス
- ✅ `TestRegisterInvalidInput` - 無効な入力（メール形式エラー）

#### メール認証
- ✅ `TestVerifySuccess_Registration` - 登録時の認証成功
- ✅ `TestVerifySuccess_EmailChange` - メールアドレス変更時の認証成功
- ✅ `TestVerifyInvalidCode` - 無効な認証コード
- ✅ `TestVerifyExpiredCode` - 期限切れ認証コード

#### 認証コード再送
- ✅ `TestResendVerificationCodeSuccess` - 再送成功
- ✅ `TestResendVerificationCodeLimitReached` - 再送制限到達

#### ログイン
- ✅ `TestLoginSuccess` - ログイン成功
- ✅ `TestLoginInvalidCredentials` - 無効な認証情報
- ✅ `TestLoginInactiveAccount` - 非アクティブアカウント

#### パスワードリセット
- ✅ `TestRequestPasswordResetSuccess` - リセット要求成功
- ✅ `TestRequestPasswordResetInactiveAccount` - 非アクティブアカウント
- ✅ `TestResetPasswordSuccess` - パスワードリセット成功
- ✅ `TestResetPasswordInvalidToken` - 無効なトークン

#### ユーザー情報取得
- ✅ `TestGetUserSuccess` - ユーザー情報取得成功
- ✅ `TestGetUserUnauthorized` - 認証なし

#### ユーザー名更新
- ✅ `TestUpdateUsernameSuccess` - 更新成功
- ✅ `TestUpdateUsernameDuplicate` - 重複ユーザー名

#### メールアドレス更新
- ✅ `TestUpdateEmailSuccess` - 更新成功
- ✅ `TestUpdateEmailDuplicateNewEmail` - 重複メールアドレス
- ✅ `TestUpdateEmailSameEmail` - 同じメールアドレス

#### パスワード更新
- ✅ `TestUpdatePasswordSuccess` - 更新成功
- ✅ `TestUpdatePasswordIncorrectCurrent` - 現在のパスワードが不正

#### ユーザー削除
- ✅ `TestDeleteUserSuccess` - 削除成功
- ✅ `TestDeleteUserUnauthorized` - 認証なし

#### トークンリフレッシュ
- ✅ `TestRefreshTokenSuccess` - リフレッシュ成功
- ✅ `TestRefreshTokenInvalid` - 無効なリフレッシュトークン

### 2. ミドルウェア（3テスト）

- ✅ `TestAuthRequiredSuccess` - 有効なトークンで認証成功
- ✅ `TestAuthRequiredNoToken` - トークンなしで401エラー
- ✅ `TestAuthRequiredInvalidToken` - 無効なトークンで401エラー

### 3. 認証ユーティリティ（6テスト）

#### JWT
- ✅ `TestGenerateJWT` - JWT生成
- ✅ `TestValidateJWT` - JWT検証
- ✅ `TestValidateJWTInvalidToken` - 無効なトークン検証
- ✅ `TestValidateJWTExpiredToken` - 期限切れトークン検証

#### パスワード
- ✅ `TestHashPassword` - パスワードハッシュ化
- ✅ `TestCheckPasswordHash` - パスワード検証

### 4. ヘルスチェック（2テスト）

- ✅ `TestPing` - Pingエンドポイント
- ✅ `TestDBCheckSuccess` - データベース接続確認

### 5. 環境設定（8テスト）

- ✅ `TestIsProduction` - 本番環境判定
- ✅ `TestIsDevelopment` - 開発環境判定
- ✅ `TestGetEnv` - 環境変数取得
- ✅ `TestGetEnvInt` - 整数型環境変数取得
- ✅ `TestGetEnvBool` - 真偽値型環境変数取得
- ✅ `TestGetEnvDuration` - 期間型環境変数取得
- ✅ `TestGetLogLevel` - ログレベル取得
- ✅ `TestGetAllowedOrigins` - CORS許可オリジン取得

### 6. 統合テスト（8テスト）

#### 開発環境テスト（4テスト）
- ✅ `TestDevelopmentEnvironment` - 開発環境設定確認
- ✅ `TestDevelopmentCORS` - 開発環境CORS設定
- ✅ `TestDevelopmentMailpitConfig` - Mailpit設定確認
- ✅ `TestDevelopmentPingEndpoint` - 開発環境エンドポイント

#### 本番環境テスト（9テスト）
- ✅ `TestProductionEnvironment` - 本番環境設定確認
- ✅ `TestProductionGinMode` - Ginリリースモード確認
- ✅ `TestProductionCORS` - 本番環境CORS設定
- ✅ `TestProductionSMTPConfig` - 本番SMTP設定確認
- ✅ `TestProductionSecurityHeaders` - セキュリティヘッダー確認
- ✅ `TestProductionJWTSecrets` - JWT強度確認
- ✅ `TestProductionPingEndpoint` - 本番環境エンドポイント
- ✅ `TestProductionCORSRestriction` - CORS制限確認
- ✅ `TestProductionLogLevel` - ログレベル確認

#### 環境切り替えテスト（4テスト）
- ✅ `TestEnvironmentSwitch` - 環境モード切り替え
- ✅ `TestSMTPConfigByEnvironment` - 環境別SMTP設定
- ✅ `TestCORSConfigByEnvironment` - 環境別CORS設定
- ✅ `TestSecurityConfigByEnvironment` - 環境別セキュリティ設定

## テスト環境の要件

テストを実行する前に、以下の環境変数が設定されている必要があります:

```env
JWT_SECRET=test_jwt_secret
JWT_REFRESH_SECRET=test_jwt_refresh_secret
APP_URL=http://localhost:3000
EMAIL_FROM=noreply@example.com
SMTP_HOST=mailpit
SMTP_PORT=1025
REDIS_HOST=redis
REDIS_PORT=6379
APP_ROOT=/usr/src
APP_LANG=en
```

### 依存サービス

- **MySQL**: ユーザーデータの保存
- **Redis**: 認証コード、セッション管理
- **Mailpit**: メール送信のテスト（localhost:8025でWebUI確認可能）

## CI/CDへの統合

GitHub ActionsなどのCI/CDパイプラインに統合する場合:

```yaml
- name: Run tests
  run: |
    docker-compose up -d
    docker-compose exec -T go sh -c "cd /usr/src/app && go test -v ./..."
```

## トラブルシューティング

### テストが失敗する場合

1. **Redisの接続エラー**: Redisコンテナが起動しているか確認
   ```bash
   docker-compose ps redis
   ```

2. **MySQLの接続エラー**: MySQLコンテナが起動しているか確認
   ```bash
   docker-compose ps mysql
   ```

3. **メール送信エラー**: Mailpitコンテナが起動しているか確認
   ```bash
   docker-compose ps mailpit
   ```

4. **環境変数の確認**: テスト環境変数が正しく設定されているか確認
   ```bash
   docker exec go-app-base-go-1 env | grep -E '(JWT|SMTP|REDIS)'
   ```

## テストのベストプラクティス

1. **テスト前のクリーンアップ**: 各テストは独立して実行できるよう、テスト前にDBとRedisをクリーンアップ
2. **モックの使用**: 外部サービス（メール送信など）はMailpitを使用してテスト
3. **エラーケースのテスト**: 正常系だけでなく、エラーケースも網羅的にテスト
4. **認証フローのテスト**: JWT生成から検証まで、完全な認証フローをテスト

