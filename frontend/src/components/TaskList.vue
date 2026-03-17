<script setup lang="ts">
import draggable from 'vuedraggable'
import { computed } from 'vue'
import { Bars3Icon } from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'
import TaskItem from '@/components/TaskItem.vue'
import type { Task } from '@/types'

const taskStore = useTaskStore()

const list = computed({
  get(): Task[] {
    return taskStore.filteredTasks
  },
  set(newOrder: Task[]) {
    taskStore.reorder(newOrder)
  },
})
</script>

<template>
  <draggable
    v-model="list"
    item-key="id"
    handle=".drag-handle"
    :animation="180"
    ghost-class="drag-ghost"
    chosen-class="drag-chosen"
    tag="ul"
    class="space-y-2"
  >
    <template #item="{ element }">
      <li class="flex items-center gap-2">
        <TaskItem :task="element" class="flex-1 min-w-0" />

        <button
          class="drag-handle flex-shrink-0 cursor-grab active:cursor-grabbing
                 text-slate-600 hover:text-slate-900 transition-colors
                 touch-none p-1 -ml-1"
          aria-label="Drag to reorder"
        >
          <Bars3Icon class="w-4 h-4" />
        </button>
      </li>
    </template>

    <template #footer>
      <li
        v-if="taskStore.filteredTasks.length === 0"
        class="text-center text-gray-700 py-16 text-sm select-none"
      >
        No tasks here.
      </li>
    </template>
  </draggable>
</template>

<style scoped>
:deep(.drag-ghost) {
  opacity: 0;
}

:deep(.drag-chosen) {
  box-shadow: 0 0 0 1px rgba(136, 223, 174, 0.5);
  border-radius: 1rem;
}
</style>
