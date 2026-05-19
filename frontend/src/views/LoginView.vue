<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { MoonIcon, SunIcon } from '@heroicons/vue/24/outline'
import { useDark, useToggle } from '@vueuse/core'

const router = useRouter()
const auth = useAuthStore()
const route = useRoute()

const isDark = useDark()
const toggleDark = useToggle(isDark)

const mode = ref<'login' | 'register'>(
  route.query.mode === 'register' ? 'register' : 'login'
)
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

/**
 * Handles user authentication by executing either login or registration based on the current mode.
 */
async function submit() {
  error.value = ''
  loading.value = true
  try {
    if (mode.value === 'login') {
      await auth.login(username.value, password.value)
    } else {
      await auth.register(username.value, password.value)
    }
    router.push('/')
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Something went wrong'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen relative bg-emerald-50/30 flex items-center justify-center p-4 dark:bg-slate-950 transition-colors">
    
    <button
      @click="toggleDark()"
      class="absolute top-4 right-4 p-2 rounded-xl text-slate-900 hover:text-emerald-600 dark:text-slate-300 dark:hover:text-emerald-400 active:scale-110 hover:cursor-pointer transform transition-transform"
      aria-label="Toggle dark mode"
    >
      <SunIcon v-if="isDark" class="w-6 h-6" />
      <MoonIcon v-else class="w-6 h-6" />
    </button>

    <div class="w-full max-w-sm">

      <h1 class="text-2xl font-bold tracking-tight text-center gap-1 dark:text-slate-100">
        Just <span class="text-emerald-500 dark:text-emerald-400">Tasks</span>
      </h1>
      <p class="text-slate-600 text-sm text-center mb-8 dark:text-slate-400">
        Tasks made simple.
      </p>

      <div class="flex bg-white rounded-2xl p-1 mb-6 border border-emerald-200 dark:bg-slate-900 dark:border-slate-800 transition-colors">
        <button
          v-for="m in (['login', 'register'] as const)"
          :key="m"
          @click="mode = m; error = ''"
          :class="[
            'flex-1 py-1.5 rounded-xl text-sm font-medium hover:cursor-pointer transition-all capitalize',
            mode === m
              ? 'bg-emerald-600 text-white shadow dark:bg-emerald-500'
              : 'text-gray-600 border border-transparent hover:border-emerald-300 dark:text-slate-400 dark:hover:border-slate-700'
          ]"
        >
          {{ m === 'login' ? 'Sign In' : 'Register' }}
        </button>
      </div>

      <form @submit.prevent="submit" class="space-y-3">
        <div>
          <label class="block text-xs text-gray-600 mb-1 dark:text-slate-400">Username</label>
          <input
            v-model="username"
            type="text"
            required
            placeholder="Username"
            class="w-full bg-white border border-emerald-200 rounded-2xl px-3 py-2.5 text-black text-sm placeholder-gray-500 focus:outline-none focus:border-emerald-500 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:placeholder-slate-500 dark:focus:border-emerald-500"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-600 mb-1 dark:text-slate-400">Password</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="Password"
            class="w-full bg-white border border-emerald-200 rounded-2xl px-3 py-2.5 text-black text-sm placeholder-gray-500 focus:outline-none focus:border-emerald-500 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:placeholder-slate-500 dark:focus:border-emerald-500"
          />
        </div>

        <p v-if="error" class="text-rose-600 text-xs pt-1 dark:text-rose-400">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-emerald-600 hover:cursor-pointer hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold rounded-xl py-2.5 text-sm transition-colors mt-2 dark:bg-emerald-600 dark:hover:bg-emerald-500"
        >
          {{ loading ? 'Loading…' : mode === 'login' ? 'Sign In' : 'Create Account' }}
        </button>
      </form>

    </div>
  </div>
</template>
