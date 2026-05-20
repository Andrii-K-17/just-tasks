import type { Category } from '@/types'

const BASE = '/api'

/**
 * Sends a generic HTTP request and handles JSON response parsing.
 */
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

/**
 * Fetch all categories.
 */
export const getCategories = (): Promise<Category[]> =>
  request<Category[]>('/categories')

/**
 * Create a new category.
 */
export const createCategory = (name: string): Promise<Category> =>
  request<Category>('/categories', {
    method: 'POST',
    body: JSON.stringify({ name }),
  })

/**
 * Delete a category by ID.
 */
export const deleteCategory = (id: number): Promise<{ deleted: boolean }> =>
  request<{ deleted: boolean }>(`/categories/${id}`, {
    method: 'DELETE',
  })
