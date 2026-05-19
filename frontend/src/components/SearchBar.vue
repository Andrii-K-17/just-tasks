<script setup lang="ts">
import { MagnifyingGlassIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'

const taskStore = useTaskStore()
</script>

<template>
  <div class="relative">
    <MagnifyingGlassIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-900 dark:text-slate-400 pointer-events-none" />

    <input
      v-model="taskStore.searchQuery"
      type="text"
      placeholder="Search tasks..."
      class="w-full border border-emerald-800 rounded-full text-xs pl-9 pr-8 py-1.5 text-black placeholder-slate-500 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-600 dark:text-slate-100 dark:placeholder-slate-400 dark:focus:border-emerald-500"
    >

    <button
      v-if="taskStore.searchQuery"
      @click="taskStore.searchQuery = ''"
      class="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-700 hover:text-gray-900 transition-colors dark:text-slate-400 dark:hover:text-slate-200"
      aria-label="Clear search"
    >
      <XMarkIcon class="w-4 h-4" />
    </button>
  </div>
</template>
