import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Task } from '@/types'
import * as tasksApi from '@/api/tasks'

export type FilterType = 'all' | 'active' | 'done'

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const filter = ref<FilterType>('all')
  const searchQuery = ref('')

  const filteredTasks = computed(() => {
    let result = tasks.value

    if (filter.value === 'active') result = result.filter(t => !t.is_completed)
    if (filter.value === 'done')   result = result.filter(t =>  t.is_completed)

    const query = searchQuery.value.trim().toLowerCase()
    if (query) result = result.filter(t => t.task_text.toLowerCase().includes(query))

    return result
  })

  const stats = computed(() => {
    const total = tasks.value.length
    const done = tasks.value.filter(t => t.is_completed).length
    const active = total - done
    const today = new Date().toISOString().slice(0, 10)
    const overdue = tasks.value.filter(t => t.deadline && !t.is_completed && t.deadline < today).length
    const byPriority = {
      low: tasks.value.filter(t => t.priority === 1).length,
      medium: tasks.value.filter(t => t.priority === 2).length,
      high: tasks.value.filter(t => t.priority === 3).length,
    }
    const rate = total > 0 ? Math.round((done / total) * 100) : 0

    return { total, done, active, overdue, byPriority, rate }
  })

  async function reorder(newOrder: Task[]): Promise<void> {
    tasks.value = newOrder
    await tasksApi.saveOrder(newOrder.map(t => t.id))
  }

  async function load(): Promise<void> {
    tasks.value = await tasksApi.fetchTasks()
  }

  async function add(payload: {
    task_text: string
    priority: number
    deadline?: string | null
  }): Promise<void> {
    const task = await tasksApi.addTask(payload)
    tasks.value.unshift(task)
  }

  async function toggle(id: number): Promise<void> {
    const task = tasks.value.find(t => t.id === id)
    if (!task) return
    const newStatus = !task.is_completed
    await tasksApi.updateTask(id, { is_completed: newStatus })
    task.is_completed = newStatus
  }

  async function editText(id: number, task_text: string): Promise<void> {
    await tasksApi.updateTask(id, { task_text })
    const task = tasks.value.find(t => t.id === id)
    if (task) task.task_text = task_text
  }

  async function editDeadline(id: number, deadline: string | null): Promise<void> {
    await tasksApi.updateTask(id, { deadline })
    const task = tasks.value.find(t => t.id === id)
    if (task) task.deadline = deadline
  }

  async function remove(id: number): Promise<void> {
    await tasksApi.removeTask(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  function reset(): void {
    tasks.value = []
    filter.value = 'all'
    searchQuery.value = ''
  }

  return {
    tasks,
    filter,
    searchQuery,
    filteredTasks,
    stats,
    load,
    add,
    toggle,
    editText,
    editDeadline,
    remove,
    reset,
    reorder
  }
})
