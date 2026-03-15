import type { User } from '@/types'

const BASE = '/api'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${BASE}${url}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  const data = await response.json()
  if (!response.ok) throw new Error(data.error ?? 'Request failed')
  return data as T
}

export const fetchMe = (): Promise<User> =>
  request<User>('/me')

export const login = (username: string, password: string): Promise<User> =>
  request<User>('/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })

export const register = (username: string, password: string): Promise<User> =>
  request<User>('/register', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })

export const logout = (): Promise<void> =>
  request<void>('/logout', {
    method: 'POST'
  })

export const deleteAccount = (): Promise<void> =>
  request<void>('/account', {
    method: 'DELETE'
  })