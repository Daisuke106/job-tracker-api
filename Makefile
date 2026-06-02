.PHONY: up down restart logs \
        docker-build docker-up docker-down docker-logs \
        run build \
        migrate-up migrate-down migrate-version migrate-force migrate-create

# ── ローカル開発（DB のみ）────────────────────────────
up:
	docker-compose up -d db

down:
	docker-compose down

restart:
	docker-compose down && docker-compose up -d db

logs:
	docker-compose logs -f db

# ── Docker フルスタック（DB + app）────────────────────
docker-build:
	docker-compose build app

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

# ── アプリ ────────────────────────────────────────────
run:
	go run ./cmd/api/main.go

build:
	go build -o bin/api ./cmd/api/main.go

# ── マイグレーション（docker-compose 経由）─────────────
# ※ DB が起動済みであること（make up 後に使用）

migrate-up:
	docker-compose run --rm migrate up

migrate-down:
	docker-compose run --rm migrate down 1

migrate-version:
	docker-compose run --rm migrate version

# 強制的にバージョンを指定（dirty状態の解除に使用）
# 例: make migrate-force V=1
migrate-force:
	docker-compose run --rm migrate force $(V)

# 新しいマイグレーションファイルを作成
# 例: make migrate-create NAME=add_questions_table
migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)
