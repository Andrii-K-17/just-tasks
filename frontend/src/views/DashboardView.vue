<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Cog6ToothIcon } from '@heroicons/vue/24/outline'
import { useAuthStore } from '@/stores/useAuthStore'
import { useTaskStore } from '@/stores/useTaskStore'
import TaskForm from '@/components/TaskForm.vue'
import TaskItem from '@/components/TaskItem.vue'
import type { FilterType } from '@/stores/useTaskStore'
import { onClickOutside } from '@vueuse/core'
import SearchBar from '@/components/SearchBar.vue'

const router = useRouter()
const auth = useAuthStore()
const taskStore = useTaskStore()

const showSettings = ref(false)
const settings = ref(null)
const filters: FilterType[] = ['all', 'active', 'done']

onMounted(() => taskStore.load())

async function logout() {
  taskStore.reset()
  await auth.logout()
  router.push('/login')
}

async function deleteAccount() {
  if (!confirm('Delete account and all tasks? This cannot be undone.')) return
  taskStore.reset()
  await auth.deleteAccount()
  router.push('/login')
}

onClickOutside(settings, () => {
  showSettings.value = false
})
</script>

<template>
  <div class="min-h-screen bg-emerald-50/30 text-black">
    <div class="sticky h-13 top-0 z-50 shadow-xs w-full border-b border-emerald-100 bg-[#f9fefb]">
      <header class="flex items-center justify-between max-w-3xl mx-auto px-4 py-1.5">
        <div>
          <h1 class="text-xl font-bold tracking-tight inline-flex items-center gap-1">
            Just <span class="text-emerald-500">Tasks</span>
          </h1>
        </div>

        <div class="relative" ref="settings">
          <button
            @click="showSettings = !showSettings"
            class="m-2 rounded-xl text-slate-900 active:rotate-30 hover:cursor-pointer transform transition-transform"
          >
            <Cog6ToothIcon class="w-6 h-6" />
          </button>

          <Transition name="fade">
            <div
              v-if="showSettings"
              class="absolute right-0 bg-slate-50 top-10 w-50 border border-emerald-400 rounded-2xl shadow-2xl z-10 overflow-hidden"
            >
              <p class="px-4 py-2.5 text-sm text-gray-800 border-b border-emerald-400 font-medium">
                {{ auth.user?.username }}
              </p>
              <button
                @click="logout"
                class="w-full text-left px-4 py-2.5 hover:cursor-pointer text-sm text-slate-900 hover:bg-emerald-200 transition-all"
              >
                Sign out
              </button>
              <button
                @click="deleteAccount"
                class="w-full text-left px-4 py-2.5 hover:cursor-pointer text-sm text-rose-600 hover:bg-emerald-200 transition-colors"
              >
                Delete account
              </button>
            </div>
          </Transition>
        </div>
      </header>
    </div>

    <div class="max-w-3xl mx-auto px-4 py-3">

      <TaskForm class="mb-1" />

      <div
        class="flex justify-between flex-wrap items-center pt-2 pb-1 mb-1 sticky top-12.5 z-10 bg-[#f9fefb] border-b border-emerald-100"
      >
        <div class="flex gap-2 mb-1">
          <button
            v-for="filter in filters"
            :key="filter"
            @click="taskStore.filter = filter"
            :class="[
              'px-3 hover:cursor-pointer py-1 rounded-full text-xs font-medium border transition-all capitalize',
              taskStore.filter === filter
                ? 'bg-emerald-600 border-emerald-600 text-white'
                : 'border-emerald-800 text-gray-800 hover:bg-emerald-100'
            ]"
          >
            {{ filter }}
          </button>
        </div>
        <SearchBar />
        <p class="rounded-full text-xs py-1.5 px-3 border border-emerald-800 text-slate-900">
          {{ taskStore.stats.done }} / {{ taskStore.stats.total }} completed
        </p>
      </div>

      <TransitionGroup name="list" tag="ul" class="space-y-2">
        <TaskItem
          v-for="task in taskStore.filteredTasks"
          :key="task.id"
          :task="task"
        />
        <li
          v-if="taskStore.filteredTasks.length === 0"
          key="empty"
          class="text-center text-gray-700 py-16 text-sm select-none"
        >
          No tasks here
        </li>
      </TransitionGroup>

    </div>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
    transition: opacity .35s ease;
}
.fade-enter-from, .fade-leave-to {
    opacity: 0;
}
.list-enter-active {
    transition: all .35s ease;
}
.list-leave-active {
    transition: all .35s ease;
}
.list-enter-from {
    opacity: 0;
    transform: translateX(-20px);
}
.list-leave-to {
    opacity: 0;
    transform: translateX(20px);
}
.list-move {
    transition: transform .35s ease;
}
.shadow-xs {
  box-shadow: 0 1.5px 1px rgba(130, 150, 133, 0.05);
}
</style>
