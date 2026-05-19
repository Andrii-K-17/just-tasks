<script setup lang="ts">
import { ref } from 'vue'
import { PlusIcon } from '@heroicons/vue/24/solid'
import { useTaskStore } from '@/stores/useTaskStore'
import { useTextareaAutosize } from '@vueuse/core'

const taskStore = useTaskStore()

const text = ref('')
const textareaRef = ref<HTMLTextAreaElement | null>(null)
useTextareaAutosize({ input: text, element: textareaRef })
const priority = ref<1 | 2 | 3>(2)
const deadline = ref('')
const error = ref('')

const priorityOptions = [
  { value: 1 as const, label: 'Low', color: 'text-gray-700 dark:text-gray-300' },
  { value: 2 as const, label: 'Med', color: 'text-yellow-700 dark:text-yellow-500' },
  { value: 3 as const, label: 'High', color: 'text-red-700 dark:text-red-400' },
]

/**
 * Validates form input, creates a new task via the store, and resets the form state.
 */
async function submit() {
  const trimmed = text.value.trim()
  if (!trimmed) return

  error.value = ''
  try {
    await taskStore.add({
      task_text: trimmed,
      priority: priority.value,
      deadline: deadline.value || null,
    })
    text.value = ''
    deadline.value = ''
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to add task'
  }
}

/**
 * Intercepts keyboard inputs to trigger form submission on Enter without the Shift modifier.
 */
function handleKeyDown(event: KeyboardEvent) {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    submit()
  }
}
</script>

<template>
  <form
    @submit.prevent="submit"
    class="bg-white border border-emerald-200 rounded-2xl p-4 space-y-3 dark:bg-slate-900 dark:border-slate-800 transition-colors"
  >
    <div class="flex gap-2 items-start">
      <textarea
        v-model="text"
        ref="textareaRef"
        rows="1"
        placeholder="New task..."
        maxlength="255"
        @keydown="handleKeyDown"
        class="flex-1 bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-2.5 text-sm text-slate-900 placeholder-gray-600 focus:outline-none focus:border-emerald-400 transition-colors resize-none overflow-hidden dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:placeholder-slate-400 dark:focus:border-emerald-500"
      ></textarea>
      <button
        type="submit"
        class="bg-emerald-500 hover:cursor-pointer hover:bg-emerald-600 text-white rounded-xl p-2 transition-colors dark:bg-emerald-600 dark:hover:bg-emerald-500"
        aria-label="Add task"
      >
        <PlusIcon class="w-6 h-6" />
      </button>
    </div>

    <div class="flex gap-2 flex-wrap">
      <div class="flex rounded-xl overflow-hidden">
        <button
          v-for="opt in priorityOptions"
          :key="opt.value"
          type="button"
          @click="priority = opt.value"
          :class="[
            'px-3 hover:cursor-pointer py-1 text-xs font-medium transition-colors',
            opt.color,
            priority === opt.value
              ? 'bg-emerald-300/80 rounded-xl border border-emerald-200 dark:bg-emerald-900/60 dark:border-emerald-700'
              : 'hover:border-emerald-200 rounded-xl border border-transparent dark:hover:border-slate-700'
          ]"
        >
          {{ opt.label }}
        </button>
      </div>

      <input
        v-model="deadline"
        type="date"
        class=" hover:cursor-pointer flex-1 min-w-0 bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors dark:bg-slate-800 dark:border-slate-700 dark:text-slate-100 dark:focus:border-emerald-500"
      >
    </div>

    <p v-if="error" class="text-rose-600 text-xs dark:text-rose-400">{{ error }}</p>
  </form>
</template>
