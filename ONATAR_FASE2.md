# Onatar — Fase 2 (v1.1 + v1.2) Product Requirements

> Extensão pós-MVP do [ONATAR_PROJECT.md](ONATAR_PROJECT.md) (v1.0.0).
> v1.1: Import de ficha D&D Beyond (RF-08). v1.2: Combat tracker (RF-09).
> Beta (auth + persistência cloud, RF-11/RF-12) é **Fase 3** — fora deste documento.
> Licença: AGPL-3.0

---

## Changelog

| # | Problema | Correção |
|---|----------|----------|
| 1 | Fase 1: `ci.yml` pina `aquasecurity/trivy-action@0.30.0`, que não existe → CI backend falha em setup | **Pin para `v0.36.0`** (latest real). |
| 2 | Fase 1: `golangci-lint-action@v6` não suporta golangci-lint v2.1 (`invalid version string`) | **Upgrade para `@v9`** (suporta v2). |
| 2 | v1.1: parser de PDF no servidor (Go) seria pesado e desalinhado com o stack (FE já usa jsPDF/html2canvas) | **Client-side com `pdfjs-dist`**; zero mudanças no backend. |
| 3 | v1.1: layout do PDF D&D Beyond muda entre revisões → parser por coordenadas seria frágil | **Parser tolerante por keywords** (STR/DEX/CON/…, nomes de spells/feats, proficiências) + **step de revisão manual** obrigatório antes de criar o draft. |
| 4 | v1.1: fixtures de PDF DDB reais são copyright WotC / dados de utilizador | **Fixtures sintéticas** geradas para teste (texto 2 colunas modelado no layout DDB, sem conteúdo proprietário). |
| 5 | v1.2: combat tracker duplicaria estado HP/conditions da ficha | **Sincronização com `Character.live`**: combatente ligado a personagem partilha HP/conditions; não-ligado guarda estado próprio. |
| 6 | v1.2: turnos e rounds confundiam-se no modelo de dados | **`round` no session + `turnIndex` no session**: wrap de `turnIndex` incrementa `round`. |
| 7 | v1.2: `save()` substitui o objeto da sessão → referências antigas ficam stale | **Views/tests re-leem via `getSession(id)`** após cada mutação (padrão documentado no store). |
| 8 | v1.1: chaves JSON de stats variam (`str`/`strength`/`constitution`) | **`abilityFromKey`** normaliza qualquer variante; `parseDDBJSON` tolera `stats|scores|abilities` aninhados. |

---

## 1. Visão & Objectivos

- **v1.1 (Sprint 6)** — RF-08: importar uma ficha exportada do D&D Beyond (PDF) e
  converter num draft editável no wizard, com revisão manual dos valores extraídos.
  Bónus: import JSON no formato D&D Beyond.
- **v1.2 (Sprint 7)** — RF-09: combat tracker local — initiative, HP, conditions
  por combate, integrado com a ficha existente.

Sem alterações ao backend (API continua read-only: `/health`, `/content`, `/build`).
Persistência continua em localStorage (consistente com RNF-03 offline parcial).

## 2. Requisitos

### 2.1 v1.1 — Import PDF D&D Beyond (RF-08)

| ID | Requisito | Critério de aceitação |
|----|-----------|------------------------|
| F2-01 | Selecionar/arrastar um PDF D&D Beyond | `file` picker + drag-drop; erro claro para ficheiros inválidos (não-PDF, > 5 MiB) |
| F2-02 | Extração de dados do PDF | Classe(s)+nível, species, background, ability scores, spells, feats, skills, HP/AC quando presentes |
| F2-03 | Revisão manual antes de criar | Preview editável dos valores extraídos; "Criar personagem" desativado até revisão ok |
| F2-04 | Conversão para draft válido | Resultado alimenta `buildDraft`/wizard sem passos em branco obrigatórios |
| F2-05 | Import JSON D&D Beyond | Adapter para JSON (mesma revisão manual) |

Não-requisitos: import automático 100% fiel (objetivo é um **ponto de partida** fiável);
suporte a PDFs de outros sistemas.

