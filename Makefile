.PHONY: up down restart run build \
        migrate-up migrate-down migrate-version migrate-force migrate-create

# ── Docker ────────────────────────────────────────────
up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose down && docker-compose up -d

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
