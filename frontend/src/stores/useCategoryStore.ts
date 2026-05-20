import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Category } from '@/types'
import * as categoriesApi from '@/api/categories'

/**
 * Global store for category management.
 */
export const useCategoryStore = defineStore('category', () => {
  /** List of all categories for the current user. */
  const categories = ref<Category[]>([])

  /** Loading state for async operations. */
  const loading = ref(false)

  /** Error message for failed requests. */
  const error = ref<string>('')

  /** Computed flag to check if categories are loaded. */
  const hasCategories = computed(() => categories.value.length > 0)

  /**
   * Loads all categories from the backend.
   */
  async function load(): Promise<void> {
    loading.value = true
    try {
      categories.value = await categoriesApi.getCategories()
      error.value = ''
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  /**
   * Creates a new category and adds it to the store.
   */
  async function add(name: string): Promise<Category> {
    const newCategory = await categoriesApi.createCategory(name)
    categories.value.push(newCategory)
    return newCategory
  }

  /**
   * Deletes a category by ID and updates the store.
   */
  async function remove(id: number): Promise<void> {
    await categoriesApi.deleteCategory(id)
    categories.value = categories.value.filter(c => c.id !== id)
  }

  /**
   * Finds a category by its ID.
   */
  function getById(id: number | null): Category | undefined {
    if (!id) return undefined
    return categories.value.find(c => c.id === id)
  }

  return {
    categories,
    loading,
    error,
    hasCategories,
    load,
    add,
    remove,
    getById
  }
})
