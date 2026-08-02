# Onatar — Plano Fase 3 (Atualizado)

> Estado do repo em 2026-08-01, 21:40 UTC.  
> Fase 2 (dívida técnica) concluída. Foco: persistência MariaDB + jornada Player completa.  
> Seed de dados disponível em `/home/hadnu/workspace/homelab/dnd-project/dnd-arts/5etools-src/`

---

## Changelog do Estado

### ✅ Concluído na Fase 2 (commits `78d72e4` + `fedc06e` + `6f865fe`)

| Bug/Feature | Estado | Commit |
|-------------|--------|--------|
| F2 — `saveCharacterFromWizard` chama `POST /build` antes de salvar | ✅ Corrigido | `78d72e4` |
| B2 — `clientIP` respeita `X-Forwarded-For` / `X-Real-IP` | ✅ Corrigido | `78d72e4` |
| B5 — Graceful shutdown (`http.Server.Shutdown`) | ✅ Corrigido | `78d72e4` |
| B4 — `abilitiesValid` verifica `POINT_BUY_MAX` (15) | ✅ Corrigido | `78d72e4` |
| Tema D&D — Tiamat font, paleta pergaminho/carmesim/ouro | ✅ Implementado | `78d72e4` |
| Favicon + meta tags + Open Graph | ✅ Implementado | `78d72e4` |
| ReviewStep — async save com loading/error states | ✅ Corrigido | `78d72e4` |
| Test fixes — vitest mock, typecheck, e2e (92/92) | ✅ Corrigido | `78d72e4` + `fedc06e` + `6f865fe` |

### ✅ Já existia antes da Fase 2

| Componente | Estado |
|------------|--------|
| Backend API (`GET /content`, `POST /build`, `GET /health`) | ✅ |
| MariaDB schema + migrations (content tables) | ✅ |
| Seed script (`cmd/seed/`) — parseia `.yaml` → structs → INSERT | ✅ |
| Rules engine (HP, AC, spell slots, features, pending choices) | ✅ |
| Rate limiting (token bucket, 10 req/min) | ✅ |
| Builder wizard — 6 steps (Class, Background, Species, Abilities, Equipment, Review) | ✅ |
| BuilderPreview — live preview sidebar (HP, AC, spell slots) | ✅ |
| CharacterView — ficha interativa (HP, spell slots, conditions, resources, features) | ✅ |
| Character store — localStorage CRUD, export JSON | ✅ |
| Campaign store — estrutura mínima (localStorage) | ✅ |
| Combat tracker (v1.2) | ✅ |
| D&D Beyond PDF import (v1.1) | ✅ |
| Suggestions engine (RF-07) | ✅ |
| Undo/redo (50 snapshots) | ✅ |
| Theme toggle (dark/light) | ✅ |
| Service Worker (cache de conteúdo) | ✅ |

---

## 1. O que falta para Fase 3

### 1.1 Persistência MariaDB (P0)

Atualmente tudo vive em **localStorage**. A Fase 3 migra para MariaDB com autenticação.

| Entidade | Onde vive agora | Onde deve viver | Auth |
|----------|-----------------|-----------------|------|
| Content (classes, spells, etc.) | MariaDB | MariaDB | Público |
| Characters | localStorage | MariaDB | GitHub OAuth |
| Campaigns | localStorage | MariaDB | GitHub OAuth |
| Drafts (builder) | Memória + undo/redo | Memória (temp) | Guest / OAuth |

### 1.2 Autenticação (P0)

| Método | Estado |
|--------|--------|
| GitHub OAuth | ❌ Não implementado |
| Session management (cookies) | ❌ Não implementado |
| Guest mode (localStorage) | ✅ Funcional — manter como fallback |

### 1.3 API Endpoints (P0)

