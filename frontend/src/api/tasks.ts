import type { Task } from '@/types'

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
 * Fetches all tasks for the current user.
 */
export const fetchTasks = (): Promise<Task[]> =>
  request<Task[]>('/tasks')

/**
 * Creates a new task (with optional category).
 */
export const addTask = (payload: {
  task_text: string
  priority: number
  deadline?: string | null
  category_id?: number | null
}): Promise<Task> =>
  request<Task>('/tasks', {
    method: 'POST',
    body: JSON.stringify(payload),
  })

/**
 * Updates specific fields of an existing task (including category).
 */
export const updateTask = (
  id: number,
  patch: Partial<Pick<Task, 'task_text' | 'is_completed' | 'priority' | 'deadline' | 'category_id'>>
): Promise<void> =>
  request<void>(`/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(patch),
  })

/**
 * Deletes a task by its ID.
 */
export const removeTask = (id: number): Promise<void> =>
  request<void>(`/tasks/${id}`, {
    method: 'DELETE'
  })

/**
 * Persists the new display order of tasks (drag-and-drop).
 */
export const saveOrder = (ids: number[]): Promise<void> =>
  request<void>('/tasks/reorder', {
    method: 'PUT',
    body: JSON.stringify({ ids }),
  })
