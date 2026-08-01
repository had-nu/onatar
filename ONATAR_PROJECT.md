# Onatar — Product Requirements Document v1

> Plataforma web de criação e gestão de fichas de personagem D&D 2024 (5.5e).  
> Inspirado no D&D Beyond, mas focado em acessibilidade para novos jogadores e simplicidade para mestres.  
> Licença: AGPL-3.0

---

## Changelog

| # | Problema | Correção |
|---|----------|----------|
| 1 | RF-05/RF-06 (campanhas/NPCs) contradiziam entrevista C ("fora do MVP") | **Campanha como entidade mínima** (apenas `id` + `name`) entra no MVP como estrutura de dados para futuro. **Gestão de campanha** (DM ver fichas, vincular, etc.) é beta. **NPCs** entram no MVP como flag `is_npc` no mesmo wizard — zero custo adicional. |
| 2 | API listava CRUD `/characters` e `/campaigns` no MVP, mas persistência é localStorage | **Removidos do MVP.** API no MVP serve apenas `GET /content` e `POST /build`. CRUD de personagens/campanhas move-se para beta (com auth + MariaDB persistence). |
| 3 | Ambiguidade entre `.md/.yaml` fonte e `data JSON` no DB | Esclarecido: `.md/.csv/.yaml` são **fonte de verdade** (source of truth). Script de seed parseia para structs Go → INSERT MariaDB. `data JSON` é o formato de armazenamento interno para flexibilidade de schema. |
| 4 | RNF-03 "100% offline" contradiz `POST /build` server-side | **Ajustado:** "Offline parcial" — frontend cacheia conteúdo via Service Worker e edita drafts em localStorage sem rede. Cálculo de build (`POST /build`) requer conectividade. Rules engine client-side é pós-MVP. |
| 5 | CI/DoD mencionava `npm run test` mas pipeline não incluía | **Adicionado:** `vitest run` (Svelte 5) ao CI desde Sprint 0. Playwright E2E permanece Sprint 5. |
| 6 | Faltava definição do engine de sugestões, schema de erro, schema de `/content`, modelo de campanha, steps do wizard | **Todos adicionados** nas secções correspondentes. |
| 7 | Ordem do wizard em §3.4 tinha SPECIES antes de BACKGROUND | **Corrigido:** CLASS → BACKGROUND → SPECIES → ABILITIES → EQUIPMENT → REVIEW. |
| 8 | Relação campanha↔personagem duplicada (`character_ids` JSON + `campaign_id FK`) | **Canónica via FK:** removido `character_ids` do ERD. |
| 9 | Roadmap Beta incluía itens v1.1/v1.2 | **Movidos** para secção "v1.1 / v1.2 — Pós-MVP (antes do Beta)". |
| 10 | Campos de sugestão sem fonte definida — risco de drift no re-seed | **§8.1** define o schema dos ficheiros fonte com `suggestedSpecies`, `suggestedBackgrounds`, `recommendedSpells` desde o início. |
| 11 | Sprint 1: `subclassLevel` e `primaryAbility` perdiam-se no MariaDB | **Migração 0002** adiciona `classes.subclass_level`; `primaryAbility` é fundido no `data JSON` pelo seed. `GET /content` devolve `subclassLevel` top-level (spec §3.5). |
| 12 | Sprint 1: `go run ./cmd/server` falhava sem `.env` exportado | **godotenv** carregado no server e no seed: `.env` é lido no working dir; env do OS prevalece (CI inalterado). |
| 13 | Sprint 1: `stop.sh` orfanava o binário de `go run` (matava só o wrapper) | **Reescrito** para matar por porta via `ss` (backend `:8090`, frontend `:5173`). |
| 14 | Sprint 1: AC do exemplo §3.5 (15 com DEX 14) inconsistente com a regra 5e | **Regra v1:** `AC = 10 + DEX mod` (exemplo é ilustrativo). Bónus de species aplicados no request (scores finais). Point-buy budget enforcement fica para o wizard (Sprint 3). |
| 15 | Sprint 2: `$state` em módulos `.ts` não é compilado (runes só em `.svelte.js/.svelte.ts`) | **Stores movidos** para `content.svelte.ts` / `characters.svelte.ts` (box helper `box.svelte.ts`). Testes de store usam `.svelte` no import. `eslint.config.js` ganha parser TS para `*.svelte.{js,ts}`. |
| 16 | Sprint 2: `export let x = $state()` é reassignado (proibido por runes) | **Padrão `box()`**: getter/setter sobre `$state` local; stores expõem `x.value`. Router já usava mutação de propriedade (válido). |
| 17 | Sprint 2: svelte-check warning "state_referenced_locally" no CharacterView | **Inicialização movida** para `onMount` (lê store e deriva `sheet`/`status`), evitando capturar `$derived` no inicializador de `$state`. |
| 18 | Sprint 3: `$state(...)` não é permitido dentro de object literal (`state_invalid_placement`) | **Helper `box()`** em `box.svelte.ts`: `$state` local + getter/setter `value`. Todos os stores usam `x.value`; testado. |
| 19 | Sprint 3: ESLint não parseava `.svelte.ts` / `new Set` em component | **Parser TS** adicionado para `*.svelte.{js,ts}` no `eslint.config.js`; `Set` substituído por array dedup em `EquipmentStep`. |
| 20 | Sprint 3: `.ts` com runes nunca é compilado (só `.svelte.ts` é) | **Stores renomeados** `content.ts`→`content.svelte.ts`, `characters.ts`→`characters.svelte.ts`; imports atualizados em views/tests. Testes de store ficam `.test.ts` (sem infix `.svelte.`). |
| 21 | Sprint 3: undo/redo perdia a posição do step | **Snapshot inclui `stepIndex`**; desfazer/refazer restaura o step em que a mudança foi feita. |
| 22 | Sprint 5: wizard E2E nunca chegava ao `Revisão` (wizard tem 6 steps, spec só percorria 5) | **Spec corrigida**: passo de Equipamento adicionado antes de Revisão; link 'Personagens' passou a `exact: true` (nav + CTA do hero colidiam em strict mode). |
| 23 | Sprint 5: `vitest run` falhava a suíte `e2e/wizard.spec.ts` (Playwright) | **`exclude: ['e2e/**', …]`** no `test` do `vite.config.ts` — Vitest e Playwright partilham glob `*.spec.ts`, agora separados. |
| 24 | Sprint 5: `prettier --check .` falhava em `test-results/.last-run.json` (artefato Playwright) | **`.prettierignore`** no frontend ignora `test-results`, `playwright-report`, `dist`, `node_modules`. |
| 25 | Sprint 5: CI só cobria unit tests Go/FE; integração e E2E não corriam na pipeline | **Job `e2e`** (MariaDB efémera → migrate → seed → server :8090 → Playwright chromium) e **trivy fs** (HIGH/CRITICAL) adicionados ao `ci.yml`; **`release.yml`** publica binários amd64/arm64 em tag `v*`. |
| 26 | Sprint 5: integração HTTP end-to-end não tinha onde correr localmente | **`internal/integration`** (health, content, build end-to-end; skip sem `TEST_DB_DSN`) + target `make test-integration`. |
| 27 | Sprint 5: docs de arquitetura/contribuição ausentes | **`ARCHITECTURE.md`** (diagramas, runes em módulos, qualidade/CI) e **`CONTRIBUTING.md`** (gate local, estilo, PR) criados; README atualizado. |

