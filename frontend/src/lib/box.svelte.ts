// Minimal reactive container for module-scope state in `.svelte.ts` stores.
// Uses a Svelte store for cross-module reactivity.
import { writable, type Writable, get } from 'svelte/store';

export function box<T>(initial: T): Writable<T> & { value: T } {
  const store = writable(initial);
  // Track the current value for synchronous access
  let currentValue = initial;
  return {
    ...store,
    get value() {
      return currentValue;
    },
    set value(next: T) {
      currentValue = next;
      store.set(next);
    },
  };
}
