# Onatar — Resolução de Dívida Técnica (Patches v2)

> 6 patches para corrigir bugs críticos e aplicar o tema visual D&D.  
> Fonte Tiamat confirmada em `/home/hadnu/workspace/homelab/dnd-project/dnd-arts/fonts/`.  
> Gerado em 2026-08-01 para o repositório `had-nu/onatar`.

---

## Aplicação dos Patches

```bash
cd /caminho/para/onatar

# 1. Verificar cada patch antes de aplicar (dry-run)
git apply --check patches/patch-01-builder-saveCharacterFromWizard.patch
git apply --check patches/patch-02-ratelimit-clientIP.patch
git apply --check patches/patch-03-main-graceful-shutdown.patch
git apply --check patches/patch-04-builder-abilitiesValid.patch
git apply --check patches/patch-05-reviewstep-async-save.patch
git apply --check patches/patch-06-app-css-theme.patch

# 2. Se todos passarem, aplicar um a um (para poder reverter individualmente)
git apply patches/patch-01-builder-saveCharacterFromWizard.patch
git apply patches/patch-02-ratelimit-clientIP.patch
git apply patches/patch-03-main-graceful-shutdown.patch
git apply patches/patch-04-builder-abilitiesValid.patch
git apply patches/patch-05-reviewstep-async-save.patch
git apply patches/patch-06-app-css-theme.patch

# 3. Commit agrupado
git add -A
git commit -m "fix(debt): resolve critical bugs F2/B2/B5/B4 + apply D&D theme

- saveCharacterFromWizard now calls POST /build before saving
- clientIP respects X-Forwarded-For and X-Real-IP
- Graceful shutdown with 10s timeout
- abilitiesValid checks POINT_BUY_MAX
- ReviewStep handles async save with loading/error states
- D&D theme: parchment, crimson, gold, Tiamat font"
```

---

## Integração de Assets Locais

### 1. Fonte Tiamat

Confirmada em `/home/hadnu/workspace/homelab/dnd-project/dnd-arts/fonts/`.

```bash
mkdir -p frontend/public/fonts/tiamat

# Copiar TODOS os ficheiros de fonte (woff2, ttf, otf, etc.)
cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/fonts/* frontend/public/fonts/tiamat/

# Verificar o que foi copiado
ls -la frontend/public/fonts/tiamat/
```

O patch 06 já adiciona o `@font-face` no `app.css`. **Ajustar os nomes dos ficheiros** conforme o que existe:

```css
/* Em frontend/src/app.css — ADICIONAR no início, antes do :root */

@font-face {
  font-family: 'Tiamat';
  src: url('/fonts/tiamat/TiamatRegular.woff2') format('woff2'),
       url('/fonts/tiamat/TiamatRegular.ttf') format('truetype');
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: 'Tiamat';
  src: url('/fonts/tiamat/TiamatBold.woff2') format('woff2'),
       url('/fonts/tiamat/TiamatBold.ttf') format('truetype');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}
```

> **Nota:** Ajustar os nomes dos ficheiros (`TiamatRegular`, `TiamatBold`, etc.) conforme os ficheiros reais na tua pasta. Se a fonte tiver outra convenção de nomes, editar o `@font-face`.

### 2. Ícones

**Estratégia:** Font Awesome 6 (já no projeto) para ícones genéricos. Assets do 5etools para classes/spells se disponíveis.

```bash
# Verificar se existem ícones no 5etools
ls /home/hadnu/workspace/homelab/dnd-project/dnd-arts/icons/ 2>/dev/null || echo "Sem pasta icons/"
ls /home/hadnu/workspace/homelab/dnd-project/5etools-img/ 2>/dev/null | head -20

# Se existirem, copiar
mkdir -p frontend/public/icons
# cp /home/hadnu/workspace/homelab/dnd-project/5etools-img/classes/*.png frontend/public/icons/
```

Se **não existirem** assets específicos, os componentes usam:
- Font Awesome para UI (sword, shield, heart, scroll, etc.)
- Iniciais da classe como placeholder (ex: "F" para Fighter em um círculo)
- Cores por spell school (já definidas no patch 06)

### 3. Favicon

```bash
# Verificar se existe favicon nos assets
ls /home/hadnu/workspace/homelab/dnd-project/dnd-arts/favicon* 2>/dev/null
ls /home/hadnu/workspace/homelab/dnd-project/5etools-src/favicon* 2>/dev/null

# Se existir, copiar
# cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/favicon.ico frontend/public/
# cp /home/hadnu/workspace/homelab/dnd-project/dnd-arts/apple-touch-icon.png frontend/public/
```

Se **não existir**, criar um favicon SVG inline simples (d20 carmesim):

```svg
<!-- frontend/public/favicon.svg -->
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <polygon points="50,5 95,25 95,75 50,95 5,75 5,25" fill="#8b0000" stroke="#c9a227" stroke-width="3"/>
  <text x="50" y="62" text-anchor="middle" fill="#f5f0e8" font-size="40" font-family="serif" font-weight="bold">20</text>
</svg>
```

E adicionar ao `frontend/index.html`:
```html
<link rel="icon" type="image/svg+xml" href="/favicon.svg" />
```

---

## Verificação Pós-Patch (Gate Completo)

### Backend

```bash
cd /caminho/para/onatar

# Compilar
go build ./cmd/server

# Testes + lint + security
go test ./...
go vet ./...
gosec ./...

# Testar graceful shutdown manualmente
./server &
PID=$!
sleep 1
kill -TERM $PID
# Esperar "server stopped" nos logs
```

### Frontend

```bash
cd frontend

# Type-check + lint + testes
npx tsc --noEmit
npm run lint
npx vitest run

# Build de produção (verifica se assets estão incluídos)
npm run build
ls dist/fonts/        # Tiamat deve estar aqui
ls dist/icons/        # Ícones devem estar aqui (se copiados)
```

### Teste Manual Crítico

1. Abrir app no browser
2. Criar personagem → passar pelos 6 steps
3. No Review, clicar "Guardar personagem"
4. **Verificar no Network tab:** `POST /api/v1/build` é chamado e retorna 200
5. **Verificar na ficha:** HP, AC, spell slots estão calculados (não vazios)
6. **Verificar tema:** fundo pergaminho, texto marrom, accent carmesim
7. **Verificar fonte:** títulos em Tiamat (serif), corpo em Source Sans 3

---

## Troubleshooting

### "error: patch does not apply"

```bash
# Verificar se o ficheiro original corresponde ao patch
git diff HEAD -- frontend/src/lib/builder.svelte.ts

# Se o agente de código modificou o ficheiro depois da minha análise,
# aplicar manualmente as mudanças (ver patch com +/- linhas)
cat patches/patch-01-builder-saveCharacterFromWizard.patch
```

### Fonte Tiamat não carrega

1. Verificar se os ficheiros estão em `frontend/public/fonts/tiamat/`
2. Verificar se os nomes no `@font-face` correspondem aos ficheiros reais
3. Verificar no DevTools → Network → Fonts se há 404
4. Se a Tiamat não tiver glifos para acentos (á, é, í, etc.), adicionar fallback:
   ```css
   font-family: 'Tiamat', 'Cinzel', Georgia, serif;
   ```

### POST /build não é chamado no save

1. Verificar se o patch 01 foi aplicado (procurar `buildDraft` em `builder.svelte.ts`)
2. Verificar se o patch 05 foi aplicado (procurar `handleSave` em `ReviewStep.svelte`)
3. Verificar se o backend está a correr (`GET /health` responde)

---

*Patches gerados por Kimi para André Ataíde. Stack: Go + Svelte 5 + MariaDB.*