---

## 1. Visão & Objectivos

- Tornar D&D 2024 acessível a quem nunca jogou — explicação contextual em cada escolha
- Ficha dinâmica e interativa (tablet/smartphone) com tracking de combate
- DM pode criar NPCs/monstros e ter visibilidade mínima das fichas dos jogadores (beta)
- 100% open source (AGPL-3), self-hostable, sem paywall de conteúdo

---

## 2. Público-alvo

| Perfil | Necessidades | Prioridade |
|--------|-------------|------------|
| Jogador novo | Wizard guiado com explicações de cada escolha; sugestões de builds | P0 |
| Jogador experiente | Criação rápida (express); ficha dinâmica; export PDF/JSON | P0 |
| Mestre (DM) | Criar NPCs/monstros (MVP); ver fichas dos jogadores (beta); gestão de campanha (beta) | P1 (MVP: NPCs) / P2 (beta: campanha) |

---

## 3. Requisitos Funcionais

### 3.1 MVP (v1.0.0)

| ID | Requisito | Prioridade |
|----|-----------|------------|
| RF-01 | Character builder wizard passo-a-passo (ver §3.4 para fluxo exato) | P0 |
| RF-02 | Visualizador de ficha web dinâmica e interativa (HP, spell slots, conditions, resources) | P0 |
| RF-03 | Lista de personagens salvos em localStorage (CRUD local) | P0 |
| RF-04 | Export de ficha para PDF e JSON | P0 |
| RF-05 | Campanha mínima (entidade: `id` + `name`) — apenas estrutura de dados, sem gestão. Relação com personagens via `characters.campaign_id` (ver §8) | P1 |
| RF-06 | DM pode criar fichas de NPCs/monstros (mesmo wizard, flag `is_npc` no draft) | P1 |
| RF-07 | Sugestões contextuais em cada step (ver §3.3) | P1 |

