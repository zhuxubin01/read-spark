# ReadSpark Backend MVP

## Prerequisites

- Go 1.24+
- Docker (for local PostgreSQL)

## Quick Start

```bash
cd /Users/zhuxubin/workspace/projects/read-spark/.worktrees/mvp/read-spark-backend

# 1) Start PostgreSQL
docker run -d --name readspark-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=readspark \
  -p 5432:5432 \
  postgres:16

# 2) Build
go build -o bin/server ./cmd/server

# 3) (Optional) Seed sample articles
go run ./scripts/seed.go

# 4) Run server
go run ./cmd/server
```

Server default address: `http://localhost:8080`

## API Verification (curl)

```bash
BASE_URL=http://localhost:8080/api/v1

# 1) Register (MVP verification code is always 123456)
REGISTER_RESP=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","code":"123456"}')
echo "$REGISTER_RESP"

ACCESS_TOKEN=$(echo "$REGISTER_RESP" | jq -r '.access_token')

# 2) Get daily articles (public)
curl -s "$BASE_URL/articles/daily" | jq

# 3) List articles (public)
curl -s "$BASE_URL/articles?page=1&page_size=10" | jq

# 4) Read one article (protected)
ARTICLE_ID=$(curl -s "$BASE_URL/articles?page=1&page_size=1" | jq -r '.articles[0].id')
curl -s "$BASE_URL/articles/$ARTICLE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq

# 5) Sync reading progress (protected)
curl -s -X POST "$BASE_URL/progress" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"article_id\":\"$ARTICLE_ID\",\"position\":320,\"percentage\":42.5}" | jq

# 6) Query reading progress list (protected)
curl -s "$BASE_URL/progress" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq

# 7) Create subscription (protected, mock receipt validation)
curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plan_type":"monthly","receipt":"mock-receipt","payment_channel":"apple"}' | jq

# 8) Query subscription status (protected)
curl -s "$BASE_URL/subscriptions/status" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq

# 9) Dictionary lookup (public)
curl -s "$BASE_URL/dictionary/hello" | jq
```

## Notes

- Current MVP uses PostgreSQL full-text search for article search.
- Real SMS verification and real receipt verification are not integrated in MVP.
