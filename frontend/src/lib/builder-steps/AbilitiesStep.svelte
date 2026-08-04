<script lang="ts">
  import { draft, setAbilityScore, setAbilityMethod } from '$lib/builder.svelte';
  import Button from '$lib/ui/Button.svelte';

  const abilities = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const;
  const abilityNames: Record<string, string> = {
    STR: 'Força', DEX: 'Destreza', CON: 'Constituição',
    INT: 'Inteligência', WIS: 'Sabedoria', CHA: 'Carisma',
  };

  const methods = [
    { id: 'standard' as const, name: 'Array Padrão', desc: '15, 14, 13, 12, 10, 8' },
    { id: 'point-buy' as const, name: 'Compra de Pontos', desc: '27 pontos para gastar' },
    { id: 'roll' as const, name: 'Rolagem 4d6', desc: 'Role os dados' },
  ];

  const standardPool = [15, 14, 13, 12, 10, 8];

  function getMod(score: number): string {
    const mod = Math.floor((score - 10) / 2);
    return mod >= 0 ? `+${mod}` : `${mod}`;
  }

  function pointBuyCost(score: number): number {
    if (score <= 13) return score - 8;
    if (score === 14) return 7;
    if (score === 15) return 9;
    return 0;
  }

  function totalPointBuy(): number {
    return Object.values(draft.abilityScores).reduce((sum, s) => sum + pointBuyCost(s), 0);
  }

  const poolChips = $derived(() => {
    if (draft.abilityMethod === 'standard') {
      const used = Object.values(draft.abilityScores);
      return standardPool.map(v => ({ value: v, used: used.includes(v) }));
    }
    if (draft.abilityMethod === 'point-buy') {
      return Array.from({ length: 15 }, (_, i) => i + 8).map(v => ({
        value: v,
        used: false,
        cost: pointBuyCost(v),
      }));
    }
    return [];
  });

  function assignFromPool(value: number) {
    const firstEmpty = abilities.find(ab => draft.abilityScores[ab] === 10);
    if (firstEmpty) setAbilityScore(firstEmpty, value);
  }

  function unassign(ab: typeof abilities[number]) {
    setAbilityScore(ab, 10);
  }

  function rollAbility(): number {
    const rolls = Array.from({ length: 4 }, () => Math.floor(Math.random() * 6) + 1);
    rolls.sort((a, b) => b - a);
    return rolls[0] + rolls[1] + rolls[2];
  }

  function rollAll() {
    for (const ab of abilities) {
      setAbilityScore(ab, rollAbility());
    }
  }
</script>

