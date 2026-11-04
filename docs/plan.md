# Monorepo Plan — budgeting-app

## Goal
Convert `budgeting-app` into a clean monorepo separating the Next.js web client, the Go API,
and shared packages, with a single source of truth for the database schema.

## Decisions
- Monorepo tooling: **Turborepo + pnpm workspaces**
- Web: **Next.js 15 (App Router) + TypeScript + Tailwind**
- API: **Go 1.23 + chi + pgx + sqlc**
- Database: **PostgreSQL**. Go side owns migrations (`apps/api/migrations`) + sqlc queries.
  Prisma in `packages/database` mirrors schema for Node/web-side type generation only.
- Shared contract: `packages/types` (Zod + TS). Web generates nothing from DB directly.
- No mobile app yet (web-first); can add `apps/mobile` later.

## Target structure

```
budgeting-app/
├── apps/
│   ├── web/                 Next.js App Router client
│   └── api/                 Go HTTP service
├── packages/
│   ├── ui/                  Shared React components
│   ├── types/               Shared TS types & Zod schemas
│   ├── database/            Prisma client (Node/web-side)
│   ├── eslint-config/       Shared ESLint flat configs
│   └── typescript-config/   Shared tsconfig presets
├── tooling/                 Scripts & codegen
├── .github/workflows/       CI
├── turbo.json
├── pnpm-workspace.yaml
└── go.work
```

## Phases

### Phase 0 — Scaffold (done in this commit)
- Root tooling files (package.json, turbo.json, pnpm-workspace.yaml, .gitignore, .npmrc, go.work)
- Shared config packages (typescript-config, eslint-config)
- Shared packages (types, ui, database)
- App shells (web, api) with health endpoint + migration + CI

### Phase 1 — Wire the contract
- Define OpenAPI spec for the Go API (swaggo or hand-written)
- Generate TS client into `packages/api-client` (or directly in web)
- Replace placeholder handlers with sqlc-generated code
- Add auth (JWT/session) end-to-end

### Phase 2 — Features
- Accounts, categories, transactions, budgets CRUD (API + web)
- Dashboard & reports
- Imports (CSV) + recurring transactions

### Phase 3 — Quality & ops
- Integration tests (API + DB via testcontainers)
- E2E (Playwright) for critical flows
- Dockerfiles + compose for local dev
- Release pipeline (build images, deploy)

## Risks / notes
- Two DB access layers (sqlc Go + Prisma Node) must stay in sync. Single source of truth
  is the SQL in `apps/api/migrations`; Prisma schema is a mirror. Add a CI check that
  diff's generated types against a regenerated Prisma client.
- pnpm `workspace:*` protocol requires pnpm; do not use npm/yarn for installs.
- Go module path is `github.com/farrelarea/budgeting-app/apps/api`; imports use full path.
- Turbo `dev` is persistent; running `pnpm dev` starts both web (3000) and api (8080).

## Verification
- `pnpm install`
- `pnpm dev` → web on :3000, api on :8080
- `curl localhost:8080/healthz` → `{"status":"ok"}`
- Home page renders API status badge using `/healthz`.
- `pnpm build && pnpm typecheck && pnpm lint`
- `cd apps/api && go build ./... && go vet ./...` (requires Go 1.23)
- CI workflow runs web + api jobs.
