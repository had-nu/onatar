# Contributing

Obrigado pelo interesse em contribuir para o Onatar. Processo leve, focado em
qualidade e documentação (AGPL-3.0).

## Setup

```bash
make migrate && make seed   # MariaDB local (ver README)
make dev                    # backend :8090 + frontend :5173
```

## Branch & PR

1. Cria uma branch a partir de `main` (`feat/`, `fix/`, `docs/`, `chore/`).
2. Faz commits pequenos e descritivos (ex.: `feat(builder): add undo/redo`).
3. Abre PR para `main`; a CI (`.github/workflows/ci.yml`) corre:
   - Backend: lint, vet, build, testes com cobertura ≥ 70%, seed idempotente,
     gosec, trivy (HIGH/CRITICAL).
   - Frontend: typecheck, vitest, eslint/prettier, build.
   - E2E: Playwright (chromium) contra MariaDB efémera + API :8090.

## Gate local obrigatório

```bash
make lint            # golangci-lint + eslint/prettier
make test            # go test + vitest
cd frontend && npm run build
make build
# opcional, com DB a correr:
TEST_DB_DSN="$DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" make test-integration
```

Regras de estilo:

- Go: `gofmt` + `golangci-lint` (inclui `gosec`); prepared statements sempre.
- Svelte/TS: `prettier` + `eslint`; runes só em `*.svelte.ts` (ver
  ARCHITECTURE.md § runes em módulos).
- Migrações: ficheiro novo em `migrations/` com up+down; nunca editar as antigas.
- Conteúdo: YAML em `data/`; seeds devem ser idempotentes.

## Documentação

- Mudanças de requisitos → atualiza `ONATAR_PROJECT.md` (changelog + PRD).
- Decisões técnicas → `ARCHITECTURE.md`.

## Reportar bugs

Issues com: passos para reproduzir, output esperado vs real, browser/OS, e logs
do dev server se relevante.