### 3.2 Pós-MVP

| ID | Requisito | Versão |
|----|-----------|--------|
| RF-08 | Import de ficha D&D Beyond (PDF parser) | v1.1 |
| RF-09 | Combate tracker (initiative, HP, conditions) | v1.2 |
| RF-10 | Integração VTT (Roll20, Foundry) | v2.0 |
| RF-11 | Gestão de campanha (DM ver fichas dos jogadores, vincular personagens, narrativa/arcos) | Beta |
| RF-12 | Autenticação (GitHub OAuth) + persistência cloud (MariaDB) | Beta |

### 3.3 Suggestions Engine (RF-07)

**Não é ML.** É um sistema de lookup baseado em regras D&D 2024:

- **Classe → Species sugeridos:** cada classe tem `suggested_species` (array de IDs) no JSON da tabela `classes`. Ex: Sorcerer Aberrant → `["kalashtar", "warforged", "tiefling"]`. Fonte: guias dos suplementos.
- **Classe → Backgrounds sugeridos:** `suggested_backgrounds` no JSON da classe.
- **Classe → Spells sugeridas:** por subclass, lista de `recommended_spells` (IDs) para cada nível.
- **Background → Skills:** descrição textual do background indica quais skills ganha.
- **Species → Traits:** exibição inline das traits ao selecionar species.

O frontend consulta estas listas no `data JSON` da entidade e exibe como "Sugestões" ou "Popular choices" em cada step. Nenhuma lógica complexa — apenas denormalização no JSON.

> **Nota de proveniência:** `suggestedSpecies`, `suggestedBackgrounds` e `recommendedSpells` vivem nos **ficheiros fonte** (`.yaml`), não no `data JSON`. O script de seed extrai-os para a DB — nunca são editados à mão. Ver §8.1.

### 3.4 Wizard Steps (RF-01) — Fluxo Exato

```
VAULT {
  CLASS(
    [CLASS FEATURES],
    [SUBCLASS — se level >= subclass_level],
    [SPELLS — se spellcaster],
    [SPECIAL TRAITS — ex: Metamagic, Invocations]
  )
  → BACKGROUND
  → SPECIES(
    [SPECIES TRAITS],
    [VARIANT — se aplicável]
  )
  → ABILITIES (
    [POINT BUY],
    [ROLLED — 4d6 drop lowest],
    [STANDARD ARRAY — 15,14,13,12,10,8]
  )
  → EQUIPMENT (
    [CLASS STARTING GEAR],
    [BACKGROUND GEAR]
  )
  → REVIEW & SAVE (SHEET — web ou PDF)
}
```

**Regras de navegação:**
- O utilizador pode avançar para o próximo step apenas se o step atual estiver válido.
- Pode retroceder a qualquer step anterior para editar.
- O step "Review" renderiza a ficha — web sheet interativa e/ou export PDF (jsPDF + html2canvas) — com preview live calculado via `POST /build`.
- "Save" persiste em localStorage (MVP) ou envia para `POST /characters` (beta).

---

## 4. Requisitos Não-funcionais

| ID | Requisito | Métrica |
|----|-----------|---------|
| RNF-01 | Tempo de resposta da API < 200ms (p95) | < 200ms |
| RNF-02 | First Contentful Paint < 1.5s no 4G | < 1.5s |
| RNF-03 | Offline parcial: frontend cacheia conteúdo (Service Worker) e edita drafts em localStorage. `POST /build` requer rede. | Sim |
| RNF-04 | Self-hostable em VPS com 1 vCPU / 1GB RAM | Sim |
| RNF-05 | Cobertura de testes backend ≥ 70% | ≥ 70% |
| RNF-06 | Sem dependências de serviços cloud proprietários | Sim |
| RNF-07 | AGPL-3 licenciamento | Sim |
| RNF-08 | Zero dados pessoais no MVP (sem GDPR) | Sim |

