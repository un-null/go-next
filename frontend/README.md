# go-next frontend

A React 19 SPA for the go-next API: a catalogue, a product page, a cart, an order history and an
account. Vite, Mantine, TanStack Query, React Router.

- Live: <https://frontend.kawabata-qppq.workers.dev>
- API: [`../backend`](../backend)

The API runs on Render's free tier, so the first request after fifteen minutes can take twenty
seconds. The UI says so after three.

## Run it locally

```sh
pnpm install
echo 'VITE_API_URL=https://go-next-f5yi.onrender.com' > .env.local
pnpm dev --port 5173 --strictPort
```

Port 5173 is not a preference: that origin is the one the API allows through CORS
(`ALLOWED_ORIGINS`). Vite would happily move to 5174 when the port is busy, and every request from
there fails as `network request failed`. `--strictPort` makes it stop instead.

`.env.local` is gitignored. The variable is read at build time, so changing it needs a restart.

## Checks

```sh
pnpm check        # Biome: lint and formatting
pnpm typecheck    # tsc --noEmit
pnpm test         # Vitest: unit tests, plus page tests in jsdom
pnpm build        # typecheck, then the production bundle
```

CI runs the same four on every pull request.

## Deploy

Static assets on Cloudflare Workers, deployed by hand:

```sh
pnpm build
pnpm exec wrangler deploy
```

`wrangler.jsonc` points at `./dist` and sets `not_found_handling`, so client-side routes survive a
reload. There is no deploy step in CI.

## Demo account

`demo@go-next.dev` / `go-next-demo` — the sign-in page has a button for it. The account is shared,
so its cart, orders and coin balance are whatever the last visitor left.