### 2.2 v1.2 — Combat Tracker (RF-09)

| ID | Requisito | Critério de aceitação |
|----|-----------|------------------------|
| F2-10 | Criar/abrir/terminar combate | Sessões persistem em localStorage; listagem de combates |
| F2-11 | Adicionar combatentes | Da lista de personagens/NPCs ou avulso (nome + HP + iniciativa) |
| F2-12 | Rolar/definir iniciativa e ordenar | Sort desc; tiebreak manual ou pelo valor inserido |
| F2-13 | Turno atual + wrap de round | next/prev; `round` incrementa ao voltar ao início |
| F2-14 | Dano/cura por combatente | Reusa o padrão do HP stepper da ficha |
| F2-15 | Condições por combatente | Reusa `CONDITIONS` (chips toggle) |
| F2-16 | Sincronização com a ficha | Combatente ligado a `Character` atualiza `live` (HP/conditions) da ficha |

## 3. Arquitetura / Decisões

### 3.1 v1.1 — Import PDF (client-side)

```
PDF (D&D Beyond)
  → pdfjs-dist (getDocument → items com texto + transform)
  → src/lib/import/ddb.ts
       extractDDDPDF(file) → ParsedDDB   (normalizado, tolerante a layout)
       toBuildRequest(parsed, content) → BuildRequest
  → Import.svelte (preview editável) → buildDraft → Character
```

- `pdfjs-dist`: worker via `?url` + `GlobalWorkerOptions.workerSrc` (Vite).
- Parser de texto: concatena linhas por coluna usando `transform[4]` (x) para
  separar colunas; procura chaves por keyword. Não depende de coordenadas exatas.
- Resolução de IDs (class/spell/feat ids): match por nome case-insensitive contra
  o `Content` cacheado (`content.svelte.ts`); desconhecidos ficam em
  `unmapped` e são mostrados na revisão.
- Ficheiros novos: `src/lib/import/ddb.ts`, `src/lib/import/ddb.test.ts`,
  `src/lib/views/Import.svelte`, `src/lib/views/Import.test.ts`.

### 3.2 v1.2 — Combat Tracker

```
src/lib/combat.svelte.ts   — store (sessões + lógica de turnos)
src/lib/views/Combat.svelte — UI
rota /combat + nav
```

- `CombatSession { id, name, round, turnIndex, combatants }`
- `Combatant { id, characterId?, name, initiative, hpCurrent, hpMax, conditions }`
- Sincronização: `linkCharacter(c)` copia `live` e escreve de volta em
  `characters.updateLive(id, { hpCurrent, conditions })` nas ações.
- `CONDITIONS`, `emptyLive` reutilizados de `types.ts`.

## 4. Roadmap de Sprints

### Sprint 6 — v1.1 Import PDF ✓
1. Dep `pdfjs-dist` + worker config no Vite.
2. `import/ddb.ts`: `extractDDDPDF` + `toBuildRequest` + testes de fixtures.
3. View `Import.svelte` + rota `/import` + link na nav.
4. Adapter JSON DDB.
5. Gate: typecheck/lint/test/build/e2e + changelog + commit.

### Sprint 7 — v1.2 Combat Tracker ✓
1. `combat.svelte.ts` + testes (rotação, wrap, sort, dano/cura, sync).
2. View `Combat.svelte` + rota `/combat` + link na nav.
3. Gate: typecheck/lint/test/build/e2e + changelog + commit.

## 5. Definition of Done

Herda o DoD do PRD v1 (§11) mais:

- v1.1: parser cobre ≥ 90% dos campos do formato DDB sintético; revisão manual
  obrigatória; unit tests do mapper; fixture sintética no repo.
- v1.2: store cobre rotação/wrap/sort/dano/conditions; sincronização com `live`
  testada; e2e smoke no fluxo principal.
- Cobertura FE e backend mantêm os limiares existentes (RNF-05 ≥ 70% backend).

---

*Documento gerado em 2026-08-01. Extensão pós-MVP do PRD v1.*
