<script setup lang="ts">
import { MagnifyingGlassIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'

const taskStore = useTaskStore()
</script>

<template>
  <div class="relative">
    <MagnifyingGlassIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-900 pointer-events-none" />

    <input
      v-model="taskStore.searchQuery"
      type="text"
      placeholder="Search tasks..."
      class="w-full border border-emerald-800 rounded-full text-xs pl-9 pr-8 py-1.5 text-black placeholder-slate-500 focus:outline-none focus:border-emerald-400 transition-colors"
    >

    <button
      v-if="taskStore.searchQuery"
      @click="taskStore.searchQuery = ''"
      class="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-700 hover:text-gray-900 transition-colors"
      aria-label="Clear search"
    >
      <XMarkIcon class="w-4 h-4" />
    </button>
  </div>
</template>