---

## 5. Stack Técnica

| Camada | Tecnologia | Justificação |
|--------|-----------|--------------|
| Backend | Go 1.23+ (chi router) | Performance, simplicidade, manutenção |
| Frontend | Svelte 5 + Vite | Menos boilerplate que React, reatividade nativa, ideal para forms complexos |
| Testes FE | Vitest + @testing-library/svelte | Testes unitários desde Sprint 0 |
| Database | MariaDB 11 | Leve, MySQL-compatível, boa performance em VPS |
| Migrations | golang-migrate | SQL puro, versionado, reversível |
| Auth (beta) | GitHub OAuth | Zero gestão de passwords |
| CI/CD | GitHub Actions | Integrado, gratuito para open source |
| Deploy | Binary + systemd (VPS) | Sem Docker no MVP — simplicidade máxima |
| Logs | Structured stdout (slog) | Simples, parseável por qualquer agente |

---

## 6. C4 Model — Context & Containers

### C1 — System Context

```
┌─────────────┐      HTTPS       ┌─────────────────────────────────────┐      SSH/systemd      ┌─────────────┐
│  Player/DM  │ ───────────────→ │           Onatar Platform           │ ←────────────────── │   Admin     │
│  (Browser)  │                  │  D&D 2024 Character Builder & Sheet │                     │ (self-host) │
└─────────────┘                  └─────────────────────────────────────┘                     └─────────────┘
                                          ↑
                                          │ (seed)
                                     ┌────────────────┐
                                     │ D&D 2024       │
                                     │ Content        │
                                     │ (.md/.csv/.yaml)│
                                     └────────────────┘
```

### C2 — Container Diagram (MVP)

```
┌─────────┐         fetch /api/*         ┌─────────────┐         SQL/TCP         ┌───────────┐
│ Browser │  ─────────────────────────→  │   Go API    │  ──────────────────→  │  MariaDB  │
│ Svelte  │                              │  REST/JSON  │                       │ Rules     │
│   5     │         GET / (static)       │             │                       │ Content   │
└─────────┘  ─────────────────────────→  └─────────────┘                       └───────────┘
       ↑                                              │
       │                                              │
       └──────────────────────────────────────────────┘
              localStorage (drafts, characters)
```

**Nota arquitetural:** No MVP, o frontend é uma SPA Svelte servida como static files pelo próprio Go binary (embed + `http.FileServer`). O Go API serve apenas `GET /content` e `POST /build`. localStorage guarda drafts e personagens. MariaDB guarda apenas conteúdo de regras. No beta, MariaDB passa a guardar também personagens e campanhas; auth via GitHub OAuth.

---

## 7. Threat Model — STRIDE

| Categoria | Ameaça | Severidade | Mitigação |
|-----------|--------|------------|-----------|
| **Spoofing** | Sem auth no MVP, qualquer um pode aceder à API se souber o IP | Média | VPS com firewall (ufw) apenas porta 443; auth OAuth no beta |
| **Tampering** | localStorage é client-side — utilizador pode manipular JSON da ficha | Média | Validar schema no backend (POST /api/build); rejeitar dados inválidos; checksum de integridade |
| **Repudiation** | Sem auth, não há accountability de quem criou/modificou o quê | Média | Logs de IP + timestamp no beta; audit trail com user_id |
| **Information Disclosure** | Stack traces em erros 500; .env exposto; directory listing | Alta | Erros genéricos em produção; validar build não inclui .env; desativar directory listing no static server |
| **Denial of Service** | POST /api/build pesado (cálculo de ficha); spam de requests | Média | Rate limiting (token bucket, 10 req/min por IP); timeout de 5s no build; input size limits (max 10 classes, max 20 spells) |
| **Elevation of Privilege** | SQL injection nos parâmetros de query; path traversal em exports | Alta | Prepared statements em 100% das queries; validar path de export; sem `fmt.Sprintf` em SQL |

### CI/CD Pipeline (obrigatória em cada PR)

