<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const router = useRouter()
const auth = useAuthStore()

const route = useRoute()

const mode = ref<'login' | 'register'>(
  route.query.mode === 'register' ? 'register' : 'login'
)
const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

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
  <div class="min-h-screen bg-emerald-50/30 flex items-center justify-center p-4">
    <div class="w-full max-w-sm">

      <h1 class="text-2xl font-bold tracking-tight text-center gap-1">
        Just <span class="text-emerald-500">Tasks</span>
      </h1>
      <p class="text-slate-600 text-sm text-center mb-8">
        Tasks made simple.
      </p>

      <div class="flex bg-white rounded-2xl p-1 mb-6 border border-emerald-200">
        <button
          v-for="m in (['login', 'register'] as const)"
          :key="m"
          @click="mode = m; error = ''"
          :class="[
            'flex-1 py-1.5 rounded-xl text-sm font-medium hover:cursor-pointer transition-all capitalize',
            mode === m
              ? 'bg-emerald-600 text-white shadow'
              : 'text-gray-600 border border-transparent hover:border-emerald-300'
          ]"
        >
          {{ m === 'login' ? 'Sign In' : 'Register' }}
        </button>
      </div>

      <form @submit.prevent="submit" class="space-y-3">
        <div>
          <label class="block text-xs text-gray-600 mb-1">Username</label>
          <input
            v-model="username"
            type="text"
            required
            placeholder="Username"
            class="w-full bg-white border border-emerald-200 rounded-2xl px-3 py-2.5 text-black text-sm placeholder-gray-500 focus:outline-none focus:border-emerald-500 transition-colors"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-600 mb-1">Password</label>
          <input
            v-model="password"
            type="password"
            required
            placeholder="Password"
            class="w-full bg-white border border-emerald-200 rounded-2xl px-3 py-2.5 text-black text-sm placeholder-gray-500 focus:outline-none focus:border-emerald-500 transition-colors"
          />
        </div>

        <p v-if="error" class="text-rose-600 text-xs pt-1">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-emerald-600 hover:cursor-pointer hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold rounded-xl py-2.5 text-sm transition-colors mt-2"
        >
          {{ loading ? 'Loading…' : mode === 'login' ? 'Sign In' : 'Create Account' }}
        </button>
      </form>

    </div>
  </div>
</template>
