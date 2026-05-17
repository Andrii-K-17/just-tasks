import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/types'
import * as authApi from '@/api/auth'

/**
 * Global store for user authentication and session management.
 */
export const useAuthStore = defineStore('auth', () => {
  /** Current authenticated user state. */
  const user = ref<User | null>(null)

  /** Computed flag to check if a user is currently logged in. */
  const isLoggedIn = computed(() => user.value !== null)

  /** Cache for the initialization request to prevent parallel duplicate calls. */
  let initPromise: Promise<void> | null = null

  /**
   * Initializes the session by fetching the current user profile.
   */
  async function init(): Promise<void> {
    if (initPromise) return initPromise

    initPromise = (async () => {
      try {
        user.value = await authApi.fetchMe()
      } catch {
        user.value = null
      } finally {
        initPromise = null
      }
    })()

    return initPromise
  }

  /**
   * Authenticates the user and sets the session state.
   */
  async function login(username: string, password: string): Promise<void> {
    user.value = await authApi.login(username, password)
  }

  /**
   * Registers a new user and automatically logs them in.
   */
  async function register(username: string, password: string): Promise<void> {
    user.value = await authApi.register(username, password)
  }

  /**
   * Terminates the current session and clears the store state.
   */
  async function logout(): Promise<void> {
    await authApi.logout()
    user.value = null
    initPromise = null
  }

  /**
   * Deletes the authenticated account and resets the store state.
   */
  async function deleteAccount(): Promise<void> {
    await authApi.deleteAccount()
    user.value = null
  }

  return {
    user,
    isLoggedIn,
    init,
    login,
    register,
    logout,
    deleteAccount
  }
})