1. `go test ./...` — unitários com cobertura ≥ 70%
2. `gosec ./...` — scan SAST Go (CWE top 25)
3. `go vet ./...` — análise estática
4. `vitest run` — testes unitários Svelte 5
5. `tsc --noEmit` — type-check frontend
6. `npm run lint` — ESLint + Prettier
7. Sem merge sem review de 1 pessoa (branch protection)

---

## 8. Entity-Relationship Diagram

### Entidades

```
┌─────────────────┐     1:N     ┌─────────────────┐
│   campaigns     │────────────→│   characters    │
│  id (PK)        │             │  id (PK)        │
│  name           │             │  name           │
│  created_at     │             │  campaign_id FK │
│  updated_at     │             │  data JSON      │
└─────────────────┘             │  is_npc boolean │
                                │  created_at     │
                                └─────────────────┘

  (MVP: estrutura            (MVP: localStorage
   de dados apenas)           apenas; beta: DB)

┌─────────────────┐     1:N     ┌─────────────────┐     N:1     ┌─────────────────┐
│    classes      │────────────→│   subclasses    │────────────→│    species      │
│  id (PK)        │             │  id (PK)        │             │  id (PK)        │
│  name           │             │  class_id FK    │             │  name           │
│  hit_die        │             │  name           │             │  data JSON      │
│  spellcaster    │             │  level_required │             └─────────────────┘
│  data JSON      │             │  data JSON      │
└─────────────────┘             └─────────────────┘

┌─────────────────┐             ┌─────────────────┐             ┌─────────────────┐
│  backgrounds    │             │     spells      │             │     feats       │
│  id (PK)        │             │  id (PK)        │             │  id (PK)        │
│  name           │             │  name           │             │  name           │
│  data JSON      │             │  level (0-9)    │             │  prerequisites  │
└─────────────────┘             │  school         │             │  data JSON      │
                                │  data JSON      │             └─────────────────┘
                                └─────────────────┘

┌─────────────────┐
│    features     │
│  id (PK)        │
│  class_id FK    │
│  subclass_id FK │
│  name           │
│  level          │
│  data JSON      │
└─────────────────┘
```

**Nota de design:**
- Tabelas de conteúdo (classes, spells, etc.) usam coluna `data JSON` para flexibilidade — o schema D&D 2024 muda entre suplementos. A normalização estrita seria impraticável.
- `.md/.csv/.yaml` são a **fonte de verdade**. Script de seed (`cmd/seed/`) parseia para structs Go → INSERT MariaDB.
- `campaigns` no MVP é apenas estrutura de dados (não há endpoints CRUD). O frontend pode criar campanhas em localStorage para agrupar personagens.
- A relação campanha↔personagem é **canónica via `characters.campaign_id FK`**. Em localStorage, o draft guarda `campaignId` no próprio personagem; no beta, a FK na MariaDB. Sem representação duplicada.
- `characters` no MVP vive apenas em localStorage. No beta, migra para MariaDB com auth.

### 8.1 Source of Truth — Schema dos ficheiros fonte

`.md/.csv/.yaml` são a fonte de verdade. Os campos de sugestão (`suggestedSpecies`, `suggestedBackgrounds`, `recommendedSpells`) vivem **nos ficheiros fonte desde o início** — nunca são editados diretamente no `data JSON` da DB.

**`data/classes/<id>.yaml`:**

```yaml
id: sorcerer
name: Sorcerer
hitDie: d6
spellcaster: true
subclassLevel: 3
primaryAbility: CHA
suggestedSpecies:   # ← AQUI, na fonte
  - kalashtar
  - tiefling
  - dragonborn
suggestedBackgrounds:   # ← AQUI
  - sage
  - charlatan
  - hermit
data:
  description: "..."
  spellcasting: { ... }
```

**`data/subclasses/<id>.yaml`:**

```yaml
id: aberrant
classId: sorcerer
name: Aberrant Mind
levelRequired: 3
recommendedSpells:   # ← AQUI
  1: [mind-sliver, arms-of-hadar]
  3: [detect-thoughts, hunger-of-hadar]
  5: [telekinesis, synaptic-static]
data:
  description: "..."
```

**`data/species/<id>.yaml`:**

```yaml
id: kalashtar
name: Kalashtar
abilityScores:
  WIS: 2
  CHA: 1
data:
  traits:
    - name: Dual Mind
      description: "..."
    - name: Mental Discipline
      description: "..."
```

