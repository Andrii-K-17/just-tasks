import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Task } from '@/types'
import * as tasksApi from '@/api/tasks'
import { useAuthStore } from '@/stores/useAuthStore'

export type FilterType = 'all' | 'active' | 'done' | 'shared'

/**
 * Global store for managing tasks state, filtering, and metrics.
 */
export const useTaskStore = defineStore('tasks', () => {
  const authStore = useAuthStore()

  /** Source list of all tasks. */
  const tasks = ref<Task[]>([])

  /** Current completion status filter. */
  const filter = ref<FilterType>('all')

  /** Search keyword for filtering tasks by text. */
  const searchQuery = ref('')

  /** Currently selected category ID for filtering tasks. */
  const selectedCategoryId = ref<number | null>(null)

  /**
   * Computes the list of tasks filtered by completion status, search query, and sharing state.
   */
  const filteredTasks = computed(() => {
    let result = tasks.value

    if (filter.value === 'active') result = result.filter(t => !t.is_completed)
    if (filter.value === 'done')   result = result.filter(t =>  t.is_completed)
    if (filter.value === 'shared') {
      result = result.filter(t => 
        (t.collaborators && t.collaborators.length > 0) || 
        (authStore.user && t.owner_name !== authStore.user.username)
      )
    }

    const query = searchQuery.value.trim().toLowerCase()
    if (query) result = result.filter(t => t.task_text.toLowerCase().includes(query))
    
    if (selectedCategoryId.value) result = result.filter(t => t.category_id === selectedCategoryId.value)

    return result
  })

  /**
  * Computes the list of tasks filtered by search query and category.
  * Used for global statistics on the main page.
  */
  const allFilteredTasks = computed(() => {
    let result = tasks.value
    
    const query = searchQuery.value.trim().toLowerCase()
    if (query) result = result.filter(t => t.task_text.toLowerCase().includes(query))

    if (selectedCategoryId.value) result = result.filter(t => t.category_id === selectedCategoryId.value)

    return result
  })

  /**
  * Computes the number of completed tasks from the globally filtered list.
  */
  const completedCount = computed(() =>
    allFilteredTasks.value.filter(t => t.is_completed).length
  )

  /**
  * Computes the total number of tasks from the globally filtered list.
  */
  const totalCount = computed(() => allFilteredTasks.value.length)

  /**
   * Computes task statistics, including completion rate and breakdown by priority.
   */
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

  /**
   * Updates the tasks order locally and persists it on the server.
   */
  async function reorder(newOrder: Task[]): Promise<void> {
    tasks.value = newOrder
    await tasksApi.saveOrder(newOrder.map(t => t.id))
  }

  /**
   * Fetches all tasks from the server and updates the store state.
   */
  async function load(): Promise<void> {
    tasks.value = await tasksApi.fetchTasks()
  }

  /**
   * Creates a new task and prepends it to the list.
   */
  async function add(payload: {
    task_text: string
    priority: number
    deadline?: string | null
    category_id?: number | null
  }): Promise<void> {
    const task = await tasksApi.addTask(payload)
    tasks.value.unshift(task)
  }

  async function editCategory(id: number, categoryId: number | null) {
    const task = tasks.value.find(t => t.id === id)
    if (task) {
      task.category_id = categoryId
      await tasksApi.updateTask(id, { category_id: categoryId })
    }
  }

  /**
   * Toggles the completion status of a specific task.
   */
  async function toggle(id: number): Promise<void> {
    const task = tasks.value.find(t => t.id === id)
    if (!task) return
    const newStatus = !task.is_completed
    await tasksApi.updateTask(id, { is_completed: newStatus })
    task.is_completed = newStatus
  }

  /**
   * Updates the text content of an existing task.
   */
  async function editText(id: number, task_text: string): Promise<void> {
    await tasksApi.updateTask(id, { task_text })
    const task = tasks.value.find(t => t.id === id)
    if (task) task.task_text = task_text
  }

  /**
   * Updates the deadline date of an existing task.
   */
  async function editDeadline(id: number, deadline: string | null): Promise<void> {
    await tasksApi.updateTask(id, { deadline })
    const task = tasks.value.find(t => t.id === id)
    if (task) task.deadline = deadline
  }

  /**
   * Deletes a task from the server and removes it from the store state.
   */
  async function remove(id: number): Promise<void> {
    await tasksApi.removeTask(id)
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  /**
   * Adds a collaborator to a task and reloads the store to get updated info.
   */
  async function addTaskCollaborator(id: number, username: string): Promise<void> {
    await tasksApi.addCollaborator(id, username)
    await load()
  }

  /**
   * Removes a collaborator from a task and reloads the store.
   */
  async function removeTaskCollaborator(id: number, collabId: number): Promise<void> {
    await tasksApi.removeCollaborator(id, collabId)
    await load()
  }

  /**
   * Resets the store state to its default empty values.
   */
  function reset(): void {
    tasks.value = []
    filter.value = 'all'
    searchQuery.value = ''
    selectedCategoryId.value = null
  }

  return {
    tasks,
    filter,
    searchQuery,
    filteredTasks,
    allFilteredTasks,
    selectedCategoryId,
    stats,
    completedCount,
    totalCount,
    load,
    add,
    toggle,
    editText,
    editDeadline,
    editCategory,
    remove,
    addTaskCollaborator,
    removeTaskCollaborator,
    reset,
    reorder
  }
})
