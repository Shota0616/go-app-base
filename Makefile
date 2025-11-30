.PHONY: test test-coverage test-verbose seed

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

# モックデータを投入
seed:
	cd go && go run cmd/seed/main.go

# Docker環境でモックデータを投入
seed-docker:
	docker exec go-app-base-go-1 sh -c "cd /usr/src/app && go run cmd/seed/main.go"