```
GET  /api/v1/content      ✅ Existe
POST /api/v1/build        ✅ Existe
GET  /health              ✅ Existe

GET    /api/v1/characters            ❌ Não existe
POST   /api/v1/characters            ❌ Não existe
GET    /api/v1/characters/:slug      ❌ Não existe
PUT    /api/v1/characters/:slug      ❌ Não existe
DELETE /api/v1/characters/:slug      ❌ Não existe
POST   /api/v1/characters/:slug/live ❌ Não existe

GET    /api/v1/campaigns             ❌ Não existe
POST   /api/v1/campaigns             ❌ Não existe
GET    /api/v1/campaigns/:slug       ❌ Não existe
PUT    /api/v1/campaigns/:slug       ❌ Não existe
DELETE /api/v1/campaigns/:slug       ❌ Não existe

GET  /auth/github                   ❌ Não existe
GET  /auth/github/callback          ❌ Não existe
POST /auth/logout                   ❌ Não existe
GET  /auth/me                       ❌ Não existe
```

### 1.4 Frontend — Auth & Sync (P0)

| Componente | Estado |
|------------|--------|
| Auth store (`auth.svelte.ts`) | ❌ Não existe |
| Login page | ❌ Não existe |
| Shell.svelte — avatar + logout | ❌ Não existe |
| Characters store — sync com API | ❌ Não existe |
| Migration localStorage → API | ❌ Não existe |

### 1.5 Export PDF/JSON (P1)

| Componente | Estado |
|------------|--------|
| Export JSON | ✅ Existe (verificar se funciona) |
| Export PDF | ⚠️ Dependências presentes (jspdf, html2canvas), não verificado end-to-end |

### 1.6 Mobile & A11y (P1)

| Componente | Estado |
|------------|--------|
| Mobile responsive | ❌ Não testado |
| A11y (ARIA, labels, focus) | ❌ Não verificado |
| Loading states (skeletons) | ❌ Não existe |

---

## 2. Schema MariaDB — Novas Tabelas

```sql
-- users (GitHub OAuth)
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    github_id BIGINT UNSIGNED NOT NULL UNIQUE,
    github_login VARCHAR(255) NOT NULL,
    github_email VARCHAR(255),
    avatar_url VARCHAR(500),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_github_id (github_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- characters (persistência cloud)
CREATE TABLE characters (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    slug VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_npc BOOLEAN DEFAULT FALSE,
    draft JSON NOT NULL,
    sheet JSON,
    live JSON,
    campaign_id BIGINT UNSIGNED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_slug (slug),
    INDEX idx_campaign (campaign_id),
    UNIQUE KEY uk_user_slug (user_id, slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- campaigns (mínimo, para agrupar)
CREATE TABLE campaigns (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    slug VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_slug (user_id, slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- sessions (stateless auth)
CREATE TABLE sessions (
    token CHAR(64) PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_expires (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

---

## 3. API Contract — Auth + Characters + Campaigns

### Auth Flow

```
GET  /auth/github          → redirect para GitHub OAuth
GET  /auth/github/callback → valida state, troca code por token, cria user, set cookie
POST /auth/logout          → invalida session, clear cookie
GET  /auth/me              → retorna { id, github_login, avatar_url } ou 401
```

**Cookie:** `session=<token>` — HttpOnly, Secure (produção), SameSite=Lax, Path=/, Max-Age=30 dias

### Characters

```
GET    /api/v1/characters
  → 200 [{ slug, name, isNpc, sheet: { level, hp: { max }, ac }, updatedAt }]

POST   /api/v1/characters
  Body: { name, draft, isNpc? }
  → 201 { slug, name, draft, sheet, live, createdAt, updatedAt }
  → Internamente chama build(draft) para calcular sheet

GET    /api/v1/characters/:slug
  → 200 { slug, name, isNpc, draft, sheet, live, campaignId, createdAt, updatedAt }

PUT    /api/v1/characters/:slug
  Body: { name?, draft?, isNpc?, campaignId? }
  → 200 { ... }
  → Se draft mudou, recalcula sheet automaticamente

DELETE /api/v1/characters/:slug
  → 204

POST   /api/v1/characters/:slug/live
  Body: { hpCurrent?, spellSlotsUsed?, conditions?, resources? }
  → 200 { live }
```

### Campaigns

```
GET    /api/v1/campaigns
  → 200 [{ slug, name, description, createdAt }]

POST   /api/v1/campaigns
  Body: { name, description? }
  → 201 { slug, name, description, createdAt }

GET    /api/v1/campaigns/:slug
  → 200 { slug, name, description, characters: [...], createdAt }