<div class="step-abilities">
  <h2 class="step-title">Atribua Atributos</h2>
  <p class="step-desc">Os seis atributos definem as capacidades f00edsicas e mentais do seu personagem.</p>

  <!-- Method tabs -->
  <div class="method-tabs">
    {#each methods as m}
      <button
        class="method-tab"
        class:active={draft.abilityMethod === m.id}
        onclick={() => setAbilityMethod(m.id)}
      >
        <span class="mt-name">{m.name}</span>
        <span class="mt-desc">{m.desc}</span>
      </button>
    {/each}
  </div>

  <!-- Pool chips -->
  {#if draft.abilityMethod === 'standard'}
    <div class="pool">
      <div class="pool-label">Valores disponíveis <span class="pool-hint">— clique para atribuir</span></div>
      <div class="pool-chips">
        {#each poolChips() as chip}
          <button
            class="pool-chip"
            class:used={chip.used}
            disabled={chip.used}
            onclick={() => assignFromPool(chip.value)}
          >
            {chip.value}
          </button>
        {/each}
      </div>
    </div>
  {:else if draft.abilityMethod === 'point-buy'}
    <div class="pool">
      <div class="pool-label">
        Pontos gastos:
        <span class="pb-total" class:pb-over={totalPointBuy() > 27}>
          {totalPointBuy()} / 27
        </span>
      </div>
      <div class="pool-chips">
        {#each poolChips() as chip}
          <button
            class="pool-chip"
            onclick={() => {
              const firstEmpty = abilities.find(ab => draft.abilityScores[ab] === 10);
              if (firstEmpty) setAbilityScore(firstEmpty, chip.value);
            }}
          >
            {chip.value}
            <span class="chip-cost">{chip.cost}pts</span>
          </button>
        {/each}
      </div>
    </div>
  {:else if draft.abilityMethod === 'roll'}
    <div class="roll-area">
      <Button variant="primary" onclick={rollAll}>🎲 Rolagem 4d6 (descarte menor)</Button>
      <p class="roll-hint">Clique para rolar os dados para todos os atributos.</p>
    </div>
  {/if}

  <!-- Ability slots -->
  <div class="abilities-grid">
    {#each abilities as ab}
      {@const score = draft.abilityScores[ab]}
      {@const mod = getMod(score)}
      <button
        class="ability-slot"
        class:filled={score !== 10}
        onclick={() => unassign(ab)}
      >
        <div class="ab-name">{abilityNames[ab]}</div>
        <div class="ab-score">{score}</div>
        <div class="ab-mod" class:ab-mod-pos={score > 10} class:ab-mod-neg={score < 10}>
          {mod}
        </div>
        {#if score !== 10}
          <span class="ab-reset">↺</span>
        {/if}
      </button>
    {/each}
  </div>

  {#if draft.abilityMethod === 'point-buy' && totalPointBuy() > 27}
    <div class="pb-warning">
      ⚠ Limite de 27 pontos excedido. Reduza alguns atributos.
    </div>
  {/if}
</div>

<style>
  .step-abilities { animation: fadeIn 0.3s ease; }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .step-title { font-size: 28px; font-weight: 500; margin: 0 0 8px 0; }
  .step-desc { font-size: 14px; color: var(--on-text-dim); margin: 0 0 24px 0; }

  .method-tabs {
    display: flex;
    gap: 8px;
    margin-bottom: 24px;
    flex-wrap: wrap;
  }
  .method-tab {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 12px 20px;
    border-radius: 10px;
    border: 1px solid var(--on-border);
    background: transparent;
    color: var(--on-text-muted);
    cursor: pointer;
    transition: all 0.15s;
    min-width: 140px;
  }
  .method-tab:hover { border-color: var(--on-border-light); color: var(--on-text); }
  .method-tab.active {
    background: var(--on-red);
    color: #fff;
    border-color: var(--on-red);
    box-shadow: 0 2px 8px rgba(197, 0, 9, 0.25);
  }
  .mt-name { font-size: 14px; font-weight: 500; }
  .mt-desc { font-size: 11px; opacity: 0.7; margin-top: 2px; }

  .pool {
    padding: 16px;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 10px;
    margin-bottom: 24px;
  }
  .pool-label {
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--on-text-dim);
    margin: 0 0 12px 0;
  }
  .pool-hint { text-transform: none; letter-spacing: 0; font-size: 12px; color: var(--on-text-dim); }
  .pool-chips { display: flex; gap: 8px; flex-wrap: wrap; }
  .pool-chip {
    position: relative;
    width: 52px;
    height: 52px;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.15s;
    border: 2px solid var(--on-border);
    background: var(--on-bg-elevated);
    color: var(--on-text);
  }
  .pool-chip:hover:not(.used) { border-color: var(--on-red); background: var(--on-bg-hover); }
  .pool-chip.used {
    opacity: 0.25;
    cursor: not-allowed;
    border-color: var(--on-border);
  }
  .chip-cost {
    font-size: 9px;
    color: var(--on-text-dim);
    font-weight: 400;
    margin-top: 2px;
  }

  .pb-total { font-weight: 600; color: var(--on-green); }
  .pb-total.pb-over { color: #ff6b6b; }
  .pb-warning {
    margin-top: 16px;
    padding: 10px 14px;
    background: rgba(197, 0, 9, 0.08);
    border: 1px solid rgba(197, 0, 9, 0.2);
    border-radius: 8px;
    font-size: 13px;
    color: #ff6b6b;
  }

  .roll-area {
    padding: 24px;
    text-align: center;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 10px;
    margin-bottom: 24px;
  }
  .roll-hint { font-size: 12px; color: var(--on-text-dim); margin: 12px 0 0 0; }

  .abilities-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
  }
  .ability-slot {
    position: relative;
    padding: 16px;
    background: var(--on-bg-surface);
    border: 1px solid var(--on-border);
    border-radius: 10px;
    text-align: center;
    transition: all 0.15s;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    color: inherit;
  }
  .ability-slot:hover { border-color: var(--on-border-light); }
  .ability-slot.filled {
    border-color: var(--on-green);
    background: rgba(0, 184, 122, 0.05);
  }
  .ab-name {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--on-text-dim);
  }
  .ab-score {
    font-size: 32px;
    font-weight: 500;
    color: var(--on-text);
    line-height: 1;
  }
  .ab-mod {
    font-size: 14px;
    font-weight: 500;
    color: var(--on-text-dim);
  }
  .ab-mod-pos { color: var(--on-green); }
  .ab-mod-neg { color: #ff6b6b; }
  .ab-reset {
    position: absolute;
    top: 6px;
    right: 8px;
    font-size: 11px;
    color: var(--on-text-dim);
    opacity: 0;
    transition: opacity 0.15s;
  }
  .ability-slot:hover .ab-reset { opacity: 1; }

  @media (max-width: 600px) {
    .abilities-grid { grid-template-columns: repeat(2, 1fr); }
  }
</style>