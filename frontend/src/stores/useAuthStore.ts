import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/types'
import * as authApi from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => user.value !== null)

  async function init(): Promise<void> {
    try {
      user.value = await authApi.fetchMe()
    } catch {
      user.value = null
    }
  }

  async function login(username: string, password: string): Promise<void> {
    user.value = await authApi.login(username, password)
  }

  async function register(username: string, password: string): Promise<void> {
    user.value = await authApi.register(username, password)
  }

  async function logout(): Promise<void> {
    await authApi.logout()
    user.value = null
  }

  async function deleteAccount(): Promise<void> {
    await authApi.deleteAccount()
    user.value = null
  }

  return { user, isLoggedIn, init, login, register, logout, deleteAccount }
})