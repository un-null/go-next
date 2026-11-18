# go-next API

A small shop API where prices are paid in coins. Go 1.23, Echo, pgx, sqlc, PostgreSQL 16.

- API: <https://go-next-f5yi.onrender.com>
- Frontend: <https://frontend.kawabata-qppq.workers.dev> (source in `../frontend`)

The free instance sleeps after fifteen minutes of inactivity, so the first request takes twenty to
twenty-five seconds. The SPA shows a hint after three.

![Backend layers](architecture.png)

## The checkout transaction

`POST /api/orders` is the part worth reading. Four things happen inside one transaction:

1. the user row is locked, so two requests cannot spend the same coins
2. the products it touches are locked in ascending id order — the order every writer uses, so two
   checkouts cannot deadlock each other
3. the order and its lines are written, stock is decremented, coins are debited and the ledger entry
   is appended
4. the cart is cleared

If any step fails, nothing is committed: a 409 for stock and a 402 for coins leave the cart, the
balance and the ledger exactly as they were.

`TestCheckout_ConcurrentLastUnit` fires two checkouts at the same last unit and asserts that one
order exists, the other call got a 409, stock is 0, the ledger has one entry and the balance was
charged once. Prices are never read from the request — the client sends product ids and quantities,
and the amounts come from the stored rows.

## API

Bearer means `Authorization: Bearer <token>`. Tokens last 24 hours and reach only their own user's
data.

| Method | Path | Auth | Notes |
| --- | --- | --- | --- |
| POST | `/api/signup` | — | `{name,email,password}` → 201. The account starts with 1,000 coins. |
| POST | `/api/login` | — | `{email,password}` → `{token,user}`. |
| GET | `/api/products` | — | `?q=&category=&page=&limit=` (defaults: page 1, limit 10). |
| GET | `/api/products/:id` | — | 404 when the product is gone. |
| GET | `/api/categories` | — | |
| GET | `/api/users/:id` | Bearer | Your own id only; the balance lives here. |
| PATCH | `/api/users/:id/name` | Bearer | `{name}`. Email takes `{email}`; password takes `{current_password,new_password}`. |
| DELETE | `/api/users/:id` | Bearer | |
| GET | `/api/cart` | Bearer | `{items,total_items,success}`. `total_items` counts lines, not units. |
| POST | `/api/cart` | Bearer | `{product_id,quantity}` **sets** the quantity; it does not add. |
| DELETE | `/api/cart?product_id=` | Bearer | Removes one line. |
| DELETE | `/api/cart/clear` | Bearer | |
| POST | `/api/orders` | Bearer | `{items:[{product_id,quantity}]}` → 201 with the order. 402 insufficient coins, 409 insufficient stock. |
| GET | `/api/orders/:id` | Bearer | Includes the line items. Another user's order is a 404. |
| GET | `/api/users/:userID/orders` | Bearer | Newest first. No items; fetch one order for those. |
| POST | `/api/coins/charge` | Bearer | Demo top-up, `{amount,description}`, 1–10,000. No payment behind it. |
| POST | `/api/coins/spend` | Bearer | `{amount,description}`; refuses to go below zero. |
| GET | `/api/coins/transactions` | Bearer | `{transactions,page,limit}`, limit capped at 100. |

A coin is one dollar, rounded up: $49.99 costs 50.

## Layout

```
cmd/service             HTTP server
cmd/migrate             migration runner
internal/entity         types and response shapes
internal/repository     SQL access (sqlc output in internal/database)
internal/usecase        business rules, including the checkout transaction
internal/delivery/http  handlers, routes, JWT middleware
migrations              schema and seed data
```

## Run it locally

```sh
docker compose -f compose.test.yml up -d          # Postgres 16 on port 5434
export DATABASE_URL='postgres://postgres:postgres@localhost:5434/gonet_test?sslmode=disable'
export JWT_SECRET='dev-secret'
go run ./cmd/migrate                              # schema and seed data
go run ./cmd/service                              # http://localhost:8080
```

The seed creates 3 categories and 20 products; `POST /api/signup` gives a new account 1,000 coins.
`GET /healthz` answers `{"status":"ok"}` without touching the database.

## Environment

| Variable | Required | Default | What it is |
| --- | --- | --- | --- |
| `JWT_SECRET` | yes | — | HMAC key for the tokens. |
| `DATABASE_URL` | yes | — | `postgres://…`. The process rewrites the scheme to `pgx5://` for the migration driver. |
| `ALLOWED_ORIGINS` | no | `http://localhost:5173` | Comma-separated list for CORS. |
| `PORT` | no | `8080` | Listen port. |
| `MIGRATIONS_DIR` | no | `migrations` | Where `cmd/migrate` looks for SQL files. |

A `.env` file is read if present, but nothing depends on it.

## Tests

```sh
go test ./...                                     # unit tests, no database
docker compose -f compose.test.yml up -d
TEST_DATABASE_URL='postgres://postgres:postgres@localhost:5434/gonet_test?sslmode=disable' \
  go test -tags=integration ./...
docker compose -f compose.test.yml down
```

The integration tests apply the migrations to a scratch database themselves, after dropping its
public schema. Do not point `TEST_DATABASE_URL` at anything you want to keep.

## Migrations

Migrations live in `migrations/` as numbered `up`/`down` pairs and are applied with:

```sh
DATABASE_URL='postgres://…' go run ./cmd/migrate
```

The command runs `up` only and logs `database is at version N (dirty=false)`. Production is a Neon
database named `gonext` — check the name before running anything by hand. The pooled connection
string works as long as `channel_binding` stays out of it; pgx resets the connection when it is
present.

Render's free plan has no pre-deploy command, so after a deploy that includes a migration, run the
command above once against the production URL. The image ships both binaries, so it is available
inside the container too.

## Deploy

The image is built from `Dockerfile` (static binaries, Alpine, runs as `nobody`) and deployed to
Render, which provides `PORT` and the environment variables above. `/healthz` is the health check.
