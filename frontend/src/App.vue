<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const router = useRouter()
const auth = useAuthStore()

onMounted(async () => {
  await auth.init()
  router.push(auth.isLoggedIn ? '/' : '/login')
})
</script>

<template>
  <RouterView />
  <footer class="flex justify-center shadow-md w-full border-t border-emerald-100 bg-emerald-50/30 text-slate-900 text-center text-sm py-4">
    <span>&copy; {{ new Date().getFullYear() }} Just Tasks. All rights reserved.</span>
  </footer>
</template>
