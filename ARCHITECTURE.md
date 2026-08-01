# Onatar — Architecture

> Actualizado após o fecho dos Sprints 0–5 (v1.0.0).

## Visão de alto nível

```
┌─────────────┐   /api/v1/* (proxy)   ┌──────────────┐   SQL   ┌─────────┐
│ Svelte 5 SPA │ ───────────────────▶ │ Go (chi) API  │ ──────▶ │ MariaDB │
│  (:5173 dev) │                      │   (:8090)     │         └─────────┘
└─────────────┘                       └──────────────┘
   hash router, localStorage            /health, /content, /build
```

- **SPA** sem framework de rotas: hash router custom (`lib/router.svelte.ts`).
- **Persistência local**: personagens, campanhas, tema, conteúdo cacheado — tudo
  em `localStorage`. A API é apenas de leitura de regras (`/content`) e cálculo
  (`/build`).
- **Service Worker** (`frontend/public/sw.js`): precache da app shell e
  stale-while-revalidate para `/api/v1/content` (RNF-03 — offline parcial).

## Backend (Go)

```
cmd/server        — entrypoint HTTP (godotenv, config, store, httpapi)
cmd/seed          — parser data/**/*.yaml → MariaDB (idempotente)
internal/config   — env: DB_*, HTTP_ADDR
internal/store    — queries + seed (prepared statements sempre)
internal/content  — structs de conteúdo /content
internal/build    — rules engine (PB, HP, AC, slots, features, validação)
internal/httpapi  — rotas chi, handlers, rate limiter (token bucket/IP)
internal/integration — testes HTTP de ponta-a-ponta contra MariaDB real
migrations/       — golang-migrate (up/down)
data/             — fonte de verdade dos ficheiros de regras (.yaml)
```

Regras chave do motor de build (`internal/build`):

| Regra | Fórmula |
|-------|---------|
| Proficiency bonus | `PB = (level + 7) / 4` (arredondado para baixo) |
| HP | `HP = dadoMax + CON mod + (avgDado + CON mod) × (level − 1)` |
| AC | `AC = 10 + DEX mod` (regra v1) |
| Spell slots | tabela full-caster por nível de conjurador |
| Escolhas pendentes | `data.choose` + `options` → `pendingChoices` |

Erros API: `{"error":{"code","message","details"}}`. `/build` é limitado a
10 req/min/IP (token bucket) com timeout de 5s e body ≤ 1 MiB.

## Frontend (Svelte 5)

```
src/lib/types.ts          — contract types (GET /content, POST /build)
src/lib/router.svelte.ts  — hash router
src/lib/theme.svelte.ts   — tema light/dark/system (data-theme)
src/lib/content.svelte.ts — cache /content (memória + localStorage, TTL)
src/lib/characters.svelte.ts — CRUD local + build + live sheet
src/lib/campaigns.svelte.ts  — campanhas mínimas (id + name)
src/lib/builder.svelte.ts — wizard store com undo/redo + validação
src/lib/box.svelte.ts     — container reactivo (`{ value }`) para runes em módulos
src/lib/export.ts         — export/import JSON + PDF (jsPDF + html2canvas)
src/lib/views/            — Landing, Characters, CharacterView, Content,
                            Builder (+ steps), Campaigns
```

### Nota — runes em módulos

`$state`/`$derived` só são compilados em ficheiros `.svelte.ts`/`.svelte.js`
(plugin `compileModule` do `@sveltejs/vite-plugin-svelte`). Ficheiros `.ts`
planos **não** são transformados. Stores com estado reactivo vivem em
`*.svelte.ts` e importam-se com `'./foo.svelte'` (o Vite resolve para
`foo.svelte.ts`). Não é permitido reassignar estado exportado — usa-se o helper
`box()` (`export const x = box(initial)` → `x.value`).

### Wizard (PRD §3.4)

CLASS → BACKGROUND → SPECIES → ABILITIES → EQUIPMENT → REVIEW.
Cada step valida antes de avançar; undo/redo restaura também a posição do step;
a sidebar de pré-visualização faz `POST /build` debounced.

## Qualidade & CI (Sprint 5)

`.github/workflows/ci.yml` (por PR + push main):

| Job | Passos |
|-----|--------|
| backend | golangci-lint, vet, build, migrate+test (coverage ≥ 70% via `internal.out`), seed idempotente, gosec, trivy fs (HIGH/CRITICAL) |
| frontend | svelte-check/tsc, vitest, eslint+prettier, vite build |
| e2e | MariaDB efémera → migrate → seed → server :8090 → Playwright (chromium) |

`.github/workflows/release.yml` (tag `v*`): `make release` (Linux amd64/arm64)
+ GitHub Release.

## Scripts locais

- `make test` — Go unit tests + vitest
- `make test-integration` — integração HTTP (requer `TEST_DB_DSN`)
- `make lint` — golangci-lint + eslint/prettier
- `make build`, `make release`
- `frontend: npm run test:e2e` — Playwright (requer backend em :8090)

## Decisões registadas

Ver Changelog em [ONATAR_PROJECT.md](ONATAR_PROJECT.md) (linhas 1–21).
Destaques: campanhas/NPCs fora do MVP (beta), offline parcial (RNF-03),
`AC = 10 + DEX`, `suggestedX` só nos ficheiros fonte, seeds idempotentes.
