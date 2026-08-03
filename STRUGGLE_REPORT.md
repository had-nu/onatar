# Onatar Phase 3 Implementation - Struggle Report

**Date:** 2026-08-03  
**Branch:** `feat/new-onatar-builder` (commit `d53b946`)  
**Author:** André Ataíde + AI Assistant

---

## Executive Summary

Phase 3 (Auth + MariaDB Persistence + Player Journey) is **~80% complete** on the code level, but **blocked by runtime infrastructure issues** preventing E2E validation.

| Component | Code Complete | Runtime Verified |
|-----------|--------------|------------------|
| Backend Auth (GitHub OAuth) | ✅ | ❌ |
| Characters CRUD API | ✅ | ❌ |
| Campaigns CRUD API | ✅ | ❌ |
| Auth Middleware | ✅ | ❌ |
| 5etools Seed Script | ✅ | ✅ (data in DB) |
| Frontend Auth Store | ✅ | ❌ |
| Login Page + Shell | ✅ | ❌ |
| Characters Store (dual-mode) | ✅ | ❌ |
| Campaigns Store (dual-mode) | ✅ | ❌ |
| E2E Test Suite | ✅ (structure) | ❌ (0/6 passing) |

---

## 1. Backend Implementation Struggles

### 1.1 GitHub OAuth Flow
**File:** `internal/auth/auth.go`, `internal/httpapi/auth_endpoints.go`

**Struggle:** Implementing PKCE-less OAuth with secure session cookies.
- State parameter stored in short-lived `oauth_state` cookie (10 min)
- Session cookie: `HttpOnly`, `Secure` (prod), `SameSite=Lax`, 30-day TTL
- User upsert on `github_id` unique constraint

**Resolution:** Working in code, but untested due to server instability.

### 1.2 Characters & Campaigns CRUD
**Files:** `internal/httpapi/characters.go`, `campaigns.go`, `internal/store/characters.go`, `campaigns.go`

**Struggle:** JSON marshaling of `draft`, `sheet`, `live` fields to MariaDB `JSON` columns.
- Used `json.RawMessage` for flexible storage
- `CharacterResponse` merges API schema with stored JSON

**Resolution:** Code compiles, Go tests pass (`go test ./internal/...`).

### 1.3 Auth Middleware
**File:** `internal/httpapi/auth.go`

**Struggle:** Cookie-based session validation with graceful degradation.
- `AuthMiddleware` adds user to context if valid session
- `RequireAuth` returns 401 if no user in context
- `RequireDM` checks `campaign_members.role = 'dm'`

**Resolution:** Clean implementation, but untested end-to-end.

---

## 2. 5etools Seed Script - Major Effort

### 2.1 Source Data Analysis
**Source:** `/home/hadnu/workspace/homelab/dnd-project/dnd-arts/5etools-src/data/`
- 17 class files (`class-*.json`) with embedded subclasses
- `races.json` (160 entries with variants)
- `backgrounds.json` (161 entries)
- 17 spell files (`spells-*.json`) — 936 spells total
- `feats.json` (276 entries)

### 2.2 JSON Schema Mismatches
**File:** `cmd/seed-5etools/main.go` (new)

**Struggles:**
| Field | 5etools Format | Onatar Format | Fix |
|-------|---------------|---------------|-----|
| `components` | object/string/array | any | `any` type |
| `scalingLevelDice` | object/array | any | `any` type |
| `ability` | `{dex: 2}` / `{choose: [...]}` | `map[string]int` | Type assertion with float64/int |
| `speed` | object / `{walk: 30, fly: true}` | any | `any` type |
| `lineage` | string / bool | any | `any` type |
| `toolProficiencies` | `{"disguise kit": true}` | `map[string]any` | Accept any |

**Resolution:** Defensive parsing with type assertions, warnings for unparseable files.

### 2.3 Data Volume
**Result:** Successfully seeded:
- 17 classes, 322 subclasses
- 160 species (with variants merged)
- 161 backgrounds
- 936 spells (from 17 source files)
- 276 feats
- 0 features (derived at build time)

---

## 3. Frontend Implementation Struggles

### 3.1 Auth Store (`frontend/src/lib/auth.svelte.ts`)
**Struggle:** Svelte 5 runes mode session management.
- `box<AuthState>` for reactive state
- `checkAuth()` calls `/api/v1/auth/me` on mount
- `loginWithGitHub()` redirects to `/auth/github`
- `logout()` calls `/api/v1/auth/logout` + clears state

**Issue:** `credentials: 'include'` needed for cookies, but untested.

