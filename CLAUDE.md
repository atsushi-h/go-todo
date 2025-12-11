# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Todo application with a Go backend API and Next.js frontend in a pnpm monorepo structure managed by Turborepo.

## Commands

### Monorepo (root level)
```bash
pnpm dev          # Start all apps (turbo dev)
pnpm build        # Build all apps (turbo build)
pnpm test         # Run all tests (turbo test)
pnpm lint         # Run all linters (turbo lint)
```

### Docker (backend development)
```bash
make dcu-dev      # Start docker compose (postgres, redis, go server)
make dcd-dev      # Stop docker compose
make dcb-dev      # Build docker containers
```

### Backend (Go API - runs inside Docker)
```bash
make test                           # Run Go tests
make lint                           # Run staticcheck
make mock                           # Generate mocks with mockery
make sqlc-generate                  # Generate sqlc code from queries
make generate                       # Full code generation (CUE → OpenAPI YAML → Go)
make migrate-diff NAME=xxx          # Create migration file
make migrate-apply                  # Apply migrations
make seed                           # Run database seeding
```

### Frontend (Next.js)
```bash
cd apps/web
pnpm dev                            # Start Next.js dev server
pnpm lint                           # Run Biome linter
pnpm format                         # Format with Biome
pnpm generate:api                   # Generate API client from OpenAPI spec (orval)
```

## Architecture

### Monorepo Structure
- `apps/api/` - Go backend (Echo framework, runs in Docker)
- `apps/web/` - Next.js 16 frontend with React 19 and TanStack Query

### Backend (apps/api)
- **Framework**: Echo v4 with oapi-codegen for OpenAPI-first development
- **Database**: PostgreSQL with pgx/v5, sqlc for type-safe queries
- **Migrations**: Atlas (`atlas.hcl`, migrations in `db/migrations/`)
- **Auth**: Google OAuth via Goth, Redis session storage
- **Code structure**:
  - `internal/gen/` - Generated API code from OpenAPI
  - `internal/handler/` - HTTP handlers implementing generated interfaces
  - `internal/service/` - Business logic with repository pattern
  - `internal/mapper/` - DTO/entity mappers
  - `db/sqlc/` - Generated database code
  - `openapi/cue/` - CUE source for OpenAPI spec

### Frontend (apps/web)
- **Framework**: Next.js 16 with App Router
- **State/Data**: TanStack Query v5 with Axios
- **Styling**: Tailwind CSS v4, Radix UI components
- **Linting**: Biome
- **API Client**: Generated via Orval from OpenAPI spec (`src/api/generated/`)

### Code Generation Flow
1. Edit OpenAPI spec in `apps/api/openapi/cue/api.cue`
2. Run `make generate` to create `openapi.yaml` and Go server code
3. Run `pnpm generate:api` in `apps/web` to generate TypeScript client

### Pre-commit Hooks (lefthook)
- Runs `pnpm lint` and `make test` before commit
