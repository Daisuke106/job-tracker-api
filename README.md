# Job Tracker API

求人票・面談管理システムの REST API です。  
Go + PostgreSQL で構築し、Layered Architecture を採用しています。

---

## 技術スタック

| 用途 | 技術 |
|------|------|
| 言語 | Go 1.26 |
| HTTP フレームワーク | [Gin](https://github.com/gin-gonic/gin) |
| DB | PostgreSQL 16 |
| DB アクセス | [sqlx](https://github.com/jmoiron/sqlx) |
| マイグレーション | [golang-migrate v4](https://github.com/golang-migrate/migrate) |
| 認証 | JWT ([golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)) |
| パスワードハッシュ | bcrypt |
| 環境変数 | [godotenv](https://github.com/joho/godotenv) |
| コンテナ | Docker / Docker Compose |

---

## アーキテクチャ

### レイヤー構成（Layered Architecture）

```
HTTPリクエスト
      │
      ▼
┌─────────────┐
│   handler   │  HTTPリクエスト/レスポンス、バリデーション
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   usecase   │  業務ロジック（認証・CRUD・所有権チェック）
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ repository  │  DB 操作（SQL 発行）
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ PostgreSQL  │
└─────────────┘
```

### ディレクトリ構成

```
job-tracker-api/
├── cmd/
│   └── api/
│       └── main.go          # エントリポイント・DI配線・ルーティング
├── internal/
│   ├── config/
│   │   └── config.go        # 環境変数の読み込み
│   ├── domain/
│   │   ├── user.go          # User エンティティ
│   │   ├── company.go       # Company エンティティ
│   │   ├── job.go           # Job エンティティ・JobStatus 定数
│   │   └── interview.go     # Interview / InterviewNote エンティティ
│   ├── handler/
│   │   ├── auth_handler.go      # 認証ハンドラ
│   │   ├── company_handler.go   # 企業 CRUD ハンドラ
│   │   ├── job_handler.go       # 求人ハンドラ
│   │   └── interview_handler.go # 面談ハンドラ（未実装）
│   ├── usecase/
│   │   ├── auth_usecase.go      # 登録・ログイン・JWT 発行
│   │   ├── company_usecase.go   # 企業 CRUD ロジック
│   │   └── job_usecase.go       # 求人 CRUD・ステータス更新ロジック
│   ├── repository/
│   │   ├── user_repository.go   # users テーブル操作
│   │   ├── company_repository.go # companies テーブル操作
│   │   └── job_repository.go    # jobs テーブル操作
│   ├── middleware/
│   │   └── auth_middleware.go   # JWT 検証・user_id をコンテキストに設定
│   └── db/
│       ├── postgres.go          # sqlx DB 接続
│       └── migrate.go           # マイグレーション実行（embed.FS 経由）
├── migrations/
│   ├── embed.go                 # SQL ファイルをバイナリに埋め込み
│   ├── 000001_init.up.sql       # 全テーブル作成
│   └── 000001_init.down.sql     # 全テーブル削除
├── Dockerfile                   # マルチステージビルド
├── docker-compose.yml           # DB / migrate / app サービス定義
├── Makefile                     # 開発用コマンド集
├── .env.example                 # 環境変数テンプレート
└── go.mod
```

---

## DB スキーマ

```
users
├── id           SERIAL PK
├── email        VARCHAR(255) UNIQUE NOT NULL
├── password_hash VARCHAR(255) NOT NULL
├── name         VARCHAR(100) NOT NULL
└── created_at   TIMESTAMP

companies
├── id           SERIAL PK
├── user_id      FK → users.id
├── name         VARCHAR(255) NOT NULL
├── industry     VARCHAR(100)
├── website_url  TEXT
├── memo         TEXT
├── created_at   TIMESTAMP
└── updated_at   TIMESTAMP

jobs
├── id           SERIAL PK
├── user_id      FK → users.id
├── company_id   FK → companies.id
├── title        VARCHAR(255) NOT NULL
├── description  TEXT
├── status       VARCHAR(50)  ← applied / screening / interview / offer / rejected / withdrawn
├── applied_at   TIMESTAMP
├── created_at   TIMESTAMP
└── updated_at   TIMESTAMP

interviews
├── id           SERIAL PK
├── job_id       FK → jobs.id
├── user_id      FK → users.id
├── scheduled_at TIMESTAMP NOT NULL
├── location     VARCHAR(255)
├── memo         TEXT
├── created_at   TIMESTAMP
└── updated_at   TIMESTAMP

interview_notes
├── id           SERIAL PK
├── interview_id FK → interviews.id
├── content      TEXT NOT NULL
└── created_at   TIMESTAMP
```

---

## API エンドポイント

### 認証（認証不要）

| メソッド | パス | 説明 |
|--------|------|------|
| POST | `/api/v1/auth/register` | ユーザー登録 |
| POST | `/api/v1/auth/login` | ログイン（JWT トークン発行） |

### 企業（要認証）

| メソッド | パス | 説明 |
|--------|------|------|
| GET | `/api/v1/companies` | 企業一覧取得 |
| POST | `/api/v1/companies` | 企業登録 |
| GET | `/api/v1/companies/:id` | 企業詳細取得 |
| PUT | `/api/v1/companies/:id` | 企業情報更新 |
| DELETE | `/api/v1/companies/:id` | 企業削除 |

### 求人（要認証）

| メソッド | パス | 説明 |
|--------|------|------|
| GET | `/api/v1/jobs` | 求人一覧取得 |
| POST | `/api/v1/jobs` | 求人登録 |
| GET | `/api/v1/jobs/:id` | 求人詳細取得 |
| PUT | `/api/v1/jobs/:id/status` | ステータス更新 |

### 面談（要認証・未実装）

| メソッド | パス | 説明 |
|--------|------|------|
| GET | `/api/v1/interviews` | 面談一覧取得 |
| POST | `/api/v1/interviews` | 面談登録 |
| GET | `/api/v1/interviews/:id` | 面談詳細取得 |
| POST | `/api/v1/interviews/:id/notes` | 面談メモ登録 |

### ヘルスチェック

| メソッド | パス | 説明 |
|--------|------|------|
| GET | `/health` | サーバー・DB 疎通確認 |

> `/health` は `feature/health-endpoint` ブランチで実装済み、マージ待ち

---

## 環境構築

### 必要なもの

- [Go 1.26+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Make](https://www.gnu.org/software/make/)（任意）

### 1. リポジトリのクローン

```bash
git clone https://github.com/Daisuke106/job-tracker-api.git
cd job-tracker-api
```

### 2. 環境変数の設定

```bash
cp .env.example .env
```

`.env` を編集して `JWT_SECRET` に任意の文字列を設定します。

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=jobuser
DB_PASSWORD=jobpass
DB_NAME=job_tracker
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-here   # ← 必須
```

### 3. 起動

#### ローカル開発（DB のみ Docker）

```bash
# DB 起動
make up        # または: docker-compose up -d db

# サーバー起動（マイグレーション自動実行）
make run       # または: go run ./cmd/api/main.go
```

#### フルスタック（全コンテナ）

```bash
make docker-up   # または: docker-compose up -d
```

### 4. 動作確認

```bash
# ユーザー登録
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"テスト太郎"}'

# ログイン
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 企業登録（取得したトークンを使用）
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"株式会社テスト","industry":"IT"}'
```

---

## マイグレーション

マイグレーションは **サーバー起動時に自動実行** されます（`embed.FS` 経由でバイナリに埋め込み済み）。

手動操作が必要な場合は Makefile コマンドを使用します。

```bash
make migrate-up               # 全マイグレーション適用
make migrate-down             # 1バージョン戻す
make migrate-version          # 現在のバージョン確認
make migrate-force V=1        # dirty 状態の強制リセット
make migrate-create NAME=xxx  # 新しいマイグレーションファイル作成
```

---

## 認証フロー

```
POST /api/v1/auth/register  →  bcrypt でパスワードをハッシュ化して保存
POST /api/v1/auth/login     →  JWT トークン（有効期限 24 時間）を発行

保護されたエンドポイント:
  Authorization: Bearer <token>  →  JWT 検証 → user_id をコンテキストに設定
```

---

## 今後の実装予定

- [ ] 面談予定の登録・取得（`/api/v1/interviews`）
- [ ] 面談メモの登録（`/api/v1/interviews/:id/notes`）
- [ ] 入力バリデーションの強化
- [ ] ユニットテスト / ハンドラテスト
- [ ] 非同期通知（面談リマインダー）
- [ ] バッチ処理（期限切れ求人の自動更新）
- [ ] OpenAPI / Swagger ドキュメント生成