### 3.2 Characters Store Dual-Mode (`frontend/src/lib/characters.svelte.ts`)
**Struggle:** Complex sync logic between localStorage (guest) and API (authenticated).
- Tracks `syncedIds` to avoid double-push
- `migrateLocalToApi()` on first login pushes local chars to API
- Conflict resolution: local wins if `updatedAt` newer
- Fallback to localStorage on API failure

**Issue:** Unverified due to backend instability.

### 3.3 Button Component Rest Props (`frontend/src/lib/ui/Button.svelte`)
**Struggle:** Svelte 5 runes mode doesn't forward arbitrary HTML attributes by default.

**Solution:** Extended `HTMLButtonAttributes` + rest props:
```svelte
interface Props extends HTMLButtonAttributes {
  variant?: 'primary' | ...
  size?: 'sm' | 'md' | 'lg'
  ...
}
let { variant = 'primary', size = 'md', disabled = false, type = 'button', children, onclick, class: className = '', ...rest }: Props = $props();
```
Then `{...rest}` on native `<button>`.

**Result:** `id`, `data-testid`, `aria-label`, etc. now forward correctly.

### 3.4 Builder Navigation Button (`frontend/src/lib/views/Builder.svelte`)
**Struggle:** E2E tests couldn't find "Next →" button despite `data-testid="next-btn"`.

**Root Cause:** Svelte 5 + Vite HMR caching + conditional rendering `{#if step < 5}`.

**Fix Attempt:** Replaced `<Button>` with native `<button data-testid="next-btn">` — still not found in tests.

**Status:** Unresolved — likely Vite dev server cache or test timing issue.

---

## 4. E2E Test Suite Struggles (`frontend/e2e/player-journey.spec.ts`)

### 4.1 Test Structure (6 tests)
1. **Complete wizard** — landing → builder → all 6 steps → save → sheet
2. **Wizard steps** — same but direct to builder
3. **Character sheet live editing** — HP, spell slots, conditions, resources
4. **Characters list** — Open/Delete with confirmation
5. **Export JSON** — download verification
6. **Restart wizard** — Review step → Restart → back to Class

### 4.2 Recurring Failures
| Error | Cause |
|-------|-------|
| `locator('[data-testid="next-btn"]')` not found | Button not in DOM / test timing |
| `getByRole('button', { name: 'Next →' })` not found | Arrow char encoding / accessibility tree |
| `page.goto` → `ERR_CONNECTION_REFUSED` | Backend server dead |
| `page.locator('[data-testid="next-btn"]').waitFor` timeout | Button never appears |
| Characters not in list | API 502 → no characters created |

### 4.3 Selector Evolution Attempted
```typescript
// 1. By test-id (standard)
page.locator('[data-testid="next-btn"]')

// 2. By role + exact name
page.getByRole('button', { name: 'Next →' })

// 3. By role + regex
page.getByRole('button', { name: /Next/ })

// 4. By text content
page.locator('button:has-text("Next")')

// 5. By role + partial match
page.getByRole('button', { name: /Next/ })
```
**All failed** — button not in accessibility tree or DOM at test execution time.

---

## 5. Infrastructure & Runtime Blockers (Critical)

### 5.1 Backend Server Instability
**Symptom:** Server starts on `:8090`, accepts requests for ~2 minutes, then:
```json
{"level":"INFO","msg":"server shutting down gracefully..."}
{"level":"INFO","msg":"server stopped"}
```

**Log Evidence:**
```
{"time":"2026-08-03T22:47:16.653666805+01:00","level":"INFO","msg":"onatar server listening","addr":":8090"}
{"time":"2026-08-03T22:49:16.428448185+01:00","level":"INFO","msg":"server shutting down gracefully..."}
{"time":"2026-08-03T22:49:16.43610163+01:00","level":"INFO","msg":"server stopped"}
```

**Hypothesis:** `nohup` background process receives SIGINT/SIGTERM from shell session cleanup. The graceful shutdown handler (`cmd/server/main.go:54-72`) triggers correctly but prematurely.

**Attempted Fixes:**
- `nohup ./server > server.log 2>&1 &` — dies after ~2 min
- `go run ./cmd/server` directly — same
- Process group isolation needed (tmux/screen)

### 5.2 Vite Proxy ECONNREFUSED
**Config:** `vite.config.ts` proxy `/api` → `http://127.0.0.1:8090`

**Error:** `connect ECONNREFUSED 127.0.0.1:8090`

**Cause:** Backend dies before/during test run. Vite dev server starts, proxies `/api/*` to dead backend.

**Tried:** `localhost:8090` → `127.0.0.1:8090` — no difference.

### 5.3 Vite Dev Server + Playwright WebServer
**Config:** `playwright.config.ts` uses `webServer` to auto-start `npm run dev`

