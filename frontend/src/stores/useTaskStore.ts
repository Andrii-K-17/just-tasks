import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Task } from '@/types'
import * as tasksApi from '@/api/tasks'

export type FilterType = 'all' | 'active' | 'done'

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const filter = ref<FilterType>('all')

  const filteredTasks = computed(() => {
    if (filter.value === 'active') return tasks.value.filter(t => !t.is_completed)
    if (filter.value === 'done') return tasks.value.filter(t =>  t.is_completed)
    return tasks.value
  })

  const stats = computed(() => ({
    total: tasks.value.length,
    done: tasks.value.filter(t => t.is_completed).length,
  }))

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

  async function remove(id: number): Promise<void> {
    await tasksApi.removeTask(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  function reset(): void {
    tasks.value = []
    filter.value = 'all'
  }

  return { tasks, filter, filteredTasks, stats, load, add, toggle, editText, remove, reset }
})