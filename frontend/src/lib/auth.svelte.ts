// Auth store (Fase 3): GitHub OAuth session management.
// Stores user, checks session on load, handles login/logout flow.
import { box } from './box.svelte'
import type { User } from './types'

export interface AuthState {
  user: User | null
  loading: boolean
  error: string | null
}

const initialState: AuthState = {
  user: null,
  loading: true,
  error: null,
}

export const auth = box<AuthState>(initialState)

const API_BASE = '/api/v1'

async function fetchJson<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'include',
  })
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const body = (await res.json()) as { error?: { message?: string } }
      if (body?.error?.message) msg = body.error.message
    } catch {
      /* ignore */
    }
    throw new Error(msg)
  }
  return res.json() as Promise<T>
}

export async function checkAuth(): Promise<User | null> {
  auth.value = { ...auth.value, loading: true, error: null }
  try {
    const user = await fetchJson<User>(`${API_BASE}/auth/me`)
    auth.value = { user, loading: false, error: null }
    return user
  } catch {
    auth.value = { user: null, loading: false, error: null }
    return null
  }
}

export function loginWithGitHub() {
  window.location.href = `${API_BASE}/auth/github`
}

export async function logout() {
  auth.value = { ...auth.value, loading: true }
  try {
    await fetch(`${API_BASE}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    })
  } catch {
    /* ignore */
  }
  auth.value = { user: null, loading: false, error: null }
}

export function getUser(): User | null {
  return auth.value.user
}

export function isAuthenticated(): boolean {
  return auth.value.user !== null
}

export function isLoading(): boolean {
  return auth.value.loading
}

export function getAuthError(): string | null {
  return auth.value.error
}

/** Test helper: reset the module singleton between tests. */
export function _resetAuth() {
  auth.value = initialState
}