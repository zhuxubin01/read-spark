# ReadSpark Backend MVP

## Prerequisites

- Go 1.24+
- Docker (for local PostgreSQL)
- `jq` (recommended for pretty-printing curl JSON output)

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

## Health / Readiness / Metrics

```bash
# liveness
curl -s http://localhost:8080/health | jq

# readiness (includes DB ping)
curl -s http://localhost:8080/ready | jq

# prometheus metrics
curl -s http://localhost:8080/metrics | head -n 20
```

## Billing Config (mock vs apple)

`configs/config.yaml`:

```yaml
billing:
  receipt_provider: "mock"      # mock | apple
  apple_environment: "sandbox" # sandbox | production
  apple_shared_secret: ""       # optional, recommended for production
  apple_verify_url: ""          # optional override for testing
```

- `receipt_provider=mock`: always uses mock verifier, suitable for local MVP.
- `receipt_provider=apple`: calls Apple verify endpoint.
- If `apple_environment=production` and Apple returns `21007`, server auto-fallbacks to sandbox verify.

## Integration Mode Switch

### Mode A: mock receipt verification (default)

1. Keep config:
   - `billing.receipt_provider: "mock"`
2. Restart server.
3. Use any non-empty `receipt` in subscription API.

### Mode B: Apple receipt verification

1. Update config:
   - `billing.receipt_provider: "apple"`
   - `billing.apple_environment: "sandbox"` (or `production`)
   - `billing.apple_shared_secret: "<your-secret>"` (optional but recommended)
2. Restart server.
3. Call subscription API with a real Apple receipt.
4. Error mapping behavior:
   - invalid/malformed receipt -> `400 INVALID_RECEIPT`
   - Apple upstream/network failure -> `502 RECEIPT_UPSTREAM_ERROR`
   - internal unexpected error -> `500 INTERNAL_ERROR`

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

# 7) Create subscription (protected)
# Mock mode: any non-empty receipt is enough
curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plan_type":"monthly","receipt":"mock-receipt","payment_channel":"apple"}' | jq

# 8) Query subscription status (protected)
curl -s "$BASE_URL/subscriptions/status" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq

# 9) Dictionary lookup (public)
curl -s "$BASE_URL/dictionary/hello" | jq

# 10) Create annotation (protected)
curl -s -X POST "$BASE_URL/annotations" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"article_id\":\"$ARTICLE_ID\",\"type\":\"highlight\",\"range_start\":1,\"range_end\":12}" | jq

# 11) List annotations (protected)
curl -s "$BASE_URL/annotations?article_id=$ARTICLE_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq

# 12) Register push token (protected, MVP mock persistence)
curl -s -X POST "$BASE_URL/push/token" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_token":"dev-token-001","platform":"ios"}' | jq
```

## Notes

- Current MVP uses PostgreSQL full-text search for article search.
- SMS verification code is configurable via `auth.verification_code` (default `123456`).
- Receipt verification is abstracted behind `ReceiptVerifier`; set `billing.receipt_provider` to `mock` or `apple`.
- Apple verifier supports auto sandbox fallback for `21007` when running in production mode.
- Real SMS verification and production-grade billing reconciliation are not fully integrated in MVP.
