import type { Task } from '@/types'

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

export const fetchTasks = (): Promise<Task[]> =>
  request<Task[]>('/tasks')

export const addTask = (payload: {
  task_text: string
  priority: number
  deadline?: string | null
}): Promise<Task> =>
  request<Task>('/tasks', {
    method: 'POST',
    body: JSON.stringify(payload),
  })

export const updateTask = (
  id: number,
  patch: Partial<Pick<Task, 'task_text' | 'is_completed' | 'priority' | 'deadline'>>
): Promise<void> =>
  request<void>(`/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify(patch),
  })

export const removeTask = (id: number): Promise<void> =>
  request<void>(`/tasks/${id}`, {
    method: 'DELETE'
  })

export const saveOrder = (ids: number[]): Promise<void> =>
  request<void>('/tasks/reorder', {
    method: 'PUT',
    body: JSON.stringify({ ids }),
  })
  