<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Cog6ToothIcon, 
  ChartBarIcon, 
  MoonIcon, 
  SunIcon, 
  TagIcon, 
  PlusIcon, 
  XMarkIcon 
} from '@heroicons/vue/24/outline'
import { useAuthStore } from '@/stores/useAuthStore'
import { useTaskStore } from '@/stores/useTaskStore'
import { useCategoryStore } from '@/stores/useCategoryStore'
import TaskForm from '@/components/TaskForm.vue'
import type { FilterType } from '@/stores/useTaskStore'
import { onClickOutside, useDark, useToggle } from '@vueuse/core'
import SearchBar from '@/components/SearchBar.vue'
import StatsModal from '@/components/StatsModal.vue'
import TaskList from '@/components/TaskList.vue'

const router = useRouter()
const auth = useAuthStore()
const taskStore = useTaskStore()
const categoryStore = useCategoryStore()

const isDark = useDark()
const toggleDark = useToggle(isDark)

const showStats = ref(false)
const showSettings = ref(false)
const settings = ref(null)

const filters: FilterType[] = ['all', 'active', 'done', 'shared']

const showNewCategoryInput = ref(false)
const newCategoryName = ref('')

onMounted(() => {
  taskStore.load()
  categoryStore.load()
})

/**
 * Creates a new category, updates the store, and resets the input state.
 */
async function createCategory() {
  const name = newCategoryName.value.trim()
  if (name) {
    await categoryStore.add(name)
    newCategoryName.value = ''
    showNewCategoryInput.value = false
  }
}

/**
 * Clears local task state, signs out the user, and redirects to the login page.
 */
async function logout() {
  taskStore.reset()
  await auth.logout()
  router.push('/login')
}

