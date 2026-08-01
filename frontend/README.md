# Onatar — Frontend

SPA Svelte 5 + Vite (TypeScript strict) servida como static files pelo backend Go.

## Comandos

- `npm run dev` — Vite dev server (:5173, proxy `/api` → Go :8090)
- `npm run typecheck` — svelte-check + tsc
- `npm run test` — Vitest + @testing-library/svelte
- `npm run lint` — ESLint + Prettier
- `npm run build` — build de produção para `dist/`

## Notas

- Svelte usa `svelte-check` para type-check (equivalente idiomático ao `tsc --noEmit` do PRD §7).
- Service worker / offline cache: Sprint 2.