O script de seed (`cmd/seed/`) extrai `suggestedSpecies`, `suggestedBackgrounds`, `recommendedSpells` diretamente do YAML e insere no `data JSON` da MariaDB. Re-seed nunca perde estes dados.

---

## 9. API Contract — MVP

### Base URL
```
https://<host>/api/v1
```

### Endpoints (MVP)

| Método | Path | Descrição | Auth |
|--------|------|-----------|------|
| GET | `/health` | Health check + versão | — |
| GET | `/content` | Todo o conteúdo de regras (classes, species, backgrounds, spells, feats, features) | — |
| POST | `/build` | Calcula ficha derivada a partir de draft | — |

### Endpoints (Beta — com auth + DB persistence)

| Método | Path | Descrição |
|--------|------|-----------|
| GET | `/characters` | Lista personagens do utilizador autenticado |
| POST | `/characters` | Cria personagem |
| GET | `/characters/:id` | Obtém personagem |
| PUT | `/characters/:id` | Atualiza personagem |
| DELETE | `/characters/:id` | Remove personagem |
| GET | `/campaigns` | Lista campanhas do utilizador |
| POST | `/campaigns` | Cria campanha |
| GET | `/campaigns/:id` | Obtém campanha + personagens vinculados |
| GET | `/campaigns/:id/characters` | Lista personagens da campanha |

### GET /content — Response Schema

```json
{
  "classes": [
    {
      "id": "sorcerer",
      "name": "Sorcerer",
      "hitDie": "d6",
      "spellcaster": true,
      "subclassLevel": 3,
      "suggestedSpecies": ["kalashtar", "tiefling"],
      "suggestedBackgrounds": ["sage", "charlatan"],
      "data": { "description": "...", "primaryAbility": "CHA", ... }
    }
  ],
  "subclasses": [
    {
      "id": "aberrant",
      "classId": "sorcerer",
      "name": "Aberrant Mind",
      "levelRequired": 3,
      "data": { "description": "...", "spells": [...] }
    }
  ],
  "species": [
    {
      "id": "kalashtar",
      "name": "Kalashtar",
      "data": { "traits": [...], "abilityScores": { "WIS": 2, "CHA": 1 } }
    }
  ],
  "backgrounds": [
    {
      "id": "sage",
      "name": "Sage",
      "data": { "skills": ["arcana", "history"], "feature": "Researcher" }
    }
  ],
  "spells": [
    {
      "id": "magic-missile",
      "name": "Magic Missile",
      "level": 1,
      "school": "evocation",
      "data": { "description": "...", "classes": ["sorcerer", "wizard"] }
    }
  ],
  "feats": [
    {
      "id": "war-caster",
      "name": "War Caster",
      "prerequisites": { "spellcasting": true },
      "data": { "description": "..." }
    }
  ],
  "features": [
    {
      "id": "font-of-magic",
      "classId": "sorcerer",
      "subclassId": null,
      "name": "Font of Magic",
      "level": 2,
      "data": { "description": "..." }
    }
  ]
}
```

### POST /build — Request/Response

**Request Body (BuildRequest):**
```json
{
  "name": "Onatar",
  "classes": [{"id": "sorcerer", "level": 6, "subclassId": "aberrant"}],
  "speciesId": "kalashtar",
  "backgroundId": "sage",
  "abilityScores": {"STR": 8, "DEX": 14, "CON": 16, "INT": 10, "WIS": 12, "CHA": 18},
  "abilityMethod": "point-buy",
  "skills": ["arcana", "insight"],
  "spells": ["magic-missile", "shield"],
  "feats": ["war-caster"],
  "isNpc": false
}
```

**Response 200 (BuildResponse):**
```json
{
  "sheet": {
    "level": 6,
    "hp": { "max": 44, "current": 44 },
    "ac": 15,
    "proficiencyBonus": 3,
    "spellSlots": [4, 3, 3, 0, 0, 0, 0, 0, 0],
    "abilities": { "STR": {"score":8,"mod":-1}, "DEX": {"score":14,"mod":+2}, ... },
    "features": [
      { "name": "Font of Magic", "level": 2, "description": "..." },
      { "name": "Metamagic", "level": 3, "description": "..." }
    ],
    "pendingChoices": [
      { "type": "metamagic", "description": "Choose 2 Metamagic options" }
    ]
  }
}
```