PUT    /api/v1/campaigns/:slug
  Body: { name?, description? }
  → 200 { ... }

DELETE /api/v1/campaigns/:slug
  → 204
```

---

## 4. Seed de Dados — 5etools-src

O utilizador confirmou que `/home/hadnu/workspace/homelab/dnd-project/dnd-arts/5etools-src/` contém tudo necessário para popular a base de dados.

### Estratégia

1. O seed script atual (`cmd/seed/`) parseia `.yaml` → structs Go → INSERT MariaDB
2. Os dados do 5etools-src provavelmente estão em formato JSON (formato nativo do 5etools)
3. Opções:
   - **A:** Converter JSON do 5etools para YAML compatível com o seed existente
   - **B:** Estender o seed script para parsear JSON do 5etools diretamente
   - **C:** Criar um script separado `cmd/seed-5etools/` que popula a partir do 5etools-src

**Recomendação:** Opção B — estender o seed existente para aceitar tanto YAML como JSON. O schema das structs Go já está definido; basta adicionar um parser JSON que mapeie os campos do 5etools para as structs.

### Mapeamento esperado (5etools → Onatar)

| 5etools | Onatar | Notas |
|---------|--------|-------|
| `class.name` | `classes.name` | |
| `class.hd.number` | `classes.hitDie` | ex: `d10` |
| `class.spellcastingAbility` | `classes.spellcaster` | true se existir |
| `class.subclassTitle` | `classes.data.subclassTitle` | ex: "Primal Path" |
| `class.subclasses[].name` | `subclasses.name` | |
| `race.name` | `species.name` | D&D 2024 usa "species" |
| `race.ability` | `species.data.abilityScores` | |
| `race.entries` | `species.data.traits` | |
| `background.name` | `backgrounds.name` | |
| `background.skillProficiencies` | `backgrounds.data.skills` | |
| `background.entries` | `backgrounds.data.feature` | |
| `spell.name` | `spells.name` | |
| `spell.level` | `spells.level` | 0-9 |
| `spell.school` | `spells.school` | |
| `spell.classes.fromClassList` | `spells.data.classes` | |
| `feat.name` | `feats.name` | |
| `feat.prerequisite` | `feats.prerequisites` | |

---

## 5. Sprints Fase 3

### Sprint 3.0 — Foundation (Backend)

| # | Tarefa | Est. |
|---|--------|------|
| 3.0.1 | Migrations: `0002_create_users.up.sql`, `0003_create_characters.up.sql`, `0004_create_campaigns.up.sql`, `0005_create_sessions.up.sql` | 1h |
| 3.0.2 | `internal/auth/` — GitHub OAuth flow (config, handlers, session management) | 4h |
| 3.0.3 | `internal/characters/` — CRUD service + handlers | 3h |
| 3.0.4 | `internal/campaigns/` — CRUD service + handlers | 2h |
| 3.0.5 | Auth middleware — proteger endpoints de characters/campaigns | 1h |
| 3.0.6 | Seed script estendido — parsear JSON do 5etools-src | 3h |
| 3.0.7 | Testes de integração para auth + characters + campaigns | 2h |

### Sprint 3.1 — Frontend Auth & Sync

| # | Tarefa | Est. |
|---|--------|------|
| 3.1.1 | `auth.svelte.ts` — login state, user data, session | 1h |
| 3.1.2 | Login page — "Login with GitHub" button | 30min |
| 3.1.3 | `Shell.svelte` — avatar + nome + logout quando logado | 1h |
| 3.1.4 | Characters store evoluído — sync com API quando logado, localStorage quando guest | 2h |
| 3.1.5 | Migration localStorage → API (quando user faz login pela primeira vez) | 2h |
| 3.1.6 | Campaign store evoluído — sync com API | 1h |

### Sprint 3.2 — Export & Polish

| # | Tarefa | Est. |
|---|--------|------|
| 3.2.1 | Export PDF end-to-end (verificar jspdf + html2canvas) | 2h |
| 3.2.2 | Export JSON end-to-end | 30min |
| 3.2.3 | Loading states (skeletons) durante fetch | 2h |
| 3.2.4 | Error boundaries + toast notifications | 2h |
| 3.2.5 | Mobile responsive (testar 375px, 768px, 1024px) | 3h |
| 3.2.6 | A11y audit (ARIA, labels, focus trap, keyboard nav) | 2h |

### Sprint 3.3 — QA & Release

| # | Tarefa | Est. |
|---|--------|------|
| 3.3.1 | Playwright E2E — jornada Player completa (login → criar → editar → ver ficha) | 4h |
| 3.3.2 | Performance audit (Lighthouse ≥ 80) | 1h |
| 3.3.3 | Documentação: `AUTH.md`, `DEPLOY.md`, `MIGRATION.md` | 2h |
| 3.3.4 | Release v1.3.0 tag + changelog | 30min |

---

## 6. Jornada do Player (User Story)

```
[LANDING PAGE]
    │
    ├── "Criar Personagem" (guest ou logado)
    │   │
    │   └── [BUILDER WIZARD] — 6 steps
    │       ├── Step 1: CLASS (selecionar classe, subclass, spells, feats)
    │       ├── Step 2: BACKGROUND (selecionar, ver skills/feature)
    │       ├── Step 3: SPECIES (selecionar, ver traits, variant)
    │       ├── Step 4: ABILITIES (standard array / point buy / rolled)
    │       ├── Step 5: EQUIPMENT (starting gear)
    │       └── Step 6: REVIEW (preview live + nome + salvar)
    │           │
    │           └── "Salvar" → POST /api/v1/characters (logado)
    │               ou localStorage (guest)
    │
    ├── "Meus Personagens"
    │   │
    │   └── [CHARACTER LIST]
    │       ├── Cards com nome, classe, nível, thumbnail
    │       ├── Filtro (NPCs, campanha, ordenar)
    │       └── Ações: Editar, Duplicar, Exportar (PDF/JSON), Apagar
    │
    └── "Ficha" (clicar num personagem)
        │
        └── [CHARACTER SHEET]
            ├── Header: Nome, Classe, Nível, AC, HP
            ├── Abilities (STR, DEX, CON, INT, WIS, CHA) com mods
            ├── Skills (proficiência highlight)
            ├── Features & Traits
            ├── Spells (preparados, slots, escola)
            ├── Equipment
            └── Live editing:
                ├── HP current/max (+/-/heal/damage)
                ├── Spell slots usados (toggle)
                ├── Conditions (add/remove)
                ├── Resources (Sorcery Points, Ki, etc.)
                └── Short Rest / Long Rest
