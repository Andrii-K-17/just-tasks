<script setup lang="ts">
import { ref, computed } from 'vue'
import { TrashIcon, CheckIcon } from '@heroicons/vue/24/outline'
import { useTaskStore } from '@/stores/useTaskStore'
import type { Task } from '@/types'
import { useTextareaAutosize } from '@vueuse/core'

const props = defineProps<{ task: Task }>()
const taskStore = useTaskStore()

const isEditing = ref(false)
const editValue = ref('')

const isEditingDeadline = ref(false)
const editDeadlineValue = ref('')

const textareaRef = ref<HTMLTextAreaElement | null>(null)
useTextareaAutosize({ element: textareaRef })

const PRIORITY: Record<number, string> = {
  1: 'bg-gray-500',
  2: 'bg-yellow-500',
  3: 'bg-red-500',
}

const today = new Date().toISOString().slice(0, 10)
const isOverdue = computed(
  () => props.task.deadline && !props.task.is_completed && props.task.deadline < today
)

function startEdit() {
  editValue.value = props.task.task_text
  isEditing.value = true
}

async function saveEdit() {
  const trimmed = editValue.value.trim()
  if (trimmed && trimmed !== props.task.task_text) {
    await taskStore.editText(props.task.id, trimmed)
  }
  isEditing.value = false
}

function startEditDeadline() {
  editDeadlineValue.value = props.task.deadline?? ''
  isEditingDeadline.value = true
}

async function saveEditDeadline() {
  const value = editDeadlineValue.value
  if (value !== props.task.deadline) {
    await taskStore.editDeadline(props.task.id, value)
  }
  isEditingDeadline.value = false
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    if (isEditing.value) saveEdit()
    if (isEditingDeadline.value) saveEditDeadline()
  }
  if (e.key === 'Escape') {
    isEditing.value = false
    isEditingDeadline.value = false
  }
}
</script>

<template>
  <li
    class="flex items-center gap-3 bg-white border border-emerald-200 rounded-2xl px-4 py-3 group hover:border-emerald-700 transition-colors"
  >
    <button
      @click="taskStore.toggle(task.id)"
      :class="[
        'w-5 h-5 rounded-full border-1 flex items-center justify-center flex-shrink-0 transition-all',
        task.is_completed
          ? 'bg-emerald-600 border-emerald-600'
          : 'border-emerald-900 hover:border-emerald-500'
      ]"
      :aria-label="task.is_completed ? 'Mark incomplete' : 'Mark complete'"
    >
      <CheckIcon v-if="task.is_completed" class="w-3 h-3 text-white" />
    </button>

    <span
      :class="['w-[0.55%] h-5 rounded-full flex-shrink-0', PRIORITY[task.priority]]"
    ></span>

    <div class="flex-1 min-w-0">
      <textarea
        v-if="isEditing"
        v-model="editValue"
        ref="textareaRef"
        @blur="saveEdit"
        @keydown="onKeydown"
        class="w-full bg-slate-50 rounded-lg px-2 py-0.5 text-sm text-slate-900 focus:outline-none focus:ring-1 focus:ring-emerald-300 resize-none overflow-hidden"
      ></textarea>
      <span
        v-else
        @click="startEdit"
        :class="[
          'text-sm cursor-text block truncate select-none',
          task.is_completed ? 'line-through text-gray-700' : 'text-gray-900'
        ]"
        :title="task.task_text"
      >
        {{ task.task_text }}
      </span>

      <input
        v-if="isEditingDeadline"
        v-model="editDeadlineValue"
        type="date"
        @blur="saveEditDeadline"
        @keydown="onKeydown"
        class="hover:cursor-pointer flex-1 min-w-0 bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors"
      >
      <span
        @click="startEditDeadline"
        v-if="task.deadline && !isEditingDeadline"
        :class="['text-xs block mt-0.5', isOverdue ? 'text-rose-600' : 'text-gray-700']"
      >
        {{ task.deadline }}
      </span>
      <button
        v-if="!task.deadline && !isEditingDeadline && isEditing"
        @click="startEditDeadline"
        class="hover:cursor-pointer flex-1 min-w-0 bg-emerald-50/30 border border-emerald-200 rounded-xl px-3 py-1 text-xs text-slate-900 focus:outline-none focus:border-emerald-400 transition-colors"
        aria-label="Set deadline"
      >
        Set deadline
      </button>
    </div>

    <button
      @click="taskStore.remove(task.id)"
      class="opacity-0 hover:cursor-pointer group-hover:opacity-100 text-gray-900 hover:text-rose-600 transition-all flex-shrink-0"
      aria-label="Delete task"
    >
      <TrashIcon class="w-5 h-5" />
    </button>
  </li>
</template>