### Error Response Format

Todas as respostas de erro seguem este schema:

```json
{
  "error": {
    "code": "POINT_BUY_EXCEEDED",
    "message": "Point buy budget exceeded: spent 29 of 27 points",
    "details": { "spent": 29, "budget": 27 }
  }
}
```

| Status | Código | Descrição |
|--------|--------|-----------|
| 400 | `INVALID_DRAFT` | Schema de draft inválido (ex: ability score > 20) |
| 400 | `POINT_BUY_EXCEEDED` | Point-buy ultrapassou 27 pontos |
| 400 | `INVALID_CLASS_COMBO` | Multiclasse com prerequisitos não cumpridos |
| 400 | `INVALID_SPELL_SELECTION` | Spell não pertence à lista de classe/nível |
| 422 | `BUILD_ERROR` | Regra de D&D impossibilita build (ex: prerequisites de feat não cumpridos) |
| 429 | `RATE_LIMITED` | Too many requests (limite: 10 req/min por IP) |
| 500 | `INTERNAL_ERROR` | Erro interno (mensagem genérica, log detalhado no servidor) |

---

## 10. Roadmap de Sprints

### Sprint 0 — Setup & Fundação
- Repo GitHub (AGPL-3), CI/CD pipeline (GitHub Actions), linting (golangci-lint, prettier)
- MariaDB schema + migrations (golang-migrate). Seed scripts (`.md/.csv/.yaml` → SQL)
- Svelte 5 + Vite setup com Vitest. Type-check strict (`tsc --noEmit`)

### Sprint 1 — Backend Core
- API REST: `GET /health`, `GET /content`, `POST /build`
- Parser de conteúdo: `.md/.csv/.yaml` → structs Go → INSERT MariaDB
- Testes unitários Go ≥ 70%. Rate limiting (token bucket). Input validation

### Sprint 2 — Frontend Core & Landing
- Svelte 5: landing page, character list (localStorage), basic character sheet viewer (read-only)
- Service worker para cache de conteúdo. Theme toggle (dark/light)
- Integração com `GET /content` — cards de classe, species, background

### Sprint 3 — Builder Wizard
- Step-by-step wizard (ver §3.4): Class → Background → Species → Abilities → Equipment → Review
- Live preview sidebar (HP, AC, spell slots, features) via `POST /build`
- Contextual help em cada escolha (tooltips com descrição + sugestões via RF-07)
- Validation por step (não avança sem completar). Undo/redo no store

### Sprint 4 — Interactive Sheet & Export
- Ficha dinâmica: editar HP current, spell slots usados, conditions, resources
- Export PDF (client-side via jsPDF + html2canvas). Export/import JSON
- Campanha mínima: criar, listar, vincular personagem por ID (localStorage)

### Sprint 5 — DevSecOps Hardening
- gosec + trivy scan na pipeline. Testes de integração na API
- Playwright E2E básico. Documentação (README, CONTRIBUTING, ARCHITECTURE.md)
- Release v1.0.0 tag. Binary cross-compile (Linux amd64/arm64)

### v1.1 / v1.2 — Pós-MVP (antes do Beta)
- v1.1: Import de ficha D&D Beyond (PDF parser) — RF-08
- v1.2: Combate tracker (initiative, HP, conditions) — RF-09

### Beta — Auth & Persistência Cloud
- GitHub OAuth. Personagens e campanhas persistem em MariaDB
- DM pode ver fichas dos jogadores na campanha. Rate limiting por user

---

## 11. Definition of Done (por User Story)

1. Código passa em `go test ./...` (backend) ou `vitest run` (frontend)
2. Type-check passa (`tsc --noEmit` no frontend; compilação Go no backend)
3. Linting passa (`golangci-lint` / `eslint`)
4. gosec não reporta HIGH/CRITICAL
5. PR revisado por 1 pessoa
6. Feature testada manualmente no browser (screenshot no PR opcional)
7. Documentação atualizada (README ou inline) se mudar API ou fluxo de UI

---

*Documento gerado em 2026-07-31. Arquiteto: André Ataíde + Kimi. Stack: Go + Svelte 5 + MariaDB.*
*Versão 1 — corrigido após revisão por agente de código.*
