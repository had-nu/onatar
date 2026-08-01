# Onatar

Plataforma web de criação e gestão de fichas de personagem D&D 2024 (5.5e).
Inspirado no D&D Beyond, mas focado em acessibilidade para novos jogadores e
simplicidade para mestres. Licença: AGPL-3.0.

> Requisitos e arquitetura: [ONATAR_PROJECT.md](ONATAR_PROJECT.md) e [ARCHITECTURE.md](ARCHITECTURE.md)

## Stack

| Camada   | Tecnologia              |
|----------|-------------------------|
| Backend  | Go 1.23+ (chi)          |
| Frontend | Svelte 5 + Vite         |
| Database | MariaDB 11              |
| Migrações| golang-migrate          |

## Quickstart (dev)

1. MariaDB local + base `onatar` (ver Makefile).
2. `cp .env.example .env`
3. `make migrate`   — aplica migrations
4. `make seed`      — popula conteúdo de regras a partir de `data/`
5. `make dev`       — backend (Go, :8090) + frontend (Vite, :5173, proxy `/api`)

O server e o seed leem `.env` automaticamente (godotenv); `go run ./cmd/server`
funciona no root do repo sem exportar variáveis. Env do OS prevalece.

## API (MVP — Sprint 1)

| Método | Path             | Descrição                                    |
|--------|------------------|----------------------------------------------|
| GET    | `/health`        | Health check + versão                        |
| GET    | `/api/v1/content`| Todo o conteúdo de regras (de MariaDB)       |
| POST   | `/api/v1/build`  | Calcula ficha derivada (draft → sheet)       |

- `/build` limitado a **10 req/min por IP** (token bucket), timeout 5s.
- Erros seguem o schema `{"error":{"code","message","details"}}` (PRD §3.5).

Exemplo:

```bash
curl http://localhost:8090/api/v1/content
curl -X POST http://localhost:8090/api/v1/build -H 'Content-Type: application/json' -d '{
  "name": "Onatar",
  "classes": [{"id": "sorcerer", "level": 6, "subclassId": "aberrant-sorcery"}],
  "speciesId": "tiefling", "backgroundId": "sage",
  "abilityScores": {"STR":8,"DEX":14,"CON":16,"INT":10,"WIS":12,"CHA":18},
  "abilityMethod": "point-buy",
  "spells": ["magic-missile","shield"], "feats": ["war-caster"]
}'
```

## Makefile

- `make migrate` — aplica migrations (up)
- `make rollback` — reverte migrations (down 1)
- `make seed` — seed idempotente de `data/**/*.yaml` → MariaDB
- `make dev` — backend + frontend em dev
- `make test` — `go test ./cmd/... ./internal/...` + `vitest run`
- `make test-integration` — testes HTTP de integração (requer `TEST_DB_DSN`)
- `make e2e` — Playwright E2E (requer backend dev em `:8090`)
- `make lint` — golangci-lint + eslint/prettier
- `make build` — build do backend
- `make release` — cross-compile Linux amd64/arm64

Testes de integração (`internal/store`) requerem `TEST_DB_DSN` (DSN MySQL);
sem ela são saltados. A CI fornece uma MariaDB efémera.

## Roadmap

MVP (v1.0.0) → v1.1 (import PDF D&D Beyond) → v1.2 (combate tracker) → Beta (auth + cloud). Detalhes no PRD §10.
