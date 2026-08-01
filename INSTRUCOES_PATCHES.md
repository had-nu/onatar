# Onatar — Resolução de Dívida Técnica (Patches)

> 6 patches para corrigir bugs críticos e aplicar o tema visual D&D.  
> Gerado em 2026-08-01 para o repositório `had-nu/onatar`.

---

## Resumo dos Patches

| # | Ficheiro | Bug | Severidade |
|---|----------|-----|------------|
| 01 | `frontend/src/lib/builder.svelte.ts` | `saveCharacterFromWizard` não chamava `POST /build` — ficha salva era draft cru | 🔴 Crítico |
| 02 | `internal/httpapi/ratelimit.go` | `clientIP` ignorava `X-Forwarded-For` — rate limit quebrado atrás de proxy | 🟡 Médio |
| 03 | `cmd/server/main.go` | Sem graceful shutdown — SIGTERM mata conexões ativas | 🟡 Médio |
| 04 | `frontend/src/lib/builder.svelte.ts` | `abilitiesValid` não verificava `POINT_BUY_MAX` (15) | 🟢 Baixo |
| 05 | `frontend/src/lib/views/builder/ReviewStep.svelte` | Não tratava `saveCharacterFromWizard` como async — sem loading nem error | 🟡 Médio |
| 06 | `frontend/src/app.css` | Tema genérico roxo → tema D&D (pergaminho, carmesim, ouro, Tiamat) | 🔴 Crítico (UX) |

---

## Como Aplicar

### Opção 1: `git apply` (recomendado)

```bash
cd /caminho/para/onatar

# Extrair o ZIP
unzip ONATAR_DEBT_PATCHES.zip -d patches/

# Aplicar um de cada vez (para rever)
git apply patches/patch-01-builder-saveCharacterFromWizard.patch
git apply patches/patch-02-ratelimit-clientIP.patch
git apply patches/patch-03-main-graceful-shutdown.patch
git apply patches/patch-04-builder-abilitiesValid.patch
git apply patches/patch-05-reviewstep-async-save.patch
git apply patches/patch-06-app-css-theme.patch

# Ou aplicar todos de uma vez
for p in patches/patch-*.patch; do git apply "$p"; done
```

### Opção 2: Copiar manualmente

Abrir cada patch e copiar as secções `+` (linhas adicionadas) e remover as `-` (linhas removidas) nos ficheiros correspondentes.

---

## Integração de Assets Locais

Tens assets em `/home/hadnu/workspace/homelab/dnd-project/dnd-arts/` incluindo a fonte **Tiamat**.

### 1. Fonte Tiamat

Copiar os ficheiros da fonte para o projeto:

```bash
# Criar diretório de fontes no frontend
mkdir -p frontend/public/fonts/tiamat

# Copiar ficheiros da fonte (ajustar extensões conforme os teus ficheiros)
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/fonts/Tiamat-*.woff2 frontend/public/fonts/tiamat/
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/fonts/Tiamat-*.ttf frontend/public/fonts/tiamat/

# Criar @font-face em frontend/src/app.css (ADICIONAR no início do ficheiro)
```

Adicionar ao `frontend/src/app.css` (antes do `:root`):

```css
/* Fonte oficial D&D — Tiamat */
@font-face {
  font-family: 'Tiamat';
  src: url('/fonts/tiamat/Tiamat-Regular.woff2') format('woff2'),
       url('/fonts/tiamat/Tiamat-Regular.ttf') format('truetype');
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: 'Tiamat';
  src: url('/fonts/tiamat/Tiamat-Bold.woff2') format('woff2'),
       url('/fonts/tiamat/Tiamat-Bold.ttf') format('truetype');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}

/* Fonte body — Source Sans 3 (Google Fonts) */
@import url('https://fonts.googleapis.com/css2?family=Source+Sans+3:wght@400;600&display=swap');
```

### 2. Ícones / Sprites

Copiar SVGs/PNGs para `frontend/public/icons/`:

```bash
mkdir -p frontend/public/icons/classes
mkdir -p frontend/public/icons/spells
mkdir -p frontend/public/icons/ui

cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/icons/classes/*.svg frontend/public/icons/classes/
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/icons/spells/*.svg frontend/public/icons/spells/
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/icons/ui/*.svg frontend/public/icons/ui/
```

Referenciar em componentes Svelte:

```svelte
<!-- Exemplo: ClassCard.svelte -->
<img src={`/icons/classes/${classId}.svg`} alt={className} />
```

### 3. Favicon

```bash
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/favicon.ico frontend/public/favicon.ico
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/apple-touch-icon.png frontend/public/apple-touch-icon.png
```

Atualizar `frontend/index.html`:

```html
<link rel="icon" type="image/x-icon" href="/favicon.ico" />
<link rel="apple-touch-icon" href="/apple-touch-icon.png" />
<meta name="description" content="Onatar — Criador de fichas de personagem D&D 2024" />
<meta property="og:title" content="Onatar" />
<meta property="og:description" content="Cria e gere fichas de personagem para D&D 2024" />
<meta property="og:image" content="/og-image.png" />
```

---

## Verificação Pós-Patch

### Backend

```bash
cd backend
go test ./...
go vet ./...
gosec ./...

# Verificar compilação
go build ./cmd/server

# Testar graceful shutdown
./server &
PID=$!
sleep 2
kill -TERM $PID
# Deve ver "server shutting down gracefully..." nos logs
```

### Frontend

```bash
cd frontend
npm run lint
npx tsc --noEmit
npm run test

# Verificar build
npm run build

# Verificar se a fonte Tiamat está no bundle
ls dist/fonts/
```

### Testes Manuais

1. Abrir o builder wizard
2. Preencher todos os 6 steps
3. No step Review, clicar "Guardar personagem"
4. Verificar no Network tab que `POST /api/v1/build` é chamado
5. Verificar que a ficha salva tem `sheet` preenchido (não `undefined`)
6. Verificar que a ficha mostra HP, AC, spell slots corretos

---

## Notas

- O patch 01 torna `saveCharacterFromWizard` **async**. O patch 05 atualiza o `ReviewStep.svelte` para tratar isso. Se tiveres outros componentes a chamar `saveCharacterFromWizard`, precisam de `await` também.
- O patch 06 muda a paleta completa. Se tiveres componentes com cores hardcoded (ex: `color: #aa3bff`), precisam de ser atualizados manualmente.
- A fonte Tiamat pode não ter suporte a todos os glifos (ex: acentos portugueses). Testar com "Espécie", "Feitiços", etc. Se necessário, adicionar `Source Sans 3` como fallback para caracteres especiais.

---

*Patches gerados por Kimi para André Ataíde. Stack: Go + Svelte 5 + MariaDB.*
