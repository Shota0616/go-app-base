.PHONY: test test-coverage test-verbose

# すべてのテストを実行
test:
	cd go && go test ./...

# カバレッジ付きでテストを実行
test-coverage:
	cd go && go test -cover ./...

# 詳細出力でテストを実行
test-verbose:
	cd go && go test -v ./...

# 特定のパッケージのテストを実行
test-controllers:
	cd go && go test -v ./cmd/api/controllers/...

test-auth:
	cd go && go test -v ./auth/...

test-middleware:
	cd go && go test -v ./middleware/...
