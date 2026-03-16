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
</template>
