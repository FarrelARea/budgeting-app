# budgeting-app

Monorepo for the budgeting app: a Next.js web client and a Go API.

## Layout

```
apps/
  web/          Next.js 15 (App Router) — UI, server components, route handlers
  api/          Go service — HTTP API, business logic, persistence
packages/
  ui/           Shared React component library (Tailwind + shadcn-style)
  types/        Shared TypeScript types & Zod schemas (contract with the API)
  database/     Prisma client + schema (Node/web-side types; Go owns migrations)
  eslint-config/ Shared ESLint flat configs
  typescript-config/ Shared tsconfig presets
tooling/        Repo scripts & codegen
.github/        CI workflows
```

## Stack

- Monorepo tooling: Turborepo + pnpm workspaces
- Web: Next.js (App Router) + TypeScript + Tailwind CSS
- API: Go (chi router) + sqlc + sqlx (PostgreSQL)
- DB source of truth: SQL migrations in `apps/api/migrations` (Go side)
- Shared TS contract: `packages/types` (+ optional Prisma client types for Node code)
- CI: GitHub Actions

## Getting started

```bash
pnpm install
pnpm dev          # runs web + api in parallel via turbo
```

Web: http://localhost:3000 — API: http://localhost:8080

## Commands

| Command | Description |
| --- | --- |
| `pnpm dev` | Start all dev servers |
| `pnpm build` | Build all apps/packages |
| `pnpm lint` | Lint everything |
| `pnpm test` | Run all tests |
| `pnpm typecheck` | Type-check TS projects |

See `docs/plan.md` (planning artifact) for the full rollout plan.
```