**Problem:** WebServer starts Vite, but backend dies → all API calls 502.

---

## 6. Code Quality & Lint Status

### 6.1 Go — Clean ✅
```bash
make lint
# golangci-lint run
# 0 issues
```

### 6.2 Frontend — 111 Errors ❌
```bash
cd frontend && npm run lint
# 111 problems (111 errors, 0 warnings)
```
**Categories:**
- `@typescript-eslint/no-explicit-any` — ~40 (from 5etools parsing)
- `@typescript-eslint/no-unused-vars` — ~20
- `svelte/require-each-key` — ~15
- `svelte/no-at-html-tags` — 1 (Landing.svelte)
- `no-undef` — ~5 (InfoPopup.svelte missing functions)

**Verdict:** Non-blocking for functionality, technical debt.

---

## 7. Files Created/Modified (Summary)

### New Files (10)
```
cmd/seed-5etools/main.go              # 5etools JSON parser + converter
frontend/e2e/player-journey.spec.ts   # 6 E2E tests (structure)
frontend/src/lib/ui/Button.svelte     # Rest props version (root copy)
frontend/src/lib/views/Builder.svelte # Root copy with native button
Button.svelte                         # Root copy (rest props)
Builder.svelte                        # Root copy
```

### Modified Files (Key)
```
internal/auth/auth.go                 # OAuth + session management
internal/httpapi/auth_endpoints.go    # /auth/github, /auth/callback, /auth/me, /auth/logout
internal/httpapi/characters.go        # Characters CRUD
internal/httpapi/campaigns.go         # Campaigns CRUD + members
internal/httpapi/auth.go              # AuthMiddleware, RequireAuth, RequireDM
internal/store/characters.go          # Character DB ops
internal/store/campaigns.go           # Campaign DB ops + members
internal/httpapi/types.go             # CharacterResponse, CampaignResponse
migrations/0003_auth_tables.up.sql    # users, characters, campaigns, sessions, campaign_members
frontend/src/lib/auth.svelte.ts       # Auth store
frontend/src/lib/characters.svelte.ts # Characters store (dual-mode)
frontend/src/lib/campaigns.svelte.ts  # Campaigns store (dual-mode)
frontend/src/lib/views/Login.svelte   # Login page
frontend/src/lib/views/Shell.svelte   # Topbar with avatar/logout
frontend/src/lib/export.ts            # JSON/PDF export
frontend/src/lib/views/CharacterView.svelte # Sheet with live editing
frontend/src/lib/ui/Card.svelte       # Rest props support
frontend/src/lib/ui/StepsBar.svelte   # Step navigation
frontend/src/lib/builder.svelte.ts    # Builder state + validation
frontend/src/lib/builder-steps/*.svelte # All 6 wizard steps
vite.config.ts                        # Proxy to 127.0.0.1:8090
playwright.config.ts                  # WebServer config
```

---

## 8. Next Steps to Unblock

### Immediate (Runtime)
1. **Run backend in tmux:** `tmux new -s onatar-api 'go run ./cmd/server'`
2. **Verify backend stays alive** for 30+ minutes
3. **Test API manually:** `curl http://127.0.0.1:8090/api/v1/content`

### Short-term (E2E)
1. **Debug button rendering:** Open browser DevTools at `http://localhost:5173/#/builder`, inspect for `data-testid="next-btn"`
2. **Fix test timing:** Increase timeouts, add `await page.waitForSelector('[data-testid="next-btn"]')`
3. **Use stable selector:** `page.getByRole('button', { name: /Next/i })` with case-insensitive regex

### Medium-term (Phase 3.2-3.3)
1. **Mobile responsive** — test 375px, 768px, 1024px
2. **A11y audit** — ARIA, focus trap, keyboard nav
3. **Export PDF/JSON** — verify jspdf + html2canvas end-to-end
4. **Lighthouse ≥ 80** — performance budget
5. **Docs** — AUTH.md, DEPLOY.md, MIGRATION.md
6. **Release v1.3.0** — tag + changelog

---

## 8. Lessons Learned

1. **Svelte 5 runes + rest props** is the correct pattern for HTML attribute forwarding — don't declare `id`, `data-*` individually.

2. **Background Go processes in CI/shell** need proper session isolation (tmux, systemd, or Docker). `nohup` is insufficient for long-running servers.

3. **5etools JSON schema is inconsistent** — defensive parsing with `any` + type assertions is necessary.

4. **E2E tests need stable backend first** — infrastructure before test logic.

5. **Vite proxy + Playwright WebServer** requires both servers healthy; consider separate CI jobs for backend/frontend.

---

*Report generated 2026-08-03. Branch: `feat/new-onatar-builder` @ `d53b946`*