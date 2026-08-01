// Minimal reactive container for module-scope state in `.svelte.ts` stores.
// Svelte runes forbid reassigning exported state, so stores expose a `value`
// getter/setter pair backed by a local `$state` variable instead.
export function box<T>(initial: T) {
  let value = $state(initial)
  return {
    get value() {
      return value
    },
    set value(next: T) {
      value = next
    },
  }
}