```

### Guest vs. Logado

| Funcionalidade | Guest (localStorage) | Logado (MariaDB) |
|----------------|----------------------|------------------|
| Criar personagem | ✅ | ✅ |
| Ver personagens | ✅ (apenas local) | ✅ (cloud + local) |
| Editar ficha | ✅ | ✅ (sync automático) |
| Exportar PDF/JSON | ✅ | ✅ |
| Aceder de outro device | ❌ | ✅ |
| Compartilhar com DM | ❌ | ✅ (via campanha, Fase 4) |
| Backup automático | ❌ | ✅ |

---

## 7. Checklist de Aceitação Fase 3

- [ ] User pode fazer login com GitHub
- [ ] User pode criar personagem passo-a-passo (6 steps)
- [ ] Cada step tem validação e não deixa avançar se inválido
- [ ] Preview live funciona no step Review (POST /build)
- [ ] Personagem é salvo em MariaDB quando logado
- [ ] Guest mode funciona (localStorage) quando não logado
- [ ] User pode ver lista de personagens (cloud quando logado, local quando guest)
- [ ] User pode editar ficha interativa (HP, spell slots, conditions)
- [ ] Export PDF funciona end-to-end
- [ ] Export JSON funciona
- [ ] Mobile responsive (testado em 375px, 768px, 1024px)
- [ ] Tema D&D mantido (Tiamat, pergaminho, carmesim, ouro)
- [ ] gosec passa sem HIGH/CRITICAL
- [ ] Testes unitários ≥ 70% backend, ≥ 50% frontend
- [ ] Playwright E2E passa (jornada Player completa)
- [ ] Lighthouse score ≥ 80

---

*Documento atualizado em 2026-08-01. Arquiteto: André Ataíde + Kimi.*
*Fase 2 concluída. Pronto para Fase 3.*