/**
 * Prompts user confirmation to permanently delete the account and all associated tasks.
 */
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
  <div class="min-h-screen bg-emerald-50/30 text-black dark:bg-slate-950 dark:text-slate-100 transition-colors">
    <div class="sticky h-13 top-0 z-50 shadow-xs w-full border-b border-emerald-100 bg-[#f9fefb] dark:bg-slate-950 dark:border-slate-800 transition-colors">
      <header class="flex items-center justify-between max-w-5xl mx-auto px-4 py-1.5">
        <div>
          <h1 class="text-xl font-bold tracking-tight inline-flex items-center gap-1 cursor-pointer" @click="router.push('/dashboard')" >
            Just <span class="text-emerald-500 dark:text-emerald-400">Tasks</span>
          </h1>
        </div>
        
        <div class="flex gap-3">
          <button
            @click="showStats = true"
            class="m-2 rounded-xl text-slate-900 hover:text-emerald-600 dark:text-slate-300 dark:hover:text-emerald-400 active:scale-x-110 hover:cursor-pointer transform transition-transform"
          >
            <ChartBarIcon class="w-6 h-6" />
          </button>
          <StatsModal v-if="showStats" @close="showStats = false" />

          <button
            @click="toggleDark()"
            class="m-2 rounded-xl text-slate-900 hover:text-emerald-600 dark:text-slate-300 dark:hover:text-emerald-400 active:scale-110 hover:cursor-pointer transform transition-transform"
          >
            <SunIcon v-if="isDark" class="w-6 h-6" />
            <MoonIcon v-else class="w-6 h-6" />
          </button>

          <div class="relative" ref="settings">
            <button
              @click="showSettings = !showSettings"
              class="m-2 rounded-xl text-slate-950 hover:text-emerald-600 dark:text-slate-300 dark:hover:text-emerald-400 active:rotate-30 hover:cursor-pointer transform transition-transform"
            >
              <Cog6ToothIcon class="w-6 h-6" />
            </button>

            <Transition name="fade">
              <div
                v-if="showSettings"
                class="absolute right-0 bg-slate-50 top-10 w-50 border border-emerald-400 rounded-2xl shadow-2xl z-10 overflow-hidden dark:bg-slate-800 dark:border-slate-700"
              >
                <p class="px-4 py-2.5 text-sm text-gray-800 border-b border-emerald-400 font-medium dark:text-slate-200 dark:border-slate-700">
                  {{ auth.user?.username }}
                </p>
                <button
                  @click="logout"
                  class="w-full text-left px-4 py-2.5 hover:cursor-pointer text-sm text-slate-900 hover:bg-emerald-200 transition-all dark:text-slate-200 dark:hover:bg-slate-700"
                >
                  Sign out
                </button>
                <button
                  @click="deleteAccount"
                  class="w-full text-left px-4 py-2.5 hover:cursor-pointer text-sm text-rose-600 hover:bg-emerald-200 transition-colors dark:text-rose-400 dark:hover:bg-slate-700"
                >
                  Delete account
                </button>
              </div>
            </Transition>
          </div>
        </div>
      </header>
    </div>

    <div class="max-w-3xl mx-auto px-4 py-3">

      <TaskForm class="mb-3" />

      <div class="flex flex-wrap items-center gap-1 mb-1">
        <button
          @click="taskStore.selectedCategoryId = null"
          :class="[
            'px-3 py-1.5 rounded-xl text-xs font-medium border transition-colors flex items-center gap-1.5 cursor-pointer',
            taskStore.selectedCategoryId === null
              ? 'bg-emerald-600 text-white border-emerald-600 dark:bg-emerald-600 dark:border-emerald-600'
              : 'border-emerald-200 text-emerald-700 hover:bg-emerald-50 dark:border-slate-700/90 dark:text-emerald-500 dark:hover:bg-slate-800/50'
          ]"
        >
          All categories
        </button>

        <div class="flex items-center gap-2 overflow-x-auto no-scrollbar">
          <button
            v-for="category in categoryStore.categories"
            :key="category.id"
            @click="taskStore.selectedCategoryId = category.id"
            :class="[
              'px-3 py-1.5 rounded-xl text-xs font-medium border transition-colors flex items-center gap-1.5 whitespace-nowrap cursor-pointer group',
              taskStore.selectedCategoryId === category.id
                ? 'bg-emerald-100 text-emerald-800 border-emerald-300 dark:bg-emerald-900/50 dark:text-emerald-300 dark:border-emerald-700'
                : 'border-emerald-200 text-emerald-700 hover:bg-emerald-50 dark:border-slate-700/90 dark:text-emerald-500 dark:hover:bg-slate-800/50'
            ]"
          >
            <TagIcon class="w-3.5 h-3.5" />
            {{ category.name }}
            
            <span
              @click.stop="categoryStore.remove(category.id)" 
              class="ml-1 opacity-0 group-hover:opacity-100 hover:text-rose-500 transition-opacity"
              title="Delete category"
            >
              <XMarkIcon class="w-3.5 h-3.5" />
            </span>
          </button>
        </div>

        <div v-if="showNewCategoryInput" class="flex items-center gap-1">
          <input 
              v-model="newCategoryName" 
              @keydown.enter="createCategory"
              @keydown.esc="showNewCategoryInput = false"
              autoFocus
              type="text" 
              placeholder="Name..." 
              class="w-24 px-2 py-1 text-xs rounded-lg border border-emerald-300 focus:outline-none focus:border-emerald-500 dark:bg-slate-800 dark:border-slate-600 dark:text-white"
            >
            <button @click="createCategory" class="p-1 text-emerald-600 hover:text-emerald-700 cursor-pointer">
              <PlusIcon class="w-4 h-4" />
            </button>
            <button @click="showNewCategoryInput = false" class="p-1 text-slate-400 hover:text-slate-600 cursor-pointer">
              <XMarkIcon class="w-4 h-4" />
            </button>
        </div>
        <button
          v-else
          @click="showNewCategoryInput = true"
          class="px-2 py-1.5 rounded-xl text-xs font-medium border border-dashed border-emerald-300 text-emerald-600 hover:bg-emerald-50 transition-colors flex items-center gap-1 cursor-pointer dark:border-emerald-700 dark:text-emerald-500 dark:hover:bg-slate-800/50"
        >
          <PlusIcon class="w-3 h-3" /> New
        </button>
      </div>
      <div
        class="flex justify-between flex-wrap items-center pt-2 pb-1 mb-1 sticky top-12.5 z-10 bg-[#f9fefb] border-b border-emerald-100 dark:bg-slate-950 dark:border-slate-800 transition-colors"
      >
        <div class="flex gap-2 mb-1">
          <button
            v-for="filter in filters"
            :key="filter"
            @click="taskStore.filter = filter"
            :class="[
              'px-3 hover:cursor-pointer py-1 rounded-full text-xs font-medium border transition-all capitalize',
              taskStore.filter === filter
                ? (filter === 'shared'
                    ? 'bg-sky-600 border-sky-600 text-white dark:bg-sky-600 dark:border-sky-600'
                    : 'bg-emerald-600 border-emerald-600 text-white dark:bg-emerald-600 dark:border-emerald-600')
                : (filter === 'shared'
                    ? 'text-gray-800 hover:bg-sky-100 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-sky-800/40'
                    : 'text-gray-800 hover:bg-emerald-100 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-emerald-800/40')
            ]"
          >
            {{ filter }}
          </button>
        </div>
        <SearchBar />
        <p class="rounded-full text-xs py-1.5 px-3 border border-emerald-800 text-slate-900 dark:border-slate-600 dark:text-slate-300">
          {{ taskStore.completedCount }} / {{ taskStore.totalCount }} completed
        </p>
      </div>

      <TaskList />

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

.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
